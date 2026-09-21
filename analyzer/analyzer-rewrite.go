package analyzer

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RewriteResult struct {
	HTML     string
	Metadata Metadata
}

type Candidate struct {
	selector *goquery.Selection
	score    int
	isEmpty  bool
}

var (
	ErrEmptyBody = errors.New("body is empty")
	ErrTooShort  = errors.New("text is too short")
	ErrNoParents = errors.New("has no parents")
)

var (
	defaultCandidates  = "article,section,h2,h3,h4,h5,h6,p,td,pre"
	unlikelyCandidates = regexp.MustCompile(
		`(?i)-ad-|ai2html|banner|breadcrumbs|combx|comment|community|cover-wrap|disqus|extra|footer|gdpr|header|legends|menu|related|remark|replies|rss|shoutbox|sidebar|skyscraper|social|sponsor|supplemental|ad-break|agegate|pagination|pager|popup|yom-remote`,
	)
	okMaybeItsACandidate = regexp.MustCompile(`(?i)and|article|body|column|content|main|mathjax|shadow`)
)

func AnalyzerRewrite(document string) (RewriteResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return RewriteResult{}, err
	}

	result := RewriteResult{}
	result.Metadata = GetMetadata(doc)

	metadata := GetMetadata(doc)
	result.Metadata = metadata

	body := doc.Find("body").First()

	if body.Length() == 0 {
		return RewriteResult{}, ErrEmptyBody
	}

	clearUnlikelyCandidates(doc)
	defaultCandidates := doc.Find(defaultCandidates)
	elementsToScore := selectionToSlice(defaultCandidates)
	elementsToScore = append(elementsToScore, selectionToSlice(defaultCandidates)...)

	// TODO: redistribute to garndparents
	candidates := decideWorthyCandidates(elementsToScore)

	idx := getTopCandidate(candidates)
	if idx != -1 {
		result.HTML = candidates[idx].selector.Text()
	}
	return result, nil
}

func decideWorthyCandidates(elementsToScore []*goquery.Selection) []Candidate {
	candidates := []Candidate{}

	for _, s := range elementsToScore {
		can, err := decideCandidate(s)
		if err != nil {
			continue
		}

		candidates = append(candidates, can)

		children := selectionToSlice(can.selector.Children())
		result := decideWorthyCandidates(children)
		candidates = append(candidates, result...)
	}

	return candidates
}

func getTopCandidate(candidates []Candidate) int {
	topIdx := -1
	topScore := 0.0

	for idx, candidate := range candidates {
		result := getLinkDensity(candidate.selector)
		scoreAfterLinks := float64(candidate.score) * (1 - result)
		if scoreAfterLinks > topScore {
			topScore = scoreAfterLinks
			topIdx = idx
		}
	}

	return topIdx
}

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

func decideCandidate(s *goquery.Selection) (Candidate, error) {
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

func getNodeTag(s *goquery.Selection) string {
	if len(s.Nodes) == 0 {
		return "invalid"
	}

	node := s.Nodes[0]
	return node.Data
}

func clearUnlikelyCandidates(doc *goquery.Document) {
	doc.Find("[class], [id]").Each(func(_ int, s *goquery.Selection) {
		className, _ := s.Attr("class")
		id, _ := s.Attr("id")

		value := strings.ToLower(className + " " + id)

		isUnlikelyClass := unlikelyCandidates.Match([]byte(value))
		mightBeCandidate := okMaybeItsACandidate.Match([]byte(value))

		tagName := getNodeTag(s)

		hasCode := hasAncestorTag(s, "code")
		hasTable := hasAncestorTag(s, "table")

		if tagName != "a" && tagName != "body" && isUnlikelyClass && !mightBeCandidate && !hasCode && !hasTable {
			s.Remove()
		}
	})
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
