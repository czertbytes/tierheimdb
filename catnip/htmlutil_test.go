package catnip

import (
	"strings"
	"testing"

	"code.google.com/p/go.net/html"
)

type ParseSelectorTest struct {
	s               string
	name, id, class string
}

var parseSelectorTests = []ParseSelectorTest{
	{"aaa", "aaa", "", ""},
	{"aaa#bbb", "aaa", "bbb", ""},
	{"aaa.bbb", "aaa", "", "bbb"},
	{"aaa#bbb.ccc", "aaa", "bbb", "ccc"},
}

func TestParseSelector(t *testing.T) {
	for _, test := range parseSelectorTests {
		n, i, c := parseSelector(test.s)
		if test.name != n || test.id != i || test.class != c {
			t.Errorf("Parsing selector '%s' failed! Name: %s [%s], Id: %s [%s], Class: %s [%s]", test.s, n, test.name, i, test.id, c, test.class)
		}
	}
}

type SelectTest struct {
	s        string
	selector []string
	nodes    []Node
}

type Node struct {
	name, data string
}

var selectTests = []SelectTest{
	{
		`<html><body><p>first</p><p>second</p></body></html>`,
		[]string{"html", "body", "p"},
		[]Node{{"p", "first"}, {"p", "second"}},
	},
	{
		`<html><body><p id="f">first</p><p class="s">second</p><p class="t">third</p></body></html>`,
		[]string{"html", "body", "p.s"},
		[]Node{{"p", "second"}},
	},
}

func TestSelect(t *testing.T) {
	for _, test := range selectTests {
		doc, _ := html.Parse(strings.NewReader(test.s))
		for i, n := range NodeSelect(doc, test.selector) {
			exp := test.nodes[i]
			name := n.Data
			data := n.FirstChild.Data
			if name != exp.name || data != exp.data {
				t.Errorf("Select failed! Got: %s Exp: %s", n, exp)
			}
		}
	}
}
