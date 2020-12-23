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
	}{
		{"", "", map[string]string{}},
		{"application/json", "application/json", map[string]string{}},
		{"*/*", "*/*", map[string]string{}},
		{"Application/JSON", "application/json", map[string]string{}},
		{" application/json ", "application/json", map[string]string{}},
		{"Application/XML;charset=utf-8", "application/xml", map[string]string{"charset": "utf-8"}},
		{"application/xml;foo=bar ", "application/xml", map[string]string{"foo": "bar"}},
		{"application/xml ; foo=bar ", "application/xml", map[string]string{"foo": "bar"}},
		{"application/xml;foo=\"bar\" ", "application/xml", map[string]string{"foo": "bar"}},
		{"application/xml;foo=\"\" ", "application/xml", map[string]string{"foo": ""}},
		{"application/xml;foo=\"\\\"b\" ", "application/xml", map[string]string{"foo": "\"b"}},
		{"a/b+c;a=b;c=d", "a/b+c", map[string]string{"a": "b", "c": "d"}},
		{"a/b;A=B", "a/b", map[string]string{"a": "b"}},
	}

	for _, table := range tables {
		request, requestError := http.NewRequest(http.MethodGet, "http://test.test", nil)
		if requestError != nil {
			log.Fatal(requestError)
		}

		if len(table.header) > 0 {
			request.Header.Set("Content-Type", table.header)
		}

		result, parameters, mediaTypeError := GetMediaType(request)
		if mediaTypeError != nil {
			t.Errorf("Unexpected error for %s", table.header)
		} else if result != table.result {
			t.Errorf("Invalid content type, got %s, exptected %s", result, table.result)
		} else if !reflect.DeepEqual(parameters, table.parameters) {
			t.Errorf("Wrong parameters, got %v, expected %v", fmt.Sprint(parameters), fmt.Sprint(table.parameters))
		}
	}
}

func TestGetMediaTypeErrors(t *testing.T) {
	tables := []struct {
		header string
		err    error
	}{
		{"Application", InvalidContentTypeError},
		{"/Application", InvalidContentTypeError},
		{"Application/", InvalidContentSubtypeError},
		{"Application/JSON/test", ExpectedParameterError},
		{"application/xml;=bar ", InvalidParameterError},
		{"application/xml; =bar ", InvalidParameterError},
		{"application/xml;foo= ", InvalidParameterError},
	}

	for _, table := range tables {
		request, requestError := http.NewRequest(http.MethodGet, "http://test.test", nil)
		if requestError != nil {
			log.Fatal(requestError)
		}

		if len(table.header) > 0 {
			request.Header.Set("Content-Type", table.header)
		}

		_, _, mediaTypeError := GetMediaType(request)
		if mediaTypeError == nil {
			t.Errorf("Expected an error for %s", table.header)
		} else if table.err != mediaTypeError {
			t.Errorf("Unexpected error \"%s\", expected \"%s\"", mediaTypeError.Error(), table.err.Error())
		}
	}
}

func TestGetAcceptedMediaType(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, "http://test.test", nil)
	if err != nil {
		log.Fatal(err)
	}
	GetAcceptedMediaType(request, []string{"application/json"})
}
