// Package rust contains shared behaviors for the rustbuddy commands
package rust

import (
	"os"
	"text/scanner"
)

// Source is a data structure representing the basic lexical structure of
// a rust source code file.
type Source struct {
	Funcs     []Fn
	RsStructs []RsStruct
	Enums     []Enum
	Traits    []Trait
	Tests     []Test
	TestBlock int
	UB        []Unsafe
}

// Span is the start end end location of a code block
type Span struct {
	Start scanner.Position
	End   scanner.Position
}

// Fn is a function in rust
type Fn struct {
	Span   Span
	Name   string
	Args   []string
	Return string
}

// Enum is an Enumeration of types in rust
type Enum struct {
	Span     Span
	Name     string
	Variants []string
}

// RsStruct is a data structure specific to rust source code. The awkward name
// is to avoid using a keyword
type RsStruct struct {
	Span    Span
	Name    string
	Methods []Fn
	Traits  []Trait
}

// Trait refers to Rust trait name. Currently, there is no use for the body of
// a trait definition.
type Trait struct {
	Span Span
	Name string
}

// Test refers to unit tests already within the source
type Test struct {
	Name string
	Span Span
}

// Unsafe are blocks of code marked unsafe
type Unsafe struct {
	Span Span
}

// Parse reads rust source code and does a simple lexical analysis
func Parse(f *os.File) (Source, error) {
	var s scanner.Scanner
	var src Source
	s.Init(f)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch s.TokenText() {
		case "!": // macros have completely unpredictable structure, so we need
			// to zip past them for sanity.
			collapseMacro(&s)
		case "#": // attribute
			attName := "#"
			for {
				c := s.Next()
				attName += string(c)
				if c == ']' {
					break
				}
			}
			switch attName {
			case "#[cfg(test)]":
				src.TestBlock = s.Pos().Line
			case "#[test]":
				t := capTest(&s)
				src.Tests = append(src.Tests, t)
			default:
				continue
			}
		// Detect trait and impl first because they can encapsulate other blocks
		case "trait":
			src.Traits = append(src.Traits, capTrait(&s))
		case "impl":
			capImpl(&src, &s)
		case "enum":
			src.Enums = append(src.Enums, capEnum(&s))
		case "struct":
			src.RsStructs = append(src.RsStructs, capStruct(&s))
		case "fn":
			fn, ubs := capFn(&s)
			src.Funcs = append(src.Funcs, fn)
			if len(ubs) > 0 {
				src.UB = append(src.UB, ubs...)
			}
		case "unsafe":
			src.UB = append(src.UB, capUB(&s))
		default:
			continue
		}
	}
	return src, nil
}
