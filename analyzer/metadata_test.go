package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

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
		{
			name:    "API Fetching",
			fixture: "api-fetching.html",
			want: Metadata{
				Title:         "This API is so Fetching!",
				Byline:        "Nikhil Marathe",
				Dir:           "",
				Lang:          "en-US",
				Excerpt:       "For more than a decade the Web has used XMLHttpRequest (XHR) to achieve asynchronous requests in JavaScript. While very useful, XHR is not a very ...",
				SiteName:      "Mozilla Hacks – the Web developer blog",
				PublishedTime: "",
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

			if got.Title != tt.want.Title {
				t.Errorf("Title = %q, want %q", got.Title, tt.want.Title)
			}
			if got.Byline != tt.want.Byline {
				t.Errorf("Byline = %q, want %q", got.Byline, tt.want.Byline)
			}
			if got.Dir != tt.want.Dir {
				t.Errorf("Dir = %q, want %q", got.Dir, tt.want.Dir)
			}
			if got.Lang != tt.want.Lang {
				t.Errorf("Lang = %q, want %q", got.Lang, tt.want.Lang)
			}
			if got.Excerpt != tt.want.Excerpt {
				t.Errorf("Excerpt = %q, want %q", got.Excerpt, tt.want.Excerpt)
			}
			if got.SiteName != tt.want.SiteName {
				t.Errorf("SiteName = %q, want %q", got.SiteName, tt.want.SiteName)
			}
			if got.PublishedTime != tt.want.PublishedTime {
				t.Errorf("PublishedTime = %q, want %q", got.PublishedTime, tt.want.PublishedTime)
			}
		})
	}
}

func TestIsSchemaJSOND(t *testing.T) {
	tests := []struct {
		name string
		data map[string]any
		want bool
	}{
		{
			name: "schema article",
			data: map[string]any{
				"@context": "https://schema.org",
				"@type":    "Article",
			},
			want: true,
		},
		{
			name: "schema context object and article type array",
			data: map[string]any{
				"@context": map[string]any{"@vocab": "http://schema.org/"},
				"@type":    []any{"Thing", "BlogPosting"},
			},
			want: true,
		},
		{
			name: "versioned schema context",
			data: map[string]any{
				"@context": "https://schema.org/version/29.0/",
				"@type":    "NewsArticle",
			},
			want: true,
		},
		{
			name: "non-schema context",
			data: map[string]any{
				"@context": "https://example.com",
				"@type":    "Article",
			},
			want: false,
		},
		{
			name: "partial article type",
			data: map[string]any{
				"@context": "https://schema.org",
				"@type":    "ArticleDraft",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isSchemaJSOND(tt.data); got != tt.want {
				t.Errorf("isSchemaJSOND() = %t, want %t", got, tt.want)
			}
		})
	}
}
