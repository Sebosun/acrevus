package analyzer

import (
	"encoding/json"
	"os"
	"regexp"
	"strings"
	"testing"
)

var whitespacePattern = regexp.MustCompile(`\s+`)

func TestAnalyzerRewrite(t *testing.T) {
	testCases := []struct {
		name                 string
		sourcePath           string
		expectedHTMLPath     string
		expectedMetadataPath string
	}{
		{
			name:                 "barebones",
			sourcePath:           "testdata/barebones/source.html",
			expectedHTMLPath:     "testdata/barebones/expected.html",
			expectedMetadataPath: "testdata/barebones/metadata.json",
		},
		{
			name:                 "barebones-ambigious",
			sourcePath:           "testdata/barebones-ambigious/source.html",
			expectedHTMLPath:     "testdata/barebones-ambigious/expected.html",
			expectedMetadataPath: "testdata/barebones-ambigious/metadata.json",
		},
		{
			name:                 "barebones-navigation",
			sourcePath:           "testdata/barebones-navigation/source.html",
			expectedHTMLPath:     "testdata/barebones-navigation/expected.html",
			expectedMetadataPath: "testdata/barebones-navigation/metadata.json",
		},
		{
			name:                 "barebones-title",
			sourcePath:           "testdata/barebones-title/source.html",
			expectedHTMLPath:     "testdata/barebones-title/expected.html",
			expectedMetadataPath: "testdata/barebones-title/metadata.json",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sourceHTML, err := os.ReadFile(testCase.sourcePath)
			if err != nil {
				t.Fatalf("read source HTML: %v", err)
			}

			expectedHTML, err := os.ReadFile(testCase.expectedHTMLPath)
			if err != nil {
				t.Fatalf("read expected HTML: %v", err)
			}

			expectedMetadata, err := os.ReadFile(testCase.expectedMetadataPath)
			if err != nil {
				t.Fatalf("read expected metadata: %v", err)
			}

			var wantMetadata Metadata
			if err := json.Unmarshal(expectedMetadata, &wantMetadata); err != nil {
				t.Fatalf("parse expected metadata: %v", err)
			}

			result, err := AnalyzerRewrite(string(sourceHTML))
			if err != nil {
				t.Fatalf("Run() failed: %v", err)
			}

			if whitespacePattern.ReplaceAllString(result.HTML, "") != whitespacePattern.ReplaceAllString(string(expectedHTML), "") {
				t.Errorf("RawHTML = %q, want %q", result.HTML, strings.TrimSpace(string(expectedHTML)))
			}
			if result.Metadata != wantMetadata {
				t.Errorf("Metadata = %+v, want %+v", result.Metadata, wantMetadata)
			}
		})
	}
}
