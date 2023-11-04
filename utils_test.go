package main

import (
	"io/ioutil"
	"log"
	"testing"
)

// ... (Your method linkCouldBePage goes here)

func TestLinkCouldBePage(t *testing.T) {
	log.SetOutput(ioutil.Discard) // Suppress logs for test

	baseUrl := "http://example.com"

	tests := []struct {
		name   string
		link   string
		result bool
	}{
		{"Non-page extension .css", "http://example.com/style.css", false},
		{"Non-page extension .jpg", "http://example.com/image.jpg", false},
		{"Valid page link", "http://example.com/page", true},
		{"Invalid URL", "http:://example", false},
		{"With query and fragment", "http://example.com/page?query=value#fragment", true},
		{"Relative URL", "/relative/path", true},
		{"With base URL prefix", "http://example.com/another/page", true},
		{"Totally different base URL", "http://another.com/page", false},
		{"Totally different base URL", "http://subdomain.example.com/page", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := linkCouldBePage(tt.link, baseUrl)
			if got != tt.result {
				t.Errorf("linkCouldBePage(%q) = %v, want %v", tt.link, got, tt.result)
			}
		})
	}
}
