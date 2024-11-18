package goapperrors

import (
	"golang.org/x/text/language"
)

// ParseAcceptLanguage parses header Accept-Language from http request.
// This func returns a list of `Tag` objects in priority order.
//
// Example:
//
//	Accept-Language: *
//	Accept-Language: fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5
func ParseAcceptLanguage(acceptLang string) ([]language.Tag, []float32, error) {
	tags, q, err := language.ParseAcceptLanguage(acceptLang)
	if err != nil {
		return nil, nil, Wrap(err)
	}
	return tags, q, nil
}

// ParseAcceptLanguageAsStr parses header Accept-Language from http request.
// This func returns a list of `Tag` strings in priority order.
//
// Example:
//
//	Accept-Language: "fr-CH, fr;q=0.9, en;q=0.8, *;q=0.5"
//	  gives result: ["fr-CH", "fr", "en", "mul"]
func ParseAcceptLanguageAsStr(acceptLang string) ([]string, error) {
	tags, _, err := ParseAcceptLanguage(acceptLang)
	if err != nil {
		return nil, Wrap(err)
	}
	langs := make([]string, 0, len(tags))
	for i := range tags {
		langs = append(langs, tags[i].String())
	}
	return langs, nil
}
