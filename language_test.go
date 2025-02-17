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
		{name: "Language and region", value: "lv-LV", result: contenttype.Language{Language: "lv", Script: "", Region: "LV"}},
		{name: "Language, script, and region", value: "lv-Latn-LV", result: contenttype.Language{Language: "lv", Script: "Latn", Region: "LV"}},
		{name: "Language, lower-case script, and region", value: "en-latn-US", result: contenttype.Language{Language: "en", Script: "Latn", Region: "US"}},
		{name: "Language only", value: "lt", result: contenttype.Language{Language: "lt", Script: "", Region: ""}},
		{name: "Language and lowercase region", value: "lv-lv", result: contenttype.Language{Language: "lv", Script: "", Region: "LV"}},
		{name: "Language and region number", value: "lv-428", result: contenttype.Language{Language: "lv", Script: "", Region: "428"}},
		{name: "Three letter language", value: "lav", result: contenttype.Language{Language: "lav", Script: "", Region: ""}},
		{name: "Language and character variant", value: "sl-rozaj", result: contenttype.Language{Language: "sl", Script: "", Region: "", Variant: "rozaj"}},
		{name: "Language, region, and digit variant", value: "de-CH-1901", result: contenttype.Language{Language: "de", Script: "", Region: "CH", Variant: "1901"}},
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
