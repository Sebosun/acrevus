package analyzer

import (
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

func (da DensityAnalyzer) bulkClean(element *rod.Element) {
	da.clean(element, "form")
	da.clean(element, "fieldset")
	da.clean(element, "object")
	da.clean(element, "embed")
	da.clean(element, "footer")
	da.clean(element, "link")
	da.clean(element, "aside")
	da.clean(element, "iframe")
	da.clean(element, "input")
	da.clean(element, "textarea")
	da.clean(element, "select")
	da.clean(element, "button")
	_ = da.cleanBr(element)
}

func (da DensityAnalyzer) clean(element *rod.Element, tag string) {
	el, err := element.Element(tag)
	if err != nil {
		return
	}

	_ = el.Remove()
}

func (da DensityAnalyzer) cleanBr(element *rod.Element) error {
	brs, err := element.Elements("br")
	if err != nil {
		return err
	}

	for _, brEl := range brs {
		next, err := brEl.Next()
		if err != nil {
			continue
		}
		tagName := next.MustEval(`() => this.tagName.toLowerCase()`).String()

		if tagName == "p" {
			brEl.Remove()
		}
	}
	return nil
}

func (da *DensityAnalyzer) cleanText(text string) string {
	// Remove extra whitespace
	re := regexp.MustCompile(`\s+`)
	cleaned := re.ReplaceAllString(strings.TrimSpace(text), " ")

	// Remove common non-content patterns
	patterns := []string{
		TextRegex[ClickHere],
		TextRegex[Dates],
		TextRegex[Emails],
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		cleaned = re.ReplaceAllString(cleaned, "")
	}

	return strings.TrimSpace(cleaned)
}


func cleanClass(text string) string {
	re := regexp.MustCompile(TextRegex[Class])
	return re.ReplaceAllString(strings.TrimSpace(text), " ")
}

func cleanStyle(text string) string {
	re := regexp.MustCompile(TextRegex[Style])
	return re.ReplaceAllString(strings.TrimSpace(text), " ")
}
