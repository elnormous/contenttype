package contenttype

import (
	"errors"
	"net/http"
	"strings"
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

func GetMediaType(request *http.Request) (string, error) {
	values := request.Header.Values("Content-Type")

	if len(values) == 0 {
		return "", nil
	}

	value := values[0]

	if value == "*/*" {
		return "", InvalidContentTypeError
	}

	slashIndex := strings.Index(value, "/")
	if slashIndex == -1 {
		return "", InvalidContentTypeError
	}

	var stringBuilder strings.Builder
	supertype, subtype := value[:slashIndex], value[slashIndex+1:]
	if !isToken(supertype) || !isToken(subtype) {
		return "", InvalidContentTypeError
	}

	stringBuilder.WriteString(strings.ToLower(supertype))
	stringBuilder.WriteByte('/')
	stringBuilder.WriteString(strings.ToLower(subtype))

	return stringBuilder.String(), nil
}
