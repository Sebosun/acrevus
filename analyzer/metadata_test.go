package analyzer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Test_getMetadata(t *testing.T) {
	sources, err := os.ReadDir(filepath.Join("metadata", "source"))
	if err != nil {
		t.Fatalf("read metadata sources: %v", err)
	}

	for _, source := range sources {
		if source.IsDir() || filepath.Ext(source.Name()) != ".html" {
			continue
		}

		name := strings.TrimSuffix(source.Name(), filepath.Ext(source.Name()))
		t.Run(name, func(t *testing.T) {
			document, err := os.ReadFile(filepath.Join("metadata", "source", source.Name()))
			if err != nil {
				t.Fatalf("read source fixture: %v", err)
			}

			expected, err := os.ReadFile(filepath.Join("metadata", "expected", name+".json"))
			if err != nil {
				t.Fatalf("read expected fixture: %v", err)
			}

			var want Metadata
			if err := json.Unmarshal(expected, &want); err != nil {
				t.Fatalf("parse expected fixture: %v", err)
			}

			got, err := GetMetadataFromString(string(document))
			if err != nil {
				t.Fatalf("GetMetadata() failed: %v", err)
			}
			if got != want {
				t.Errorf("GetMetadata() = %+v, want %+v", got, want)
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
