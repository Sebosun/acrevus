package repl

func (m model) View() string {
	// The header
	var view string

	switch m.state {
	case StateReading:
		view = m.ViewArticle()
	case StateMain:
		view = m.PreviewView()
	}

	m.updateViewHeight(len(view))

	return view
}
