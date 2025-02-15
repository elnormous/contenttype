package contenttype_test

import (
	"testing"

	"github.com/elnormous/contenttype"
)

func TestNewLanguage(t *testing.T) {
	testCases := []struct {
		name   string
		value  string
		result contenttype.Language
	}{
		{name: "Empty string", value: "", result: contenttype.Language{}},
		{name: "Language and region", value: "lv-LV", result: contenttype.Language{Language: "lv", Script: "", Region: "lv"}},
		{name: "Language, script, and region", value: "en-latin-US", result: contenttype.Language{Language: "en", Script: "latin", Region: "us"}},
		{name: "Language only", value: "lt", result: contenttype.Language{Language: "lt", Script: "", Region: ""}},
		{name: "Language and lowercase region", value: "lv-lv", result: contenttype.Language{Language: "lv", Script: "", Region: "lv"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := contenttype.NewLanguage(testCase.value)

			if result.Language != testCase.result.Language || result.Script != testCase.result.Script || result.Region != testCase.result.Region {
				t.Fatalf("Invalid language, got %s, exptected %s for %s", result, testCase.result, testCase.value)
			}
		})
	}
}
