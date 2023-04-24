package main

import "strings"

//go:generate gonum -types=ExternalLanguages -output=external_languages.gen.go -input=external_languages.go

// Start:ExternalLanguages

// Name: PHP
// PHP bindings

// Name: JS
// JS bindings

func compareExternalLanguages(a, b any) bool {
	if a == nil || b == nil {
		return a == b
	}

	if a, ok := a.(ExternalLanguage); ok {
		if b, ok := b.(ExternalLanguage); ok {
			return strings.ToLower(a.Value) == strings.ToLower(b.Value)
		}

		if b, ok := b.(string); ok {
			return strings.ToLower(a.Value) == strings.ToLower(b)
		}
	}

	return false
}

// End:ExternalLanguages
