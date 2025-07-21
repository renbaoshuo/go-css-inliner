package cssinliner

import (
	"go.baoshuo.dev/cssparser"
)

type InlinerOption func(*Inliner)

// WithAllowLoadRemoteStylesheets allows the inliner to fetch remote stylesheets.
func WithAllowLoadRemoteStylesheets(allow bool) InlinerOption {
	return func(inliner *Inliner) {
		inliner.allowLoadRemoteStylesheets = allow
	}
}

// WithAllowReadLocalFiles allows the inliner to fetch local stylesheets.
func WithAllowReadLocalFiles(allow bool, path string) InlinerOption {
	return func(inliner *Inliner) {
		inliner.allowReadLocalFiles = allow
		inliner.path = path
	}
}

func WithParserOptions(parserOptions ...cssparser.ParserOption) InlinerOption {
	return func(inliner *Inliner) {
		inliner.parserOptions = append(inliner.parserOptions, parserOptions...)
	}
}
