// go-crawl-redux implements a concurrent same-domain web crawler
//
// domains are crawled using a worker pool pattern, and results are
// separated by whether they are internal (same domain) or external links
//
package main

import (
	"fmt"
)

func main() {

	URL := "https://monzo.com"
	crawlResults := Crawl(URL)

	for _, r := range crawlResults {
		fmt.Printf("\n\n|--- %s\n", r.requestURL)

		fmt.Printf("|--- [%d internal links]\n", len(r.internalLinks))
		for l := range r.internalLinks {
			fmt.Printf("| |--- %s\n", l)
		}

		fmt.Printf("|--- [%d external links]\n", len(r.externalLinks))
		for l := range r.externalLinks {
			fmt.Printf("| |--- %s\n", l)
		}
	}
}
