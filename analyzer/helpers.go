package analyzer

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

func (da *DensityAnalyzer) redistributeToParents(blocks *[]ContentBlock) {
	for _, block := range *blocks {
		parent, err := block.Element.Parent()
		if err == nil {
			for i, bl := range *blocks {
				if da.isSameElement(parent, bl.Element) {
					(*blocks)[i].Density += block.Density * 0.5
				}

				grandPar, err := parent.Parent()
				if err != nil {
					continue
				}

				if da.isSameElement(grandPar, bl.Element) {
					(*blocks)[i].Density += block.Density * 0.25
				}

				grandGrandPar, err := grandPar.Parent()
				if err != nil {
					continue
				}
				if da.isSameElement(grandGrandPar, bl.Element) {
					(*blocks)[i].Density += block.Density * 0.15
				}
			}
		}
	}
}
