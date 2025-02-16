package contenttype

import (
	"testing"
)

func TestDigitChar(t *testing.T) {
	testCases := []struct {
		name   string
		value  byte
		result bool
	}{
		{name: "Zero", value: '0', result: true},
		{name: "One", value: '1', result: true},
		{name: "Two", value: '2', result: true},
		{name: "Three", value: '3', result: true},
		{name: "Nine", value: '9', result: true},
		{name: "Lower-case letter", value: 'a', result: false},
		{name: "Upper-case letter", value: 'A', result: false},
		{name: "Space", value: ' ', result: false},
		{name: "Slash", value: '/', result: false},
		{name: "Colon", value: ':', result: false},
		{name: "NUL", value: 0x00, result: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := isDigitChar(testCase.value)

			if result != testCase.result {
				t.Fatalf("Invalid digit, got %t, exptected %t for %s", result, testCase.result, string(testCase.value))
			}
		})
	}
}

func TestAlphaChar(t *testing.T) {
	testCases := []struct {
		name   string
		value  byte
		result bool
	}{
		{name: "Lower-case letter", value: 'a', result: true},
		{name: "Upper-case letter", value: 'A', result: true},
		{name: "Lower-case letter", value: 'z', result: true},
		{name: "Upper-case letter", value: 'Z', result: true},
		{name: "Digit", value: '0', result: false},
		{name: "Space", value: ' ', result: false},
		{name: "Slash", value: '/', result: false},
		{name: "Colon", value: ':', result: false},
		{name: "Underscore", value: '_', result: false},
		{name: "Hyphen", value: '-', result: false},
		{name: "Period", value: '.', result: false},
		{name: "Apostrophe", value: '\'', result: false},
		{name: "Exclamation mark", value: '!', result: false},
		{name: "Brace", value: '}', result: false},
		{name: "Bracket", value: '[', result: false},
		{name: "At", value: '@', result: false},
		{name: "Grave accent", value: '`', result: false},
		{name: "NUL", value: 0x00, result: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := isAlphaChar(testCase.value)

			if result != testCase.result {
				t.Fatalf("Invalid alpha, got %t, exptected %t for %s", result, testCase.result, string(testCase.value))
			}
		})
	}
}

func TestVisibleChar(t *testing.T) {
	testCases := []struct {
		name   string
		value  byte
		result bool
	}{
		{name: "Lower-case letter", value: 'a', result: true},
		{name: "Upper-case letter", value: 'A', result: true},
		{name: "Exclamation mark", value: '!', result: true},
		{name: "Underscore", value: '_', result: true},
		{name: "Digit", value: '0', result: true},
		{name: "Space", value: ' ', result: false},
		{name: "NUL", value: 0x00, result: false},
		{name: "US", value: 0x1F, result: false},
		{name: "DEL", value: 0x7F, result: false},
		{name: "Tilde", value: '~', result: true},
		{name: "Out of ASCII range", value: 0x80, result: false},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := isVisibleChar(testCase.value)

			if result != testCase.result {
				t.Fatalf("Invalid visible, got %t, exptected %t for %s", result, testCase.result, string(testCase.value))
			}
		})
	}
}

func TestSkipChar(t *testing.T) {
	testCases := []struct {
		name      string
		value     string
		skip      byte
		remaining string
		result    bool
	}{
		{name: "Empty string", value: "", skip: ' ', remaining: "", result: false},
		{name: "Skip one letter", value: "A", skip: 'A', remaining: "", result: true},
		{name: "Skip a different letter", value: "A", skip: 'B', remaining: "A", result: false},
		{name: "Skip string with same letters", value: "AA", skip: 'A', remaining: "A", result: true},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			remaining, skipped := skipCharacter(testCase.value, testCase.skip)

			if remaining != testCase.remaining {
				t.Fatalf("Invalid remaining, got %s, exptected %s for %s", remaining, testCase.remaining, testCase.value)
			}

			if skipped != testCase.result {
				t.Fatalf("Invalid result, got %v, exptected %v for %s", skipped, testCase.result, testCase.value)
			}
		})
	}
}
