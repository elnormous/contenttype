package contenttype_test

import (
	"errors"
	"log"
	"net/http"
	"reflect"
	"testing"

	"github.com/elnormous/contenttype"
)

func TestNewMediaType(t *testing.T) {
	testCases := []struct {
		name   string
		value  string
		result contenttype.MediaType
	}{
		{"Empty string", "", contenttype.MediaType{}},
		{"Type and subtype", "application/json", contenttype.MediaType{"application", "json", contenttype.Parameters{}}},
		{"Type, subtype, parameter", "a/b;c=d", contenttype.MediaType{"a", "b", contenttype.Parameters{"c": "d"}}},
		{"Subtype only", "/b", contenttype.MediaType{}},
		{"Type only", "a/", contenttype.MediaType{}},
		{"Type, subtype, invalid parameter", "a/b;c", contenttype.MediaType{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := contenttype.NewMediaType(testCase.value)

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
		value  contenttype.MediaType
		result string
	}{
		{"Empty media type", contenttype.MediaType{}, ""},
		{"Type and subtype", contenttype.MediaType{"application", "json", contenttype.Parameters{}}, "application/json"},
		{"Type, subtype, parameter", contenttype.MediaType{"a", "b", contenttype.Parameters{"c": "d"}}, "a/b;c=d"},
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
		result contenttype.MediaType
	}{
		{"Empty header", "", contenttype.MediaType{}},
		{"Type and subtype", "application/json", contenttype.MediaType{"application", "json", contenttype.Parameters{}}},
		{"Wildcard", "*/*", contenttype.MediaType{"*", "*", contenttype.Parameters{}}},
		{"Capital subtype", "Application/JSON", contenttype.MediaType{"application", "json", contenttype.Parameters{}}},
		{"Space in front of type", " application/json ", contenttype.MediaType{"application", "json", contenttype.Parameters{}}},
		{"Capital and parameter", "Application/XML;charset=utf-8", contenttype.MediaType{"application", "xml", contenttype.Parameters{"charset": "utf-8"}}},
		{"White space after parameter", "application/xml;foo=bar ", contenttype.MediaType{"application", "xml", contenttype.Parameters{"foo": "bar"}}},
		{"White space after subtype and before parameter", "application/xml ; foo=bar ", contenttype.MediaType{"application", "xml", contenttype.Parameters{"foo": "bar"}}},
		{"Quoted parameter", "application/xml;foo=\"bar\" ", contenttype.MediaType{"application", "xml", contenttype.Parameters{"foo": "bar"}}},
		{"Quoted empty parameter", "application/xml;foo=\"\" ", contenttype.MediaType{"application", "xml", contenttype.Parameters{"foo": ""}}},
		{"Quoted pair", "application/xml;foo=\"\\\"b\" ", contenttype.MediaType{"application", "xml", contenttype.Parameters{"foo": "\"b"}}},
		{"Whitespace after quoted parameter", "application/xml;foo=\"\\\"B\" ", contenttype.MediaType{"application", "xml", contenttype.Parameters{"foo": "\"b"}}},
		{"Plus in subtype", "a/b+c;a=b;c=d", contenttype.MediaType{"a", "b+c", contenttype.Parameters{"a": "b", "c": "d"}}},
		{"Capital parameter", "a/b;A=B", contenttype.MediaType{"a", "b", contenttype.Parameters{"a": "b"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "http://test.test", nil)
			if err != nil {
				log.Fatal(err)
			}

			if len(testCase.header) > 0 {
				request.Header.Set("Content-Type", testCase.header)
			}

			result, err := contenttype.GetMediaType(request)
			if err != nil {
				t.Errorf("Unexpected error \"%w\" for %s", err, testCase.header)
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
		{"Type only", "Application", contenttype.ErrInvalidMediaType},
		{"Subtype only", "/Application", contenttype.ErrInvalidMediaType},
		{"Type with slash", "Application/", contenttype.ErrInvalidMediaType},
		{"Invalid token character", "a/b\x19", contenttype.ErrInvalidMediaType},
		{"Invalid character after subtype", "Application/JSON/test", contenttype.ErrInvalidMediaType},
		{"No parameter name", "application/xml;=bar ", contenttype.ErrInvalidParameter},
		{"Whitespace and no parameter name", "application/xml; =bar ", contenttype.ErrInvalidParameter},
		{"No value and whitespace", "application/xml;foo= ", contenttype.ErrInvalidParameter},
		{"Invalid character in value", "a/b;c=\x19", contenttype.ErrInvalidParameter},
		{"Invalid character in quoted string", "a/b;c=\"\x19\"", contenttype.ErrInvalidParameter},
		{"Invalid character in quoted pair", "a/b;c=\"\\\x19\"", contenttype.ErrInvalidParameter},
		{"No assignment after parameter", "a/b;c", contenttype.ErrInvalidParameter},
		{"No semicolon before paremeter", "a/b e", contenttype.ErrInvalidMediaType},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "http://test.test", nil)
			if err != nil {
				log.Fatal(err)
			}

			if len(testCase.header) > 0 {
				request.Header.Set("Content-Type", testCase.header)
			}

			_, err = contenttype.GetMediaType(request)
			if err == nil {
				t.Errorf("Expected an error for %s", testCase.header)
			} else if !errors.Is(err, testCase.err) {
				t.Errorf("Unexpected error \"%w\", expected \"%v\" for %s", err, testCase.err, testCase.header)
			}
		})
	}
}

func TestGetAcceptableMediaType(t *testing.T) {
	testCases := []struct {
		name                string
		header              string
		availableMediaTypes []contenttype.MediaType
		result              contenttype.MediaType
		extensionParameters contenttype.Parameters
	}{
		{"Empty header", "", []contenttype.MediaType{
			{"application", "json", contenttype.Parameters{}},
		}, contenttype.MediaType{"application", "json", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Type and subtype", "application/json", []contenttype.MediaType{
			{"application", "json", contenttype.Parameters{}},
		}, contenttype.MediaType{"application", "json", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Capitalized type and subtype", "Application/Json", []contenttype.MediaType{
			{"application", "json", contenttype.Parameters{}},
		}, contenttype.MediaType{"application", "json", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Multiple accept types", "text/plain,application/xml", []contenttype.MediaType{
			{"text", "plain", contenttype.Parameters{}},
		}, contenttype.MediaType{"text", "plain", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Multiple accept types, second available", "text/plain,application/xml", []contenttype.MediaType{
			{"application", "xml", contenttype.Parameters{}},
		}, contenttype.MediaType{"application", "xml", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Accept weight", "text/plain;q=1.0", []contenttype.MediaType{
			{"text", "plain", contenttype.Parameters{}},
		}, contenttype.MediaType{"text", "plain", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Wildcard", "*/*", []contenttype.MediaType{
			{"application", "json", contenttype.Parameters{}},
		}, contenttype.MediaType{"application", "json", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Wildcard subtype", "application/*", []contenttype.MediaType{
			{"application", "json", contenttype.Parameters{}},
		}, contenttype.MediaType{"application", "json", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Weight with dot", "a/b;q=1.", []contenttype.MediaType{
			{"a", "b", contenttype.Parameters{}},
		}, contenttype.MediaType{"a", "b", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Multiple weights", "a/b;q=0.1,c/d;q=0.2", []contenttype.MediaType{
			{"a", "b", contenttype.Parameters{}},
			{"c", "d", contenttype.Parameters{}},
		}, contenttype.MediaType{"c", "d", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Multiple weights and default weight", "a/b;q=0.2,c/d;q=0.2", []contenttype.MediaType{
			{"a", "b", contenttype.Parameters{}},
			{"c", "d", contenttype.Parameters{}},
		}, contenttype.MediaType{"a", "b", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Wildcard subtype and weight", "a/*;q=0.2,a/c", []contenttype.MediaType{
			{"a", "b", contenttype.Parameters{}},
			{"a", "c", contenttype.Parameters{}},
		}, contenttype.MediaType{"a", "c", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Different accept order", "a/b,a/a", []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "b", contenttype.Parameters{}},
		}, contenttype.MediaType{"a", "b", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Wildcard subtype with multiple available types", "a/*", []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "b", contenttype.Parameters{}},
		}, contenttype.MediaType{"a", "a", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Wildcard subtype against weighted type", "a/a;q=0.2,a/*", []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "b", contenttype.Parameters{}},
		}, contenttype.MediaType{"a", "b", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Media type parameter", "a/a;q=0.2,a/a;c=d", []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "a", contenttype.Parameters{"c": "d"}},
		}, contenttype.MediaType{"a", "a", contenttype.Parameters{"c": "d"}}, contenttype.Parameters{}},
		{"Weight and media type parameter", "a/b;q=1;e=e", []contenttype.MediaType{
			{"a", "b", contenttype.Parameters{}},
		}, contenttype.MediaType{"a", "b", contenttype.Parameters{}}, contenttype.Parameters{"e": "e"}},
		{"", "a/*,a/a;q=0", []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "b", contenttype.Parameters{}},
		}, contenttype.MediaType{"a", "b", contenttype.Parameters{}}, contenttype.Parameters{}},
		{"Maximum length weight", "a/a;q=0.001,a/b;q=0.002", []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "b", contenttype.Parameters{}},
		}, contenttype.MediaType{"a", "b", contenttype.Parameters{}}, contenttype.Parameters{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "http://test.test", nil)
			if err != nil {
				log.Fatal(err)
			}

			if len(testCase.header) > 0 {
				request.Header.Set("Accept", testCase.header)
			}

			result, extensionParameters, err := contenttype.GetAcceptableMediaType(request, testCase.availableMediaTypes)

			if err != nil {
				t.Errorf("Unexpected error \"%w\" for %s", err, testCase.header)
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
		availableMediaTypes []contenttype.MediaType
		err                 error
	}{
		{"No available type", "", []contenttype.MediaType{}, contenttype.ErrNoAvailableTypeGiven},
		{"No acceptable type", "application/xml", []contenttype.MediaType{{"application", "json", contenttype.Parameters{}}}, contenttype.ErrNoAcceptableTypeFound},
		{"Invalid character after subtype", "application/xml/", []contenttype.MediaType{{"application", "json", contenttype.Parameters{}}}, contenttype.ErrInvalidMediaRange},
		{"Comma after subtype with no parameter", "application/xml,", []contenttype.MediaType{{"application", "json", contenttype.Parameters{}}}, contenttype.ErrInvalidMediaType},
		{"Subtype only", "/xml", []contenttype.MediaType{{"application", "json", contenttype.Parameters{}}}, contenttype.ErrInvalidMediaType},
		{"Type with comma and without subtype", "application/,", []contenttype.MediaType{{"application", "json", contenttype.Parameters{}}}, contenttype.ErrInvalidMediaType},
		{"Invalid character", "a/b c", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidMediaRange},
		{"No value for parameter", "a/b;c", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidParameter},
		{"Wildcard type only", "*/b", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidMediaType},
		{"Invalid character in weight", "a/b;q=a", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidWeight},
		{"Weight bigger than 1.0", "a/b;q=11", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidWeight},
		{"More than 3 digitas after dot", "a/b;q=1.0000", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidWeight},
		{"Invalid character after dot", "a/b;q=1.a", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidWeight},
		{"Invalid digit after dot", "a/b;q=1.100", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidWeight},
		{"Type with weight zero only", "a/b;q=0", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrNoAcceptableTypeFound},
		{"No value for extension parameter", "a/a;q=1;ext=", []contenttype.MediaType{{"a", "a", contenttype.Parameters{}}}, contenttype.ErrInvalidParameter},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request, err := http.NewRequest(http.MethodGet, "http://test.test", nil)
			if err != nil {
				log.Fatal(err)
			}

			if len(testCase.header) > 0 {
				request.Header.Set("Accept", testCase.header)
			}

			_, _, err = contenttype.GetAcceptableMediaType(request, testCase.availableMediaTypes)
			if err == nil {
				t.Errorf("Expected an error for %s", testCase.header)
			} else if !errors.Is(err, testCase.err) {
				t.Errorf("Unexpected error \"%w\", expected \"%v\" for %s", err, testCase.err, testCase.header)
			}
		})
	}
}
