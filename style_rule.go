package css_inliner

import (
	"fmt"
	"regexp"

	css_parser "go.baoshuo.dev/css-parser"
)

const (
	inlineFakeSelector = "*INLINE*"

	cssIdentifierRegexp = `-?([_a-zA-Z]|[\x{00A0}-\x{FFFF}]|(\\[^\r\n\f0-9a-fA-F]))([_a-zA-Z0-9-]|[\x{00A0}-\x{FFFF}]|(\\[^\r\n\f0-9a-fA-F]))*`
)

var (
	idSelectorRegexp            = regexp.MustCompile(`#` + cssIdentifierRegexp)
	classSelectorRegExp         = regexp.MustCompile(`\.` + cssIdentifierRegexp)
	attributeSelectorRegexp     = regexp.MustCompile(`\[\s*([_a-zA-Z0-9-]+)\s*([~|^$*]?=)?\s*("([^"]*)"|'([^']*)'|([_a-zA-Z0-9-]+))?\s*\]`)
	pseudoClassSelectorRegExp   = regexp.MustCompile(`:([_a-zA-Z-]|[\x{00A0}-\x{FFFF}]|\\.)([_a-zA-Z0-9-]|[\x{00A0}-\x{FFFF}]|\\.)*(?:\((?:[^"']|\"[^\"]*\"|'[^']*')+\))?`)
	pseudoElementSelectorRegExp = regexp.MustCompile(`::(-?([_a-zA-Z]|[\x{00A0}-\x{FFFF}]|\\.)([_a-zA-Z0-9-]|[\x{00A0}-\x{FFFF}]|\\.)*)`)

	// Because of negative lookbehind is not supported in golang, we cannot use `(?<![.#:_a-zA-Z0-9-])`
	// in the regexp, so we removes selectors which are previously matched by idSelectorRegexp,
	// classSelectorRegExp, pseudoClassSelectorRegExp and pseudoElementSelectorRegExp to ensure we
	// only match type selectors.
	//
	// typeSelectorRegExp          = regexp.MustCompile(`(?<![.#:_a-zA-Z0-9-])((?:([_a-zA-Z]|[\x{00A0}-\x{FFFF}]|\\.)([_a-zA-Z0-9-]|[\x{00A0}-\x{FFFF}]|\\.)*\|)?([_a-zA-Z]|[\x{00A0}-\x{FFFF}]|\\.)([_a-zA-Z0-9-]|[\x{00A0}-\x{FFFF}]|\\.)*)`)
	typeSelectorRegExp = regexp.MustCompile(`((?:([_a-zA-Z]|[\x{00A0}-\x{FFFF}]|\\.)([_a-zA-Z0-9-]|[\x{00A0}-\x{FFFF}]|\\.)*\|)?([_a-zA-Z]|[\x{00A0}-\x{FFFF}]|\\.)([_a-zA-Z0-9-]|[\x{00A0}-\x{FFFF}]|\\.)*)`)
)

// StyleRule represents a Qualifier Rule for a uniq selector
type StyleRule struct {
	Selector     string                       // The style rule selector
	Declarations []*css_parser.CssDeclaration // The style rule properties
	Specificity  int                          // Selector specificity
}

func NewStyleRule(selector string, declarations []*css_parser.CssDeclaration) *StyleRule {
	return &StyleRule{
		Selector:     selector,
		Declarations: declarations,
		Specificity:  ComputeSpecificity(selector),
	}
}

func (styleRule *StyleRule) String() string {
	result := ""

	result += styleRule.Selector

	if len(styleRule.Declarations) == 0 {
		result += ";"
	} else {
		result += " {\n"

		for _, decl := range styleRule.Declarations {
			result += fmt.Sprintf("  %s\n", decl.String())
		}

		result += "}"
	}

	return result
}

// ComputeSpecificity computes style rule specificity
//
// cf. http://www.w3.org/TR/selectors/#specificity
func ComputeSpecificity(selector string) int {
	result := 0

	if selector == inlineFakeSelector {
		result += 1000
	}

	idSelectors := idSelectorRegexp.FindAllStringSubmatch(selector, -1)
	selector = idSelectorRegexp.ReplaceAllString(selector, "")

	classSelectors := classSelectorRegExp.FindAllStringSubmatch(selector, -1)
	selector = classSelectorRegExp.ReplaceAllString(selector, "")

	attributeSelectors := attributeSelectorRegexp.FindAllStringSubmatch(selector, -1)
	selector = attributeSelectorRegexp.ReplaceAllString(selector, "")

	pseudoClassSelectors := pseudoClassSelectorRegExp.FindAllStringSubmatch(selector, -1)
	selector = pseudoClassSelectorRegExp.ReplaceAllString(selector, "")

	pseudoElementSelectors := pseudoElementSelectorRegExp.FindAllStringSubmatch(selector, -1)
	selector = pseudoElementSelectorRegExp.ReplaceAllString(selector, "")

	typeSelectors := typeSelectorRegExp.FindAllStringSubmatch(selector, -1)
	// selector = typeSelectorRegExp.ReplaceAllString(selector, "")

	a := len(idSelectors)
	b := len(classSelectors) + len(attributeSelectors) + len(pseudoClassSelectors)
	c := len(typeSelectors) + len(pseudoElementSelectors)

	result += a*100 + b*10 + c

	return result
}
