package contenttype

import (
	"errors"
	"net/http"
	"strings"
)

var InvalidMediaTypeError = errors.New("Invalid media type")
var InvalidMediaRangeError = errors.New("Invalid media range")
var InvalidParameterError = errors.New("Invalid parameter")
var InvalidExtensionParameterError = errors.New("Invalid extension parameter")
var NoAvailableTypeGivenError = errors.New("No available type given")
var NoAcceptableTypeFoundError = errors.New("No acceptable type found")
var InvalidWeightError = errors.New("Invalid wieght")

type Parameters = map[string]string

type MediaType struct {
	Type       string
	Subtype    string
	Parameters Parameters
}

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

func skipWhiteSpaces(s string) string {
	for i := 0; i < len(s); i++ {
		if !isWhiteSpaceChar(s[i]) {
			return s[i:]
		}
	}

	return ""
}

func consumeToken(s string) (token, remaining string, consumed bool) {
	for i := 0; i < len(s); i++ {
		if !isTokenChar(s[i]) {
			return strings.ToLower(s[:i]), s[i:], i > 0
		}
	}

	return strings.ToLower(s), "", len(s) > 0
}

func consumeQuotedString(s string) (token, remaining string, consumed bool) {
	var stringBuilder strings.Builder

	index := 0
	for ; index < len(s); index++ {
		if s[index] == '\\' {
			index++
			if len(s) <= index || !isQuotedPairChar(s[index]) {
				return "", s, false
			}

			stringBuilder.WriteByte(s[index])
		} else if isQuotedTextChar(s[index]) {
			stringBuilder.WriteByte(s[index])
		} else {
			break
		}
	}

	return strings.ToLower(stringBuilder.String()), s[index:], true
}

func consumeType(s string) (string, string, string, bool) {
	// RFC 7231, 3.1.1.1. Media Type
	s = skipWhiteSpaces(s)

	var t, subt string
	var ok bool
	t, s, ok = consumeToken(s)
	if !ok {
		return "", "", s, false
	}

	if len(s) == 0 || s[0] != '/' {
		return "", "", s, false
	}

	s = s[1:] // skip the slash

	subt, s, ok = consumeToken(s)
	if !ok {
		return "", "", s, false
	}

	if t == "*" && subt != "*" {
		return "", "", s, false
	}

	s = skipWhiteSpaces(s)

	return t, subt, s, true
}

func consumeParameter(s string) (string, string, string, bool) {
	// RFC 7231, 3.1.1.1. Media Type
	s = skipWhiteSpaces(s)

	var ok bool
	var key string
	key, s, ok = consumeToken(s)
	if !ok {
		return "", "", s, false
	}

	if len(s) == 0 || s[0] != '=' {
		return "", "", s, false
	}

	s = s[1:] // skip the equal sign

	var value string
	if len(s) > 0 && s[0] == '"' {
		s = s[1:] // skip the opening quote

		value, s, ok = consumeQuotedString(s)
		if !ok {
			return "", "", s, false
		}

		if len(s) == 0 || s[0] != '"' {
			return "", "", s, false
		}

		s = s[1:] // skip the closing quote

	} else {
		value, s, ok = consumeToken(s)
		if !ok {
			return "", "", s, false
		}
	}

	s = skipWhiteSpaces(s)

	return key, value, s, true
}

func checkWeight(s string) bool {
	// RFC 7231, 5.3.1. Quality Values
	if len(s) == 0 || (s[0] != '0' && s[0] != '1') {
		return false
	}

	if len(s) > 1 {
		if s[1] != '.' || len(s) > 5 {
			return false
		}

		for index := 2; index < len(s); index++ {
			if !isDigitChar(s[index]) ||
				(s[0] == '1' && s[index] != '0') { // weight can not be greater than 1
				return false
			}
		}
	}

	return true
}

func GetMediaType(request *http.Request) (MediaType, error) {
	// RFC 7231, 3.1.1.5. Content-Type
	contentTypeHeaders := request.Header.Values("Content-Type")

	if len(contentTypeHeaders) == 0 {
		return MediaType{}, nil
	}

	s := contentTypeHeaders[0]
	mediaType := MediaType{}
	var consumed bool

	mediaType.Type, mediaType.Subtype, s, consumed = consumeType(s)

	if !consumed {
		return MediaType{}, InvalidMediaTypeError
	}

	mediaType.Parameters = make(Parameters)

	for len(s) > 0 && s[0] == ';' {
		s = s[1:] // skip the semicolon

		key, value, remaining, consumed := consumeParameter(s)

		if !consumed {
			return MediaType{}, InvalidParameterError
		}

		s = remaining

		mediaType.Parameters[key] = value
	}

	if len(s) > 0 {
		return MediaType{}, InvalidMediaTypeError
	}

	return mediaType, nil
}

func GetAcceptableMediaType(request *http.Request, availableMediaTypes []MediaType) (MediaType, error) {
	// RFC 7231, 5.3.2. Accept
	if len(availableMediaTypes) == 0 {
		return MediaType{}, NoAvailableTypeGivenError
	}

	acceptHeaders := request.Header.Values("Accept")

	if len(acceptHeaders) == 0 {
		return availableMediaTypes[0], nil
	}

	s := acceptHeaders[0]

	resultMediaType := MediaType{}
	resultWeight := ""
	acceptableTypeFound := false

	for mediaTypeCount := 0; len(s) > 0; mediaTypeCount++ {
		if mediaTypeCount > 0 {
			// every media type after the first one must start with a comma
			if s[0] != ',' {
				break
			}
			s = s[1:] // skip the comma
		}

		mediaType := MediaType{}
		var consumed bool
		mediaType.Type, mediaType.Subtype, s, consumed = consumeType(s)
		if !consumed {
			return MediaType{}, InvalidMediaTypeError
		}

		parameters := make(Parameters)
		currentWeight := "1"

		// media type parameters
		for len(s) > 0 && s[0] == ';' {
			s = s[1:] // skip the semicolon

			var key, value string
			key, value, s, consumed = consumeParameter(s)

			if !consumed {
				return MediaType{}, InvalidParameterError
			}

			parameters[key] = value

			if key == "q" {
				if !checkWeight(value) {
					return MediaType{}, InvalidWeightError
				}

				currentWeight = value
				break // "q" parameter separates media type parameters from Accept extension parameters
			}
		}

		// extension parameters
		for len(s) > 0 && s[0] == ';' {
			s = s[1:] // skip the semicolon

			key, value, remaining, consumed := consumeParameter(s)

			if !consumed {
				return MediaType{}, InvalidParameterError
			}

			s = remaining

			parameters[key] = value // TODO: store in a separate map
		}

		for _, availableMediaType := range availableMediaTypes {
			if (mediaType.Type == "*" || mediaType.Type == availableMediaType.Type) &&
				(mediaType.Subtype == "*" || mediaType.Subtype == availableMediaType.Subtype) {

				if currentWeight > "0" && // 0 means "not acceptable"
					(!acceptableTypeFound || currentWeight > resultWeight) {
					resultMediaType = availableMediaType
					resultMediaType.Parameters = parameters
					resultWeight = currentWeight
					acceptableTypeFound = true
				}
			}
		}

		s = skipWhiteSpaces(s)
	}

	if len(s) > 0 {
		return MediaType{}, InvalidMediaRangeError
	}

	if !acceptableTypeFound {
		return MediaType{}, NoAcceptableTypeFoundError
	}

	return resultMediaType, nil
}
