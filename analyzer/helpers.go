package analyzer

import (
	"fmt"
	"time"

	"github.com/go-rod/rod"
)

func (da *DensityAnalyzer) weighScoreByTag(blocks *[]ContentBlock) {
	for i, block := range *blocks {
		score := block.Density

		switch block.TagName {
		case "article", "main":
			score *= 1.5
		case "p":
			score += 5
		case "div":
			score += 3
		case "option", "select", "li", "ul", "ol", "dl", "form":
			score *= 0.01
		case "INVALID":
			score *= 0.01
		}

		switch {
		case block.TextLength > 1000:
			score *= 1.2
		case block.TextLength > 5000:
			score *= 2
		}

		(*blocks)[i].Density = score
	}
}

var fingerprintJS = `() =>
		this.id + "|"
		+ this.tagName + '|'
		+ this.className + '|'
		+ this.childElementCount + '|'
		+ this.clientHeight + '|'
		+ this.clientWidth + '|'
		+ this.textContent.substring(0, 50)`

func elementFingerprint(el *rod.Element) string {
	return el.MustEval(fingerprintJS).String()
}

func (da *DensityAnalyzer) redistributeToParents(blocks *[]ContentBlock) {
	start := time.Now()

	fpIndex := make(map[string][]int, len(*blocks))
	for i := range *blocks {
		bl := &(*blocks)[i]
		if bl.CachedFingerprintString == "" {
			bl.CachedFingerprintString = elementFingerprint(bl.Element)
		}
		fpIndex[bl.CachedFingerprintString] = append(fpIndex[bl.CachedFingerprintString], i)
	}

	loopStart := time.Now()
	for _, block := range *blocks {
		parent, err := block.Element.Parent()
		if err != nil {
			continue
		}

		parentFP := elementFingerprint(parent)
		for _, i := range fpIndex[parentFP] {
			(*blocks)[i].Density += block.Density * 0.5
		}

		grandPar, err := parent.Parent()
		if err != nil {
			continue
		}
		grandFP := elementFingerprint(grandPar)
		for _, i := range fpIndex[grandFP] {
			(*blocks)[i].Density += block.Density * 0.25
		}

		grandGrandPar, err := grandPar.Parent()
		if err != nil {
			continue
		}
		grandGrandFP := elementFingerprint(grandGrandPar)
		for _, i := range fpIndex[grandGrandFP] {
			(*blocks)[i].Density += block.Density * 0.15
		}
	}
	fmt.Printf("  [timer] parent redistribution loop: %v\n", time.Since(loopStart))
	fmt.Printf("  [timer] redistributeToParents total: %v\n", time.Since(start))
}
