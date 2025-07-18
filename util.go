package css_inliner

import (
	"strings"

	css_parser "go.baoshuo.dev/css-parser"
)

// TODO: Get list of supported pseudo selectors from "github.com/andybalholm/cascadia" package
// https://github.com/andybalholm/cascadia/blob/5263deb988702df34b4de5b8cd2fe53add4bea3d/parser.go#L473-L607

var unsupportedSelectors = []string{
	":active", ":after", ":before", ":checked", ":disabled", ":enabled",
	":first-line", ":first-letter", ":focus", ":hover", ":invalid", ":in-range",
	":lang", ":link", ":root", ":selection", ":target", ":valid", ":visited"}

func Inlinable(selector string) bool {
	if strings.Contains(selector, "::") {
		return false
	}

	for _, badSel := range unsupportedSelectors {
		if strings.Contains(selector, badSel) {
			return false
		}
	}

	return true
}

func computeStyleValue(declarations []*css_parser.CssDeclaration) string {
	result := ""

	// set style attribute value
	for _, declaration := range declarations {
		if result != "" {
			result += " "
		}

		result += declaration.StringWithImportant(false)
	}

	return result
}

func mergeStyleDeclarations(styleRules []*StyleRule, output map[string]*StyleDeclaration) {
	for _, styleRule := range styleRules {
		for _, declaration := range styleRule.Declarations {
			styleDecl := NewStyleDeclaration(styleRule, declaration)

			if (output[declaration.Property] == nil) || (styleDecl.Specificity() >= output[declaration.Property].Specificity()) {
				output[declaration.Property] = styleDecl
			}
		}
	}
}
