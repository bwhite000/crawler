Go Crawler
==========

A web crawler package that can be used to traverse through a website and scrape
data from each webpage's Document as the crawl progresses concurrently.

~~~go
package main

import (
    "log"

    "github.com/bwhite000/crawler"
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

Installation
------------

In the terminal, please type:

~~~bash
go get -u github.com/bwhite000/crawler
~~~

Scraper Methods
---------------

Create the scraper by providing it with a `goquery` Document pointer during instantiation. These methods can then be called on that Document.

~~~go
// Scraper is a tool to help with scraping data.
type Scraper struct {
    Document *goquery.Document
}
~~~

#### `Exists(selector string) bool`

Checks if the selector matches an Element in the Document.

~~~go
if scraper.Exists("[itemtype='http://schema.org/Product']") {
    // ...
}
~~~

#### `Float(selector string) float64`

Gets the text content from the matched Element, then parses a float from the string.

~~~go
percentage := scraper.Float("#percentage-box")
~~~

#### `GetAttr(selector string, attrName string) string`

Gets the attribute value from the matched Element.

~~~go
title := scraper.GetAttr("meta[property='og:title']", "content")
~~~

#### `Html(selector string) string`

Get the inner HTML value from the matched Element.

~~~go
divHTML := scraper.Html("div.elm-with-text")
~~~

#### `Int(selector string) int`

Gets the text content from the matched Element, then parses an integer from the string.

~~~go
year := scraper.Int("div.year-container")
~~~

#### `Text(selector string) string`

Gets the text from the matched Element.

~~~go
bank := scraper.Text("#bank-title-elm")
~~~

#### `ToFloat(input string) float64`

Parses a float value from the provided string. It is okay to have non-numeric values on either side of the expected float in the string.

~~~go
stockPrice := scraper.ToFloat(scraper.GetAttr("meta[property='og:stock']", "content"))
~~~

Crawler Methods
---------------

~~~go
// Crawler crawls a website until the specified number of webpages have been crawled.
type Crawler struct {
    MaxFetches           int
    RequestSpaceMs       int
    IgnoreQueryParams    bool
    StartURL             string
    OnPageLoadedListener chan<- *PageData
}
~~~

### Properties

#### `MaxFetches int`

The maximum number of pages to crawl before completing and exiting the crawl. Useful for testing purposes when setting the value to a lower amount, or for setting an upper limit to prevent very deep crawls.

#### `RequestSpaceMs int`

The amount of time in milliseconds to wait between page fetches. This can be used to avoid rate throttling and limiting situations from the website that the script is crawling.

#### `IgnoreQueryParams bool`

Specifies if the crawler should differentiate between already crawled pages based on the query parameter strings. The default value is `true` to ignore query string parameters. Some websites may have links to the same webpage multiple times with query string parameters to do not affect the page's contents and that would result in rescraping the same page repeatedly.

#### `StartURL string`

The URL to begin the crawling process at.

#### `OnPageLoadedListener chan<- *PageData`

The channel that fetched page information structs will be pushed to when found. Listen on this channel in a loop to process crawled webpages.

### Methods

#### `Begin() void`

Starts the crawl process using the specified values.

#### `WasIndexed(url string) bool`

Tests if the provided URL parameter has already been indexed by this crawl already.

PageData Structure
------------------

~~~go
// PageData contains the details about the page that has been crawled.
type PageData struct {
    URL      string
    Document *goquery.Document
}
~~~

#### `URL string`

The URL of the webpage that this struct represents.

#### `Document *goquery.Document`

The `goquery` Document that this struct represents the DOM tree for.

Dependencies
------------

* __html__ [https://golang.org/x/net/html] - Creates the parse tree from the fetched HTML string.
* __goquery__ [https://github.com/PuerkitoBio/goquery] - Adds the ability to query Elements from the `html` package's parse tree using CSS selector syntax.
    * Depends on:
        * __cascadia__ [https://github.com/andybalholm/cascadia] - Implements CSS selectors for use with the parse trees produced by the `html` package.
        * __html__ [https://golang.org/x/net/html] - Please see above.
* __Go Standard Library__ - Used for all remaining functionality built into the language.
