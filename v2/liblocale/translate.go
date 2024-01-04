package liblocale

import (
	"regexp"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	regexTemplate      = regexp.MustCompile(`^{{(.+)}}$`) // Regex for detect {{syntax}}
	regexChainTemplate = regexp.MustCompile(`^<T (.+)>$`) // Regex for detect chain template <T syntax>
)

// Translate
func Translate(loc *i18n.Localizer, syntax string, data interface{}) string {
	// Translate template data
	templateData := translateTemplateData(loc, data)

	// Translate syntax
	s, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID:    syntax,
		TemplateData: templateData,
	})
	if err != nil {
		return fallbackTranslation(loc, s)
	}

	// Check translated syntax is chained syntax <<SYNTAX>>
	isChain, syntaxChain := checkSyntax(s, regexChainTemplate)
	if isChain {
		s = Translate(loc, syntaxChain, data)
	}
	return s
}

// translateTemplateData provide translation if template data is given
func translateTemplateData(loc *i18n.Localizer, data interface{}) map[string]string {
	templateData := map[string]string{}
	if t, ok := data.(map[string]string); ok {
		for k, v := range t {
			t[k] = translateWithModifiers(loc, v)
		}
		templateData = t
	}
	return templateData
}

// translateWithModifiers translates and applies any modifiers to the text
func translateWithModifiers(loc *i18n.Localizer, text string) string {
	isSyntax, s := checkSyntax(text, regexTemplate)
	if !isSyntax {
		return text
	}

	txt := strings.Fields(s)
	txtLen := len(txt)
	localeTxt := Translate(loc, txt[0], nil)
	switch {
	case txtLen > 1 && txt[1] == "upper":
		return strings.ToUpper(localeTxt)
	case txtLen > 1 && txt[1] == "lower":
		return strings.ToLower(localeTxt)
	case txtLen > 1 && txt[1] == "title":
		return cases.Title(language.English).String(localeTxt)
	default:
		return localeTxt
	}
}

// fallbackTranslation provides a translation for when the main translation fails
func fallbackTranslation(loc *i18n.Localizer, syntax string) string {
	s, err := loc.Localize(&i18n.LocalizeConfig{
		MessageID: "common.error.i18n.syntax",
		TemplateData: map[string]interface{}{
			"syntax": syntax,
		},
	})
	if err != nil {
		return syntax // Return the original syntax as a last resort
	}
	return s
}

// checkSyntax verify syntax with regex
func checkSyntax(input string, reg *regexp.Regexp) (bool, string) {
	// Check if the input matches the pattern
	matches := reg.FindStringSubmatch(input)

	// If it's a match, return true and the extracted key
	if len(matches) > 1 {
		return true, matches[1]
	}

	// If it's not a match, return false and an empty string
	return false, ""
}
