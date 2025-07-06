// Package repl, handling displaying the app info on terminal
package repl

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"fmt"
	"os"
	"sebosun/acrevus-go/storage"

	tea "github.com/charmbracelet/bubbletea"
)

type Sizes struct {
	width  int
	height int
}

type model struct {
	sizes      Sizes
	scroll     int      // 0-1000 whatever
	choices    []string // items on the to-do list
	cursor     int      // which to-do list item our cursor is pointing at
	error      string
	curArticle storage.Entry
	selected   int // which article is being read
	articles   []storage.Entry
	isReading  bool
}

func initialModel(fileData storage.FileData) model {
	choices := []string{}

	for _, v := range fileData.Entries {
		choices = append(choices, v.Title)
	}

	return model{
		// Our to-do list is a grocery list
		cursor:  0,
		choices: choices,
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected:  0,
		articles:  fileData.Entries,
		isReading: false,
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return tea.SetWindowTitle("Acrevus")
}

func InitTea() {
	fileData, err := storage.GetFileData()
	if err != nil {
		fmt.Printf("Error initalizing Tea from filedata: %v", err)
		os.Exit(1)
	}
	p := tea.NewProgram(initialModel(fileData), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
