# Go CSS Inliner

## Installation

```bash
go get go.baoshuo.dev/cssinliner
```

## Usage

```go
package main

import (
  "go.baoshuo.dev/cssinliner"
)

func main() {
  html := `<html><head><style>body { color: red; }</style></head><body>Hello World</body></html>`

  result, err := cssinliner.Inline(html)
  if err != nil {
    panic(err)
  }

  println(result)
}
```

This package provides these functions:

- `Inline(html string, options... InlinerOption) (string, error)`<br />
  Inlines CSS styles into the provided HTML string.
- `InlineFile(path string, options... InlinerOption) (string, error)`<br />
  Reads an HTML file from the specified path and inlines CSS styles into it.

The available options include:

- `WithAllowLoadRemoteStylesheets(allow bool)`<br />
  Allows the inliner to fetch remote stylesheets.
- `WithAllowReadLocalFiles(allow bool, path string)`<br />
  Allows the inliner to fetch local stylesheets from the specified path.

## Credits

- https://github.com/aymerick/douceur
- https://github.com/PuerkitoBio/goquery

## Author

**go-css-inliner** © [Baoshuo](https://baoshuo.ren), Released under the [MIT](./LICENSE) License.

> [Personal Homepage](https://baoshuo.ren) · [Blog](https://blog.baoshuo.ren) · GitHub [@renbaoshuo](https://github.com/renbaoshuo)
