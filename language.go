package contenttype

import "strings"

type Language struct {
	Language string
	Script   string
	Region   string
}

// NewLanguage parses the string and returns an instance of Language struct.
func NewLanguage(s string) Language {
	language, err := ParseLanguage(s)
	if err != nil {
		return Language{}
	}

	return language
}

// ParseLanguage parses the given string as a language and returns it as a Language.
// If the string cannot be parsed an appropriate error is returned.
func ParseLanguage(s string) (Language, error) {
	// RFC 4647, 2.1 Basic Language Range
	language := Language{}
	var consumed bool
	if language.Language, language.Script, language.Region, s, consumed = consumeLanguageTags(skipWhitespaces(s)); !consumed {
		return Language{}, ErrInvalidLanguage
	}

	// there must not be anything left after parsing the header
	if len(s) > 0 {
		return Language{}, ErrInvalidMediaType
	}

	return language, nil
}

func consumeLanguage(s string) (language, remaining string, consumed bool) {
	// RFC 5646, 2.1. Syntax
	for i := 0; i < len(s) && i < 8; i++ {
		if !isAlphaChar(s[i]) {
			if len(s) >= 2 {
				return strings.ToLower(s[:i]), s[i:], true
			} else {
				return "", s, false
			}
		}
	}

	if len(s) >= 2 {
		return strings.ToLower(s), "", len(s) >= 2
	} else {
		return "", s, false
	}
}

func consumeScript(s string) (script, remaining string, consumed bool) {
	// RFC 5646, 2.1. Syntax
	for i := 0; i < len(s) && i < 4; i++ {
		if !isAlphaChar(s[i]) {
			if len(s) >= 2 {
				return strings.ToLower(s[:i]), s[i:], true
			} else {
				return "", s, false
			}
		}
	}

	if len(s) >= 2 {
		return strings.ToLower(s), "", true
	} else {
		return "", s, false
	}
}

func consumeRegion(s string) (region, remaining string, consumed bool) {
	// RFC 5646, 2.1. Syntax
	for i := 0; i < len(s) && i < 3; i++ {
		if !isAlphaChar(s[i]) {
			return strings.ToLower(s[:i]), s[i:], len(s) > 0
		}
	}

	return strings.ToLower(s), "", len(s) > 0
}

func consumeLanguageTags(s string) (language, script, region, remaining string, consumed bool) {
	language, s, consumed = consumeLanguage(s)

	if !consumed {
		return "", "", "", "", false
	}

	if len(s) == 0 {
		return language, "", "", "", true
	}

	if s[0] != '-' {
		return "", "", "", "", false
	}

	/*tag1, s, consumed := consumeTag(s[1:])

	if len(s) == 0 {
		return language, "", tag1, "", true
	}

	if s[0] != '-' {
		return "", "", "", "", false
	}

	tag2, s, consumed := consumeTag(s[1:])

	return language, tag1, tag2, s, true*/

	return "", "", "", "", false
}
