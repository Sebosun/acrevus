package analyzer

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

func (a *AnalyzerGoquery) getLinkDensity(s *goquery.Selection) float64 {
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

func (a *AnalyzerGoquery) getNodeTag(s *goquery.Selection) string {
	if len(s.Nodes) == 0 {
		return "invalid"
	}

	node := s.Nodes[0]
	return node.Data
}

func (a *AnalyzerGoquery) hasAncestorTag(s *goquery.Selection, tag string) bool {
	maxDepth := 3

	cur := s
	for range maxDepth {
		parent := cur.Parent()

		// means element is empty == parent doesnt exist
		if parent.Length() == 0 {
			return false
		}
		tagName := a.getNodeTag(s.Parent())
		if tagName == tag {
			return true
		}

		cur = s.Parent()
	}

	return false
}

func (a *AnalyzerGoquery) printSliceSelection(items []Candidate) {
	for _, c := range items {
		fmt.Printf("Selector %s - score %.2f ", a.getNodeTag(c.selector), c.score)
		fmt.Println(c.selector.Text())
		fmt.Printf("\n")
	}
}

func (a *AnalyzerGoquery) printH(s *goquery.Selection) {
	fmt.Println(s.Html())
}

func (a *AnalyzerGoquery) printT(s *goquery.Selection) {
	fmt.Println(s.Text())
}

func (a *AnalyzerGoquery) selectionToSlice(s *goquery.Selection) []*goquery.Selection {
	acc := []*goquery.Selection{}

	s.Each(func(_ int, s *goquery.Selection) {
		acc = append(acc, s)
	})
	return acc
}
