// Package cmd contains all core components of rustbuddy.
/*
Copyright © 2019 Samuel Kreimeyer

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/skreimeyer/rustbuddy/rust"
	"github.com/spf13/cobra"
)

// mktestCmd represents the mktest command
var mktestCmd = &cobra.Command{
	Use:   "mktest [file]",
	Short: "Generate templates for table-based unit tests",
	Args:  cobra.MinimumNArgs(1),
	Long: `Mktest performs a simple lexical analysis of a rust code file and
	identifies signatures for functions and methods. Each becomes its own test
	function, which requires filling out one or more "Case" structs. The basic
	structure of the test is:
	for each case:
		assert(
			my_function(test_arguments) == what-I-want,
			show-this-note-about-the-testcase-on-error
			)`,
	Run: func(cmd *cobra.Command, args []string) {
		makeTest(args)
	},
}

var out string
var app bool

func init() {
	rootCmd.AddCommand(mktestCmd)
	mktestCmd.Flags().BoolVar(&app, "append", false, "Append the output to the source file")
	mktestCmd.Flags().StringVar(&out, "output", "", "Name of file to write output. Defaults to stdout")
}

func makeTest(args []string) {
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

	files := args
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
		if out != "" {
			destination, err = os.Create(out)
			if err != nil {
				fmt.Println("Unable to create file:", err)
				return
			}
		}
		if app == true {
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
		if val.Parent.Struct == "" && len(val.Parent.Trait) > 1 {
			ignore = append(ignore, i)
		}
	}
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
