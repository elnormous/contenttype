package contenttype

import (
	"errors"
	"net/http"
	"strings"
	"unicode"
)

var InvalidContentTypeError = errors.New("Invalid content type")

func isTSpecial(r rune) bool {
	return strings.ContainsRune(`()<>@,;:\"/[]?=`, r)
}

func isTokenChar(r rune) bool {
	return r > 0x20 && r < 0x7f && !isTSpecial(r)
}

func isNotTokenChar(r rune) bool {
	return !isTokenChar(r)
}

func isToken(s string) bool {
	if len(s) == 0 {
		return false
	}
	return strings.IndexFunc(s, isNotTokenChar) == -1
}

func GetMediaType(request *http.Request) (string, map[string]string, error) {
	values := request.Header.Values("Content-Type")

	if len(values) == 0 {
		return "", map[string]string{}, nil
	}

	value := values[0]

	if value == "*/*" {
		return "", nil, InvalidContentTypeError
	}

	slashIndex := strings.IndexByte(value, '/')
	if slashIndex == -1 {
		return "", nil, InvalidContentTypeError
	}

	parameters := make(map[string]string)

	endIndex := strings.Index(value, ";")
	if endIndex == -1 {
		endIndex = len(value)
	} else {
		// TODO: parse parameters in while loop
	}

	var stringBuilder strings.Builder
	supertype := strings.TrimLeftFunc(value[:slashIndex], unicode.IsSpace)
	subtype := strings.TrimRightFunc(value[slashIndex+1:endIndex], unicode.IsSpace)

	if !isToken(supertype) || !isToken(subtype) {
		return "", nil, InvalidContentTypeError
	}

	stringBuilder.WriteString(strings.ToLower(supertype))
	stringBuilder.WriteByte('/')
	stringBuilder.WriteString(strings.ToLower(subtype))

	return stringBuilder.String(), parameters, nil
}
