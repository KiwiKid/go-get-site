package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"regexp"
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
			_, err := io.WriteString(w, fmt.Sprintf(`<a class="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300"  href="/sites/%s">%s</a>`, escapedUrl, text))
			return err
		case Post:
			_, err := io.WriteString(w, fmt.Sprintf(`<button  hx-post="/sites/%s" hx-target="#container"  class="bg-blue-500 text-white py-2 px-4 rounded hover:bg-blue-600 focus:ring focus:ring-opacity-50 focus:ring-blue-300 focus:border-blue-300">%s</button> `, escapedUrl, text))
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

func splitQuestion(questionAndAnswer string, questionMode bool) string {
	splitString := strings.Split(questionAndAnswer, "?")

	if len(questionAndAnswer) == 0 {
		log.Printf("Error split question, empty question", questionAndAnswer)
		return ""
	}
	if len(splitString) == 1 {
		log.Printf("Error split question, empty question", questionAndAnswer)
		return questionAndAnswer
	}
	if len(splitString) == 2 {
		if questionMode {
			return splitString[0] + "?"
		} else {
			return splitString[1]
		}
	}

	log.Fatalf("Multiple ? in the question?!?: \n%s", questionAndAnswer)
	panic("woah")
}

func (w Website) getTidyTitle(pageTitle string) string {
	return strings.Replace(pageTitle, w.TitleReplace, "", -1)
}

func splitIntoBlocks(content string) []string {
	// Regular expression to match two or more newline characters, possibly surrounded by other whitespace
	re := regexp.MustCompile(`\.\s*\n\s*\n+`)

	// Split the content by the regular expression
	blocks := re.Split(content, -1)

	// Trim whitespace and filter out any empty blocks
	var filteredBlocks []string
	for _, block := range blocks {
		trimmedBlock := strings.TrimSpace(block)
		if trimmedBlock != "" {
			filteredBlocks = append(filteredBlocks, trimmedBlock)
		}
	}

	return filteredBlocks
}

func substring(s string, start uint, end uint) string {
	// Convert string to rune slice for Unicode safety
	runeSlice := []rune(s)

	// Check for valid start and end indices
	if start < 0 || int(end) > len(runeSlice) || start > end {
		return ""
	}

	// Return the substring
	return string(runeSlice[start:end])
}

type Progress struct {
	Total int
	Done  int
}

func (p Progress) getProgress(withWidth bool) string {
	percentage := float64(p.Done) / float64(p.Total) * 100.0

	if withWidth {
		return fmt.Sprintf("width:%.2f%%", percentage)
	}
	return fmt.Sprintf("%.2f%%", percentage)
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

func stripAnchors(link string) (string, error) {
	parsed, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	// Clear out the fragment (which is the anchor tag part of a URL)
	parsed.Fragment = ""
	return parsed.String(), nil
}

func linkCouldBePage(s string, baseUrl string) bool {
	nonPageExtensions := map[string]bool{
		".css":         true,
		".js":          true,
		".ico":         true,
		".png":         true,
		".jpg":         true,
		".jpeg":        true,
		".gif":         true,
		".bmp":         true,
		".tif":         true,
		".tiff":        true,
		".svg":         true,
		".webp":        true,
		".mp3":         true,
		".wav":         true,
		".mp4":         true,
		".mov":         true,
		".avi":         true,
		".mkv":         true,
		".pdf":         true,
		".xml":         true,
		".txt":         true,
		".less":        true,
		".webmanifest": true,
		".zip":         true,
		".rar":         true,
		".tar":         true,
		".gz":          true,
		".json":        true,
		".woff":        true,
		".woff2":       true,
		".ttf":         true,
		".eot":         true,
		".otf":         true,
		".flv":         true,
		".swf":         true,
		".iso":         true,
	}

	// Parse the href to a URL structure
	u, err := url.Parse(s)
	if err != nil {
		log.Printf("failed to parse:check %s %v\n", s, err)
		return false
	}
	if len(u.String()) == 0 {
		log.Printf("failed to parse:check (empty string) %s\n", s)

		return false
	}
	bu, buErr := url.Parse(baseUrl)
	if buErr != nil {
		return false
	}
	// Extract only the path, ignoring query and fragment
	path := u.EscapedPath()

	log.Printf("linkCouldBePage:check %s\n", s)
	// Check the extension
	for ext := range nonPageExtensions {
		if strings.HasSuffix(strings.ToLower(path), ext) {
			return false
		}
	}

	baseDomain := bu.Host
	pathDomain := u.Host

	if baseDomain == pathDomain || strings.HasPrefix(s, "/") {
		log.Print("linkCouldBePage:VALID\n")
		return true
	} else {
		log.Print("linkCouldBePage:INVALID\n")
		return false
	}
}
