package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/skreimeyer/rustbuddy/rust"
)

func main() {
	// template setup
	const tmpl = `
#[cfg(test)]
mod tests {
	use super::*;

	{{range skipTraits .Funcs | skipMain}}#[test]
	fn test_{{.Parent.Trait}}{{.Parent.Struct}}{{.Name}}() {
		{{if skipSelf .Args | len | ne 0 }}struct Input {
			{{range skipSelf .Args}}{{.}},
			{{end}}};{{end}}
		struct Output {
			R: {{.Return}}
		};
		struct Case { {{if ne .Parent.Struct ""}}
			obj: {{.Parent.Struct}}{{end}}
			inpt: Input,
			out: Output,
			comment: string,
		};
		// __TEST CASES GO HERE__
		let cases = vec![
			// FIXME
			// Case { {{if ne .Parent.Struct ""}}
			//	obj: {{.Parent.Struct}}{}{{end}}
			// 	inpt: Input {},
			// 	out: Output{},
			// 	comment: "",
			// },
		]
		// __END TEST CASES__
		for c in cases.iter() {
			assert!({{if ne .Parent.Struct ""}}c.obj.{{end}}{{.Name}}({{range skipSelf .Args}}c.inpt.{{stripType .}},{{end}}) == c.out, c.comment)
		}
	}
	{{end}}
}
`
	fmap := template.FuncMap{
		"stripType":  stripType,
		"skipSelf":   skipSelf,
		"skipMain":   skipMain,
		"skipTraits": skipTraits,
	}
	testTemp := template.Must(template.New("testTemp").Funcs(fmap).Parse(tmpl))
	flag.Usage = func() {
		fmt.Printf("Usage: mktest [-Options ...] [File1.rs ...FileN.rs]\n")
		flag.PrintDefaults()
	}

	// flag setup

	app := flag.Bool("a", false, "Append the output to the source file")
	out := flag.String("o", "", "Name of file to write output. Defaults to stdout")
	// This is a bad design.
	flag.Parse()

	// the meat

	files := flag.Args()
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Println("File Read error:", err)
			continue
		}
		tree, err := rust.Parse(f)
		if err != nil {
			fmt.Println("Parsing error:", err)
			continue
		}
		// This needs a redesign
		destination := os.Stdout
		if *out != "" {
			destination, err = os.Create(*out)
			if err != nil {
				fmt.Println("Unable to create file:", err)
				return
			}
		}
		if *app == true {
			b, _ := ioutil.ReadFile(fname) // this err is redundant
			destination.Write(b)
		}
		err = testTemp.Execute(destination, tree)
		if err != nil {
			fmt.Println("Template error:", err)
			continue
		}
	}
}

// skipSelf is a helper function to omit the &self argument in methods.
func skipSelf(args []string) []string {
	if len(args) > 0 && args[0] == "&self" {
		return args[1:]
	}
	return args
}

// stripType is a helper function that takes a function argument (ie,
//	name: type) and returns a string with only the argument name.
func stripType(s string) string {
	if strings.HasPrefix(s, "&self") || len(s) == 0 {
		return ""
	}
	i := strings.Index(s, ":")
	return strings.TrimSpace(s[:i])
}

// skipMain ignores the "main" function
func skipMain(funcs []rust.Fn) []rust.Fn {
	for i, val := range funcs {
		if val.Name == "main" {
			return append(funcs[:i], funcs[i+1:]...)
		}
	}
	return funcs
}

// skipTraits ignores those function which are only associated with trait
// definitions. There aren't meaningful tests that we can perform for them.
// without a LOT of introspection we don't want right now.
func skipTraits(funcs []rust.Fn) []rust.Fn {
	ignore := []int{}
	for i, val := range funcs {
		fmt.Println("Parent struct", val.Parent.Struct)
		fmt.Println("Parent Trait", val.Parent.Trait)
		if val.Parent.Struct == "" && len(val.Parent.Trait) > 1 {
			ignore = append(ignore, i)
		}
	}
	fmt.Println("ignore list:", ignore)
	if len(ignore) == 0 {
		return funcs
	}
	output := []rust.Fn{}
	for j, n := range ignore {
		if j == len(ignore)-1 {
			output = append(output, funcs[n+1:]...)
		}
		if n == 0 {
			continue
		}
		next := ignore[j+1]
		output = append(output, funcs[n+1:next]...)
	}
	return output
}
