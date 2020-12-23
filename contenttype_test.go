package contenttype

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
)

func TestGetMediaType(t *testing.T) {
	testCases := []struct {
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

	for _, testCase := range testCases {
		request, requestError := http.NewRequest(http.MethodGet, "http://test.test", nil)
		if requestError != nil {
			log.Fatal(requestError)
		}

		if len(testCase.header) > 0 {
			request.Header.Set("Content-Type", testCase.header)
		}

		result, parameters, mediaTypeError := GetMediaType(request)
		if mediaTypeError != nil {
			t.Errorf("Unexpected error for %s: %s", testCase.header, mediaTypeError.Error())
		} else if result != testCase.result {
			t.Errorf("Invalid content type, got %s, exptected %s", result, testCase.result)
		} else if !reflect.DeepEqual(parameters, testCase.parameters) {
			t.Errorf("Wrong parameters, got %v, expected %v", fmt.Sprint(parameters), fmt.Sprint(testCase.parameters))
		}
	}
}

func TestGetMediaTypeErrors(t *testing.T) {
	testCases := []struct {
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

	for _, testCase := range testCases {
		request, requestError := http.NewRequest(http.MethodGet, "http://test.test", nil)
		if requestError != nil {
			log.Fatal(requestError)
		}

		if len(testCase.header) > 0 {
			request.Header.Set("Content-Type", testCase.header)
		}

		_, _, mediaTypeError := GetMediaType(request)
		if mediaTypeError == nil {
			t.Errorf("Expected an error for %s", testCase.header)
		} else if testCase.err != mediaTypeError {
			t.Errorf("Unexpected error \"%s\", expected \"%s\"", mediaTypeError.Error(), testCase.err.Error())
		}
	}
}

func TestGetAcceptableMediaType(t *testing.T) {
	testCases := []struct {
		header              string
		availableMediaTypes []string
		result              string
		parameters          map[string]string
	}{
		{"", []string{"application/json"}, "application/json", map[string]string{}},
		{"application/json", []string{"application/json"}, "application/json", map[string]string{}},
		{"Application/Json", []string{"application/json"}, "application/json", map[string]string{}},
		{"application/json,application/xml", []string{"application/json"}, "application/json", map[string]string{}},
	}

	for _, testCase := range testCases {
		request, requestError := http.NewRequest(http.MethodGet, "http://test.test", nil)
		if requestError != nil {
			log.Fatal(requestError)
		}

		if len(testCase.header) > 0 {
			request.Header.Set("Accept", testCase.header)
		}

		result, parameters, mediaTypeError := GetAcceptableMediaType(request, testCase.availableMediaTypes)

		if mediaTypeError != nil {
			t.Errorf("Unexpected error for %s: %s", testCase.header, mediaTypeError.Error())
		} else if result != testCase.result {
			t.Errorf("Invalid content type, got %s, exptected %s", result, testCase.result)
		} else if !reflect.DeepEqual(parameters, testCase.parameters) {
			t.Errorf("Wrong parameters, got %v, expected %v", fmt.Sprint(parameters), fmt.Sprint(testCase.parameters))
		}
	}
}

func TestGetAcceptableMediaTypeErrors(t *testing.T) {
	testCases := []struct {
		header              string
		availableMediaTypes []string
		err                 error
	}{
		{"", []string{}, NoAvailableTypeGivenError},
		{"application/xml", []string{"application/json"}, NoAcceptableTypeFoundError},
	}

	for _, testCase := range testCases {
		request, requestError := http.NewRequest(http.MethodGet, "http://test.test", nil)
		if requestError != nil {
			log.Fatal(requestError)
		}

		if len(testCase.header) > 0 {
			request.Header.Set("Accept", testCase.header)
		}

		_, _, mediaTypeError := GetAcceptableMediaType(request, testCase.availableMediaTypes)
		if mediaTypeError == nil {
			t.Errorf("Expected an error for %s", testCase.header)
		} else if testCase.err != mediaTypeError {
			t.Errorf("Unexpected error \"%s\", expected \"%s\"", mediaTypeError.Error(), testCase.err.Error())
		}
	}
}
