package mediatype

import (
	"errors"
	"net/http"
	"strings"
)

var InvalidContentTypeError = errors.New("Invalid content type")

func GetContentType(request *http.Request) (string, error) {
	values := request.Header.Values("Content-Type")

	if len(values) == 0 {
		return "", nil
	}

	value := values[0]

	if value == "*/*" {
		return "", InvalidContentTypeError
	}

	return strings.ToLower(value), nil
}
