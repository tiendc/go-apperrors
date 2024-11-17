package goapperrors

// Language represents a language.
// Language values can be anything and can be set by client code.
// For example, Language("string") or Language(language.Tag from "golang.org/x/text/language").
type Language any

const (
	LanguageEn = "en"
	LanguageFr = "fr"
	LanguageDe = "de"
	LanguageEs = "es"
	LanguageIt = "it"
	LanguagePt = "pt"
	LanguageRu = "ru"
	LanguageZh = "zh"
	LanguageJa = "ja"
	LanguageKo = "ko"
	LanguageAr = "ar"
	LanguageHi = "hi"
)

type TranslationFunc func(lang Language, key string, params map[string]any) (string, error)
