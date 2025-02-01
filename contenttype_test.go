package contenttype_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/elnormous/contenttype"
)

var (
	// selection of pre-defined values used for testing
	instEmpty        contenttype.MediaType = contenttype.MediaType{}
	instSimple       contenttype.MediaType = contenttype.NewMediaType("text/plain")
	instWildcard     contenttype.MediaType = contenttype.NewMediaType("*/*")
	instTextWildcard contenttype.MediaType = contenttype.NewMediaType("text/*")
	instParams       contenttype.MediaType = contenttype.NewMediaType("application/json; q=0.001; charset=utf-8")
	instJSON         contenttype.MediaType = contenttype.NewMediaType("application/json")
	instJSON2        contenttype.MediaType = contenttype.NewMediaType("application/json; charset=utf-8")
	instAppWildcard  contenttype.MediaType = contenttype.NewMediaType("application/*")
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

func TestParseMediaType(t *testing.T) {
	testCases := []struct {
		name   string
		value  string
		result contenttype.MediaType
	}{
		{name: "Type and subtype", value: "application/json", result: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}},
		{name: "Type and subtype with whitespaces", value: "application/json   ", result: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}},
		{name: "Type, subtype, parameter", value: "a/b;c=d", result: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{"c": "d"}}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := contenttype.ParseMediaType(testCase.value)
			if err != nil {
				t.Errorf("Unexpected error \"%v\" for %s", err, testCase.value)
			} else if result.Type != testCase.result.Type || result.Subtype != testCase.result.Subtype {
				t.Fatalf("Invalid content type, got %s/%s, exptected %s/%s for %s", result.Type, result.Subtype, testCase.result.Type, testCase.result.Subtype, testCase.value)
			} else if !reflect.DeepEqual(result.Parameters, testCase.result.Parameters) {
				t.Fatalf("Wrong parameters, got %v, expected %v for %s", result.Parameters, testCase.result.Parameters, testCase.value)
			}
		})
	}
}

func TestParseMediaTypeErrors(t *testing.T) {
	testCases := []struct {
		name  string
		value string
		err   error
	}{
		{name: "Empty string", value: "", err: contenttype.ErrInvalidMediaType},
		{name: "Subtype only", value: "/b", err: contenttype.ErrInvalidMediaType},
		{name: "Type only", value: "a/", err: contenttype.ErrInvalidMediaType},
		{name: "Type, subtype, invalid parameter", value: "a/b;c", err: contenttype.ErrInvalidParameter},
		{name: "Type and parameter without subtype", value: "a/;c", err: contenttype.ErrInvalidMediaType},
		{name: "Type and subtype with remaining data", value: "a/b c", err: contenttype.ErrInvalidMediaType},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			_, err := contenttype.ParseMediaType(testCase.value)
			if err == nil {
				t.Errorf("Expected an error for %s", testCase.value)
			} else if !errors.Is(err, testCase.err) {
				t.Errorf("Unexpected error \"%v\", expected \"%v\" for %s", err, testCase.err, testCase.value)
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

func TestMediaType_MIME(t *testing.T) {
	testCases := []struct {
		name   string
		value  contenttype.MediaType
		result string
	}{
		{name: "Empty media type", value: contenttype.MediaType{}, result: ""},
		{name: "Type and subtype", value: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}, result: "application/json"},
		{name: "Type, subtype, parameter", value: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{"c": "d"}}, result: "a/b"},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.value.MIME()

			if result != testCase.result {
				t.Errorf("Invalid result type, got %s, exptected %s", result, testCase.result)
			}
		})
	}
}

func TestMediaType_IsWildcard(t *testing.T) {
	testCases := []struct {
		name   string
		value  contenttype.MediaType
		result bool
	}{
		{name: "Empty media type", value: contenttype.MediaType{}, result: false},
		{name: "Type and subtype", value: contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{}}, result: false},
		{name: "Type, subtype, parameter", value: contenttype.MediaType{Type: "a", Subtype: "b", Parameters: contenttype.Parameters{"c": "d"}}, result: false},
		{name: "text/*", value: contenttype.MediaType{Type: "text", Subtype: "*"}, result: true},
		{name: "application/*; charset=utf-8", value: contenttype.MediaType{Type: "application", Subtype: "*", Parameters: contenttype.Parameters{"charset": "utf-8"}}, result: true},
		{name: "*/*", value: contenttype.MediaType{Type: "*", Subtype: "*"}, result: true},
		// invalid MIME type, but will return true
		{name: "*/json", value: contenttype.MediaType{Type: "*", Subtype: "json"}, result: true},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := testCase.value.IsWildcard()
			if result != testCase.result {
				t.Errorf("Invalid result type, got %v, expected %v", result, testCase.result)
			}
		})
	}
}

// ExampleMediaType_MIME comparing to MIME types
func ExampleMediaType_MIME() {
	mt := contenttype.NewMediaType("application/json; charset=utf-8")
	fmt.Printf("MIME(): %s\n", mt.MIME())
	fmt.Printf("matches: application/json: %v\n", mt.MIME() == "application/json")
	fmt.Printf("matches: application/*: %v\n", mt.MIME() == "application/*")
	fmt.Printf("matches: text/plain: %v\n", mt.MIME() == "text/plain")
	// Output: MIME(): application/json
	// matches: application/json: true
	// matches: application/*: false
	// matches: text/plain: false
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
			request := httptest.NewRequest(http.MethodGet, "http://test.test", nil)

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
			request := httptest.NewRequest(http.MethodGet, "http://test.test", nil)

			if len(testCase.header) > 0 {
				request.Header.Set("Content-Type", testCase.header)
			}

			_, err := contenttype.GetMediaType(request)
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
		{name: "Spaces around comma", header: "a/a;q=0.1 , a/b , a/c", availableMediaTypes: []contenttype.MediaType{
			{"a", "a", contenttype.Parameters{}},
		}, result: contenttype.MediaType{Type: "a", Subtype: "a", Parameters: contenttype.Parameters{}}, extensionParameters: contenttype.Parameters{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://test.test", nil)

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
		{"More than 3 digits after dot", "a/b;q=1.0000", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidWeight},
		{"Invalid character after dot", "a/b;q=1.a", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidWeight},
		{"Invalid digit after dot", "a/b;q=1.100", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidWeight},
		{"Weight with two dots", "a/b;q=0..1", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrInvalidWeight},
		{"Type with weight zero only", "a/b;q=0", []contenttype.MediaType{{"a", "b", contenttype.Parameters{}}}, contenttype.ErrNoAcceptableTypeFound},
		{"No value for extension parameter", "a/a;q=1;ext=", []contenttype.MediaType{{"a", "a", contenttype.Parameters{}}}, contenttype.ErrInvalidParameter},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "http://test.test", nil)

			if len(testCase.header) > 0 {
				request.Header.Set("Accept", testCase.header)
			}

			_, _, err := contenttype.GetAcceptableMediaType(request, testCase.availableMediaTypes)
			if err == nil {
				t.Errorf("Expected an error for %s", testCase.header)
			} else if !errors.Is(err, testCase.err) {
				t.Errorf("Unexpected error \"%v\", expected \"%v\" for %s", err, testCase.err, testCase.header)
			}
		})
	}
}

func TestMediaType_Equal(t *testing.T) {
	// create a map of items to turn into a permutation, these should all be
	// different
	mediaTypes := map[string]contenttype.MediaType{
		"empty":        instEmpty,
		"simple":       instSimple,
		"wildcard":     instWildcard,
		"textwildcard": instTextWildcard,
		"params":       instParams,
		"json":         instJSON,
		"json2":        instJSON2,
		"appwildcard":  instAppWildcard,
	}

	type test struct {
		name string
		a    contenttype.MediaType
		b    contenttype.MediaType
		want bool
	}
	tests := []test{}

	// create permutation
	for outerName, outerMt := range mediaTypes {
		for innerName, innerMt := range mediaTypes {
			tests = append(tests,
				test{
					fmt.Sprintf("%s vs %s", outerName, innerName),
					outerMt,
					innerMt,
					outerName == innerName,
				})
		}
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Equal(tt.b); got != tt.want {
				t.Errorf("MediaType.Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ExampleMediaType_MIME comparing two media types with their parameters
func ExampleMediaType_Equal() {
	base := contenttype.NewMediaType("application/json; charset=utf-8")
	noMatch := contenttype.NewMediaType("application/json")
	match := contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{"charset": "utf-8"}}

	fmt.Printf("matches exactly: %v\n", base.Equal(base))
	fmt.Printf("matches exactly: %v\n", base.Equal(noMatch))
	fmt.Printf("matches exactly: %v\n", base.Equal(match))
	// Output: matches exactly: true
	// matches exactly: false
	// matches exactly: true
}

func TestMediaType_EqualsMIME(t *testing.T) {
	// create a map of items to turn into a permutation, these should all be
	// different
	mtut := map[string]contenttype.MediaType{
		"empty":        instEmpty,
		"simple":       instSimple,
		"wildcard":     instWildcard,
		"textwildcard": instTextWildcard,
		"appwildcard":  instAppWildcard,
		"params":       instParams,
	}

	type test struct {
		name string
		a    contenttype.MediaType
		b    contenttype.MediaType
		want bool
	}
	tests := []test{
		// all of these are equal
		{"params vs json", instParams, instJSON, true},
		{"params vs json2", instParams, instJSON2, true},
		{"json vs params", instJSON, instParams, true},
		{"json2 vs params", instJSON2, instParams, true},
		{"json vs json", instJSON, instJSON, true},
		{"json vs json2", instJSON, instJSON2, true},
		{"json2 vs json", instJSON2, instJSON, true},
		{"json2 vs json2", instJSON2, instJSON2, true},
	}

	// create permutation of the remaining tests from the map which are only equal
	// to themselves
	for outerName, outerMt := range mtut {
		for innerName, innerMt := range mtut {
			tests = append(tests,
				test{
					fmt.Sprintf("%s vs %s", outerName, innerName),
					outerMt,
					innerMt,
					outerName == innerName,
				})
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.EqualsMIME(tt.b); got != tt.want {
				t.Errorf("MediaType.EqualsMIME() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ExampleMediaType_EqualsMIME comparing only on the MIME media type which must
// match both the type and subtype
func ExampleMediaType_EqualsMIME() {
	base := contenttype.NewMediaType("application/json; q=0.01; charset=utf-8")
	noMatch := contenttype.NewMediaType("text/json")
	partialWildcard := contenttype.NewMediaType("application/*")
	diffParams := contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{"charset": "utf-8"}}
	match := contenttype.MediaType{Type: "application", Subtype: "json"}

	fmt.Printf("matches exactly: %v\n", base.EqualsMIME(base))
	fmt.Printf("matches exactly: %v\n", base.EqualsMIME(noMatch))
	fmt.Printf("matches exactly: %v\n", base.EqualsMIME(partialWildcard))
	fmt.Printf("matches exactly: %v\n", base.EqualsMIME(diffParams))
	fmt.Printf("matches exactly: %v\n", base.EqualsMIME(match))
	// Output: matches exactly: true
	// matches exactly: false
	// matches exactly: false
	// matches exactly: true
	// matches exactly: true
}

func TestMediaType_Matches(t *testing.T) {
	tests := []struct {
		name string
		a    contenttype.MediaType
		b    contenttype.MediaType
		want bool
	}{
		{"empty matches empty", instEmpty, instEmpty, true},
		{"text/plain matches text/plain", instSimple, instSimple, true},
		{"text/* matches text/plain", instTextWildcard, instSimple, true},
		{"*/* matches text/plain", instWildcard, instSimple, true},
		{"text/plain matches text/*", instSimple, instTextWildcard, true},
		{"text/plain matches */*", instSimple, instWildcard, true},
		{"text/plain doesn't match application/*", instSimple, instAppWildcard, false},
		{"text/* doesn't match application/*", instTextWildcard, instAppWildcard, false},
		{"*/* matches application/*", instWildcard, instAppWildcard, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.Matches(tt.b); got != tt.want {
				t.Errorf("MediaType.Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

// ExampleMediaType_Matches comparing only on the MIME media type which must
// either match exactly or be the wildcard token
func ExampleMediaType_Matches() {
	base := contenttype.NewMediaType("application/json; q=0.01; charset=utf-8")
	noMatch := contenttype.NewMediaType("text/json")
	partialWildcard := contenttype.NewMediaType("application/*")
	fullWildcard := contenttype.NewMediaType("*/*")
	diffParams := contenttype.MediaType{Type: "application", Subtype: "json", Parameters: contenttype.Parameters{"charset": "utf-8"}}
	match := contenttype.MediaType{Type: "application", Subtype: "json"}

	fmt.Printf("matches exactly: %v\n", base.Matches(base))
	fmt.Printf("matches exactly: %v\n", base.Matches(noMatch))
	fmt.Printf("matches exactly: %v\n", base.Matches(partialWildcard))
	fmt.Printf("matches exactly: %v\n", base.Matches(fullWildcard))
	fmt.Printf("matches exactly: %v\n", base.Matches(diffParams))
	fmt.Printf("matches exactly: %v\n", base.Matches(match))
	// Output: matches exactly: true
	// matches exactly: false
	// matches exactly: true
	// matches exactly: true
	// matches exactly: true
	// matches exactly: true
}

func TestMediaType_MatchesAny(t *testing.T) {
	tests := []struct {
		name string
		a    contenttype.MediaType
		bs   []contenttype.MediaType
		want bool
	}{
		{"vs no list", instEmpty, nil, false},
		{"vs empty list", instEmpty, []contenttype.MediaType{}, false},
		{"empty vs matching single", instEmpty, []contenttype.MediaType{instEmpty}, true},
		{"empty vs non-matching single", instEmpty, []contenttype.MediaType{instJSON}, false},
		{"empty vs second match", instEmpty, []contenttype.MediaType{instJSON, instEmpty}, true},
		{"specific vs wildcard only", instSimple, []contenttype.MediaType{instTextWildcard}, true},
		{"specific vs second item wildcard", instSimple, []contenttype.MediaType{instJSON, instTextWildcard}, true},
		{"wildcard vs anything", instWildcard, []contenttype.MediaType{instJSON}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.MatchesAny(tt.bs...); got != tt.want {
				t.Errorf("MediaType.MatchesAny() = %v, want %v", got, tt.want)
			}
		})
	}
}
