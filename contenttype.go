package contenttype

import (
	"errors"
	"net/http"
	"strings"
)

var InvalidContentTypeError = errors.New("Invalid content type")
var InvalidParameterError = errors.New("Invalid parameter")

func isWhiteSpaceChar(r rune) bool {
	// RFC 7230, 3.2.3. Whitespace
	return r == 0x09 || r == 0x20 // HTAB or SP
}

func isTokenChar(r rune) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return strings.ContainsRune("!#$%&'*+-.^_`|~", r) ||
		// RFC 5234, Appendix B.1. Core Rules
		(r >= 0x30 && r <= 0x39) || // DIGIT
		(r >= 0x41 && r <= 0x5A) || (r >= 0x61 && r <= 0x7A) // ALPHA
}

func isNotTokenChar(r rune) bool {
	return !isTokenChar(r)
}

func isVisibleChar(r rune) bool {
	// RFC 5234, Appendix B.1. Core Rules
	return r >= 0x21 && r <= 0x7E
}

func isObsoleteTextChar(r rune) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return r >= 0x80 && r <= 0xFF
}

func isQuotedTextChar(r rune) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return r == 0x09 || r == 0x20 || // HTAB or SP
		r == 0x21 ||
		(r >= 0x23 && r <= 0x5B) ||
		(r >= 0x5D && r <= 0x7E) ||
		isObsoleteTextChar(r)
}

func isQuotedPairChar(r rune) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return r == 0x09 || r == 0x20 || // HTAB or SP
		isVisibleChar(r) ||
		isObsoleteTextChar(r)
}

func consumeToken(s string) (token, remaining string, consumed bool) {
	index := strings.IndexFunc(s, isNotTokenChar)
	if index == -1 {
		return s, "", len(s) > 0
	} else {
		return s[:index], s[index:], index > 0
	}
}

func consumeQuotedString(s string) (token, remaining string, consumed bool) {
	if len(s) == 0 || s[0] != '"' {
		return "", s, false
	}

	var stringBuilder strings.Builder

	for index := 1; index < len(s); index++ {
		if s[index] == '"' {
			return stringBuilder.String(), s[index+1:], true
		}

		if s[index] == '\\' {
			index++
			if len(s) <= index || !isQuotedPairChar(rune(s[index])) {
				return "", s, false
			}

			stringBuilder.WriteByte(s[index])
		} else {
			if !isQuotedTextChar(rune(s[index])) {
				return "", s, false
			}
			stringBuilder.WriteByte(s[index])
		}
	}

	return "", s, false
}

func skipWhiteSpaces(s string) string {
	return strings.TrimLeftFunc(s, isWhiteSpaceChar)
}

func GetMediaType(request *http.Request) (string, map[string]string, error) {
	// RFC 7231, 3.1.1.1. Media Type
	contentTypes := request.Header.Values("Content-Type")

	if len(contentTypes) == 0 {
		return "", map[string]string{}, nil
	}

	s := skipWhiteSpaces(contentTypes[0])

	var ok bool
	var supertype string
	supertype, s, ok = consumeToken(s)
	if !ok {
		return "", nil, InvalidContentTypeError
	}

	if len(s) == 0 || s[0] != '/' {
		return "", nil, InvalidContentTypeError
	}

	s = s[1:] // skip the slash

	var subtype string
	subtype, s, ok = consumeToken(s)
	if !ok {
		return "", nil, InvalidContentTypeError
	}

	s = skipWhiteSpaces(s)

	parameters := make(map[string]string)

	for len(s) != 0 {
		if s[0] != ';' {
			return "", nil, InvalidParameterError
		}

		s = s[1:] // skip the semicolon
		s = skipWhiteSpaces(s)

		var key string
		key, s, ok = consumeToken(s)
		if !ok {
			return "", nil, InvalidParameterError
		}

		if len(s) == 0 || s[0] != '=' {
			return "", nil, InvalidParameterError
		}

		s = s[1:] // skip the equal sign

		var value string
		if len(s) != 0 && s[0] == '"' { // opening quote
			value, s, ok = consumeQuotedString(s)

			if !ok {
				return "", nil, InvalidParameterError
			}

		} else {
			value, s, ok = consumeToken(s)
			if !ok {
				return "", nil, InvalidParameterError
			}
		}

		parameters[key] = value

		s = skipWhiteSpaces(s)
	}

	var stringBuilder strings.Builder
	stringBuilder.WriteString(strings.ToLower(supertype))
	stringBuilder.WriteByte('/')
	stringBuilder.WriteString(strings.ToLower(subtype))

	return stringBuilder.String(), parameters, nil
}
