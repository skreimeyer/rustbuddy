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
	"github.com/spf13/cobra"
)

var mkerrCmd = &cobra.Command{
	Use:   "mkerr",
	Short: "Generate a custom error for a single file",
	Long: `mkerr uses the file or module name to template out a custom error
	identical to that shown in the "Defining and Error Type" from the
	Rust by Example book. Output written to stdout by default`,
	Run: func(cmd *cobra.Command, args []string) {
		makeErr()
	},
}

var writeout bool

func init() {
	rootCmd.AddCommand(mkerrCmd)
	mkerrCmd.Flags().BoolVar(&writeout, "write", false, "write into the existing file")
}

func makeErr() {
	const tmpl = `
use std::error::Error;
use std::fmt;

#[derive(Debug)]
struct {{.name}} {
    message: String
}

impl {{.name}} {
    fn new(message: &str) -> {{.name}} {
        {{.name}}{message: message.to_string()}
    }
}

impl fmt::Display for {{.name}} {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f,"{}",self.message)
    }
}

impl Error for {{.name}} {
    fn description(&self) -> &str {
        &self.message
    }
}`
}
