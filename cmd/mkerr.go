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
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var mkerrCmd = &cobra.Command{
	Use:   "mkerr [flags]",
	Short: "Generate a custom error for a single file",
	Long: `mkerr uses the file or module name to template out a custom error
	identical to that shown in the "Defining and Error Type" from the
	Rust by Example book. w written to stdout by default`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dest := ""
		src, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		if writeout == true {
			dest = args[0]
		}
		if name == "" {
			name = makeName(args[0]) + "Error"
		}
		makeErr(name, src, dest)
	},
}

var writeout bool
var name string

func init() {
	rootCmd.AddCommand(mkerrCmd)
	mkerrCmd.Flags().BoolVar(&writeout, "write", false, "write w in place of the existing file")
	mkerrCmd.Flags().StringVar(&name, "name", "", "name for custom error. Defaults to module name.")
}

// Templates a typical error declaration block. Uses bufio scanner because we
// actually need the content of comment lines, so we're not omitting any of the
// content of the original file.
func makeErr(errName string, source *os.File, outFile string) {
	const eTmpl = `
use std::error::Error;
use std::fmt;

#[derive(Debug)]
struct {{.E}} {
    message: String
}

impl {{.E}} {
    fn new(message: &str) -> {{.E}} {
        {{.E}}{message: message.to_string()}
    }
}

impl fmt::Display for {{.E}} {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f,"{}",self.message)
    }
}

impl Error for {{.E}} {
    fn description(&self) -> &str {
        &self.message
    }
}

`
	eTemp := template.Must(template.New("eTemp").Parse(eTmpl))

	type customErr struct {
		E string
	}
	ce := customErr{E: errName}
	writeComplete := false
	scn := bufio.NewScanner(source)
	var b bytes.Buffer
	for scn.Scan() {
		if strings.HasPrefix(scn.Text(), "//") {
			b.Write(scn.Bytes())
			b.Write([]byte("\n"))
			continue
		}
		if strings.HasPrefix(scn.Text(), "/*") {
			b.Write(scn.Bytes())
			b.Write([]byte{'\n'}) // avoid trimming newline
			for {
				scn.Scan()
				b.Write(scn.Bytes())
				b.Write([]byte("\n"))
				if strings.Contains(scn.Text(), "*/") {
					scn.Scan() // so we don't catch the trailing slash
					break
				}
			}
		}
		if writeComplete == false {
			err := eTemp.Execute(&b, ce)
			writeComplete = true
			if err != nil {
				fmt.Println("Failed to write error template:", err)
			}
		}
		b.Write(scn.Bytes())
		b.Write([]byte("\n"))
	}
	source.Close()
	if outFile != "" {
		ioutil.WriteFile(outFile, b.Bytes(), 0644)
	} else {
		os.Stdout.Write(b.Bytes())

	}
}

func makeName(fname string) string {
	result := path.Base(fname)
	if result == "mod.rs" {
		dir, _ := path.Split(fname)
		if dir == "" { // User supplied a file mod.rs in working dir
			wd, err := os.Getwd()
			if err != nil {
				result = "Custom"
			} else {
				result = path.Base(wd)
			}
		} else {
			dir = strings.Trim(dir, "/") // remove trailing slash
			i := strings.LastIndex(dir, "/")
			result = dir[i+1:]
		}

	}
	return strings.TrimRight(result, ".rs")
}
