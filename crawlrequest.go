package main

import (
	"fmt"
	"net/http"
)

// CrawlResult represents the response to a page RequestLinks
type CrawlResult struct {
	requestURL                   string
	internalLinks, externalLinks StringSet
}

// RequestLinks finds all href links from a web page with a given url
func RequestLinks(link string, baseURL string) (CrawlResult, error) {
	resp, err := http.Get(link)

	if err != nil {
		return CrawlResult{}, fmt.Errorf("get request to url: %s failed", link)
	}

	links, err := FindLinks(baseURL, resp.Body)

	if err != nil {
		return CrawlResult{}, fmt.Errorf("failed to find links in response body for url: %s", link)
	}

	links.requestURL = link

	err = resp.Body.Close()

	if err != nil {
		return CrawlResult{}, fmt.Errorf("failed to close response body for url: %s", link)
	}

	return links, nil
}
