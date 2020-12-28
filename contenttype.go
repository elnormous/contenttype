package contenttype

import (
	"errors"
	"net/http"
	"strings"
)

var ErrInvalidMediaType = errors.New("Invalid media type")
var ErrInvalidMediaRange = errors.New("Invalid media range")
var ErrInvalidParameter = errors.New("Invalid parameter")
var ErrInvalidExtensionParameter = errors.New("Invalid extension parameter")
var ErrNoAcceptableTypeFound = errors.New("No acceptable type found")
var ErrNoAvailableTypeGiven = errors.New("No available type given")
var ErrInvalidWeight = errors.New("Invalid wieght")

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
	// RFC 7230, 3.2.3. Whitespace
	for i := 0; i < len(s); i++ {
		if !isWhiteSpaceChar(s[i]) {
			return s[i:]
		}
	}

	return ""
}

func consumeToken(s string) (token, remaining string, consumed bool) {
	// RFC 7230, 3.2.6. Field Value Components
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
	var consumed bool
	t, s, consumed = consumeToken(s)
	if !consumed {
		return "", "", s, false
	}

	if len(s) == 0 || s[0] != '/' {
		return "", "", s, false
	}

	s = s[1:] // skip the slash

	subt, s, consumed = consumeToken(s)
	if !consumed {
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

	var consumed bool
	var key string
	key, s, consumed = consumeToken(s)
	if !consumed {
		return "", "", s, false
	}

	if len(s) == 0 || s[0] != '=' {
		return "", "", s, false
	}

	s = s[1:] // skip the equal sign

	var value string
	if len(s) > 0 && s[0] == '"' {
		s = s[1:] // skip the opening quote

		value, s, consumed = consumeQuotedString(s)
		if !consumed {
			return "", "", s, false
		}

		if len(s) == 0 || s[0] != '"' {
			return "", "", s, false
		}

		s = s[1:] // skip the closing quote

	} else {
		value, s, consumed = consumeToken(s)
		if !consumed {
			return "", "", s, false
		}
	}

	s = skipWhiteSpaces(s)

	return key, value, s, true
}

func getWeight(s string) (int, bool) {
	// RFC 7231, 5.3.1. Quality Values
	result := 0
	multiplier := 1000
	for i := 0; i < len(s); i++ {
		if i == 0 {
			// the first character must be 0 or 1
			if s[i] != '0' && s[i] != '1' {
				return 0, false
			}

			result = int(s[i]-'0') * multiplier
			multiplier /= 10
		} else if i == 1 {
			// the second character must be a dot
			if s[i] != '.' {
				return 0, false
			}
		} else if i > 4 { // the string can not be longer than 5 characters
			return 0, false
		} else {
			// the remaining characters must be digits and the value can not be greater than 1.000
			if (s[0] == '1' && s[i] != '0') ||
				(s[i] < '0' || s[i] > '9') {
				return 0, false
			}

			result += int(s[i]-'0') * multiplier
			multiplier /= 10
		}
	}

	return result, true
}

func compareMediaTypes(checkMediaType, mediaType MediaType) bool {
	if (checkMediaType.Type == "*" || checkMediaType.Type == mediaType.Type) &&
		(checkMediaType.Subtype == "*" || checkMediaType.Subtype == mediaType.Subtype) {

		for checkKey, checkValue := range checkMediaType.Parameters {
			if value, found := mediaType.Parameters[checkKey]; !found || value != checkValue {
				return false
			}
		}

		return true
	}

	return false
}

func getPrecedence(checkMediaType, mediaType MediaType) bool {
	if len(mediaType.Type) == 0 || len(mediaType.Subtype) == 0 { // not set
		return true
	}

	if (mediaType.Type == "*" && checkMediaType.Type != "*") ||
		(mediaType.Subtype == "*" && checkMediaType.Subtype != "*") ||
		(len(mediaType.Parameters) < len(checkMediaType.Parameters)) {
		return true
	}

	return false
}

func NewMediaType(s string) MediaType {
	mediaType := MediaType{}
	var consumed bool
	mediaType.Type, mediaType.Subtype, s, consumed = consumeType(s)
	if !consumed {
		return MediaType{}
	}

	mediaType.Parameters = make(Parameters)

	for len(s) > 0 && s[0] == ';' {
		s = s[1:] // skip the semicolon

		key, value, remaining, consumed := consumeParameter(s)
		if !consumed {
			return MediaType{}
		}

		s = remaining

		mediaType.Parameters[key] = value
	}

	return mediaType
}

func (mediaType *MediaType) String() string {
	var stringBuilder strings.Builder

	if len(mediaType.Type) > 0 || len(mediaType.Subtype) > 0 {
		stringBuilder.WriteString(mediaType.Type)
		stringBuilder.WriteByte('/')
		stringBuilder.WriteString(mediaType.Subtype)
	}

	for key, value := range mediaType.Parameters {
		stringBuilder.WriteByte(';')
		stringBuilder.WriteString(key)
		stringBuilder.WriteByte('=')
		stringBuilder.WriteString(value)
	}

	return stringBuilder.String()
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
		return MediaType{}, ErrInvalidMediaType
	}

	mediaType.Parameters = make(Parameters)

	for len(s) > 0 && s[0] == ';' {
		s = s[1:] // skip the semicolon

		key, value, remaining, consumed := consumeParameter(s)
		if !consumed {
			return MediaType{}, ErrInvalidParameter
		}

		s = remaining

		mediaType.Parameters[key] = value
	}

	// there must not be anything left after parsing the header
	if len(s) > 0 {
		return MediaType{}, ErrInvalidMediaType
	}

	return mediaType, nil
}

func GetAcceptableMediaType(request *http.Request, availableMediaTypes []MediaType) (MediaType, Parameters, error) {
	// RFC 7231, 5.3.2. Accept
	if len(availableMediaTypes) == 0 {
		return MediaType{}, Parameters{}, ErrNoAvailableTypeGiven
	}

	acceptHeaders := request.Header.Values("Accept")
	if len(acceptHeaders) == 0 {
		return availableMediaTypes[0], Parameters{}, nil
	}

	s := acceptHeaders[0]

	weights := make([]struct {
		mediaType           MediaType
		extensionParameters Parameters
		weight              int
		order               int
	}, len(availableMediaTypes))

	for mediaTypeCount := 0; len(s) > 0; mediaTypeCount++ {
		if mediaTypeCount > 0 {
			// every media type after the first one must start with a comma
			if s[0] != ',' {
				break
			}
			s = s[1:] // skip the comma
		}

		acceptableMediaType := MediaType{}
		var consumed bool
		acceptableMediaType.Type, acceptableMediaType.Subtype, s, consumed = consumeType(s)
		if !consumed {
			return MediaType{}, Parameters{}, ErrInvalidMediaType
		}

		acceptableMediaType.Parameters = make(Parameters)
		weight := 1000 // 1.000

		// media type parameters
		for len(s) > 0 && s[0] == ';' {
			s = s[1:] // skip the semicolon

			var key, value string
			key, value, s, consumed = consumeParameter(s)
			if !consumed {
				return MediaType{}, Parameters{}, ErrInvalidParameter
			}

			if key == "q" {
				weight, consumed = getWeight(value)
				if !consumed {
					return MediaType{}, Parameters{}, ErrInvalidWeight
				}
				break // "q" parameter separates media type parameters from Accept extension parameters
			}

			acceptableMediaType.Parameters[key] = value
		}

		extensionParameters := make(Parameters)
		for len(s) > 0 && s[0] == ';' {
			s = s[1:] // skip the semicolon

			key, value, remaining, consumed := consumeParameter(s)
			if !consumed {
				return MediaType{}, Parameters{}, ErrInvalidParameter
			}

			s = remaining

			extensionParameters[key] = value
		}

		for i := 0; i < len(availableMediaTypes); i++ {
			if compareMediaTypes(acceptableMediaType, availableMediaTypes[i]) &&
				getPrecedence(acceptableMediaType, weights[i].mediaType) {
				weights[i].mediaType = acceptableMediaType
				weights[i].extensionParameters = extensionParameters
				weights[i].weight = weight
				weights[i].order = mediaTypeCount
			}
		}

		s = skipWhiteSpaces(s)
	}

	// there must not be anything left after parsing the header
	if len(s) > 0 {
		return MediaType{}, Parameters{}, ErrInvalidMediaRange
	}

	resultIndex := -1
	for i := 0; i < len(availableMediaTypes); i++ {
		if resultIndex != -1 {
			if weights[i].weight > weights[resultIndex].weight ||
				(weights[i].weight == weights[resultIndex].weight && weights[i].order < weights[resultIndex].order) {
				resultIndex = i
			}
		} else if weights[i].weight > 0 {
			resultIndex = i
		}
	}

	if resultIndex == -1 {
		return MediaType{}, Parameters{}, ErrNoAcceptableTypeFound
	}

	return availableMediaTypes[resultIndex], weights[resultIndex].extensionParameters, nil
}
