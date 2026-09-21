package analyzer

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func getLinkDensity(s *goquery.Selection) float64 {
	text := s.Text()
	if len(text) == 0 {
		return 0
	}

	textLength := float64(len(text))
	linksLength := 0.0

	s.Find("a").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if ok {
			linksLength += float64(len(href))
		}
	})

	if linksLength <= 0.0 {
		return 0.0
	}

	return textLength / linksLength
}

func getNodeTag(s *goquery.Selection) string {
	if len(s.Nodes) == 0 {
		return "invalid"
	}

	node := s.Nodes[0]
	return node.Data
}

func hasAncestorTag(s *goquery.Selection, tag string) bool {
	maxDepth := 3

	cur := s
	for range maxDepth {
		parent := cur.Parent()

		// means element is empty == parent doesnt exist
		if parent.Length() == 0 {
			return false
		}
		tagName := getNodeTag(s.Parent())
		if tagName == tag {
			return true
		}

		cur = s.Parent()
	}

	return false
}

func printSliceSelection(items []Candidate) {
	for _, c := range items {
		fmt.Printf("Selector %s - score %d ", getNodeTag(c.selector), c.score)
		fmt.Println(c.selector.Text())
		fmt.Printf("\n")
	}
}

func printH(s *goquery.Selection) {
	fmt.Println(s.Html())
}

func printT(s *goquery.Selection) {
	fmt.Println(s.Text())
}

func selectionToSlice(s *goquery.Selection) []*goquery.Selection {
	acc := []*goquery.Selection{}

	s.Each(func(_ int, s *goquery.Selection) {
		acc = append(acc, s)
	})
	return acc
}

