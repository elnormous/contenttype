package contenttype

import (
	"errors"
	"net/http"
	"strings"
)

var InvalidMediaTypeError = errors.New("Invalid media type")
var InvalidMediaRangeError = errors.New("Invalid media range")
var InvalidParameterError = errors.New("Invalid parameter")
var NoAvailableTypeGivenError = errors.New("No available type given")
var NoAcceptableTypeFoundError = errors.New("No acceptable type found")

func isWhiteSpaceChar(c byte) bool {
	// RFC 7230, 3.2.3. Whitespace
	return c == 0x09 || c == 0x20 // HTAB or SP
}

func isDigitChar(c byte) bool {
	// RFC 5234, Appendix B.1. Core Rules
	return c >= 0x30 && c <= 0x39
}

func isAlphaChar(c byte) bool {
	// RFC 5234, Appendix B.1. Core Rules
	return (c >= 0x41 && c <= 0x5A) || (c >= 0x61 && c <= 0x7A)
}

func isTokenChar(c byte) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return c == '!' || c == '#' || c == '$' || c == '%' || c == '&' || c == '\'' || c == '*' ||
		c == '+' || c == '-' || c == '.' || c == '^' || c == '_' || c == '`' || c == '|' || c == '~' ||
		isDigitChar(c) ||
		isAlphaChar(c)
}

func isVisibleChar(c byte) bool {
	// RFC 5234, Appendix B.1. Core Rules
	return c >= 0x21 && c <= 0x7E
}

func isObsoleteTextChar(c byte) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return c >= 0x80 && c <= 0xFF
}

func isQuotedTextChar(c byte) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return c == 0x09 || c == 0x20 || // HTAB or SP
		c == 0x21 ||
		(c >= 0x23 && c <= 0x5B) ||
		(c >= 0x5D && c <= 0x7E) ||
		isObsoleteTextChar(c)
}

func isQuotedPairChar(c byte) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return c == 0x09 || c == 0x20 || // HTAB or SP
		isVisibleChar(c) ||
		isObsoleteTextChar(c)
}

func consumeToken(s string) (token, remaining string, consumed bool) {
	for i := 0; i < len(s); i++ {
		if !isTokenChar(s[i]) {
			return s[:i], s[i:], i > 0
		}
	}

	return s, "", len(s) > 0
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
			if len(s) <= index || !isQuotedPairChar(s[index]) {
				return "", s, false
			}

			stringBuilder.WriteByte(s[index])
		} else {
			if !isQuotedTextChar(s[index]) {
				return "", s, false
			}
			stringBuilder.WriteByte(s[index])
		}
	}

	return "", s, false
}

func skipWhiteSpaces(s string) string {
	for i := 0; i < len(s); i++ {
		if !isWhiteSpaceChar(s[i]) {
			return s[i:]
		}
	}

	return ""
}

func consumeMediaType(s string) (string, string, bool) {
	// RFC 7231, 3.1.1.1. Media Type
	s = skipWhiteSpaces(s)

	var ok bool
	var supertype string
	supertype, s, ok = consumeToken(s)
	if !ok {
		return "", s, false
	}
	supertype = strings.ToLower(supertype)

	if len(s) == 0 || s[0] != '/' {
		return "", s, false
	}

	s = s[1:] // skip the slash

	var subtype string
	subtype, s, ok = consumeToken(s)
	if !ok {
		return "", s, false
	}
	subtype = strings.ToLower(subtype)

	s = skipWhiteSpaces(s)

	return supertype + "/" + subtype, s, true
}

func GetMediaType(request *http.Request) (string, map[string]string, error) {
	// RFC 7231, 3.1.1.5. Content-Type
	contentTypeHeaders := request.Header.Values("Content-Type")

	if len(contentTypeHeaders) == 0 {
		return "", map[string]string{}, nil
	}

	mediaType, s, consumed := consumeMediaType(contentTypeHeaders[0])

	if !consumed {
		return "", map[string]string{}, InvalidMediaTypeError
	}

	parameters := make(map[string]string)

	for len(s) > 0 && s[0] == ';' {
		s = s[1:] // skip the semicolon
		s = skipWhiteSpaces(s)

		var ok bool
		var key string
		key, s, ok = consumeToken(s)
		if !ok {
			return "", map[string]string{}, InvalidParameterError
		}
		key = strings.ToLower(key)

		if len(s) == 0 || s[0] != '=' {
			return "", map[string]string{}, InvalidParameterError
		}

		s = s[1:] // skip the equal sign

		var value string
		if len(s) > 0 && s[0] == '"' { // opening quote
			value, s, ok = consumeQuotedString(s)

			if !ok {
				return "", map[string]string{}, InvalidParameterError
			}

		} else {
			value, s, ok = consumeToken(s)
			if !ok {
				return "", map[string]string{}, InvalidParameterError
			}
		}
		value = strings.ToLower(value)

		parameters[key] = value

		s = skipWhiteSpaces(s)
	}

	if len(s) > 0 {
		return "", map[string]string{}, InvalidMediaTypeError
	}

	return mediaType, parameters, nil
}

func GetAcceptableMediaType(request *http.Request, availableMediaTypes []string) (string, map[string]string, error) {
	// RFC 7231, 5.3.2. Accept
	if len(availableMediaTypes) == 0 {
		return "", map[string]string{}, NoAvailableTypeGivenError
	}

	acceptHeaders := request.Header.Values("Accept")

	if len(acceptHeaders) == 0 {
		return availableMediaTypes[0], map[string]string{}, nil
	}

	s := acceptHeaders[0]

	for len(s) > 0 {
		mediaType, remaining, consumed := consumeMediaType(s)

		s = remaining

		if !consumed {
			return "", map[string]string{}, InvalidMediaTypeError
		}

		for _, availableMediaType := range availableMediaTypes {
			if mediaType == availableMediaType {
				return mediaType, map[string]string{}, nil
			}
		}

		if len(s) > 0 {
			if s[0] != ',' {
				return "", map[string]string{}, InvalidMediaRangeError
			}

			s = s[1:] // skip the comma
		}
	}

	if len(s) > 0 {
		return "", map[string]string{}, InvalidMediaTypeError
	}

	return "", map[string]string{}, NoAcceptableTypeFoundError
}
