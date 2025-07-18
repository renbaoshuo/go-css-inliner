package css_inliner

type InlinerOption func(*Inliner)

// WithAllowRemoteContent allows the inliner to fetch remote stylesheets.
func WithAllowRemoteContent(allow bool) InlinerOption {
	return func(inliner *Inliner) {
		inliner.allowRemoteContent = allow
	}
}

// WithAllowLocalFiles allows the inliner to fetch local stylesheets.
func WithAllowLocalFiles(allow bool, path string) InlinerOption {
	return func(inliner *Inliner) {
		inliner.allowLocalFiles = allow
		inliner.path = path
	}
}
