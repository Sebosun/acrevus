package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractJSOND(t *testing.T) {
	data := parseJSOND("\n<![CDATA[\n{\"headline\": \"Example article\"}\n]]>\n")

	if data["headline"] != "Example article" {
		t.Errorf("extractJSOND() headline = %v, want %q", data["headline"], "Example article")
	}
}

func Test_getMetadata(t *testing.T) {
	tests := []struct {
		name    string
		fixture string
		want    Metadata
		wantErr bool
	}{
		{
			name:    "ACLU article",
			fixture: "aclu.html",
			want: Metadata{
				Title:         "Facebook Is Tracking Me Even Though I’m Not on Facebook",
				Byline:        "Daniel Kahn Gillmor",
				Dir:           "ltr",
				Lang:          "en",
				Excerpt:       "Facebook collects data about people who have never even opted in. But there are ways these non-users can protect themselves.",
				SiteName:      "American Civil Liberties Union",
				PublishedTime: "2018-04-05T06:00",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			document, err := os.ReadFile(filepath.Join("metadata", tt.fixture))
			if err != nil {
				t.Fatalf("read fixture %q: %v", tt.fixture, err)
			}

			got, gotErr := GetMetadata(string(document))
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getMetadata() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getMetadata() succeeded unexpectedly")
			}

			if got != tt.want {
				t.Errorf("getMetadata() = %v, want %v", got, tt.want)
			}
		})
	}
}
