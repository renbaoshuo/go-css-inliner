package css_inliner

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	css_parser "go.baoshuo.dev/css-parser"
	"golang.org/x/net/html"
)

const elementMarkerAttr = "data-inliner-marker"

type Inliner struct {
	html          string                      // Raw HTML content
	path          string                      // Path to the HTML file
	doc           *goquery.Document           // Parsed HTML document
	stylesheets   []*css_parser.CssStylesheet // Parsed CSS stylesheets
	elements      map[string]*Element         // HTML elements matching collected inlinable style rules
	rawRules      []fmt.Stringer              // CSS rules that are not inlinable but that must be inserted in output document
	elementMarker int                         // current element marker value

	allowRemoteContent bool // Whether to allow remote content (e.g., <link rel="stylesheet" href="http://example.com/style.css" />)
	allowLocalFiles    bool // Whether to allow local files (e.g., <link rel="stylesheet" href="/path/to/local/file.css" />)
}

func NewInliner(html string, options ...InlinerOption) *Inliner {
	inliner := &Inliner{
		html:     html,
		elements: make(map[string]*Element),
	}

	for _, option := range options {
		option(inliner)
	}

	return inliner
}

// Inline processes the HTML content and inlines the CSS styles.
func Inline(html string, options ...InlinerOption) (string, error) {
	result, err := NewInliner(html, options...).Inline()
	if err != nil {
		return "", err
	}

	return result, nil
}

func (inliner *Inliner) Inline() (string, error) {
	if err := inliner.parseHTML(); err != nil {
		return "", err
	}

	if inliner.allowRemoteContent {
		if err := inliner.fetchExternalStylesheets(); err != nil {
			return "", fmt.Errorf("failed to fetch external stylesheets: %w", err)
		}
	}

	if err := inliner.parseStylesheets(); err != nil {
		return "", err
	}

	inliner.collectElementsAndRules()

	if err := inliner.inlineStyleRules(); err != nil {
		return "", err
	}

	inliner.insertRawStylesheet()

	return inliner.genHTML()
}

func (inliner *Inliner) parseHTML() error {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(inliner.html))
	if err != nil {
		return err
	}

	inliner.doc = doc

	return nil
}

func (inliner *Inliner) fetchExternalStylesheets() error {
	// TODO: Implement fetching of external stylesheets
	return errors.New("fetching external stylesheets is not implemented")
}

func (inliner *Inliner) parseStylesheets() error {
	var result error

	inliner.doc.Find("style").EachWithBreak(func(i int, s *goquery.Selection) bool {
		stylesheet, err := css_parser.ParseStylesheet(s.Text())
		if err != nil {
			result = err
			return false
		}

		inliner.stylesheets = append(inliner.stylesheets, stylesheet)

		// removes parsed stylesheet
		s.Remove()

		return true
	})

	return result
}

func (inliner *Inliner) collectElementsAndRules() {
	for _, stylesheet := range inliner.stylesheets {
		for _, rule := range stylesheet.Rules {
			if rule.Kind == css_parser.QualifiedRule {
				inliner.handleQualifiedRule(rule)
			} else {
				inliner.rawRules = append(inliner.rawRules, rule)
			}
		}
	}
}

func (inliner *Inliner) handleQualifiedRule(rule *css_parser.CssRule) {
	for _, selector := range rule.Selectors {
		if Inlinable(selector) {
			inliner.doc.Find(selector).Each(func(i int, s *goquery.Selection) {
				// get marker
				eltMarker, exists := s.Attr(elementMarkerAttr)
				if !exists {
					// mark element
					eltMarker = strconv.Itoa(inliner.elementMarker)
					s.SetAttr(elementMarkerAttr, eltMarker)
					inliner.elementMarker++

					// add new element
					inliner.elements[eltMarker] = NewElement(s)
				}

				// add style rule for element
				inliner.elements[eltMarker].addStyleRule(NewStyleRule(selector, rule.Declarations))
			})
		} else {
			// Keep it 'as is'
			inliner.rawRules = append(inliner.rawRules, NewStyleRule(selector, rule.Declarations))
		}
	}
}

func (inliner *Inliner) inlineStyleRules() error {
	for _, element := range inliner.elements {
		// remove marker
		element.element.RemoveAttr(elementMarkerAttr)

		// inline element
		err := element.inline()
		if err != nil {
			return err
		}
	}

	return nil
}

func (inliner *Inliner) computeRawCSS() string {
	result := ""

	for _, rawRule := range inliner.rawRules {
		result += rawRule.String()
		result += "\n"
	}

	return result
}

func (inliner *Inliner) insertRawStylesheet() {
	head := inliner.doc.Find("head")

	// create a new head element if it doesn't exist
	if head.Length() == 0 {
		head = inliner.doc.Find("html").PrependHtml("<head></head>").End()
	} else if head.Length() > 1 {
		head = head.First() // ensure only one head element
	}

	rawCss := inliner.computeRawCSS()
	if rawCss != "" {
		cssNode := &html.Node{
			Type: html.TextNode,
			Data: "\n" + rawCss,
		}

		styleNode := &html.Node{
			Type: html.ElementNode,
			Data: "style",
			Attr: []html.Attribute{{Key: "type", Val: "text/css"}},
		}

		styleNode.AppendChild(cssNode)
		head.AppendNodes(styleNode)
	}
}

func (inliner *Inliner) genHTML() (string, error) {
	return inliner.doc.Html()
}
