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

var (
	ErrEmptyBody = errors.New("body is empty")
	ErrTooShort  = errors.New("text is too short")
	ErrNoParents = errors.New("has no parents")
)

var defaultCandidates = "article,section,h2,h3,h4,h5,h6,p,td,pre"

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

	defaultCandidates := doc.Find(defaultCandidates)
	elementsToScore := selectionToSlice(defaultCandidates)
	elementsToScore = append(elementsToScore, selectionToSlice(defaultCandidates)...)

	candidates := decideWorthyCandidates(elementsToScore)

	// printSliceSelection(candidates)

	idx := getTopCandidates(candidates)
	if idx != -1 {
		result.HTML = candidates[idx].selector.Text()
	}

	result.Metadata = GetMetadata(doc)

	fmt.Println(result.HTML)
	return result, nil
}

type Candidate struct {
	selector *goquery.Selection
	score    int
	isEmpty  bool
}

func decideWorthyCandidates(elementsToScore []*goquery.Selection) []Candidate {
	candidates := []Candidate{}

	for _, s := range elementsToScore {
		can, err := isNodeACandidate(s)
		if err != nil {
			// fmt.Println("Error - ", err.Error())
			continue
		}

		candidates = append(candidates, can)

		children := selectionToSlice(can.selector.Children())
		result := decideWorthyCandidates(children)
		candidates = append(candidates, result...)
	}

	return candidates
}

func getTopCandidates(candidates []Candidate) int {
	topIdx := -1
	topScore := 0

	for idx, v := range candidates {
		if v.score > topScore {
			topIdx = idx
		}
	}

	return topIdx
}

func isNodeACandidate(s *goquery.Selection) (Candidate, error) {
	text := s.Text()
	// nodes with too short of a text gets skipped
	if len(s.Text()) < 25 {
		return Candidate{}, ErrTooShort
	}

	// elems with no parents get skipped
	if s.Parent().Length() == 0 {
		return Candidate{}, ErrNoParents
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

	return Candidate{selector: s, score: contentScore, isEmpty: true}, nil
}

func selectionToSlice(s *goquery.Selection) []*goquery.Selection {
	acc := []*goquery.Selection{}

	s.Each(func(_ int, s *goquery.Selection) {
		acc = append(acc, s)
	})
	return acc
}

func getNodeType(s *goquery.Selection) string {
	if len(s.Nodes) == 0 {
		return "invalid"
	}

	node := s.Nodes[0]
	return node.Data
}

func printSliceSelection(items []Candidate) {
	for _, c := range items {
		fmt.Printf("Selector %s - score %d", getNodeType(c.selector), c.score)
		fmt.Printf("\n")
		fmt.Println(c.selector.Text())
	}
}

func printH(s *goquery.Selection) {
	fmt.Println(s.Html())
}

func printT(s *goquery.Selection) {
	fmt.Println(s.Text())
}
