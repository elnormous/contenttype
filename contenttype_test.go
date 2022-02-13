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
		{name: "Empty string", value: "", result: contenttype.MediaType{}},
		{name: "Type and subtype", value: "application/json", result: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}},
		{name: "Type, subtype, parameter", value: "a/b;c=d", result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{"c": "d"}}},
		{name: "Subtype only", value: "/b", result: contenttype.MediaType{}},
		{name: "Type only", value: "a/", result: contenttype.MediaType{}},
		{name: "Type, subtype, invalid parameter", value: "a/b;c", result: contenttype.MediaType{}},
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
		{name: "Empty media type", value: contenttype.MediaType{}, result: ""},
		{name: "Type and subtype", value: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}, result: "application/json"},
		{name: "Type, subtype, parameter", value: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{"c": "d"}}, result: "a/b;c=d"},
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
		{name: "Empty header", header: "", result: contenttype.MediaType{}},
		{name: "Type and subtype", header: "application/json", result: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}},
		{name: "Wildcard", header: "*/*", result: contenttype.MediaType{Type: "*", Subtype: "*", Parameters: contenttype.Parameters{}}},
		{name: "Capital subtype", header: "Application/JSON", result: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}},
		{name: "Space in front of type", header: " application/json ", result: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}},
		{name: "Capital and parameter", header: "Application/XML;charset=utf-8", result: contenttype.MediaType{Type: "application", Subtype: "xml", Parameters: contenttype.Parameters{"charset": "utf-8"}}},
		{name: "Spaces around semicolon", header: "a/b ; c=d", result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{"c": "d"}}},
		{name: "Spaces around semicolons", header: "a/b ; c=d ; e=f", result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{"c": "d", "e": "f"}}},
		{name: "Two spaces around semicolons", header: "a/b  ;  c=d  ;  e=f", result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{"c": "d", "e": "f"}}},
		{name: "White space after parameter", header: "application/xml;foo=bar ", result: contenttype.MediaType{Type: "application", Subtype: "xml", Parameters: contenttype.Parameters{"foo": "bar"}}},
		{name: "White space after subtype and before parameter", header: "application/xml ; foo=bar ", result: contenttype.MediaType{Type: "application", Subtype: "xml", Parameters: contenttype.Parameters{"foo": "bar"}}},
		{name: "Quoted parameter", header: "application/xml;foo=\"bar\" ", result: contenttype.MediaType{Type: "application", Subtype: "xml", Parameters: contenttype.Parameters{"foo": "bar"}}},
		{name: "Quoted empty parameter", header: "application/xml;foo=\"\" ", result: contenttype.MediaType{Type: "application", Subtype: "xml", Parameters: contenttype.Parameters{"foo": ""}}},
		{name: "Quoted pair", header: "application/xml;foo=\"\\\"b\" ", result: contenttype.MediaType{Type: "application", Subtype: "xml", Parameters: contenttype.Parameters{"foo": "\"b"}}},
		{name: "Whitespace after quoted parameter", header: "application/xml;foo=\"\\\"B\" ", result: contenttype.MediaType{Type: "application", Subtype: "xml", Parameters: contenttype.Parameters{"foo": "\"b"}}},
		{name: "Plus in subtype", header: "a/b+c;a=b;c=d", result: contenttype.MediaType{Type: "a", Subtype: "b+c", Parameters: contenttype.Parameters{"a": "b", "c": "d"}}},
		{name: "Capital parameter", header: "a/b;A=B", result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{"a": "b"}}},
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
				t.Errorf("Unexpected error \"%v\" for %s", err, testCase.header)
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
		{"No semicolon before parameter", "a/b e", contenttype.ErrInvalidMediaType},
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
				t.Errorf("Unexpected error \"%v\", expected \"%v\" for %s", err, testCase.err, testCase.header)
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
		{name: "Empty header", availableMediaTypes: []contenttype.MediaType{
			{"application", "json", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Type and subtype", header: "application/json", availableMediaTypes: []contenttype.MediaType{
			{"application", "json", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Capitalized type and subtype", header: "Application/Json", availableMediaTypes: []contenttype.MediaType{
			{"application", "json", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Multiple accept types", header: "text/plain,application/xml", availableMediaTypes: []contenttype.MediaType{
			{"text", "plain", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "text", Subtype: "plain", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Multiple accept types, second available", header: "text/plain,application/xml", availableMediaTypes: []contenttype.MediaType{
			{"application", "xml", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "application", Subtype: "xml", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Accept weight", header: "text/plain;q=1.0", availableMediaTypes: []contenttype.MediaType{
			{"text", "plain", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "text", Subtype: "plain", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Wildcard", header: "*/*", availableMediaTypes: []contenttype.MediaType{
			{"application", "json", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Wildcard subtype", header: "application/*", availableMediaTypes: []contenttype.MediaType{
			{"application", "json", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Weight with dot", header: "a/b;q=1.", availableMediaTypes: []contenttype.MediaType{
			{"a", "b", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Multiple weights", header: "a/b;q=0.1,c/d;q=0.2", availableMediaTypes: []contenttype.MediaType{
			{"a", "b", contenttype.Parameters{}},
			{"c", "d", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "c", Subtype: "d", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Multiple weights and default weight", header: "a/b;q=0.2,c/d;q=0.2", availableMediaTypes: []contenttype.MediaType{
			{"a", "b", contenttype.Parameters{}},
			{"c", "d", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Wildcard subtype and weight", header: "a/*;q=0.2,a/c", availableMediaTypes: []contenttype.MediaType{
			{"a", "b", contenttype.Parameters{}},
			{"a", "c", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "a", Subtype: "c", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Different accept order", header: "a/b,a/a", availableMediaTypes: []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "b", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Wildcard subtype with multiple available types", header: "a/*", availableMediaTypes: []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "b", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "a", Subtype: "a", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Wildcard subtype against weighted type", header: "a/a;q=0.2,a/*", availableMediaTypes: []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "b", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Media type parameter", header: "a/a;q=0.2,a/a;c=d", availableMediaTypes: []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "a", contenttype.Parameters{"c": "d"}},
		}, result: contenttype.MediaType{Type: "a", Subtype: "a", Parameters: contenttype.Parameters{"c": "d"}}, extensionParameters: contenttype.Parameters{}},
		{name: "Weight and media type parameter", header: "a/b;q=1;e=e", availableMediaTypes: []contenttype.MediaType{
			{"a", "b", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{"e": "e"}},
		{header: "a/*,a/a;q=0", availableMediaTypes: []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "b", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
		{name: "Maximum length weight", header: "a/a;q=0.001,a/b;q=0.002", availableMediaTypes: []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
			{"a", "b", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
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
				t.Errorf("Unexpected error \"%v\" for %s", err, testCase.header)
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
				t.Errorf("Unexpected error \"%v\", expected \"%v\" for %s", err, testCase.err, testCase.header)
			}
		})
	}
}
