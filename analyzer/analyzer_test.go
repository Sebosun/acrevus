package analyzer

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var staticServer *httptest.Server

func TestMain(m *testing.M) {
	staticServer = httptest.NewServer(http.FileServer(http.Dir("testdata")))

	exitCode := m.Run()
	staticServer.Close()

	os.Exit(exitCode)
}

const isSkip = true

func TestAnalyzer(t *testing.T) {
	t.Run("Extracts barebones article", func(t *testing.T) {
		skipCI(t)
		fixtureURL := staticServer.URL + "/barebones.html"
		result, err := Run(fixtureURL)

		wantedTag := "article"
		wantedTitle := "Barebones article"

		gotTag := result.Content.TagName
		gotTitle := result.Metadata.Title

		if err != nil {
			t.Errorf("Error shouldn't happen here %s", err.Error())
		}

		if gotTag != wantedTag {
			t.Errorf("Should correct extract tag name. Got %s - wanted %s", gotTag, wantedTag)
		}

		if gotTitle != wantedTitle {
			t.Errorf("Should correct extract title name. Got %s - wanted %s", gotTitle, wantedTitle)
		}
	})

	t.Run("Ignores navigation tags", func(t *testing.T) {
		skipCI(t)
		fixtureURL := staticServer.URL + "/barebones-navigation.html"
		result, err := Run(fixtureURL)

		wantedTag := "article"
		wantedTitle := "Barebones-navigation article"

		gotTag := result.Content.TagName
		gotTitle := result.Metadata.Title

		if err != nil {
			t.Errorf("Error shouldn't happen here %s", err.Error())
		}

		if gotTag != wantedTag {
			t.Errorf("Should correct extract tag name. Got %s - wanted %s", gotTag, wantedTag)
		}

		if gotTitle != wantedTitle {
			t.Errorf("Should correct extract title name. Got %s - wanted %s", gotTitle, wantedTitle)
		}
	})

	t.Run("Prefers title from og:title", func(t *testing.T) {
		skipCI(t)
		fixtureURL := staticServer.URL + "/barebones-title.html"
		result, err := Run(fixtureURL)

		wantedTag := "article"
		wantedTitle := "Different title"

		gotTag := result.Content.TagName
		gotTitle := result.Metadata.Title

		if err != nil {
			t.Errorf("Error shouldn't happen here %s", err.Error())
		}

		if gotTag != wantedTag {
			t.Errorf("Should correct extract tag name. Got %s - wanted %s", gotTag, wantedTag)
		}

		if gotTitle != wantedTitle {
			t.Errorf("Should correct extract title name. Got %s - wanted %s", gotTitle, wantedTitle)
		}
	})

	t.Run("Finds articles with ambigious html structure ", func(t *testing.T) {
		skipCI(t)
		fixtureURL := staticServer.URL + "/barebones-ambigious.html"
		result, err := Run(fixtureURL)

		wantedTag := "div"
		wantedTitle := "Barebones article"
		wantedTagID := "searched-for"

		gotTag := result.Content.TagName
		gotTitle := result.Metadata.Title
		gotTagID := getHTMLElemID(result.Content.Element)

		if err != nil {
			t.Errorf("Error shouldn't happen here %s", err.Error())
		}

		if gotTag != wantedTag {
			t.Errorf("Should correct extract tag name. Got %s - wanted %s", gotTag, wantedTag)
		}

		if gotTitle != wantedTitle {
			t.Errorf("Should correct extract title name. Got %s - wanted %s", gotTitle, wantedTitle)
		}

		if !strings.Contains(result.RawHTML, wantedTagID) {
			t.Errorf("Element has wrong id. Got %s - wanted %s", gotTagID, wantedTagID)
		}
	})

	t.Run("Real life example of Edward Fesers blog", func(t *testing.T) {
		skipCI(t)
		fixtureURL := staticServer.URL + "/fesers-blog.html"
		result, err := Run(fixtureURL)

		wantedTag := "div"
		wantedTitle := "Berkeley’s God"
		wantedTagID := "post-body entry-content"

		gotTag := result.Content.TagName
		gotTitle := result.Metadata.Title
		gotTagID := getHTMLElemID(result.Content.Element)

		if err != nil {
			t.Errorf("Error shouldn't happen here %s", err.Error())
		}

		if gotTag != wantedTag {
			t.Errorf("Should correct extract tag name. Got %s - wanted %s", gotTag, wantedTag)
		}

		if gotTitle != wantedTitle {
			t.Errorf("Should correct extract title name. Got %s - wanted %s", gotTitle, wantedTitle)
		}

		if !strings.Contains(result.RawHTML, wantedTagID) {
			t.Errorf("Element has wrong id. Got %s - wanted %s", gotTagID, wantedTagID)
		}
	})

	// t.Run("handles a page without article content", func(t *testing.T) {
	// 	// Add analyzer assertions here.
	// })
}

func skipCI(t *testing.T) {
	if isSkip {
		t.Skip("Skipping test")
	}
}
