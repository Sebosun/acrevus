package repl

import (
	"fmt"
	"os"
	"path"

	"sebosun/acrevus-go/storage"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) updatePreview(msg tea.Msg) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			newFile, err := m.deleteArticle()
			if err != nil {
				// likely should have some codes for various errors etc
				m.error = fmt.Sprintf("Error %v", err)
			}
			m.articles = newFile

		case "ctrl+c", "q":
			return m, tea.Quit
		case "right", "l":
			m.selected = m.cursor
			m.cursor = 0
			m.state = StateReading
			newM, err := m.getRawArticle()
			if err == nil {
				m = newM
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.articles)-1 {
				m.cursor++
			}
		}
	}
	return m, nil
}

func (m model) getRawArticle() (model, error) {
	art := m.articles[m.selected]
	artPath, err := storage.GetArticlesPath()
	if err != nil {
		return model{}, fmt.Errorf("error reading article file %w", err)
	}

	artPath = path.Join(artPath, art.Path)
	article, err := os.ReadFile(artPath)
	if err != nil {
		return model{}, fmt.Errorf("error reading article file %w", err)
	}
	m.articleRawHTML = string(article)

	return m, nil
}
