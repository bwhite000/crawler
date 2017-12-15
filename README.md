Go Crawler
==========

A web crawler package that can be used to traverse through a website and scrape
data from each webpage's Document as the crawl progresses concurrently.

~~~go
package main

import (
    "log"

    "github.com/bwhite000/go-crawler"
)

var onPageLoadedChan = make(chan *crawler.PageData)

func onPageLoaded() {
    for {
        // Listen for the next incoming webpage data from the crawler.
        data := <-onPageLoadedChan
        doc := data.Document

        // Implement the included scraping helper methods.
        scraper := &crawler.Scraper{Document: doc}

        // Process the webpage's Document to scrape useful data.
        log.Println("Title: ", scraper.GetAttr("meta[property='og:title']", "content"))
    }
}

func init() {
    go onPageLoaded()
}

func main() {
    log.Println("Beginning crawl...")

    // Initialize the crawler.
    crawlerObj := &crawler.Crawler{
        StartURL:             "https://example.com/photos/switzerland",
        OnPageLoadedListener: onPageLoadedChan,
    }

    // Begin the crawl.
    crawlerObj.Begin()
}
~~~
