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
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/skreimeyer/rustbuddy/rust"
	"github.com/spf13/cobra"
)

// vetCmd represents the vet command
var vetCmd = &cobra.Command{
	Use:   "vet [FILENAME]",
	Short: "Summarize unsafe blocks in a file or crate",
	Long: `Vet identifies blocks of "unsafe" rust and gives a summary of how
much of a file is marked "unsafe." Vet is an auditing tool to get a simplistic 
look at potential risk areas in code.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runVet(args)
	},
}

func init() {
	rootCmd.AddCommand(vetCmd)
}

func runVet(args []string) {
	for _, fname := range args {
		f, err := os.Open(fname)
		if err != nil {
			fmt.Println("cannot open file", err)
			return
		}
		src, err := rust.Parse(f)
		if err != nil {
			fmt.Println("cannot parse file", err)
			return
		}
		f.Close()
		uLoc := 0
		for _, u := range src.UB {
			uLoc += u.Span.End.Line - u.Span.Start.Line
		}
		fAgain, _ := os.Open(fname)
		sLoc, err := lineCounter(fAgain)
		if err != nil {
			fmt.Println("error getting linecount", err)
			return
		}
		upc := float64(uLoc) / float64(sLoc) * 100.0
		fmt.Printf("%s summary:\nSLOC:%d\tunsafe:%.2f%%\tunsafe blocks:%d\n",
			fname,
			sLoc,
			upc,
			len(src.UB),
		)
	}
}

func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
