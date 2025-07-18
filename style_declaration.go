package cssinliner

import (
	cssparser "go.baoshuo.dev/cssparser"
)

type StyleDeclaration struct {
	StyleRule   *StyleRule
	Declaration *cssparser.Declaration
}

func NewStyleDeclaration(styleRule *StyleRule, declaration *cssparser.Declaration) *StyleDeclaration {
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
