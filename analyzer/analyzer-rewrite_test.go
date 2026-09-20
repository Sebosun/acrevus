package analyzer

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAnalyzerRewrite(t *testing.T) {
	cases, err := os.ReadDir("testdata")
	if err != nil {
		t.Fatalf("read analyzer test cases: %v", err)
	}

	for i, testCase := range cases {
		if !testCase.IsDir() {
			continue
		}

		if i != 1 {
			continue
		}
		t.Run(testCase.Name(), func(t *testing.T) {
			sourceHTML, err := os.ReadFile(filepath.Join("testdata", testCase.Name(), "source.html"))
			if err != nil {
				t.Fatalf("read source HTML: %v", err)
			}

			expectedHTML, err := os.ReadFile(filepath.Join("testdata", testCase.Name(), "expected.html"))
			if err != nil {
				t.Fatalf("read expected HTML: %v", err)
			}

			expectedMetadata, err := os.ReadFile(filepath.Join("testdata", testCase.Name(), "metadata.json"))
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

			if result.HTML != strings.TrimSpace(string(expectedHTML)) {
				t.Errorf("RawHTML = %q, want %q", result.HTML, strings.TrimSpace(string(expectedHTML)))
			}
			if result.Metadata != wantMetadata {
				t.Errorf("Metadata = %+v, want %+v", result.Metadata, wantMetadata)
			}
		})
	}
}
