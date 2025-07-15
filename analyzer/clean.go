package analyzer

import (
	"regexp"
	"strings"

	"github.com/go-rod/rod"
)

func (da DensityAnalyzer) clean(element *rod.Element, tag string) {
	el, err := element.Element(tag)
	if err != nil {
		return
	}

	el.Remove()
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

func cleanClass(text string) string {
	re := regexp.MustCompile(TextRegex[Class])
	return re.ReplaceAllString(strings.TrimSpace(text), " ")
}

func cleanStyle(text string) string {
	re := regexp.MustCompile(TextRegex[Style])
	return re.ReplaceAllString(strings.TrimSpace(text), " ")
}
