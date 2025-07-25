package cssinliner

import (
	"slices"
	"sort"

	"github.com/PuerkitoBio/goquery"
	cssparser "go.baoshuo.dev/cssparser"
)

type Element struct {
	element       *goquery.Selection // The goquery handler
	styleRules    []*StyleRule       // The style rules to apply on that element
	parserOptions []cssparser.ParserOption
}

type AttrToStyleRule struct {
	styleName string   // The name of the style property
	elements  []string // The elements that have this style property
}

var attrToStyle = map[string][]*AttrToStyleRule{
	"align": {
		{
			"float",
			[]string{"img"},
		},
		{
			"text-align",
			[]string{"h1", "h2", "h3", "h4", "h5", "h6", "p", "div", "blockquote", "tr", "th", "td"},
		},
	},
	"bgcolor": {{
		"background-color",
		[]string{"body", "table", "tr", "th", "td"},
	}},
	"background": {{
		"background-image",
		[]string{"table"},
	}},
	"valign": {{
		"vertical-align",
		[]string{"th", "td"},
	}},
	"width": {{
		"width",
		[]string{"img", "table", "th", "td"},
	}},
	"height": {{
		"height",
		[]string{"img", "table", "th", "td"},
	}},
}

func NewElement(element *goquery.Selection, parserOptions ...cssparser.ParserOption) *Element {
	return &Element{
		element:       element,
		parserOptions: parserOptions,
	}
}

func (element *Element) addStyleRule(styleRule *StyleRule) {
	element.styleRules = append(element.styleRules, styleRule)
}

func (element *Element) inline() error {
	// compute declarations
	declarations, err := element.computeDeclarations()
	if err != nil {
		return err
	}

	// set style attribute
	styleValue := computeStyleValue(declarations)
	if styleValue != "" {
		element.element.SetAttr("style", styleValue)
	}

	return nil
}

func (element *Element) computeDeclarations() ([]*cssparser.Declaration, error) {
	result := []*cssparser.Declaration{}

	styles := make(map[string]*StyleDeclaration)

	// First: parsed stylesheets rules
	mergeStyleDeclarations(element.styleRules, styles)

	// Second: attributes
	attrRules, err := element.parseAttributes()
	if err != nil {
		return result, err
	}

	// Then: inline rules
	inlineRules, err := element.parseInlineStyle()
	if err != nil {
		return result, err
	}

	mergeStyleDeclarations(inlineRules, styles)
	mergeStyleDeclarations(attrRules, styles)

	// map to array
	for _, styleDecl := range styles {
		result = append(result, styleDecl.Declaration)
	}

	// sort declarations by property name
	sort.Sort(cssparser.DeclarationsByProperty(result))

	return result, nil
}

func (element *Element) parseAttributes() ([]*StyleRule, error) {
	result := []*StyleRule{}
	declarations := []*cssparser.Declaration{}

	for attr, rules := range attrToStyle {
		value, exists := element.element.Attr(attr)
		if !exists || value == "" {
			continue
		}

		for _, rule := range rules {
			if slices.Contains(rule.elements, element.element.Nodes[0].Data) {
				declarations = append(declarations, &cssparser.Declaration{
					Property: rule.styleName,
					Value:    value,
				})
			}
		}
	}

	if len(declarations) > 0 {
		result = append(result, NewStyleRule(inlineFakeSelector, declarations))
	}

	return result, nil
}

func (element *Element) parseInlineStyle() ([]*StyleRule, error) {
	result := []*StyleRule{}

	styleValue, exists := element.element.Attr("style")
	if !exists || styleValue == "" {
		return result, nil
	}

	declarations, err := cssparser.ParseDeclarations(styleValue, element.parserOptions...)
	if err != nil {
		return result, err
	}

	result = append(result, NewStyleRule(inlineFakeSelector, declarations))

	return result, nil
}
