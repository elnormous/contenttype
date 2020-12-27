package contenttype

import (
	"log"
	"net/http"
	"reflect"
	"testing"
)

func TestNewMediaType(t *testing.T) {
	testCases := []struct {
		value  string
		result MediaType
	}{
		{"", MediaType{}},
		{"application/json", MediaType{"application", "json", Parameters{}}},
		{"a/b;c=d", MediaType{"a", "b", Parameters{"c": "d"}}},
		{"/b", MediaType{}},
		{"a/", MediaType{}},
		{"a/b;c", MediaType{}},
	}

	for _, testCase := range testCases {
		result := NewMediaType(testCase.value)

		if result.Type != testCase.result.Type || result.Subtype != testCase.result.Subtype {
			t.Errorf("Invalid content type, got %s/%s, exptected %s/%s for %s", result.Type, result.Subtype, testCase.result.Type, testCase.result.Subtype, testCase.value)
		} else if !reflect.DeepEqual(result.Parameters, testCase.result.Parameters) {
			t.Errorf("Wrong parameters, got %v, expected %v for %s", result.Parameters, testCase.result.Parameters, testCase.value)
		}
	}
}

func TestGetMediaType(t *testing.T) {
	testCases := []struct {
		header string
		result MediaType
	}{
		{"", MediaType{}},
		{"application/json", MediaType{"application", "json", Parameters{}}},
		{"*/*", MediaType{"*", "*", Parameters{}}},
		{"Application/JSON", MediaType{"application", "json", Parameters{}}},
		{" application/json ", MediaType{"application", "json", Parameters{}}},
		{"Application/XML;charset=utf-8", MediaType{"application", "xml", Parameters{"charset": "utf-8"}}},
		{"application/xml;foo=bar ", MediaType{"application", "xml", Parameters{"foo": "bar"}}},
		{"application/xml ; foo=bar ", MediaType{"application", "xml", Parameters{"foo": "bar"}}},
		{"application/xml;foo=\"bar\" ", MediaType{"application", "xml", Parameters{"foo": "bar"}}},
		{"application/xml;foo=\"\" ", MediaType{"application", "xml", Parameters{"foo": ""}}},
		{"application/xml;foo=\"\\\"b\" ", MediaType{"application", "xml", Parameters{"foo": "\"b"}}},
		{"application/xml;foo=\"\\\"B\" ", MediaType{"application", "xml", Parameters{"foo": "\"b"}}},
		{"a/b+c;a=b;c=d", MediaType{"a", "b+c", Parameters{"a": "b", "c": "d"}}},
		{"a/b;A=B", MediaType{"a", "b", Parameters{"a": "b"}}},
	}

	for _, testCase := range testCases {
		request, requestError := http.NewRequest(http.MethodGet, "http://test.test", nil)
		if requestError != nil {
			log.Fatal(requestError)
		}

		if len(testCase.header) > 0 {
			request.Header.Set("Content-Type", testCase.header)
		}

		result, mediaTypeError := GetMediaType(request)
		if mediaTypeError != nil {
			t.Errorf("Unexpected error \"%s\" for %s", mediaTypeError.Error(), testCase.header)
		} else if result.Type != testCase.result.Type || result.Subtype != testCase.result.Subtype {
			t.Errorf("Invalid content type, got %s/%s, exptected %s/%s for %s", result.Type, result.Subtype, testCase.result.Type, testCase.result.Subtype, testCase.header)
		} else if !reflect.DeepEqual(result.Parameters, testCase.result.Parameters) {
			t.Errorf("Wrong parameters, got %v, expected %v for %s", result.Parameters, testCase.result.Parameters, testCase.header)
		}
	}
}

func TestGetMediaTypeErrors(t *testing.T) {
	testCases := []struct {
		header string
		err    error
	}{
		{"Application", InvalidMediaTypeError},
		{"/Application", InvalidMediaTypeError},
		{"Application/", InvalidMediaTypeError},
		{"a/b\x19", InvalidMediaTypeError},
		{"Application/JSON/test", InvalidMediaTypeError},
		{"application/xml;=bar ", InvalidParameterError},
		{"application/xml; =bar ", InvalidParameterError},
		{"application/xml;foo= ", InvalidParameterError},
		{"a/b;c=\x19", InvalidParameterError},
		{"a/b;c=\"\x19\"", InvalidParameterError},
		{"a/b;c=\"\\\x19\"", InvalidParameterError},
		{"a/b;c", InvalidParameterError},
		{"a/b e", InvalidMediaTypeError},
	}

	for _, testCase := range testCases {
		request, requestError := http.NewRequest(http.MethodGet, "http://test.test", nil)
		if requestError != nil {
			log.Fatal(requestError)
		}

		if len(testCase.header) > 0 {
			request.Header.Set("Content-Type", testCase.header)
		}

		_, mediaTypeError := GetMediaType(request)
		if mediaTypeError == nil {
			t.Errorf("Expected an error for %s", testCase.header)
		} else if testCase.err != mediaTypeError {
			t.Errorf("Unexpected error \"%s\", expected \"%s\" for %s", mediaTypeError.Error(), testCase.err.Error(), testCase.header)
		}
	}
}

func TestGetAcceptableMediaType(t *testing.T) {
	testCases := []struct {
		header              string
		availableMediaTypes []MediaType
		result              MediaType
	}{
		{"", []MediaType{{"application", "json", Parameters{}}}, MediaType{"application", "json", Parameters{}}},
		{"application/json", []MediaType{{"application", "json", Parameters{}}}, MediaType{"application", "json", Parameters{}}},
		{"Application/Json", []MediaType{{"application", "json", Parameters{}}}, MediaType{"application", "json", Parameters{}}},
		{"text/plain,application/xml", []MediaType{{"text", "plain", Parameters{}}}, MediaType{"text", "plain", Parameters{}}},
		{"text/plain,application/xml", []MediaType{{"application", "xml", Parameters{}}}, MediaType{"application", "xml", Parameters{}}},
		{"text/plain;q=1.0", []MediaType{{"text", "plain", Parameters{}}}, MediaType{"text", "plain", Parameters{}}},
		{"*/*", []MediaType{{"application", "json", Parameters{}}}, MediaType{"application", "json", Parameters{}}},
		{"application/*", []MediaType{{"application", "json", Parameters{}}}, MediaType{"application", "json", Parameters{}}},
		{"a/b;q=1.", []MediaType{{"a", "b", Parameters{}}}, MediaType{"a", "b", Parameters{}}},
		{"a/b;q=0.1,c/d;q=0.2", []MediaType{
			{"a", "b", Parameters{}},
			{"c", "d", Parameters{}},
		}, MediaType{"c", "d", Parameters{}}},
		{"a/b;q=0.2,c/d;q=0.2", []MediaType{
			{"a", "b", Parameters{}},
			{"c", "d", Parameters{}},
		}, MediaType{"a", "b", Parameters{}}},
		{"a/*;q=0.2,a/c", []MediaType{
			{"a", "b", Parameters{}},
			{"a", "c", Parameters{}},
		}, MediaType{"a", "c", Parameters{}}},
		{"a/b,a/a", []MediaType{
			{"a", "a", Parameters{}},
			{"a", "b", Parameters{}},
		}, MediaType{"a", "b", Parameters{}}},
		{"a/*", []MediaType{
			{"a", "a", Parameters{}},
			{"a", "b", Parameters{}},
		}, MediaType{"a", "a", Parameters{}}},
		{"a/a;q=0.2,a/*", []MediaType{
			{"a", "a", Parameters{}},
			{"a", "b", Parameters{}},
		}, MediaType{"a", "b", Parameters{}}},
		{"a/a;q=0.2,a/a;c=d", []MediaType{
			{"a", "a", Parameters{}},
			{"a", "a", Parameters{"c": "d"}},
		}, MediaType{"a", "a", Parameters{"c": "d"}}},
	}

	for _, testCase := range testCases {
		request, requestError := http.NewRequest(http.MethodGet, "http://test.test", nil)
		if requestError != nil {
			log.Fatal(requestError)
		}

		if len(testCase.header) > 0 {
			request.Header.Set("Accept", testCase.header)
		}

		result, _, mediaTypeError := GetAcceptableMediaType(request, testCase.availableMediaTypes)

		if mediaTypeError != nil {
			t.Errorf("Unexpected error \"%s\" for %s", mediaTypeError.Error(), testCase.header)
		} else if result.Type != testCase.result.Type || result.Subtype != testCase.result.Subtype {
			t.Errorf("Invalid content type, got %s/%s, exptected %s/%s for %s", result.Type, result.Subtype, testCase.result.Type, testCase.result.Subtype, testCase.header)
		} else if !reflect.DeepEqual(result.Parameters, testCase.result.Parameters) {
			t.Errorf("Wrong parameters, got %v, expected %v for %s", result.Parameters, testCase.result.Parameters, testCase.header)
		}
	}
}

func TestGetAcceptableMediaTypeErrors(t *testing.T) {
	testCases := []struct {
		header              string
		availableMediaTypes []MediaType
		err                 error
	}{
		{"", []MediaType{}, NoAvailableTypeGivenError},
		{"application/xml", []MediaType{{"application", "json", Parameters{}}}, NoAcceptableTypeFoundError},
		{"application/xml/", []MediaType{{"application", "json", Parameters{}}}, InvalidMediaRangeError},
		{"application/xml,", []MediaType{{"application", "json", Parameters{}}}, InvalidMediaTypeError},
		{"/xml", []MediaType{{"application", "json", Parameters{}}}, InvalidMediaTypeError},
		{"application/,", []MediaType{{"application", "json", Parameters{}}}, InvalidMediaTypeError},
		{"a/b c", []MediaType{{"a", "b", Parameters{}}}, InvalidMediaRangeError},
		{"a/b;c", []MediaType{{"a", "b", Parameters{}}}, InvalidParameterError},
		{"*/b", []MediaType{{"a", "b", Parameters{}}}, InvalidMediaTypeError},
		{"a/b;q=a", []MediaType{{"a", "b", Parameters{}}}, InvalidWeightError},
		{"a/b;q=11", []MediaType{{"a", "b", Parameters{}}}, InvalidWeightError},
		{"a/b;q=1.0000", []MediaType{{"a", "b", Parameters{}}}, InvalidWeightError},
		{"a/b;q=1.a", []MediaType{{"a", "b", Parameters{}}}, InvalidWeightError},
		{"a/b;q=1.100", []MediaType{{"a", "b", Parameters{}}}, InvalidWeightError},
		{"a/b;q=0", []MediaType{{"a", "b", Parameters{}}}, NoAcceptableTypeFoundError},
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
			t.Errorf("Unexpected error \"%s\", expected \"%s\" for %s", mediaTypeError.Error(), testCase.err.Error(), testCase.header)
		}
	}
}
