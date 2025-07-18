package css_inliner

import (
	css_parser "go.baoshuo.dev/css-parser"
)

type StyleDeclaration struct {
	StyleRule   *StyleRule
	Declaration *css_parser.CssDeclaration
}

func NewStyleDeclaration(styleRule *StyleRule, declaration *css_parser.CssDeclaration) *StyleDeclaration {
	return &StyleDeclaration{
		StyleRule:   styleRule,
		Declaration: declaration,
	}
}

func (styleDecl *StyleDeclaration) Specificity() int {
	if styleDecl.Declaration.Important {
		return 10000
	}

	return styleDecl.StyleRule.Specificity
}
