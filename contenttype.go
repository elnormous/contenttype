package contenttype

import (
	"errors"
	"log"
	"net/http"
	"strings"
)

var InvalidContentTypeError = errors.New("Invalid content type")

func isTokenChar(r rune) bool {
	// RFC 7230, 3.2.6. Field Value Components
	return strings.ContainsRune("!#$%&'*+-.^_`|~", r) ||
		(r >= 0x30 && r <= 0x39) || // DIGIT
		(r >= 0x41 && r <= 0x5A) || (r >= 0x61 && r <= 0x7A) // ALPHA
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

func isWhiteSpaceChar(r rune) bool {
	// RFC 7230, 3.2.3. Whitespace
	return r == 0x09 || r == 0x20 // HTAB or SP
}

func GetMediaType(request *http.Request) (string, map[string]string, error) {
	contentTypes := request.Header.Values("Content-Type")

	if len(contentTypes) == 0 {
		return "", map[string]string{}, nil
	}

	contentType := strings.TrimFunc(contentTypes[0], isWhiteSpaceChar)

	slashIndex := strings.Index(contentType, "/")
	if slashIndex == -1 {
		return "", nil, InvalidContentTypeError
	}

	parameters := make(map[string]string)

	endIndex := strings.Index(contentType, ";")
	if endIndex == -1 {
		endIndex = len(contentType)
	} else {
		parameterIndex := endIndex
		parameterString := contentType

		for parameterIndex != -1 {
			parameterString := parameterString[parameterIndex+1:]

			equalIndex := strings.Index(contentType, "=")
			key := contentType[:equalIndex]
			log.Println(key)

			parameterIndex = strings.Index(parameterString, "/")
		}
	}

	var stringBuilder strings.Builder
	supertype := contentType[:slashIndex]
	subtype := contentType[slashIndex+1 : endIndex]

	if !isToken(supertype) || !isToken(subtype) {
		return "", nil, InvalidContentTypeError
	}

	stringBuilder.WriteString(strings.ToLower(supertype))
	stringBuilder.WriteByte('/')
	stringBuilder.WriteString(strings.ToLower(subtype))

	return stringBuilder.String(), parameters, nil
}
