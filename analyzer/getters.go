package analyzer

import (
	"regexp"
	"slices"
	"strings"

	"github.com/go-rod/rod"
)

// Loose adapatation from mozilla/readibility
func (da *DensityAnalyzer) getTitle() string {
	dupa, err := da.page.Search("og:title")
	if err == nil {
		el := dupa.First
		dupa, err := el.Attribute("content")
		if err == nil {
			return *dupa
		}
	}

	title := da.page.MustInfo().Title
	origTitle := title

	if title == "" {
		elTitle, err := da.page.Element("title")
		if err == nil {
			title = elTitle.MustText()
			origTitle = title
		}
	}

	// TODO: Title separators

	if strings.Contains(title, ": ") {
		hElems := da.page.MustElements("h1, h2")
		trimmed := strings.TrimSpace(title)

		match := slices.ContainsFunc(hElems, func(el *rod.Element) bool {
			elText := strings.TrimSpace(el.MustText())
			return trimmed == elText
		})
		if !match {
			last := strings.LastIndex(origTitle, ":")
			// this won't break since we're already checking if : is within title string
			title = title[:last+1]

			isTooShort := len(strings.TrimSpace(title)) < 3
			if isTooShort {
				first := strings.Index(origTitle, ":")
				title = title[:first+1]
			}
		}
	} else if len(title) > 150 || len(title) < 15 {
		el := da.page.MustElements("h1")
		// There technically should only be one h1 on the page, if we go by standards
		if len(el) == 1 {
			title = el[0].MustText()
		}
	}

	// whitespaces
	re := regexp.MustCompile(`\s{2,}`)
	title = re.ReplaceAllString(strings.TrimSpace(title), " ")

	return title
}
