package contenttype

import (
	"log"
	"net/http"
	"reflect"
	"testing"
)

func TestNewMediaType(t *testing.T) {
	testCases := []struct {
		name   string
		value  string
		result MediaType
	}{
		{"Empty string", "", MediaType{}},
		{"Type and subtype", "application/json", MediaType{"application", "json", Parameters{}}},
		{"Type, subtype, parameter", "a/b;c=d", MediaType{"a", "b", Parameters{"c": "d"}}},
		{"Subtype only", "/b", MediaType{}},
		{"Type only", "a/", MediaType{}},
		{"Type, subtype, invalid parameter", "a/b;c", MediaType{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := NewMediaType(testCase.value)

			if result.Type != testCase.result.Type || result.Subtype != testCase.result.Subtype {
				t.Fatalf("Invalid content type, got %s/%s, exptected %s/%s for %s", result.Type, result.Subtype, testCase.result.Type, testCase.result.Subtype, testCase.value)
			} else if !reflect.DeepEqual(result.Parameters, testCase.result.Parameters) {
				t.Fatalf("Wrong parameters, got %v, expected %v for %s", result.Parameters, testCase.result.Parameters, testCase.value)
			}
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		name   string
		value  MediaType
		result string
	}{
		{"Empty media type", MediaType{}, ""},
		{"Type and subtype", MediaType{"application", "json", Parameters{}}, "application/json"},
		{"Type, subtype, parameter", MediaType{"a", "b", Parameters{"c": "d"}}, "a/b;c=d"},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.value.String()

			if result != testCase.result {
				t.Errorf("Invalid result type, got %s, exptected %s", result, testCase.result)
			}
		})
	}
}

func TestGetMediaType(t *testing.T) {
	testCases := []struct {
		name   string
		header string
		result MediaType
	}{
		{"Empty header", "", MediaType{}},
		{"Type and subtype", "application/json", MediaType{"application", "json", Parameters{}}},
		{"Wildcard", "*/*", MediaType{"*", "*", Parameters{}}},
		{"Capital subtype", "Application/JSON", MediaType{"application", "json", Parameters{}}},
		{"Space in front of type", " application/json ", MediaType{"application", "json", Parameters{}}},
		{"Capital and parameter", "Application/XML;charset=utf-8", MediaType{"application", "xml", Parameters{"charset": "utf-8"}}},
		{"White space after parameter", "application/xml;foo=bar ", MediaType{"application", "xml", Parameters{"foo": "bar"}}},
		{"White space after subtype and before parameter", "application/xml ; foo=bar ", MediaType{"application", "xml", Parameters{"foo": "bar"}}},
		{"Quoted parameter", "application/xml;foo=\"bar\" ", MediaType{"application", "xml", Parameters{"foo": "bar"}}},
		{"Quoted empty parameter", "application/xml;foo=\"\" ", MediaType{"application", "xml", Parameters{"foo": ""}}},
		{"Quoted pair", "application/xml;foo=\"\\\"b\" ", MediaType{"application", "xml", Parameters{"foo": "\"b"}}},
		{"Whitespace after quoted parameter", "application/xml;foo=\"\\\"B\" ", MediaType{"application", "xml", Parameters{"foo": "\"b"}}},
		{"Plus in subtype", "a/b+c;a=b;c=d", MediaType{"a", "b+c", Parameters{"a": "b", "c": "d"}}},
		{"Capital parameter", "a/b;A=B", MediaType{"a", "b", Parameters{"a": "b"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
		})
	}
}

func TestGetMediaTypeErrors(t *testing.T) {
	testCases := []struct {
		name   string
		header string
		err    error
	}{
		{"Type only", "Application", ErrInvalidMediaType},
		{"Subtype only", "/Application", ErrInvalidMediaType},
		{"Type with slash", "Application/", ErrInvalidMediaType},
		{"Invalid token character", "a/b\x19", ErrInvalidMediaType},
		{"Invalid character after subtype", "Application/JSON/test", ErrInvalidMediaType},
		{"No parameter name", "application/xml;=bar ", ErrInvalidParameter},
		{"Whitespace and no parameter name", "application/xml; =bar ", ErrInvalidParameter},
		{"No value and whitespace", "application/xml;foo= ", ErrInvalidParameter},
		{"Invalid character in value", "a/b;c=\x19", ErrInvalidParameter},
		{"Invalid character in quoted string", "a/b;c=\"\x19\"", ErrInvalidParameter},
		{"Invalid character in quoted pair", "a/b;c=\"\\\x19\"", ErrInvalidParameter},
		{"No assignment after parameter", "a/b;c", ErrInvalidParameter},
		{"No semicolon before paremeter", "a/b e", ErrInvalidMediaType},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
		})
	}
}

func TestGetAcceptableMediaType(t *testing.T) {
	testCases := []struct {
		name                string
		header              string
		availableMediaTypes []MediaType
		result              MediaType
		extensionParameters Parameters
	}{
		{"Empty header", "", []MediaType{{"application", "json", Parameters{}}}, MediaType{"application", "json", Parameters{}}, Parameters{}},
		{"Type and subtype", "application/json", []MediaType{{"application", "json", Parameters{}}}, MediaType{"application", "json", Parameters{}}, Parameters{}},
		{"Capitalized type and subtype", "Application/Json", []MediaType{{"application", "json", Parameters{}}}, MediaType{"application", "json", Parameters{}}, Parameters{}},
		{"Multiple accept types", "text/plain,application/xml", []MediaType{{"text", "plain", Parameters{}}}, MediaType{"text", "plain", Parameters{}}, Parameters{}},
		{"Multiple accept types, second available", "text/plain,application/xml", []MediaType{{"application", "xml", Parameters{}}}, MediaType{"application", "xml", Parameters{}}, Parameters{}},
		{"Accept weight", "text/plain;q=1.0", []MediaType{{"text", "plain", Parameters{}}}, MediaType{"text", "plain", Parameters{}}, Parameters{}},
		{"Wildcard", "*/*", []MediaType{{"application", "json", Parameters{}}}, MediaType{"application", "json", Parameters{}}, Parameters{}},
		{"Wildcard subtype", "application/*", []MediaType{{"application", "json", Parameters{}}}, MediaType{"application", "json", Parameters{}}, Parameters{}},
		{"Weight with dot", "a/b;q=1.", []MediaType{{"a", "b", Parameters{}}}, MediaType{"a", "b", Parameters{}}, Parameters{}},
		{"Multiple weights", "a/b;q=0.1,c/d;q=0.2", []MediaType{
			{"a", "b", Parameters{}},
			{"c", "d", Parameters{}},
		}, MediaType{"c", "d", Parameters{}}, Parameters{}},
		{"Multiple weights and default weight", "a/b;q=0.2,c/d;q=0.2", []MediaType{
			{"a", "b", Parameters{}},
			{"c", "d", Parameters{}},
		}, MediaType{"a", "b", Parameters{}}, Parameters{}},
		{"Wildcard subtype and weight", "a/*;q=0.2,a/c", []MediaType{
			{"a", "b", Parameters{}},
			{"a", "c", Parameters{}},
		}, MediaType{"a", "c", Parameters{}}, Parameters{}},
		{"Different accept order", "a/b,a/a", []MediaType{
			{"a", "a", Parameters{}},
			{"a", "b", Parameters{}},
		}, MediaType{"a", "b", Parameters{}}, Parameters{}},
		{"Wildcard subtype with multiple available types", "a/*", []MediaType{
			{"a", "a", Parameters{}},
			{"a", "b", Parameters{}},
		}, MediaType{"a", "a", Parameters{}}, Parameters{}},
		{"Wildcard subtype against weighted type", "a/a;q=0.2,a/*", []MediaType{
			{"a", "a", Parameters{}},
			{"a", "b", Parameters{}},
		}, MediaType{"a", "b", Parameters{}}, Parameters{}},
		{"Media type parameter", "a/a;q=0.2,a/a;c=d", []MediaType{
			{"a", "a", Parameters{}},
			{"a", "a", Parameters{"c": "d"}},
		}, MediaType{"a", "a", Parameters{"c": "d"}}, Parameters{}},
		{"Weight and media type parameter", "a/b;q=1;e=e", []MediaType{{"a", "b", Parameters{}}}, MediaType{"a", "b", Parameters{}}, Parameters{"e": "e"}},
		{"", "a/*,a/a;q=0", []MediaType{
			{"a", "a", Parameters{}},
			{"a", "b", Parameters{}},
		}, MediaType{"a", "b", Parameters{}}, Parameters{}},
		{"Maximum length weight", "a/a;q=0.001,a/b;q=0.002", []MediaType{
			{"a", "a", Parameters{}},
			{"a", "b", Parameters{}},
		}, MediaType{"a", "b", Parameters{}}, Parameters{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, requestError := http.NewRequest(http.MethodGet, "http://test.test", nil)
			if requestError != nil {
				log.Fatal(requestError)
			}

			if len(testCase.header) > 0 {
				request.Header.Set("Accept", testCase.header)
			}

			result, extensionParameters, mediaTypeError := GetAcceptableMediaType(request, testCase.availableMediaTypes)

			if mediaTypeError != nil {
				t.Errorf("Unexpected error \"%s\" for %s", mediaTypeError.Error(), testCase.header)
			} else if result.Type != testCase.result.Type || result.Subtype != testCase.result.Subtype {
				t.Errorf("Invalid content type, got %s/%s, exptected %s/%s for %s", result.Type, result.Subtype, testCase.result.Type, testCase.result.Subtype, testCase.header)
			} else if !reflect.DeepEqual(result.Parameters, testCase.result.Parameters) {
				t.Errorf("Wrong parameters, got %v, expected %v for %s", result.Parameters, testCase.result.Parameters, testCase.header)
			} else if !reflect.DeepEqual(extensionParameters, testCase.extensionParameters) {
				t.Errorf("Wrong extension parameters, got %v, expected %v for %s", extensionParameters, testCase.extensionParameters, testCase.header)
			}
		})
	}
}

func TestGetAcceptableMediaTypeErrors(t *testing.T) {
	testCases := []struct {
		name                string
		header              string
		availableMediaTypes []MediaType
		err                 error
	}{
		{"No available type", "", []MediaType{}, ErrNoAvailableTypeGiven},
		{"No acceptable type", "application/xml", []MediaType{{"application", "json", Parameters{}}}, ErrNoAcceptableTypeFound},
		{"Invalid character after subtype", "application/xml/", []MediaType{{"application", "json", Parameters{}}}, ErrInvalidMediaRange},
		{"Comma after subtype with no parameter", "application/xml,", []MediaType{{"application", "json", Parameters{}}}, ErrInvalidMediaType},
		{"Subtype only", "/xml", []MediaType{{"application", "json", Parameters{}}}, ErrInvalidMediaType},
		{"Type with comma and without subtype", "application/,", []MediaType{{"application", "json", Parameters{}}}, ErrInvalidMediaType},
		{"Invalid character", "a/b c", []MediaType{{"a", "b", Parameters{}}}, ErrInvalidMediaRange},
		{"No value for parameter", "a/b;c", []MediaType{{"a", "b", Parameters{}}}, ErrInvalidParameter},
		{"Wildcard type only", "*/b", []MediaType{{"a", "b", Parameters{}}}, ErrInvalidMediaType},
		{"Invalid character in weight", "a/b;q=a", []MediaType{{"a", "b", Parameters{}}}, ErrInvalidWeight},
		{"Weight bigger than 1.0", "a/b;q=11", []MediaType{{"a", "b", Parameters{}}}, ErrInvalidWeight},
		{"More than 3 digitas after dot", "a/b;q=1.0000", []MediaType{{"a", "b", Parameters{}}}, ErrInvalidWeight},
		{"Invalid character after dot", "a/b;q=1.a", []MediaType{{"a", "b", Parameters{}}}, ErrInvalidWeight},
		{"Invalid digit after dot", "a/b;q=1.100", []MediaType{{"a", "b", Parameters{}}}, ErrInvalidWeight},
		{"Type with weight zero only", "a/b;q=0", []MediaType{{"a", "b", Parameters{}}}, ErrNoAcceptableTypeFound},
		{"No value for extension parameter", "a/a;q=1;ext=", []MediaType{{"a", "a", Parameters{}}}, ErrInvalidParameter},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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
		})
	}
}
