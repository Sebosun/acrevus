package analyzer

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type RewriteResult struct {
	HTML     string
	Metadata Metadata
}

type Candidate struct {
	selector *goquery.Selection
	score    float64
	isEmpty  bool
	depth    int
}

var (
	ErrEmptyBody = errors.New("body is empty")
	ErrTooShort  = errors.New("text is too short")
	ErrNoParents = errors.New("has no parents")
)

var (
	defaultCandidates  = "body,article,section,h2,h3,h4,h5,h6,p,td,pre"
	unlikelyCandidates = regexp.MustCompile(
		`(?i)-ad-|ai2html|banner|breadcrumbs|combx|comment|community|cover-wrap|disqus|extra|footer|gdpr|header|legends|menu|related|remark|replies|rss|shoutbox|sidebar|skyscraper|social|sponsor|supplemental|ad-break|agegate|pagination|pager|popup|yom-remote`,
	)
	okMaybeItsACandidate = regexp.MustCompile(`(?i)and|article|body|column|content|main|mathjax|shadow`)
	presentationalAttributes = []string{
    "align",
    "background",
    "bgcolor",
    "border",
    "cellpadding",
    "cellspacing",
    "frame",
    "hspace",
    "rules",
    "style",
    "valign",
    "vspace",
  }
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

	cleanScripts(doc)
	clearUnlikelyCandidates(doc)
	defaultCandidates := doc.Find(defaultCandidates)
	elementsToScore := selectionToSlice(defaultCandidates)
	elementsToScore = append(elementsToScore, selectionToSlice(defaultCandidates)...)

	// TODO: redistribute to garndparents
	candidates := decideWorthyCandidates(elementsToScore, 0)

	idx := getTopCandidate(candidates)
	if idx != -1 {
		result.HTML = candidates[idx].selector.Text()
	}

	for _, v := range candidates {
		fmt.Println(v.score, len(v.selector.Text()))
	}

	return result, nil
}

func decideWorthyCandidates(elementsToScore []*goquery.Selection, depth int) []Candidate {
	candidates := []Candidate{}

	for _, s := range elementsToScore {
		can, err := decideCandidate(s)
		if err != nil {
			continue
		}
		can.depth = depth

		candidates = append(candidates, can)

		children := selectionToSlice(can.selector.Children())
		result := decideWorthyCandidates(children, depth+1)

		for _, child := range result {

			divider := 0.0

			switch depth {
			case 0:
				divider = 1
			case 1:
				divider = 2
			default:
				divider = float64(depth * 3)
			}

			can.score += child.score / divider
		}

		candidates = append(candidates, result...)
	}

	return candidates
}

func getTopCandidate(candidates []Candidate) int {
	for idx, candidate := range candidates {
		result := getLinkDensity(candidate.selector)
		scoreAfterLinks := float64(candidate.score) * (1 - result)

		candidates[idx].score = scoreAfterLinks
	}

	slices.SortFunc(candidates, func(a Candidate, b Candidate) int {
		return int(a.score - b.score)
	})

	return 0
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

	contentScore := 0.0

	// Add a point for the paragraph itself as a base.
	contentScore += 1

	// Add points for any commas within this paragraph.
	contentScore += float64(len(strings.Split(text, ",")))

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

func cleanScripts(doc *goquery.Document) {
	doc.Find("script").Each(func(_ int, s *goquery.Selection) {
		s.Remove()
	})

	doc.Find("noscript").Each(func(_ int, s *goquery.Selection) {
		s.Remove()
	})
}
