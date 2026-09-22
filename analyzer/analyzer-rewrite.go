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
	unlikelyRoles = []string{
		"menu", "menubar",
		"complementary",
		"navigation",
		"alert",
		"alertdialog",
		"dialog",
	}
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
	positive      = regexp.MustCompile(`(?i)article|body|content|entry|hentry|h-entry|main|page|pagination|post|text|blog|story`)
	negative      = regexp.MustCompile(`(?i)-ad-|hidden|^hid$| hid$| hid |^hid |banner|combx|comment|com-|contact|footer|gdpr|masthead|media|meta|outbrain|promo|related|scroll|share|shoutbox|sidebar|skyscraper|sponsor|shopping|tags|widget`)
	shareElements = regexp.MustCompile(`(?i)(\b|_)(share|sharedaddy)(\b|_),`)
)

const (
	defaultMaxParentDepth = 5
	defaultTreshold       = 500
)

type AnalyzerGoquery struct {
	MaxParentDepth int
	TextTreshold   int
}

func NewAnalyzerGoquery() *AnalyzerGoquery {
	return &AnalyzerGoquery{MaxParentDepth: defaultMaxParentDepth, TextTreshold: defaultTreshold}
}

func AnalyzerRewrite(document string) (RewriteResult, error) {
	return NewAnalyzerGoquery().Parse(document)
}

func (a *AnalyzerGoquery) Parse(document string) (RewriteResult, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(document))
	if err != nil {
		return RewriteResult{}, err
	}

	parseResult := RewriteResult{}
	parseResult.Metadata = GetMetadata(doc)

	body := doc.Find("body").First()

	if body.Length() == 0 {
		return RewriteResult{}, ErrEmptyBody
	}

	a.clearUnlikelyCandidates(doc)
	a.clearEmptyElements(doc)

	defaultCandidates := doc.Find(defaultCandidates)
	elementsToScore := a.selectionToSlice(defaultCandidates)
	elementsToScore = append(elementsToScore, a.selectionToSlice(defaultCandidates)...)

	candidates := a.decideWorthyCandidates(elementsToScore, 0)

	a.getTopCandidate(candidates)

	winner := candidates[0]

	a.cleanShare(winner.selector)
	a.cleanScripts(winner.selector)
	a.cleanUnecessary(winner.selector)
	a.cleanPresentational(winner.selector)

	html, err := goquery.OuterHtml(winner.selector)
	if err != nil {
		return RewriteResult{}, err
	}

	parseResult.HTML = html
	return parseResult, nil
}

func (a *AnalyzerGoquery) decideWorthyCandidates(elementsToScore []*goquery.Selection, depth int) []Candidate {
	candidates := []Candidate{}

	for _, s := range elementsToScore {
		candidate, err := a.decideCandidate(s)
		if err != nil {
			continue
		}

		parents := a.getParents(s)
		if len(parents) <= 0 {
			continue
		}

		// We're skipping the initial items we selected
		// We're looking up h1, h2s etc - they are not likely to be the article itself
		// but vital part of the article
		candidates = append(candidates, candidate)

		for parentDepth, parent := range parents {
			parentCandidate, err := a.decideCandidate(parent)
			if err != nil {
				continue
			}

			divider := 0.0
			switch parentDepth {
			case 0:
				divider = 1
			case 1:
				divider = 2
			default:
				divider = float64(depth * 3)
			}

			parentCandidate.score += candidate.score / divider
			candidates = append(candidates, parentCandidate)
		}

	}

	return candidates
}

func (a *AnalyzerGoquery) getParents(s *goquery.Selection) []*goquery.Selection {
	curDepth := 0
	cur := s.Parent()
	parents := []*goquery.Selection{}
	for a.MaxParentDepth > curDepth {
		if cur.Length() == 0 {
			break
		}
		parents = append(parents, cur)
		curDepth += 1

		cur = cur.Parent()

	}
	return parents
}

func (a *AnalyzerGoquery) getTopCandidate(candidates []Candidate) {
	for idx, candidate := range candidates {
		result := a.getLinkDensity(candidate.selector)
		scoreAfterLinks := float64(candidate.score) * (1 - result)

		candidates[idx].score = scoreAfterLinks
	}

	slices.SortFunc(candidates, func(a Candidate, b Candidate) int {
		return int(b.score - a.score)
	})
}

func (a *AnalyzerGoquery) decideCandidate(s *goquery.Selection) (Candidate, error) {
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

	tagName := a.getNodeTag(s)

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

	contentScore += float64(a.getClassWeight(s))

	return Candidate{selector: s, score: contentScore, isEmpty: true}, nil
}

func (a *AnalyzerGoquery) getClassWeight(s *goquery.Selection) int {
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

func (a *AnalyzerGoquery) clearUnlikelyCandidates(doc *goquery.Document) {
	doc.Find("[class], [id]").Each(func(_ int, s *goquery.Selection) {
		className, _ := s.Attr("class")
		id, _ := s.Attr("id")

		value := strings.ToLower(className + " " + id)

		isUnlikelyClass := unlikelyCandidates.Match([]byte(value))
		mightBeCandidate := okMaybeItsACandidate.Match([]byte(value))

		tagName := a.getNodeTag(s)

		hasCode := a.hasAncestorTag(s, "code")
		hasTable := a.hasAncestorTag(s, "table")

		if tagName != "a" && tagName != "body" && isUnlikelyClass && !mightBeCandidate && !hasCode && !hasTable {
			s.Remove()
		}
	})

	for _, role := range unlikelyRoles {
		roleQuery := fmt.Sprintf(`[role="%s"]`, role)
		doc.Find(roleQuery).Remove()
	}
}

func (a *AnalyzerGoquery) clearEmptyElements(doc *goquery.Document) {
	tagNames := []string{"div", "section", "header", "h1", "h2", "h3", "h4", "h5", "h6"}

	for _, tag := range tagNames {
		doc.Find(tag).Each(func(_ int, s *goquery.Selection) {
			if len(s.Nodes) > 0 {
				// https://developer.mozilla.org/en-US/docs/Web/API/Node/nodeType
				// golang for whatever reason decided ElementNode is 3
				isElementNode := s.Nodes[0].Type == html.ElementNode
				hasNoChildren := s.Children().Length() == 0
				noText := len(s.Text()) == 0

				isEmpty := isElementNode && hasNoChildren && noText
				if (isEmpty) {
					s.Remove()
				}
			}
		})
	}
}

func (a *AnalyzerGoquery) removeAllNodesWithSelector(doc *goquery.Selection, selector string) {
	doc.Find(selector).Remove()
}

func (a *AnalyzerGoquery) cleanScripts(doc *goquery.Selection) {
	a.removeAllNodesWithSelector(doc, "script")
	a.removeAllNodesWithSelector(doc, "noscript")
	a.removeAllNodesWithSelector(doc, "style")
}

func (a *AnalyzerGoquery) cleanUnecessary(doc *goquery.Selection) {
	a.removeAllNodesWithSelector(doc, "object")
	a.removeAllNodesWithSelector(doc, "embed")
	a.removeAllNodesWithSelector(doc, "footer")
	a.removeAllNodesWithSelector(doc, "link")
	a.removeAllNodesWithSelector(doc, "aside")

	a.removeAllNodesWithSelector(doc, "iframe")
	a.removeAllNodesWithSelector(doc, "input")
	a.removeAllNodesWithSelector(doc, "textarea")
	a.removeAllNodesWithSelector(doc, "select")
	a.removeAllNodesWithSelector(doc, "button")

	a.removeAllNodesWithSelector(doc, "field")
	a.removeAllNodesWithSelector(doc, "fieldset")
}

func (a *AnalyzerGoquery) cleanPresentational(s *goquery.Selection) {
	s.RemoveAttr("class")
	s.RemoveAttr("width")
	s.RemoveAttr("height")

	for _, attribute := range presentationalAttributes {
		s.RemoveAttr(attribute)
	}
	s.Children().Each(func(_ int, child *goquery.Selection) {
		a.cleanPresentational(child)
	})
}

// Have to use it before we clean classes otherwise we're done brah
func (a *AnalyzerGoquery) cleanShare(s *goquery.Selection) {
	s.Children().Each(func(_ int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		id, _ := s.Attr("id")

		string := class + id
		// if we have any "share" element lets remove the child that contains it
		tresholdMet := len(s.Text()) < a.TextTreshold
		if tresholdMet && shareElements.Match([]byte(string)) {
			s.Remove()
		}
	})
}
