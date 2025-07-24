package repl

import (
	"regexp"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type DisplayNode struct {
	NodeType    atom.Atom
	TextContent string
	Children    []DisplayNode
}

func getParsedArticles(rawHTML string) (DisplayNode, error) {
	n, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return DisplayNode{}, err
	}

	nodes, _ := runTextParser(n, 0)

	return nodes, nil
}

func runTextParser(n *html.Node, depth int) (DisplayNode, bool) {
	ds, shouldAppend := textParser(n)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		childResult, shouldAppndChild := runTextParser(child, depth+1)
		if shouldAppndChild {
			ds.Children = append(ds.Children, childResult)
		}
	}

	shouldAppend = shouldAppend || len(ds.Children) > 0
	return ds, shouldAppend
}

func textParser(node *html.Node) (DisplayNode, bool) {
	// replace every two+ spaces with 1 space
	re := regexp.MustCompile(`\s{2,}`)

	switch node.Type {
	case html.TextNode:
		par := node.Parent.DataAtom
		stripped := re.ReplaceAll([]byte(node.Data), []byte(" "))
		return DisplayNode{
			NodeType:    par,
			TextContent: string(stripped),
		}, true
	case html.ElementNode:
		return DisplayNode{NodeType: node.DataAtom}, true
	}
	return DisplayNode{}, false
}
