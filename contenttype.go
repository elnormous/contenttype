package contenttype

import (
	"errors"
	"net/http"
	"strings"
)

var InvalidContentTypeError = errors.New("Invalid content type")
var InvalidParameterError = errors.New("Invalid parameter")

func isTokenChar(r rune) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return strings.ContainsRune("!#$%&'*+-.^_`|~", r) ||
		(r >= 0x30 && r <= 0x39) || // DIGIT
		(r >= 0x41 && r <= 0x5A) || (r >= 0x61 && r <= 0x7A) // ALPHA
}

func isNotTokenChar(r rune) bool {
	return !isTokenChar(r)
}

func consumeToken(s string) (token, remaining string) {
	index := strings.IndexFunc(s, isNotTokenChar)
	if index == -1 {
		return s, ""
	} else {
		return s[:index], s[index:]
	}
}

func isWhiteSpaceChar(r rune) bool {
	// RFC 7230, 3.2.3. Whitespace
	return r == 0x09 || r == 0x20 // HTAB or SP
}

func skipWhiteSpaces(s string) string {
	return strings.TrimLeftFunc(s, isWhiteSpaceChar)
}

func GetMediaType(request *http.Request) (string, map[string]string, error) {
	contentTypes := request.Header.Values("Content-Type")

	if len(contentTypes) == 0 {
		return "", map[string]string{}, nil
	}

	s := skipWhiteSpaces(contentTypes[0])

	var supertype string
	supertype, s = consumeToken(s)

	if len(supertype) == 0 {
		return "", nil, InvalidContentTypeError
	}

	if len(s) == 0 || s[0] != '/' {
		return "", nil, InvalidContentTypeError
	}

	s = s[1:] // skip the slash

	var subtype string
	subtype, s = consumeToken(s)

	if len(subtype) == 0 {
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
		key, s = consumeToken(s)

		if len(s) == 0 || s[0] != '=' {
			return "", nil, InvalidParameterError
		}

		s = s[1:] // skip the equal sign
		var value string
		value, s = consumeToken(s)

		parameters[key] = value

		s = skipWhiteSpaces(s)
	}

	var stringBuilder strings.Builder
	stringBuilder.WriteString(strings.ToLower(supertype))
	stringBuilder.WriteByte('/')
	stringBuilder.WriteString(strings.ToLower(subtype))

	return stringBuilder.String(), parameters, nil
}
