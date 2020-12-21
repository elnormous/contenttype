package mediatype

import (
	"log"
	"net/http"
	"testing"
)

func TestGetMediaType(t *testing.T) {
	tables := []struct {
		header string
		result string
		err    error
	}{
		{"", "", nil},
		{"application/json", "application/json", nil},
		{"*/*", "", InvalidContentTypeError},
		{"Application/JSON", "application/json", nil},
	}

	for _, table := range tables {
		request, err := http.NewRequest(http.MethodGet, "http://test.test", nil)
		if err != nil {
			log.Fatal(err)
		}

		if len(table.header) > 0 {
			request.Header.Set("Content-Type", table.header)
		}
		result, err := GetMediaType(request)
		if table.err != nil {
			if err == nil {
				t.Errorf("Expected an error")
			} else if table.err != err {
				t.Errorf("Unexpected error \"" + err.Error() + "\", expected \"" + table.err.Error() + "\"")
			}
		} else if table.err == nil && err != nil {
			t.Errorf("Got an unexpected error " + err.Error())
		} else if result != table.result {
			t.Errorf("Invalid content type, go " + result + ", expected " + table.result)
		}
	}
}
