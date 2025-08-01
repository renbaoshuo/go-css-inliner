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

// WithParserOptions allows setting custom CSS parser options.
// This can be used to customize the behavior of the CSS parser.
func WithParserOptions(parserOptions ...cssparser.ParserOption) InlinerOption {
	return func(inliner *Inliner) {
		inliner.parserOptions = append(inliner.parserOptions, parserOptions...)
	}
}

type HtmlPreprocessor func(html, path string) (string, error)

// WithHtmlPreprocessor allows setting a custom HTML preprocessor function.
// This function can be used to modify the HTML before processing.
func WithHtmlPreprocessor(preprocessor HtmlPreprocessor) InlinerOption {
	return func(inliner *Inliner) {
		inliner.htmlPreprocessor = preprocessor
	}
}

type CssFilePreprocessor func(css, path string) (string, error)

// WithCssPreprocessor allows setting a custom CSS preprocessor function.
// This function can be used to modify the CSS before inlining.
//
// NOTE: This function only applies to CSS files, not <style> tags.
// If you want to preprocess CSS in <style> tags, use the `WithHtmlPreprocessor` option
// to modify the HTML to include the preprocessed CSS directly in the <style> tags.
func WithCssFilePreprocessor(preprocessor CssFilePreprocessor) InlinerOption {
	return func(inliner *Inliner) {
		inliner.cssFilePreprocessor = preprocessor
	}
}
