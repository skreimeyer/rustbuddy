package rust

import (
	"strings"
	"text/scanner"
)

// auxilliary functions are at the bottom.

// capture a function body
func capFn(s *scanner.Scanner) (Fn, []Unsafe) {
	var UBs []Unsafe
	var fnBody string
	var fnSig string
	var sp Span
	sp.Start = s.Pos()
	for {
		c := s.Next()
		if c == ';' { // This is a function without a body.
			break
		}
		if c == '{' {
			fnBody = collapse(c, s)
			break
		}
		if c == '<' {
			collapse(c, s)
			continue
		}
		fnSig += string(c)
	}
	fnSig = strings.TrimSpace(fnSig)
	f := parseFnSig(fnSig)
	sp.End = s.Pos()
	f.Span = sp
	// Go back through the function body to check for `unsafe`. This really
	// should be refactored completely to find these blocks in the first pass.
	if len(fnBody) > 0 {
		fnBody += string(scanner.EOF) // yes, this is actually necessary
		var ubscan scanner.Scanner
		ubscan.Init(strings.NewReader(fnBody))
		for tok := ubscan.Scan(); tok != scanner.EOF; tok = ubscan.Scan() {
			if ubscan.TokenText() == "unsafe" {
				UBs = append(UBs, capUB(&ubscan))
			}
		}
	}
	return f, UBs
}

// parseFn extracts the meaningful components from a function signature string.
func parseFnSig(s string) Fn {
	argBegin := strings.Index(s, "(")
	argEnd := strings.LastIndex(s, ")")
	retBegin := strings.LastIndex(s, "->")
	r := ""
	if retBegin != -1 {
		//argEnd will get tripped up by a unit return "()" we fix it here
		argEnd = strings.LastIndex(s[:retBegin], ")")
		r = s[retBegin+3:]
	}
	return Fn{
		Span:   Span{},
		Name:   s[0:argBegin],
		Args:   strings.Split(s[argBegin+1:argEnd], ","),
		Return: r,
	}
}

// capture a trait definition body, ignoring child functions
func capTrait(s *scanner.Scanner) Trait {
	t := ""
	start := s.Pos()
	for {
		c := s.Next()
		if c == '<' {
			collapse(c, s)
			continue
		}
		if c == '{' || c == '(' {
			collapse(c, s)
			break
		}
		t += string(c)
	}
	end := s.Pos()
	spn := Span{
		Start: start,
		End:   end,
	}
	t = strings.TrimSpace(t)
	return Trait{
		Name: t,
		Span: spn,
	}
}

// Capture a test block, ignoring everything but the function name
func capTest(s *scanner.Scanner) Test {
	start := s.Pos()
	name := ""
	for {
		c := s.Next()
		if c == '(' {
			c = advTo('{', s)
			collapse(c, s)
			break
		}
		name += string(c)
	}
	name = strings.Split(name, "fn ")[1]
	end := s.Pos()
	spn := Span{
		Start: start,
		End:   end,
	}
	return Test{
		Name: name,
		Span: spn,
	}
}

// Capture a struct block, ignoring fields
func capStruct(s *scanner.Scanner) RsStruct {
	start := s.Pos()
	name := ""
	for {
		c := s.Next()
		if c == ';' {
			break
		}
		if c == '{' {
			collapse(c, s)
			break
		}
		if c == '(' || c == '<' {
			collapse(c, s)
			continue
		}
		name += string(c)
	}
	name = strings.TrimSpace(name)
	end := s.Pos()
	spn := Span{
		Start: start,
		End:   end,
	}
	return RsStruct{
		Span:    spn,
		Name:    name,
		Methods: []Fn{},
		Traits:  []Trait{},
	}
}

// Capture the enum block and variants
func capEnum(s *scanner.Scanner) Enum {
	start := s.Pos()
	vars := []string{}
	name := ""
	depth := 0
	for {
		c := s.Next()
		if c == '{' {
			depth++
			break
		}
		name += string(c)
	}
	name = strings.TrimSpace(name)
	endEnum := false
	for {
		variant := ""
		for {
			c := s.Next()
			if c == ',' {
				break
			}
			if c == '{' || c == '(' {
				variant += string(c)
				variant += collapse(c, s)
				break
			}
			if c == '}' {
				endEnum = true
				break
			}
			if c == '/' {
				advTo('\n', s)
				continue
			}
			variant += string(c)
		}
		variant = strings.TrimSpace(variant)
		if variant != "" {
			vars = append(vars, variant)
		}
		if endEnum == true {
			break
		}

	}
	end := s.Pos()
	spn := Span{
		Start: start,
		End:   end,
	}
	return Enum{
		Span:     spn,
		Name:     name,
		Variants: vars,
	}

}

// impl signatures can be highly varied. One-pass procedural handling does not
// seem to have an obvious, practical implementation.
func capImpl(src *Source, s *scanner.Scanner) {
	var (
		sig        string
		traitName  string
		structName string
	)
	for {
		c := s.Next()
		if c == '<' || c == '(' {
			collapse(c, s)
			continue
		}
		if c == '{' {
			break
		}
		sig += string(c)
	}
	if strings.Contains(sig, " for ") {
		parts := strings.Split(sig, " for ")
		i := strings.LastIndex(parts[0], " ")
		traitName = parts[0][i:]
		j := strings.IndexRune(parts[1], '<')
		if j == -1 {
			j = strings.IndexAny(parts[1], " \n")
		}
		structName = parts[1][:j]
	} else {
		parts := strings.Split(sig, ">")
		skip := strings.Count(parts[0], "<")
		structName = parts[skip]
		structName = strings.Split(structName, "<")[0]
		structName = strings.TrimSpace(structName)
	}
	// Index struct & trait within the existing parse tree and create new items
	// if they don't already exist
	if traitName != "" {
		exists := false
		for _, t := range src.Traits {
			if t.Name == traitName {
				exists = true
				break
			}
		}
		if !exists {
			newTrait := Trait{
				Name: traitName,
				Span: Span{},
			}
			src.Traits = append(src.Traits, newTrait)
		}
	}
	m := -1
	for j, s := range src.RsStructs {
		if s.Name == structName {
			m = j
			break
		}
	}
	if m == -1 {
		newStruct := RsStruct{
			Name:    structName,
			Span:    Span{},
			Methods: []Fn{},
			Traits:  []Trait{},
		}
		src.RsStructs = append(src.RsStructs, newStruct)
		m = len(src.RsStructs) - 1
	}
	// capture all child functions and append to methods array

	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		switch s.TokenText() {
		case "fn":
			f, ubs := capFn(s)
			src.RsStructs[m].Methods = append(src.RsStructs[m].Methods, f)
			if len(ubs) > 0 {
				src.UB = append(src.UB, ubs...)
			}
		case "}":
			goto stopCapture // There should be no nested brackets outside of functions
		}
	}
stopCapture:
}

func capUB(s *scanner.Scanner) Unsafe {
	var sp Span
	advTo('{', s)
	sp.Start = s.Pos()
	collapse('{', s)
	sp.End = s.Pos()
	return Unsafe{Span: sp}
}

func collapse(current rune, s *scanner.Scanner) string {
	var content string
	left := current
	right := '}'
	switch left {
	case '(':
		right = ')'
	case '<':
		right = '>'
	default:
	}
	open := 1
	for {
		c := s.Next()
		content += string(c)
		if c == right {
			open--
		}
		if c == left {
			open++
		}
		if open == 0 {
			break
		}
	}
	return content
}

// Call from exclamation point. Will peek and then advance to first opening block
// which may be ( or {, then calls collapse.
func collapseMacro(s *scanner.Scanner) {
cmTop:
	switch s.Peek() {
	case '=':
		break // false alarm
	case '{':
		s.Next()
		collapse('{', s)
		break
	case '(':
		s.Next()
		collapse('(', s)
		break
	case ' ', '\n':
		s.Next()
		goto cmTop
	default: // some kind of char. ! is for negation and not a macro
		break
	}
	return
}

func advTo(target rune, s *scanner.Scanner) rune {
	var c rune
	for {
		c := s.Next()
		if c == target {
			break
		}
	}
	return c
}
