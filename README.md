# go-css-inliner

## Installation

```bash
go get go.baoshuo.dev/css-inliner
```

## Usage

```go
package main

import (
  "go.baoshuo.dev/css-inliner"
)

func main() {
  html := `<html><head><style>body { color: red; }</style></head><body>Hello World</body></html>`
  
  result, err := css_inliner.Inline(html)
  if err != nil {
    panic(err)
  }
  
  println(result)
}
```

## Credits

- https://github.com/aymerick/douceur
- https://github.com/PuerkitoBio/goquery

## Author

**go-css-inliner** © [Baoshuo](https://baoshuo.ren), Released under the [MIT](./LICENSE) License.

> [Personal Homepage](https://baoshuo.ren) · [Blog](https://blog.baoshuo.ren) · GitHub [@renbaoshuo](https://github.com/renbaoshuo)
