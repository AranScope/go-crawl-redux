# go-crawl-redux
🕷️ A same-domain concurrent web crawler written in Go 🕷️

## Usage
```go
go run .
```
Crawls the URL supplied in crawler.go, with the default channel sizes specified in the same file. Prints results to stdout.

## Test
```go
go test
```

## Caveats

We make several caveats and assumptions made in the implementation of the crawler, these are listed for convenience below

| Caveat | Explanation |
| ---------- | ----------- |
| External domains will not be crawled | As we aim for completeness i.e. crawling all pages for a given domain, it is not practical to crawl external domains
| Very large sites **may** not be completely crawled | In the case of very large sites, with large numbers of links, to avoid deadlocks and excessive memory usage we discard URLs. This happens when the results and tasks channels is full, and so we can scale this by increasing their respective sizes
| Does not respect robots.txt | As this is a test exercise rather than production-ready code, we do not respect robots.txt, and so this should **only be used on pre-authorized sites** |
| There is not extensive test coverage | As current there is only one full system integration test, in a production system there would be significantly higher functionality coverage via. additional unit and integration tests.
