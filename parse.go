package main

import (
	"golang.org/x/net/html"
	"io"
	"net/url"
	"strings"
)

// isVisitableLink checks if a given href links to a different document
func isVisitableLink(link string) bool {
	return !(link == "" || strings.HasPrefix(link, "#"))
}

func removeTrailingSlash(URL string) string {
	return strings.TrimRight(URL, "/")
}

// FindLinks finds all of the href links in a HTML stream provided by htmlReader
func FindLinks(rawBaseURL string, HTMLReader io.Reader) (CrawlResult, error) {
	links, err := extractRawHrefs(HTMLReader)

	if err != nil {
		return CrawlResult{}, err
	}

	baseURL, err := url.Parse(rawBaseURL)

	if err != nil {
		return CrawlResult{}, err
	}

	parsedLinks := make(map[string]struct{})

	var response = CrawlResult{
		"",
		StringSet{},
		StringSet{},
	}

	for _, link := range links {
		if isVisitableLink(link) {
			linkUrl, err := url.Parse(removeTrailingSlash(link))

			if err != nil {
				continue
			}

			// remove the hashbang as this is a same-page link
			linkUrl.Fragment = ""

			if !linkUrl.IsAbs() {
				linkUrl = baseURL.ResolveReference(linkUrl)
			}

			if !(linkUrl.Scheme == "http" || linkUrl.Scheme == "https") {
				continue
			}

			_, parsed := parsedLinks[linkUrl.String()]

			if !parsed {
				if linkUrl.Hostname() == baseURL.Hostname() {
					response.internalLinks.Add(linkUrl.String())
					parsedLinks[linkUrl.String()] = struct{}{}
				} else {
					response.externalLinks.Add(linkUrl.String())
				}
			}
		}
	}

	return response, nil
}

// extractRawHrefs extracts all href values from anchor tags in a HTML stream provided by htmlReader
func extractRawHrefs(htmlReader io.Reader) ([]string, error) {
	tokenizer := html.NewTokenizer(htmlReader)
	var links []string

	for {
		tt := tokenizer.Next()

		switch tt {

		case html.ErrorToken:
			err := tokenizer.Err()

			if err == io.EOF {
				return links, nil
			}

			return nil, err

		case html.StartTagToken:
			token := tokenizer.Token()

			if token.Data == "a" {
				for _, a := range token.Attr {
					if a.Key == "href" {
						links = append(links, a.Val)
						break
					}
				}
			}
		}
	}
}
