package main

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

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

func stringToUint(s string) (uint, error) {
	i, err := strconv.Atoi(s) // Convert string to int
	if err != nil {
		return 0, err
	}
	if i < 0 {
		return 0, fmt.Errorf("negative value not allowed: %d", i)
	}
	return uint(i), nil
}

func linkCouldBePage(s string) bool {
	nonPageExtensions := map[string]bool{
		".css":   true,
		".js":    true,
		".ico":   true,
		".png":   true,
		".jpg":   true,
		".jpeg":  true,
		".gif":   true,
		".bmp":   true,
		".tif":   true,
		".tiff":  true,
		".svg":   true,
		".webp":  true,
		".mp3":   true,
		".wav":   true,
		".mp4":   true,
		".mov":   true,
		".avi":   true,
		".mkv":   true,
		".pdf":   true,
		".xml":   true,
		".txt":   true,
		".zip":   true,
		".rar":   true,
		".tar":   true,
		".gz":    true,
		".json":  true,
		".woff":  true,
		".woff2": true,
		".ttf":   true,
		".eot":   true,
		".otf":   true,
		".flv":   true,
		".swf":   true,
		".iso":   true,
	}

	// Parse the href to a URL structure
	u, err := url.Parse(s)
	if err != nil {
		return false
	}

	// Extract only the path, ignoring query and fragment
	path := u.EscapedPath()

	// Check the extension
	for ext := range nonPageExtensions {
		if strings.HasSuffix(strings.ToLower(path), ext) {
			return false
		}
	}
	return true
}
