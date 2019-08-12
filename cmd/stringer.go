package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"text/template"

	"github.com/skreimeyer/rustbuddy/rust"
	"github.com/spf13/cobra"
)

// stringerCmd represents the stringer command
var stringerCmd = &cobra.Command{
	Use:   "stringer [FLAGS] [SOURCE FILE] [ENUMS, ...]",
	Short: "Create a string representation method for enums",
	Long: `Stringer automates creating simple string representations of your enum
	types. This is done by implementing the following methods and trait:
	
	to_str()
	to_string()
	Display

	The string will be the name of the particular variant, and nothing more. No
	attempt is made to infer meaning from more complex types. The following would
	be the results from this example enum definition:

	enum MyEnum{
		First,
		Second(i32),
		Third{a:i32, b:char c:Something::Complicated},
	}

	MyEnum::First.to_str() == "First" // &str
	MyEnum::Second.to_string() == "Second" // String
	println!("{}",MyEnum::Third) == "Third" // fmt::Result
	`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		stringify(args)
	},
}

var allEnum bool

func init() {
	rootCmd.AddCommand(stringerCmd)
	stringerCmd.Flags().BoolVar(&allEnum, "all", false, "impl to_string for all enums")
	stringerCmd.Flags().BoolVar(&writeout, "write", false, "write output into source file")
}

func stringify(args []string) {
	const stringerTemplate = `
//GENERATED CODE DO NOT EDIT
const {{.Name}}_STR: &str = "{{$c := concat .Variants}}{{$c}}";
impl {{.Name}} {
	fn to_str(&self) -> &str {
		match &self {
			{{$e := .Name}}{{range $_,$v := .Variants}}{{$e}}::{{$v}} => &{{$e}}_STR{{slicer $c $v}},
			{{end}}
		}
	}

	fn to_string(&self) -> String {
		String::from(self.to_str())
	}
}

use std::fmt;
impl fmt::Display for {{.Name}} {
	fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
		write!(f, "{}", self.to_str())
	}
}
// END GENERATED CODE`
	fmap := template.FuncMap{
		"concat":    concat,
		"slicer":    slicer,
		"stripTail": stripTail,
		"plusOne":   plusOne,
	}
	var buf bytes.Buffer
	var q enumQueue
	tmpl := template.Must(template.New("stringerTemplate").Funcs(fmap).Parse(stringerTemplate))
	f, err := os.Open(args[0])
	if err != nil {
		fmt.Println("Can't open", args[0], err)
		return
	}
	defer f.Close()
	src, err := rust.Parse(f)
	if err != nil {
		fmt.Println("Cannot parse", args[0], err)
		return
	}
	if allEnum == true {
		q = src.Enums
	} else {
		for _, n := range args[1:] {
			for _, m := range src.Enums {
				if m.Name == n {
					q = append(q, m)
				}
			}
		}

	}
	sort.Sort(enumQueue(q))
	fBytes, err := ioutil.ReadFile(args[0]) // this is a poor implementation
	if err != nil {
		fmt.Println(err)
		return
	}
	lastByte := 0
	for _, e := range q {
		end := e.Span.End.Offset
		buf.Write(fBytes[lastByte:end])
		err = tmpl.Execute(&buf, e)
		if err != nil {
			fmt.Println("template error:", err)
			return
		}
		lastByte = end
	}
	buf.Write(fBytes[lastByte:])
	f.Close()
	if writeout == true {
		err := ioutil.WriteFile(args[0], buf.Bytes(), 0655)
		if err != nil {
			fmt.Println("failed to write", err)
			return
		}
		return
	}
	w := bufio.NewWriter(os.Stdout)
	w.Write(buf.Bytes())
	w.Flush()
	return
}

func concat(s []string) string {
	output := ""
	for _, a := range s {
		output += stripTail(a)
	}
	return output
}

func stripTail(s string) string {
	i := strings.IndexAny(s, ":{(")
	if i == -1 {
		return s
	}
	return s[:i]
}

func plusOne(i int) int {
	return i + 1
}

func slicer(s string, v string) string {
	v = stripTail(v)
	i := strings.Index(s, v)
	j := i + len(v)
	return fmt.Sprintf("[%d..%d]", i, j)
}

type enumQueue []rust.Enum

func (q enumQueue) Len() int      { return len(q) }
func (q enumQueue) Swap(i, j int) { q[i], q[j] = q[j], q[i] }
func (q enumQueue) Less(i, j int) bool {
	return q[i].Span.End.Offset < q[j].Span.End.Offset
}
