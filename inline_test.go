package css_inliner

import (
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
