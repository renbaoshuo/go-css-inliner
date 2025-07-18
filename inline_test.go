package css_inliner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInline(t *testing.T) {
	source := `<html><head><style>body { color: red; }</style></head><body>Hello World</body></html>`
	expected := `<html><head></head><body style="color: red;">Hello World</body></html>`

	inliner := NewInliner(source)

	result, err := inliner.Inline()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestInlineFile(t *testing.T) {
	source := `<html><head><link rel="stylesheet" href="./style.css" /></head><body>Hello File</body></html>`
	stylesheet := `body { color: blue; }`
	expected := `<html><head></head><body style="color: blue;">Hello File</body></html>`

	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "inline-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Write the HTML file
	htmlPath := filepath.Join(tempDir, "index.html")
	err = os.WriteFile(htmlPath, []byte(source), 0644)
	if err != nil {
		t.Fatalf("Failed to write HTML file: %v", err)
	}

	// Write the stylesheet to a file
	cssPath := filepath.Join(tempDir, "style.css")
	err = os.WriteFile(cssPath, []byte(stylesheet), 0644)
	if err != nil {
		t.Fatalf("Failed to write stylesheet: %v", err)
	}

	// Run the inliner
	result, err := InlineFile(htmlPath)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
