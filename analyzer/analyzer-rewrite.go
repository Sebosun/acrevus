package analyzer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RewriteResult struct {
	HTML     string
	Metadata Metadata
}

var ErrEmptyBody = errors.New("body is empty")

var defaultCandidates = "section,h2,h3,h4,h5,h6,p,td,pre"


func AnalyzerRewrite(document string) (RewriteResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return RewriteResult{}, err
	}

	result := RewriteResult{}

	metadata := GetMetadata(doc)
	result.Metadata = metadata

	body := doc.Find("body").First()

	if body.Length() == 0 {
		printH(body)
		printT(body)
		return RewriteResult{}, ErrEmptyBody
	}

	candidates := doc.Find(defaultCandidates)
	
	walk(body)

	return result, nil
}

func walk(s *goquery.Selection) []*goquery.Selection {
	elementsToScore := []*goquery.Selection{}

	s.Each(func(_ int, s *goquery.Selection) {
		nodeType := getNodeType(s)

		// nodes with too short of a text get excluded
		if len(s.Text()) < 25 { }
	})
	return elementsToScore
}

func isNodeACandidate(s *goquery.Selection) bool {
	text := s.Text()
	// nodes with too short of a text gets skipped
	if len(s.Text()) < 25 {
		return false
	}

	// elems with no parents get skipped
	if s.Parent().Length() == 0 {
		return false
	}

	contentScore := 0

	// Add a point for the paragraph itself as a base.
	contentScore += 1

	// Add points for any commas within this paragraph.
	contentScore += len(strings.Split(text, ","))

	switch {
	case len(text) >= 300:
		contentScore += 3
	case len(text) >= 200:
		contentScore += 2
	case len(text) >= 100:
		contentScore += 1
	}

	return true
}

func getNodeType(s *goquery.Selection) string {
	if len(s.Nodes) == 0 {
		return "invalid"
	}

	node := s.Nodes[0]
	return node.Data
}

func printH(s *goquery.Selection) {
	fmt.Println(s.Html())
}

func printT(s *goquery.Selection) {
	fmt.Println(s.Text())
}
