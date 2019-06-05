package main

import "fmt"

// worker completes a crawl of a single URL, finding the contained links and
// sending these back through the results channel
func worker(url string, tasks <-chan string, results chan<- CrawlResult) {
	for t := range tasks {
		links, err := RequestLinks(t, url)

		if err != nil {
			fmt.Println(err)
		} else {
			select {
			case results <- links:
			default:
				fmt.Printf("results channel full, discarding results for url: %s\n", url)
			}
		}
	}
}

// Crawl completes a full same domain crawl of a given URL
// this means a complete site map will be generated, this includes
// external URLs which are *not* crawled
func Crawl(URL string) []CrawlResult {
	numWorkers := 100
	tasksChanSize := 1000
	resultsChanSize := 10

	fmt.Printf("crawling %s\n", URL)

	visited := StringSet{}

	var crawlResults []CrawlResult

	tasks := make(chan string, tasksChanSize)
	defer close(tasks)

	results := make(chan CrawlResult, resultsChanSize)
	defer close(results)

	crawlsInProgress := 0

	for i := 0; i < numWorkers; i++ {
		go worker(URL, tasks, results)
	}

	visited.Add(URL)
	crawlsInProgress++

	tasks <- URL

	for r := range results {

		crawlsInProgress--

		crawlResults = append(crawlResults, r)

		for link := range r.internalLinks {
			if !visited.Contains(link) {
				visited.Add(link)
				crawlsInProgress++

				select {
				case tasks <- link:
				default:
					fmt.Printf("failed to issue task for URL: %s\n", link)
				}
			}
		}

		for link := range r.externalLinks {
			visited.Add(link)
		}

		fmt.Printf("\033[2K\rcrawled %d links", len(visited))

		if crawlsInProgress == 0 {
			break
		}
	}

	return crawlResults
}
