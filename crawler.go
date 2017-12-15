package crawler

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

var urlIndex = make(map[string]bool)
var alreadyIndexedUrls = make(map[string]bool)

// Set the config that it is okay to handle webpages with invalid HTTPS certificates.
var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

// Initialize a new HTTP Client.
var client = &http.Client{Transport: tr}

// Fetch fetches the remote resource at the provided url.
func Fetch(respChan chan<- []byte, url string) {
	// Fetch the remote webpage.
	if resp, respErr := client.Get(url); respErr == nil {
		// Close the body buffer and connection when done.
		defer resp.Body.Close()

		// Read the contents of the HTTP response body.
		if body, readBodyErr := ioutil.ReadAll(resp.Body); readBodyErr == nil {
			respChan <- body
		} else {
			fmt.Println("fetch(): ioutil.ReadAll() error", readBodyErr)
			respChan <- []byte("")
		}
	} else {
		fmt.Println("fetch(): client.Get() error:", respErr)
		respChan <- []byte("")
	}
}

// PageData contains the details about the page that was crawled.
type PageData struct {
	URL      string
	Document *goquery.Document
}

// Crawler crawls a webpage until a specified number of webpages has been crawled.
type Crawler struct {
	MaxFetches           int
	Origin               string
	RequestSpaceMs       int
	IgnoreQueryParams    bool
	StartURL             string
	OnPageLoadedListener chan<- *PageData

	crawlDepth  int
	homepageURL *url.URL
}

// Begin starts the crawl process.
func (cr *Crawler) Begin() {
	cr.Start(cr.StartURL)
}

// Start begins the crawl of the provided URL.
func (cr *Crawler) Start(urlParam string) {
	// Check if the crawl depth has exceeded that max allowed amount.
	if cr.crawlDepth >= cr.MaxFetches {
		return
	}

	// Do not index a URL that has already been indexed.
	if _, exists := alreadyIndexedUrls[urlParam]; exists {
		return
	}

	// Check if this is the first crawl
	if cr.crawlDepth <= 0 {
		if homepageURL, urlParseErr := url.Parse(urlParam); urlParseErr == nil {
			// Assign the value of the parsed URL.
			cr.homepageURL = homepageURL

			// Create the origin for resolving in relative paths.
			cr.Origin = (cr.homepageURL.Scheme + "://" + cr.homepageURL.Host)
		} else {
			log.Fatal(urlParseErr)
		}
	}

	// Increment the crawl depth
	cr.crawlDepth++

	fmt.Println("("+strconv.Itoa(cr.crawlDepth)+" of "+strconv.Itoa(cr.MaxFetches)+") Fetching for:", urlParam)

	// Pause and sleep, if set to.
	if cr.RequestSpaceMs > 0 {
		time.Sleep(time.Duration(cr.RequestSpaceMs) * time.Millisecond)
	}

	// Create a channel.
	pageContentsBufChan := make(chan []byte)

	// Fetch the remote webpage.
	go Fetch(pageContentsBufChan, urlParam)

	// Wait for the goroutine to complete the webpage fetch.
	pageContentsBuf := <-pageContentsBufChan
	close(pageContentsBufChan) // Close the channel.

	if len(pageContentsBuf) > 0 {
		// Parse the HTML into [Node]s.
		if node, htmlParseErr := html.Parse(bytes.NewBuffer(pageContentsBuf)); htmlParseErr == nil {
			// Create a new goquery [Document] from the parsed [Nodes].
			doc := goquery.NewDocumentFromNode(node)

			// Send the parsed document through the event listener.
			cr.OnPageLoadedListener <- &PageData{URL: urlParam, Document: doc}

			// Set the default canonicalURL value.
			canonicalURL := urlParam

			// Attempt to get the canonicalURL value from the webpage; check if the Element is in the DOM.
			if canonicalElm := doc.Find("link[rel='canonical']"); canonicalElm.Size() > 0 {
				// Check if the matched Element has an "href" attribute.
				if attrVal, hrefExists := canonicalElm.Attr("href"); hrefExists {
					canonicalURL = (cr.Origin + attrVal)
				}
			}

			// Add this base URL to the list of already visited pages.
			alreadyIndexedUrls[urlParam] = true

			// Add the parsed webpage canonical URL to the list of already indexed URLs.
			alreadyIndexedUrls[canonicalURL] = true

			// Query for all HTML Anchor Elements on this webpage to detect more links to crawl.
			doc.Find("a").Each(func(i int, s *goquery.Selection) {
				if hrefVal, hrefExists := s.Attr("href"); hrefExists {
					if strings.HasPrefix(hrefVal, "/") {
						hrefVal = (cr.Origin + hrefVal)
					} else if strings.HasPrefix(hrefVal, "http") == false {
						return
					}

					if hrefURL, err := url.Parse(hrefVal); err == nil {
						if cr.homepageURL.Host == hrefURL.Host {
							// Check if the query parameters should be removed from the URL.
							if cr.IgnoreQueryParams {
								// Empty the query String parameters.
								hrefURL.RawQuery = ""

								urlIndex[hrefURL.String()] = true
							} else {
								urlIndex[hrefVal] = true
							}
						}
					} else {
						log.Fatal(err)
					}
				} else {
					return
				}
			})

			// Clear the memory used by the Nodes and Document.
			node = nil
			doc = nil

			// Check if the crawler depth has reached the max depth amount yet.
			if cr.crawlDepth < cr.MaxFetches {
				for k := range urlIndex {
					cr.Start(k)
				}
			} else { // Complete the crawl recursion and make the final log.
				fmt.Println("Number of urls reachable for indexing:", len(urlIndex))
			}
		} else {
			log.Println(htmlParseErr)
		}
	}
}

// WasIndexed tests if the provided URL has already been indexed before.
func (cr *Crawler) WasIndexed(url string) bool {
	if _, existsAlready := alreadyIndexedUrls[url]; existsAlready == false {
		return true
	}

	return false
}
