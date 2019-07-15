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
	"fmt"
	"html/template"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

var mkerrCmd = &cobra.Command{
	Use:   "mkerr",
	Short: "Generate a custom error for a single file",
	Long: `mkerr uses the file or module name to template out a custom error
	identical to that shown in the "Defining and Error Type" from the
	Rust by Example book. Output written to stdout by default`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dest := os.Stdout
		src, err := os.Open(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		if writeout == true {
			dest, err = os.Create(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
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
	mkerrCmd.Flags().BoolVar(&writeout, "write", false, "write output in place of the existing file")
	mkerrCmd.Flags().StringVar(&name, "name", "", "name for custom error. Defaults to module name.")
}

func makeErr(errName string, source *os.File, output *os.File) {
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
}`
	eTemp := template.Must(template.New("eTemp").Parse(eTmpl))

	type customErr struct {
		E string
	}
	ce := customErr{E: errName}
	writeComplete := false
	scn := bufio.NewScanner(source)
	for scn.Scan() {
		if strings.HasPrefix(scn.Text(), "//") {
			output.Write(scn.Bytes())
			output.Write([]byte("\n"))
			continue
		}
		if strings.HasPrefix(scn.Text(), "/*") {
			for {
				scn.Scan()
				output.Write(scn.Bytes())
				output.Write([]byte("\n"))
				if strings.Contains(scn.Text(), "*/") {
					break
				}
			}
		}
		if writeComplete == false {
			err := eTemp.Execute(output, ce)
			writeComplete = true
			if err != nil {
				fmt.Println("Failed to write error template:", err)
			}
		}
		output.Write(scn.Bytes())
		output.Write([]byte("\n"))
	}
	source.Close()
	output.Close()
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
