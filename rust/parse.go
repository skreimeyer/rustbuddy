// Package rust contains shared behaviors for the rustbuddy commands
package rust

import (
	"os"
	"strings"
	"text/scanner"
)

// SyntaxTree is a data structure representing the basic lexical structure of
// a rust source code file.
type SyntaxTree struct {
	Funcs     []Fn
	Tests     []Test
	TestBlock int
}

// Fn is a function in rust
type Fn struct {
	Name   string
	Args   []string
	Return string
	Parent Method
	Line   int
}

// Method are all functions explicitly implemented for a particular struct
type Method struct {
	Struct string
	Trait  string
}

// Test refers to unit tests already within the source
type Test struct {
	Name string
	Line int
}

// Parse reads rust source code and does a simple lexical analysis
func Parse(f *os.File) (SyntaxTree, error) {
	var s scanner.Scanner
	var st SyntaxTree
	s.Init(f)
	depth := 0
	parent := Method{}
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch s.TokenText() {
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
				st.TestBlock = s.Pos().Line
			case "#[test]":
				localDepth := 0
				line := s.Pos().Line
				fnStr := ""
				foundBlock := false
				for {
					if localDepth == 0 && foundBlock == true {
						break
					}
					c := s.Next()
					fnStr += string(c)
					switch c {
					case '{':
						localDepth++
						foundBlock = true
					case '}':
						localDepth--
					default:
						continue
					}
				}
				st.Tests = append(st.Tests, Test{
					Name: parseTestName(fnStr),
					Line: line,
				})
			default:

			}

		case "trait":
			traitSig := ""
			for {
				c := s.Next()
				if c == '{' {
					depth++
					break
				}
				traitSig += string(c)
			}
			traitSig = strings.TrimSpace(traitSig)
			parent.Trait = traitSig
			parent.Struct = ""
		case "impl":
			methodSig := "impl"
			for {
				c := s.Next()
				if c == '{' {
					depth++
					break
				}
				methodSig += string(c)
			}
			parent = parseMethod(methodSig)
		case "fn":
			var fnSig string
			line := s.Pos().Line
			for {
				c := s.Next()
				if c == '{' {
					depth++
					break
				}
				if c == ';' {
					break
				}
				fnSig += string(c)
			}
			fnSig = strings.TrimSpace(fnSig)
			f := parseFn(fnSig, line)
			if depth > 0 {
				f.Parent = parent
			}
			st.Funcs = append(st.Funcs, f)
		case "{":
			depth++
		case "}":
			depth--

		default:
			continue
		}
	}
	return st, nil
}

// parseFn extracts the meaningful components from a function signature string.
func parseFn(s string, lineNum int) Fn {
	argBegin := strings.Index(s, "(")
	argEnd := strings.LastIndex(s, ")")
	retBegin := strings.LastIndex(s, "->")
	return Fn{
		Name:   s[0:argBegin],
		Args:   strings.Split(s[argBegin+1:argEnd], ","),
		Return: s[retBegin+3:],
		Line:   lineNum,
	}
}

// parseMethod takes a method string and breaks down the name of the parent
// struct and associated trait, if any.
func parseMethod(s string) Method {
	var m Method
	s = strings.TrimSpace(s)
	words := strings.Split(s, " ")
	if len(words) == 4 {
		m.Trait = words[1]
	}
	m.Struct = words[len(words)-1]
	return m
}

// parseTestName extracts the name of a test function
func parseTestName(s string) string {
	iFn := strings.Index(s, "fn")
	iArg := strings.Index(s, "(")
	return strings.TrimSpace(s[iFn+2 : iArg-1])
}
