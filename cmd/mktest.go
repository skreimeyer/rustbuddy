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
	const mktestTemplate = `
#[cfg(test)]
mod tests {
	use super::*;

	// generated code. Edit only test cases!
	// functions
	{{range .Funcs | skipMain}}#[test]
	fn test_{{.Name}}() {
		{{if len .Args | ne 0 }}#[derive(PartialEq)]
		struct Input {
			{{range .Args}}{{.}},
			{{end}}};{{end}}
		#[derive(PartialEq)]
		struct Output {
			r: {{orUnit .Return}},
		};
		#[derive(PartialEq)]
		struct Case { 
			input: Input,
			out: Output,
			comment: String,
		};
		// start test cases
		for c in vec![
			// make your test cases here
			// Case {
			// 	input: Input {},
			// 	out: Output{r: },
			// 	comment: String::from(""),
			// },
		]
		// end of test cases
		.into_iter() { {{$total := skipSelf .Args | len | lessOne}}
			assert!(
				Output{r: {{.Name}}({{range $i,$x := .Args}}c.input.{{stripType $x}}{{if ne $i $total}}, {{end}}{{end}}) }== c.out, c.comment
		)
		}
	}
	{{end}}
	{{if len .RsStructs | ne 0}}// methods{{end}}
	{{range .RsStructs}}{{$parent := .Name}}{{range .Methods}}#[test]
	fn test_{{$parent}}_{{.Name}}() {
		{{if skipSelf .Args | len | ne 0 }}#[derive(PartialEq)]struct Input {
			{{range skipSelf .Args}}{{.}},
			{{end}}};{{end}}
		#[derive(PartialEq)]
		struct Output {
			r: {{orUnit .Return}}
		};
		#[derive(PartialEq)]
		struct Case { 
			obj:		{{$parent}},{{if skipSelf .Args | len | ne 0 }}
			input:		Input,{{end}}
			out:		Output,
			comment:	String,
		};
		// __TEST CASES GO HERE__
		let cases = vec![
			// FIXME
			// Case {
			//	obj: 	{{$parent}}{},{{if skipSelf .Args | len | ne 0 }}
			//	input:	Input{},{{end}}
			// 	out: 	Output{r: },
			// 	comment:String::from(""),
			// },
		]
		// __END TEST CASES__
		.into_iter() { {{$total := skipSelf .Args | len | lessOne}}
				assert!(
					Output{r: c.obj.{{.Name}}({{range $i, $x := skipSelf .Args}}c.input.{{stripType $x}}{{if ne $i $total}}, {{end}}{{end}}) == c.out, c.comment
				)
			}
		}

	}
	{{end}}{{end}}}//End generated code
`
	fmap := template.FuncMap{
		"stripType": stripType,
		"skipSelf":  skipSelf,
		"skipMain":  skipMain,
		"lessOne":   lessOne,
		"orUnit":    orUnit,
	}
	testTemp := template.Must(template.New("testTemp").Funcs(fmap).Parse(mktestTemplate))

	files := args
	for _, fname := range files {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Println("File Read error:", err)
			continue
		}
		source, err := rust.Parse(f)
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
			f.Close()
			destination, err = os.OpenFile(fname, os.O_RDWR|os.O_APPEND, 0660)
			if err != nil {
				fmt.Println("Cannot write to source file:", err)
				return
			}
		}
		err = testTemp.Execute(destination, source)
		if err != nil {
			fmt.Println("Template error:", err)
			continue
		}
		destination.Close()
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
	if i == -1 {
		return strings.TrimSpace(s)
	}
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

func lessOne(i int) int {
	return i - 1
}

func orUnit(s string) string {
	if s == "" {
		return "()"
	}
	return s
}
