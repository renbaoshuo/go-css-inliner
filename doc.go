/*
Package `cssinliner` provides functionality to inline CSS styles into HTML documents.

It allows for inlining styles from remote stylesheets, local files, and inline styles within the HTML document itself.

It supports options to control whether remote stylesheets and local files should be loaded, and it processes the HTML to apply styles directly to elements.

Here's a brief overview of the main functions:
- Inline(html string, options... InlinerOption) (string, error): Inlines CSS styles into the provided HTML string.
- InlineFile(path string, options... InlinerOption) (string, error): Reads an HTML file from the specified path and inlines CSS styles into it.

The available options include:
- WithAllowLoadRemoteStylesheets(allow bool): Allows the inliner to fetch remote stylesheets.
- WithAllowReadLocalFiles(allow bool, path string): Allows the inliner to fetch local stylesheets from the specified path.

The source code of this package is hosted on GitHub: https://github.com/renbaoshuo/go-css-inliner
*/
package cssinliner
