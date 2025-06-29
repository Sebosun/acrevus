// Package repl, handling displaying the app info on terminal
package repl

// These imports will be used later on the tutorial. If you save the file
// now, Go might complain they are unused, but that's fine.
// You may also need to run `go mod tidy` to download bubbletea and its
// dependencies.
import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type Sizes struct {
	width  int
	height int
}

type model struct {
	sizes       Sizes
	totalHeight int              // 0-1000 whatever
	scroll      int              // 0-1000 whatever
	choices     []string         // items on the to-do list
	cursor      int              // which to-do list item our cursor is pointing at
	selected    map[int]struct{} // which to-do items are selected
}

func initialModel() model {
	return model{
		// Our to-do list is a grocery list
		cursor:      20,
		totalHeight: 500,
		choices:     []string{"Buy carrots", "Buy celery", "Buy kohlrabi", "lmaoooooo"},
		// A map which indicates which choices are selected. We're using
		// the  map like a mathematical set. The keys refer to the indexes
		// of the `choices` slice, above.
		selected: make(map[int]struct{}),
	}
}

func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return tea.SetWindowTitle("Acrevus")
}

func InitTea() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
