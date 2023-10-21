package main

import (
	"context"
	"fmt"
	"io"
	"net/url"

	"github.com/a-h/templ"
)

type TagType string

const (
	AHref TagType = "ahref"
	Post          = "post"
)

func linkGenerator(tagType TagType, postfix string, text string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		escapedUrl := url.QueryEscape(url.QueryEscape(postfix))

		switch tagType {
		case AHref:
			_, err := io.WriteString(w, fmt.Sprintf(`<a class="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300"  href="/site/%s">%s</a>`, escapedUrl, text))
			return err
		case Post:
			_, err := io.WriteString(w, fmt.Sprintf(`<button  hx-post="/site/%s" hx-target="#container"  class="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300">%s</button> `, escapedUrl, text))
			return err
		}

		return fmt.Errorf("unsupported tag type: %s", tagType)
	})
}

func doubleEscape(str string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		escapedStr := url.QueryEscape(url.QueryEscape(str))
		_, err := io.WriteString(w, escapedStr)
		return err
	})

}

func addQueryParam(originalURL, fromValue string) (string, error) {
	// Parse the original URL
	u, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	// Get existing query parameters
	q := u.Query()

	// Add or update the 'from' query parameter
	q.Set("from", fromValue)

	// Set the updated query parameters back to the URL
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func removeHTTPScheme(rawURL string) (string, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	// Clear the scheme
	parsedURL.Scheme = ""
	if parsedURL.Host != "" {
		return fmt.Sprintf("%s%s", parsedURL.Host, parsedURL.Path), nil
	}
	return parsedURL.Path, nil
}
