package cssinliner

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"go.baoshuo.dev/cssparser"
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

func TestInlineWithRemoteStylesheet(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/css")
		w.Write([]byte("body { color: green; }"))
	}))
	defer server.Close()

	source := fmt.Sprintf(`<html><head><link rel="stylesheet" href="%s" /></head><body>Hello Remote</body></html>`, server.URL)
	expected := `<html><head></head><body style="color: green;">Hello Remote</body></html>`

	result, err := Inline(source, WithAllowLoadRemoteStylesheets(true))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestInlineWithComplexSelectors(t *testing.T) {
	source := `<html><head><style>
		div.container > p { color: red; }
		ul li:first-child { color: blue; }
		a[target="_blank"] { text-decoration: none; }
		h1, h2 { font-weight: bold; }
		div ~ span { background-color: yellow; }
	</style></head><body>
		<div class="container">
			<p>This is a paragraph.</p>
			<div>
				<p>Not a direct child</p>
				<span>A span</span>
			</div>
		</div>
		<p>This paragraph is not inside the container.</p>
		<ul>
			<li>First item</li>
			<li>Second item</li>
		</ul>
		<a href="#" target="_blank">Link</a>
		<h1>Heading 1</h1>
		<h2>Heading 2</h2>
		<span>A sibling span</span>
	</body></html>`
	expected := `<html><head></head><body>
		<div class="container">
			<p style="color: red;">This is a paragraph.</p>
			<div>
				<p>Not a direct child</p>
				<span>A span</span>
			</div>
		</div>
		<p>This paragraph is not inside the container.</p>
		<ul>
			<li style="color: blue;">First item</li>
			<li>Second item</li>
		</ul>
		<a href="#" target="_blank" style="text-decoration: none;">Link</a>
		<h1 style="font-weight: bold;">Heading 1</h1>
		<h2 style="font-weight: bold;">Heading 2</h2>
		<span style="background-color: yellow;">A sibling span</span>
	</body></html>`

	inliner := NewInliner(source)

	result, err := inliner.Inline()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestInlineWithBrokenCSS(t *testing.T) {
	source := `<html><head><style>
	.name {
		color: red;
		background-co
		border-radius: 5px;
		border: 1px solid #fff;
	}
	几个无效字符
	</style></head><body>
		<div class="name" style="color: blue;">Hello</div>
	</body></html>`
	expected := `<html><head></head><body>
		<div class="name" style="border: 1px solid #fff; color: blue;">Hello</div>
	</body></html>`

	result, err := Inline(source, WithParserOptions(cssparser.WithLooseParsing(true)))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}
