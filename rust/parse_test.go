package rust

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestEnum(t *testing.T) {
	f, _ := os.Open("cases/sample_enum.rs")
	expected := Enum{
		Name: "FlashMessage",
		Variants: []string{
			"Success",
			"Warning{ category: i32, message: String }",
			"Error(String)",
		},
	}
	src, _ := Parse(f)
	if cmpall(src.Enums[0].Variants, expected.Variants) != true {
		got := strings.Join(src.Enums[0].Variants, ", ")
		t.Errorf("Invalid Enum parse. Values are the following:\n%s", got)
	}
	if src.Enums[0].Name != expected.Name {
		t.Errorf("Invalid Enum parse. Names do not match.\nexpected: %s\tgot: %s", expected.Name, src.Enums[0].Name)
	}

}

func TestStruct(t *testing.T) {
	f, _ := os.Open("cases/sample_struct.rs")
	expected := []string{
		"Person",
		"Nil",
		"Pair",
		"Point",
		"Rectangle",
	}
	src, _ := Parse(f)
	found := []string{}
	for _, n := range src.RsStructs {
		found = append(found, n.Name)
	}
	if cmpall(found, expected) != true {
		got := strings.Join(found, "|")
		t.Errorf("Invalid struct parse. Values are the following:\n%s", got)
	}
}

func TestUnsafe(t *testing.T) {
	f, _ := os.Open("cases/sample_unsafe.rs")
	expected := 1
	src, _ := Parse(f)
	found := len(src.UB)
	if found != expected {
		t.Errorf("Expected %d blocks. Found %d", expected, found)
	}
}

func TestFn(t *testing.T) {
	f, _ := os.Open("cases/sample_fn.rs")
	expectedNames := []string{
		"main",
		"is_divisible_by",
		"fizzbuzz",
		"fizzbuzz_to",
	}
	expectedReturns := []string{
		"",
		"bool",
		"()",
		"",
	}
	src, _ := Parse(f)
	foundNames := []string{}
	foundReturns := []string{}
	for _, n := range src.Funcs {
		foundNames = append(foundNames, n.Name)
		foundReturns = append(foundReturns, n.Return)
	}
	if cmpall(foundNames, expectedNames) != true {
		t.Errorf("Invalid fn parse. Names are the following:\n%v", foundNames)
	}
	if cmpall(foundReturns, expectedReturns) != true {
		t.Errorf("Invalid fn parse. Returns are the following:\n%v", foundReturns)
	}
	fmt.Println("last line")
}

func TestImpl(t *testing.T) {
	f, _ := os.Open("cases/sample_impl.rs")
	exMap := make(map[string][]string)
	fndMap := make(map[string][]string)
	exMap["Val"] = []string{"value"}
	exMap["GenVal"] = []string{"value"}
	src, _ := Parse(f)
	for _, n := range src.RsStructs {
		methods := []string{}
		for _, m := range n.Methods {
			methods = append(methods, m.Name)
		}
		fndMap[n.Name] = methods
	}
	for k, v := range exMap {
		if len(fndMap[k]) == 0 {
			t.Errorf("Invalid impl parse. Did not find: %s", k)
		}
		if cmpall(fndMap[k], v) != true {
			t.Errorf("invalid impl parse.\nSearch:%v\nFound:%v", v, fndMap[k])
		}
	}

}

func cmpall(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for _, m := range a {
		for j, n := range b {
			if m == n {
				b = append(b[:j], b[j+1:]...)
				break
			}
		}
	}
	if len(b) > 0 {
		return false
	}
	return true
}

func contains(a []string, b string) bool {
	for _, n := range a {
		if n == b {
			return true
		}
	}
	return false
}
