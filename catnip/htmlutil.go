package catnip

import (
	"regexp"
	"strings"

	"code.google.com/p/go.net/html"
)

var (
	selectorsRegExp *regexp.Regexp
)

func init() {
	var err error
	selectorsRegExp, err = regexp.Compile(`([a-zA-Z0-9]+)(?:#([a-zA-Z0-9-_]+))?(?:\.([a-zA-Z0-9-_]+))?`)
	if err != nil {
		panic(err)
	}
}

func NodeAttribute(n *html.Node, attrName string) string {
	for _, attr := range n.Attr {
		if attr.Key == attrName {
			return attr.Val
		}
	}

	return ""
}

func NodeSelect(n *html.Node, path []string) []*html.Node {
	nodes := []*html.Node{}
	elementName, elementId, elementClass := parseSelector(path[0])
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.DataAtom.String() == elementName {
			var cId string
			var cClasses []string
			for _, attr := range c.Attr {
				if attr.Key == "id" {
					cId = attr.Val
				}
				if attr.Key == "class" {
					cClasses = strings.Split(attr.Val, " ")
				}
			}

			hasSameId := true
			if len(elementId) > 0 {
				if elementId != cId {
					hasSameId = false
				}
			}

			hasSameClass := true
			if len(elementClass) > 0 {
				hasSameClass = false
				for _, cClass := range cClasses {
					if elementClass == cClass {
						hasSameClass = true
					}
				}
			}

			if hasSameId == true && hasSameClass == true {
				if len(path) > 1 {
					nodes = append(nodes, NodeSelect(c, path[1:])...)
				} else {
					nodes = append(nodes, c)
				}
			}
		}
	}

	return nodes
}

func parseSelector(selector string) (string, string, string) {
	selectorsRes := selectorsRegExp.FindAllSubmatch([]byte(selector), -1)
	if len(selectorsRes) != 1 || len(selectorsRes[0]) != 4 {
		return "", "", ""
	}

	return string(selectorsRes[0][1]), string(selectorsRes[0][2]), string(selectorsRes[0][3])
}
