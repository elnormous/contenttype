package contenttype

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func TestGetMediaType(t *testing.T) {
	tables := []struct {
		header     string
		result     string
		parameters map[string]string
		err        error
	}{
		{"", "", map[string]string{}, nil},
		{"application/json", "application/json", map[string]string{}, nil},
		{"*/*", "", nil, InvalidContentTypeError},
		{"Application/JSON", "application/json", map[string]string{}, nil},
		{"Application", "", nil, InvalidContentTypeError},
		{"Application/JSON/test", "", nil, InvalidContentTypeError},
		{" application/json ", "application/json", map[string]string{}, nil},
		// {"Application/XML;charset=utf-8", "application/xml", map[string]string{
		// 	"charset": "utf-8",
		// }, nil},
	}

	for _, table := range tables {
		request, err := http.NewRequest(http.MethodGet, "http://test.test", nil)
		if err != nil {
			log.Fatal(err)
		}

		if len(table.header) > 0 {
			request.Header.Set("Content-Type", table.header)
		}
		result, parameters, err := GetMediaType(request)
		if table.err != nil {
			if err == nil {
				t.Errorf("Expected an error")
			} else if table.err != err {
				t.Errorf("Unexpected error \"%s\", expected \"%s\"", err.Error(), table.err.Error())
			}
		} else if table.err == nil && err != nil {
			t.Errorf("Got an unexpected error \"%s\"", err.Error())
		} else if result != table.result {
			t.Errorf("Invalid content type, got %s, exptected %s", result, table.result)
		} else if !reflect.DeepEqual(parameters, table.parameters) {

			t.Errorf("Wrong parameters, got %v, expected %v", fmt.Sprint(parameters), fmt.Sprint(table.parameters))
		}
	}
}
