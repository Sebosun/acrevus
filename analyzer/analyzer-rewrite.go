package analyzer

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

type RewriteResult struct {
	HTML     string
	Metadata Metadata
}

type Candidate struct {
	selector    *goquery.Selection
	score       float64
	isEmpty     bool
	depth       int
	fingerprint string
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
	okMaybeItsACandidate     = regexp.MustCompile(`(?i)and|article|body|column|content|main|mathjax|shadow`)
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
	positive = regexp.MustCompile(`(?i)article|body|content|entry|hentry|h-entry|main|page|pagination|post|text|blog|story`)
	negative = regexp.MustCompile(`(?i)-ad-|hidden|^hid$| hid$| hid |^hid |banner|combx|comment|com-|contact|footer|gdpr|masthead|media|meta|outbrain|promo|related|scroll|share|shoutbox|sidebar|skyscraper|sponsor|shopping|tags|widget`)
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

	cleanPresentational(doc)

	clearUnlikelyCandidates(doc)
	defaultCandidates := doc.Find(defaultCandidates)
	elementsToScore := selectionToSlice(defaultCandidates)
	elementsToScore = append(elementsToScore, selectionToSlice(defaultCandidates)...)

	candidates := decideWorthyCandidates(elementsToScore, 0)

	getTopCandidate(candidates)

	winner := candidates[0]

	cleanScripts(winner.selector)
	cleanUnecessary(winner.selector)

	html, err := winner.selector.Html()
	if err != nil {
		return RewriteResult{}, err
	}

	result.HTML = html
	return result, nil
}

func decideWorthyCandidates(elementsToScore []*goquery.Selection, depth int) []Candidate {
	seen := make(map[*html.Node]struct{})
	candidates := []Candidate{}

	for _, s := range elementsToScore {
		can, err := decideCandidate(s)
		if err != nil {
			continue
		}
		can.depth = depth

		parents := getParents(s, 5)
		if len(parents) <= 0 {
			continue
		}

		// We're skipping the initial items we selected
		// We're looking up h1, h2s etc - they are not likely to be the article itself
		// but vital part of the article
		if depth > 0 {
			candidates = append(candidates, can)
			node := s.Get(0)
			seen[node] = struct{}{}
		}

		// children := selectionToSlice(can.selector.Children())
		result := decideWorthyCandidates(parents, depth+1)

		for _, par := range result {

			divider := 0.0

			switch depth {
			case 0:
				divider = 1
			case 1:
				divider = 2
			default:
				divider = float64(depth * 3)
			}

			par.score += par.score / divider
		}


		for _, v := range result {
			node := s.Get(0)
			_, exists := seen[node]
			if (!exists) {
				seen[node] = struct{}{}
				candidates = append(candidates, v)
			}
		}
	}

	return candidates
}

func getParents(s *goquery.Selection, maxDepth int) []*goquery.Selection {
	curDepth := 0
	cur := s.Parent()
	parents := []*goquery.Selection{}
	for maxDepth > curDepth {
		if cur.Length() == 0 {
			break
		}
		parents = append(parents, cur)
		curDepth += 1

		cur = cur.Parent()

	}
	return parents
}

func getTopCandidate(candidates []Candidate) {
	for idx, candidate := range candidates {
		result := getLinkDensity(candidate.selector)
		scoreAfterLinks := float64(candidate.score) * (1 - result)

		candidates[idx].score = scoreAfterLinks
	}

	slices.SortFunc(candidates, func(a Candidate, b Candidate) int {
		return int(a.score - b.score)
	})
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

	tagName := getNodeTag(s)

	switch tagName {
	case "div":
		contentScore = 5
	case "PRE":
	case "TD":
	case "BLOCKQUOTE":
		contentScore += 3

	case "ADDRESS":
	case "OL":
	case "UL":
	case "DL":
	case "DD":
	case "DT":
	case "LI":
	case "FORM":
		contentScore -= 3

	case "H1":
	case "H2":
	case "H3":
	case "H4":
	case "H5":
	case "H6":
	case "TH":
		contentScore -= 6
	}

	contentScore += float64(getClassWeight(s))

	return Candidate{selector: s, score: contentScore, isEmpty: true}, nil
}

func getClassWeight(s *goquery.Selection) int {
	weight := 0

	class, ok := s.Attr("class")
	if ok {
		hasNegatives := negative.Match([]byte(class))
		if hasNegatives {
			weight -= 25
		}

		hasPositives := positive.Match([]byte(class))
		if hasPositives {
			weight += 25
		}
	}

	id, ok := s.Attr("id")

	if ok {
		hasNegatives := negative.Match([]byte(id))
		if hasNegatives {
			weight -= 25
		}

		hasPositives := positive.Match([]byte(id))
		if hasPositives {
			weight += 25
		}
	}

	return weight
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

func removeAllNodesWithSelector(doc *goquery.Selection, selector string) {
	doc.Find(selector).Remove()
}

func cleanScripts(doc *goquery.Selection) {
	removeAllNodesWithSelector(doc, "script")
	removeAllNodesWithSelector(doc, "noscript")
}

func cleanUnecessary(doc *goquery.Selection) {
	removeAllNodesWithSelector(doc, "object")
	removeAllNodesWithSelector(doc, "embed")
	removeAllNodesWithSelector(doc, "footer")
	removeAllNodesWithSelector(doc, "link")
	removeAllNodesWithSelector(doc, "aside")

	removeAllNodesWithSelector(doc, "iframe")
	removeAllNodesWithSelector(doc, "input")
	removeAllNodesWithSelector(doc, "textarea")
	removeAllNodesWithSelector(doc, "select")
	removeAllNodesWithSelector(doc, "button")

	removeAllNodesWithSelector(doc, "field")
	removeAllNodesWithSelector(doc, "fieldset")
}

func cleanPresentational(doc *goquery.Document) {
	doc.Find("style").Remove()
}
