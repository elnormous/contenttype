// Package contenttype implements HTTP Content-Type, Accept and Accept-Language header parsers.
package contenttype

import (
	"errors"
	"net/http"
	"reflect"
	"strings"
)

var (
	// ErrInvalidMediaType is returned when the media type in the Content-Type or Accept header is syntactically invalid.
	ErrInvalidMediaType = errors.New("invalid media type")
	// ErrInvalidMediaRange is returned when the range of media types in the Content-Type or Accept header is syntactically invalid.
	ErrInvalidMediaRange = errors.New("invalid media range")
	// ErrInvalidParameter is returned when the media type parameter in the Content-Type or Accept header is syntactically invalid.
	ErrInvalidParameter = errors.New("invalid parameter")
	// ErrInvalidExtensionParameter is returned when the media type extension parameter in the Content-Type or Accept header is syntactically invalid.
	ErrInvalidExtensionParameter = errors.New("invalid extension parameter")
	// ErrNoAcceptableTypeFound is returned when Accept header contains only media types that are not in the acceptable media type list.
	ErrNoAcceptableTypeFound = errors.New("no acceptable type found")
	// ErrNoAvailableTypeGiven is returned when the acceptable media type list is empty.
	ErrNoAvailableTypeGiven = errors.New("no available type given")
	// ErrInvalidWeight is returned when the media type weight in Accept header is syntactically invalid.
	ErrInvalidWeight = errors.New("invalid weight")
	// ErrInvalidLanguage is returned when the language is syntactically invalid.
	ErrInvalidLanguage = errors.New("invalid language")
)

// Parameters represents media type parameters as a key-value map.
type Parameters = map[string]string

// MediaType holds the type, subtype and parameters of a media type.
type MediaType struct {
	Type       string
	Subtype    string
	Parameters Parameters
}

// NewMediaType parses the string and returns an instance of MediaType struct.
func NewMediaType(s string) MediaType {
	mediaType, err := ParseMediaType(s)
	if err != nil {
		return MediaType{}
	}

	return mediaType
}

// Converts the MediaType to string.
func (mediaType MediaType) String() string {
	var stringBuilder strings.Builder

	if len(mediaType.Type) > 0 || len(mediaType.Subtype) > 0 {
		stringBuilder.WriteString(mediaType.Type)
		stringBuilder.WriteByte('/')
		stringBuilder.WriteString(mediaType.Subtype)

		for key, value := range mediaType.Parameters {
			stringBuilder.WriteByte(';')
			stringBuilder.WriteString(key)
			stringBuilder.WriteByte('=')
			stringBuilder.WriteString(value)
		}
	}

	return stringBuilder.String()
}

// MIME returns the MIME type without any of the parameters
func (mediaType MediaType) MIME() string {
	var stringBuilder strings.Builder

	if len(mediaType.Type) > 0 || len(mediaType.Subtype) > 0 {
		stringBuilder.WriteString(mediaType.Type)
		stringBuilder.WriteByte('/')
		stringBuilder.WriteString(mediaType.Subtype)
	}

	return stringBuilder.String()
}

// Equal checks whether the provided MIME media type matches this one
// including all parameters
func (mediaType MediaType) Equal(mt MediaType) bool {
	return reflect.DeepEqual(mediaType, mt)
}

// EqualsMIME checks whether the base MIME types match
func (mediaType MediaType) EqualsMIME(mt MediaType) bool {
	return (mediaType.Type == mt.Type) && (mediaType.Subtype == mt.Subtype)
}

// Matches checks whether the MIME media types match handling wildcards in either
func (mediaType MediaType) Matches(mt MediaType) bool {
	t := mediaType.Type == mt.Type || (mediaType.Type == "*") || (mt.Type == "*")
	st := mediaType.Subtype == mt.Subtype || mediaType.Subtype == "*" || mt.Subtype == "*"
	return t && st
}

// MatchesAny checks whether the MIME media types matches any of the specified
// list of mediatype handling wildcards in any of them
func (mediaType MediaType) MatchesAny(mts ...MediaType) bool {
	for _, mt := range mts {
		if mediaType.Matches(mt) {
			return true
		}
	}
	return false
}

// IsWildcard returns true if either the Type or Subtype are the wildcard character '*'
func (mediaType MediaType) IsWildcard() bool {
	return mediaType.Type == `*` || mediaType.Subtype == `*`
}

// GetMediaType gets the content of Content-Type header, parses it, and returns the parsed MediaType.
// If the request does not contain the Content-Type header, an empty MediaType is returned.
func GetMediaType(request *http.Request) (MediaType, error) {
	// RFC 7231, 3.1.1.5. Content-Type
	contentTypeHeaders := request.Header.Values("Content-Type")
	if len(contentTypeHeaders) == 0 {
		return MediaType{}, nil
	}

	return ParseMediaType(contentTypeHeaders[0])
}

// ParseMediaType parses the given string as a MIME media type (with optional parameters) and returns it as a MediaType.
// If the string cannot be parsed an appropriate error is returned.
func ParseMediaType(s string) (MediaType, error) {
	// RFC 7231, 3.1.1.1. Media Type
	mediaType := MediaType{
		Parameters: Parameters{},
	}
	var consumed bool
	if mediaType.Type, mediaType.Subtype, s, consumed = consumeType(skipWhitespaces(s)); !consumed {
		return MediaType{}, ErrInvalidMediaType
	}

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

// GetAcceptableMediaType chooses a media type from available media types according to the Accept header.
// Returns the most suitable media type or an error if no type can be selected.
func GetAcceptableMediaType(request *http.Request, availableMediaTypes []MediaType) (MediaType, Parameters, error) {
	// RFC 7231, 5.3.2. Accept
	if len(availableMediaTypes) == 0 {
		return MediaType{}, Parameters{}, ErrNoAvailableTypeGiven
	}

	acceptHeaders := request.Header.Values("Accept")
	if len(acceptHeaders) == 0 {
		return availableMediaTypes[0], Parameters{}, nil
	}

	return GetAcceptableMediaTypeFromHeader(acceptHeaders[0], availableMediaTypes)
}

// GetAcceptableMediaTypeFromHeader chooses a media type from available media types according to the specified Accept header value.
// Returns the most suitable media type or an error if no type can be selected.
func GetAcceptableMediaTypeFromHeader(headerValue string, availableMediaTypes []MediaType) (MediaType, Parameters, error) {
	s := headerValue

	weights := make([]struct {
		mediaType           MediaType
		extensionParameters Parameters
		weight              uint
		order               uint
	}, len(availableMediaTypes))

	for mediaTypeCount := uint(0); len(s) > 0; mediaTypeCount++ {
		if mediaTypeCount > 0 {
			// every media type after the first one must start with a comma
			var skipped bool
			s, skipped = skipCharacter(s, ',')
			if !skipped {
				break
			}
		}

		acceptableMediaType := MediaType{
			Parameters: Parameters{},
		}
		var consumed bool
		if acceptableMediaType.Type, acceptableMediaType.Subtype, s, consumed = consumeType(skipWhitespaces(s)); !consumed {
			return MediaType{}, Parameters{}, ErrInvalidMediaType
		}

		weight := uint(1000) // 1.000

		// media type parameters
		for len(s) > 0 && s[0] == ';' {
			s = s[1:] // skip the semicolon

			var key, value string
			if key, value, s, consumed = consumeParameter(s); !consumed {
				return MediaType{}, Parameters{}, ErrInvalidParameter
			}

			if key == "q" {
				if weight, consumed = getWeight(value); !consumed {
					return MediaType{}, Parameters{}, ErrInvalidWeight
				}
				break // "q" parameter separates media type parameters from Accept extension parameters
			}

			acceptableMediaType.Parameters[key] = value
		}

		extensionParameters := Parameters{}
		for len(s) > 0 && s[0] == ';' {
			s = s[1:] // skip the semicolon

			var key, value, remaining string
			if key, value, remaining, consumed = consumeParameter(s); !consumed {
				return MediaType{}, Parameters{}, ErrInvalidParameter
			}

			s = remaining

			extensionParameters[key] = value
		}

		for i, availableMediaType := range availableMediaTypes {
			if compareMediaTypes(acceptableMediaType, availableMediaType) &&
				getPrecedence(acceptableMediaType, weights[i].mediaType) {
				weights[i].mediaType = acceptableMediaType
				weights[i].extensionParameters = extensionParameters
				weights[i].weight = weight
				weights[i].order = mediaTypeCount
			}
		}

		s = skipWhitespaces(s)
	}

	// there must not be anything left after parsing the header
	if len(s) > 0 {
		return MediaType{}, Parameters{}, ErrInvalidMediaRange
	}

	resultIndex := -1
	for i, weight := range weights {
		if resultIndex != -1 {
			if weight.weight > weights[resultIndex].weight ||
				(weight.weight == weights[resultIndex].weight && weight.order < weights[resultIndex].order) {
				resultIndex = i
			}
		} else if weight.weight > 0 {
			resultIndex = i
		}
	}

	if resultIndex == -1 {
		return MediaType{}, Parameters{}, ErrNoAcceptableTypeFound
	}

	return availableMediaTypes[resultIndex], weights[resultIndex].extensionParameters, nil
}

func isWhitespaceChar(c byte) bool {
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
	return c >= 0x80 // c is always less than or equal to 0xFF
}

func isQuotedTextChar(c byte) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return isWhitespaceChar(c) ||
		c == 0x21 ||
		(c >= 0x23 && c <= 0x5B) ||
		(c >= 0x5D && c <= 0x7E) ||
		isObsoleteTextChar(c)
}

func isQuotedPairChar(c byte) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return isWhitespaceChar(c) ||
		isVisibleChar(c) ||
		isObsoleteTextChar(c)
}

func skipWhitespaces(s string) string {
	// RFC 7230, 3.2.3. Whitespace
	for i := 0; i < len(s); i++ {
		if !isWhitespaceChar(s[i]) {
			return s[i:]
		}
	}

	return ""
}

func skipCharacter(s string, c byte) (remaining string, consumed bool) {
	if len(s) == 0 || s[0] != c {
		return s, false
	}

	return s[1:], true
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
	// RFC 7230, 3.2.6. Field Value Components
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
	t, s, consumed := consumeToken(s)
	if !consumed {
		return "", "", s, false
	}

	s, skipped := skipCharacter(s, '/')
	if !skipped {
		return "", "", s, false
	}

	st, s, consumed := consumeToken(s)
	if !consumed {
		return "", "", s, false
	}

	if t == "*" && st != "*" {
		return "", "", s, false
	}

	return t, st, skipWhitespaces(s), true
}

func consumeParameter(s string) (string, string, string, bool) {
	// RFC 7231, 3.1.1.1. Media Type
	var consumed bool
	var key string
	if key, s, consumed = consumeToken(skipWhitespaces(s)); !consumed {
		return "", "", s, false
	}

	var skipped bool
	s, skipped = skipCharacter(s, '=')
	if !skipped {
		return "", "", s, false
	}

	var value string
	if s, skipped = skipCharacter(s, '"'); skipped {
		if value, s, consumed = consumeQuotedString(s); !consumed {
			return "", "", s, false
		}

		if s, skipped = skipCharacter(s, '"'); !skipped { // skip the closing quote
			return "", "", s, false
		}
	} else {
		if value, s, consumed = consumeToken(s); !consumed {
			return "", "", s, false
		}
	}

	return key, value, skipWhitespaces(s), true
}

func getWeight(s string) (uint, bool) {
	// RFC 7231, 5.3.1. Quality Values
	result := uint(0)
	multiplier := uint(1000)

	// the string must not have more than three digits after the decimal point
	if len(s) > 5 {
		return 0, false
	}

	for i := 0; i < len(s); i++ {
		if i == 0 {
			// the first character must be 0 or 1
			if s[i] != '0' && s[i] != '1' {
				return 0, false
			}

			result = uint(s[i]-'0') * multiplier
			multiplier /= 10
		} else if i == 1 {
			// the second character must be a dot
			if s[i] != '.' {
				return 0, false
			}
		} else {
			// the remaining characters must be digits and the value can not be greater than 1.000
			if (s[0] == '1' && s[i] != '0') ||
				!isDigitChar(s[i]) {
				return 0, false
			}

			result += uint(s[i]-'0') * multiplier
			multiplier /= 10
		}
	}

	return result, true
}

func compareMediaTypes(checkMediaType, mediaType MediaType) bool {
	// RFC 7231, 5.3.2. Accept
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
	// RFC 7231, 5.3.2. Accept
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

func consumeLanguage(s string) (string, string, bool) {
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

func consumeScript(s string) (string, string, bool) {
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

func consumeRegion(s string) (string, string, bool) {
	// RFC 5646, 2.1. Syntax
	for i := 0; i < len(s) && i < 3; i++ {
		if !isAlphaChar(s[i]) {
			return strings.ToLower(s[:i]), s[i:], len(s) > 0
		}
	}

	return strings.ToLower(s), "", len(s) > 0
}

func consumeLanguageTags(s string) (string, string, string, string, bool) {
	language, s, consumed := consumeLanguage(s)

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
