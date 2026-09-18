package analyzer

import (
	"regexp"
	"strings"
)

// TODO: Add getTitle tests for missing OG metadata/content, missing <title>, and heading fallbacks.
// Loose adapatation from mozilla/readibility
func (da *DensityAnalyzer) getTitle() string {
	// Easy way first
	ogTitles, err := da.page.Elements("meta[property='og:title'], meta[name='og:title']")

	if err == nil {
		for _, ogTitle := range ogTitles {
			// <meta property="og:title" content="Your Title"> <---- here
			content, err := ogTitle.Attribute("content")
			if err == nil && content != nil && strings.TrimSpace(*content) != "" {
				return *content
			}
		}
	}

	// Harder, getting info from
	title := ""
	if info, err := da.page.Info(); err == nil {
		title = info.Title
	}

	origTitle := title

	// looking up title elems
	if title == "" {
		elTitles, err := da.page.Elements("title")
		hasTitleElems := err == nil && len(elTitles) > 0
		if hasTitleElems {
			text, err := elTitles[0].Text()
			if err == nil {
				title = text
				origTitle = title
			}
		}
	}

	// TODO: Title separators

	if strings.Contains(title, ": ") {
		hElems, err := da.page.Elements("h1, h2")
		if err == nil {
			trimmed := strings.TrimSpace(title)
			match := false
			for _, el := range hElems {
				elText, err := el.Text()
				if err == nil && trimmed == strings.TrimSpace(elText) {
					match = true
					break
				}
			}

			if !match {
				separatorIndex := strings.LastIndex(origTitle, ":")
				// this won't break since we're already checking if : is within title string
				title = title[:separatorIndex+1]

				isTooShort := len(strings.TrimSpace(title)) < 3
				if isTooShort {
					first := strings.Index(origTitle, ":")
					title = title[:first+1]
				}
			}
		}
	} else if len(title) > 150 || len(title) < 15 {
		el, err := da.page.Elements("h1")
		// There technically should only be one h1 on the page, if we go by standards
		if err == nil && len(el) == 1 {
			text, err := el[0].Text()
			if err == nil {
				title = text
			}
		}
	}

	// cleanup
	re := regexp.MustCompile(`\s{2,}`)
	title = re.ReplaceAllString(strings.TrimSpace(title), " ")

	return title
}

func (da *DensityAnalyzer) getAuthor() string {
	return ""
}

// func (da *DensityAnalyzer) getDescription() string {}
//
// func (da *DensityAnalyzer) getPublishDate() string {}
//
// func (da *DensityAnalyzer) getModifiedDate() string {}
