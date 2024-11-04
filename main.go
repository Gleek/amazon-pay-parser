package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/net/html"
)

// Extract text content from a node
func extractText(n *html.Node) string {
	var buf bytes.Buffer
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	normalized := regexp.MustCompile(`\s+`).ReplaceAllString(buf.String(), " ")
	return strings.TrimSpace(normalized)
}

// hasExactClasses checks if a node has exactly the given classes
func hasExactClasses(n *html.Node, classNames []string) bool {
	for _, attr := range n.Attr {
		if attr.Key == "class" {
			classes := strings.Fields(attr.Val)
			if len(classes) != len(classNames) {
				return false
			}
			sort.Strings(classes)
			sort.Strings(classNames)
			for i := range classes {
				if classes[i] != classNames[i] {
					return false
				}
			}
			return true
		}
	}
	return false
}

// Recursively search for nodes with the exact classes and extract their text
func extractClassContent(n *html.Node, classMap map[string]*string) {
	if n.Type == html.ElementNode {
		for className, content := range classMap {
			classParts := strings.Fields(className)
			if hasExactClasses(n, classParts) {
				*content = extractText(n)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractClassContent(c, classMap)
	}
}

// Process each span with id=itemDetailExpandedView and extract text content for final classes
func processHTML(n *html.Node) {
	spans := findNodes(n, "itemDetailExpandedView", true)

	for _, span := range spans {
		aRowDivs := findNodes(span, "a-row", false)
		if len(aRowDivs) == 0 {
			continue
		}
		classMap := map[string]*string{
			"a-section pad-header-text":                                 new(string), // Party
			"a-section payment-details-desktop":                         new(string), // Medium
			"a-size-base a-color-tertiary":                              new(string), // Date
			"a-column a-span3 a-text-right pad-header-text a-span-last": new(string), // Amount
		}

		extractClassContent(aRowDivs[0], classMap)

		party := *classMap["a-section pad-header-text"]
		medium := *classMap["a-section payment-details-desktop"]
		date := *classMap["a-size-base a-color-tertiary"]
		amount := *classMap["a-column a-span3 a-text-right pad-header-text a-span-last"]

		// Print CSV row if all fields are filled
		if party != "" && medium != "" && date != "" && amount != "" {
			fmt.Printf("%q,%q,%q,%q\n", party, medium, date, cleanAmount(amount))
		}
	}
}

// only keep digits and decimal point and negative sign
func cleanAmount(amount string) string {
	return regexp.MustCompile(`[^0-9.-]`).ReplaceAllString(amount, "")
}

// Find nodes with specific classes
func findNodes(n *html.Node, class string, isID bool) []*html.Node {
	var result []*html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, attr := range n.Attr {
				if isID && attr.Key == "id" && attr.Val == class {
					result = append(result, n)
				} else if !isID && attr.Key == "class" && strings.TrimSpace(attr.Val) == class {
					result = append(result, n)
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(n)
	return result
}

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <html-file-path>", os.Args[0])
	}

	filePath := os.Args[1]
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	processHTML(doc)
}
