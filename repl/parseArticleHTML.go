package repl

import (
	"regexp"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type DisplayNode struct {
	ParentNode  atom.Atom
	TextContent string
	Children    []DisplayNode
}

func runTextParser(n *html.Node, depth int) DisplayNode {
	node, _ := textParser(n)

	for child := n.FirstChild; child != nil; child = child.NextSibling {
		childResult := runTextParser(child, depth+1)
		node.Children = append(node.Children, childResult)
	}

	return node
}

func textParser(node *html.Node) (DisplayNode, bool) {

	re := regexp.MustCompile(`\s{2,}`)

	switch node.Type {
	case html.TextNode:
		par := node.Parent.DataAtom
		stripped := re.ReplaceAll([]byte(node.Data), []byte(" "))
		return DisplayNode{
			ParentNode:  par,
			TextContent: string(stripped),
		}, true
	}
	return DisplayNode{}, false
}
