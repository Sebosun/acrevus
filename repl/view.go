package repl

func (m model) View() string {
	// The header
	if m.isReading {
		return m.ViewArticle()
	}

	return m.PreviewView()
}
