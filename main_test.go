package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var serverURL string

func TestMain(m *testing.M) {
	ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {

		var URI string

		if req.Method != "GET" {
			panic(fmt.Errorf("unexpected request type %s", req.Method))
		}

		switch req.URL.String() {
		case "/about":
			URI = "stubs/site/about.html"
		case "/contact":
			URI = "stubs/site/contact.html"
		case "/":
			URI = "stubs/site/index.html"
		default:
			panic(fmt.Errorf("request domain %s not set up in test server", req.URL.String()))
		}

		body, err := ioutil.ReadFile(URI)

		if err != nil {
			panic(err)
		}

		res.Header().Set("Content-Type", "text/html")
		_, err = res.Write(body)

		if err != nil {
			panic(err)
		}
	}))

	serverURL = ts.URL

	defer ts.Close()
	os.Exit(m.Run())
}

func FindCrawlResultByRequestURL(results []CrawlResult, requestURL string) (CrawlResult, bool) {
	for _, result := range results {
		if result.requestURL == requestURL {
			return result, true
		}
	}

	return CrawlResult{}, false
}

func AssertStringSetsEqual(actual StringSet, expected StringSet, t *testing.T) {
	if len(actual) != len(expected) {
		t.Errorf("actual set size %d, expected %d", len(actual), len(expected))
	}

	for v := range expected {
		_, ok := actual[v]

		if !ok {
			t.Errorf("actual set missing expected element %s", v)
		}
	}
}

func AssertEqualCrawlResults(actual []CrawlResult, expected []CrawlResult, t *testing.T) {
	if len(actual) != len(expected) {
		t.Errorf("got %d crawl results, expected %d", len(actual), len(expected))
	}

	for _, e := range expected {
		a, ok := FindCrawlResultByRequestURL(actual, e.requestURL)

		if !ok {
			t.Errorf("expected crawl result with requestURL: %s", e.requestURL)
			continue
		}

		AssertStringSetsEqual(a.externalLinks, e.externalLinks, t)
		AssertStringSetsEqual(a.internalLinks, e.internalLinks, t)
	}
}

func TestGivenValidDomain_WhenCrawl_ReturnsExpectedLinks(t *testing.T) {
	fmt.Println("server URL: " + serverURL)
	actualCrawlResults := Crawl(serverURL)

	expectedResults := []CrawlResult{
		{
			requestURL:    serverURL,
			internalLinks: NewStringSet([]string{serverURL + "/about"}),
			externalLinks: NewStringSet([]string{"https://monzo.com/contact", "https://monzo.com/about", "https://monzo.com/contact?name=Aran"}),
		},
		{
			requestURL:    serverURL + "/about",
			internalLinks: NewStringSet([]string{serverURL, serverURL + "/contact"}),
			externalLinks: NewStringSet([]string{"https://aran.site/hello/world", "https://aran.site/hello/world?name=aran"}),
		},
		{
			requestURL:    serverURL + "/contact",
			internalLinks: NewStringSet([]string{serverURL}),
			externalLinks: NewStringSet([]string{"https://youtube.com/hello/world"}),
		},
	}

	AssertEqualCrawlResults(actualCrawlResults, expectedResults, t)
}
