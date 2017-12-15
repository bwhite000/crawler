package crawler

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// Scraper is a tool to help with scraping data.
type Scraper struct {
	Document *goquery.Document
}

var floatRegExp = regexp.MustCompile(`\d+(\.\d{2})?`)
var intRegExp = regexp.MustCompile(`\d+`)

// Float gets the [Text] value of the Element matching the provided selector, then
// parses the floating point value from the string.
func (sc *Scraper) Float(selector string) float64 {
	if elm := sc.Document.Find(selector); elm.Size() > 0 {
		text := elm.Text()
		// Remove all commas in the float string.
		text = strings.Replace(text, ",", "", -1)

		if floatRegExp.MatchString(text) {
			foundVal := floatRegExp.FindString(text)

			if floatVal, err := strconv.ParseFloat(foundVal, 64); err == nil {
				return floatVal
			}
		}
	}

	return 0.0
}

// ToFloat converts the provided string to a floating point value.
func (sc *Scraper) ToFloat(input string) float64 {
	// Remove all commas in the float string.
	input = strings.Replace(input, ",", "", -1)

	if floatRegExp.MatchString(input) {
		foundVal := floatRegExp.FindString(input)

		if floatVal, err := strconv.ParseFloat(foundVal, 64); err == nil {
			return floatVal
		}
	}

	return 0.0
}

// Int gets the [Text] value of the Element matching the provided selector, then
// parses the integer value from the string.
func (sc *Scraper) Int(selector string) int {
	if elm := sc.Document.Find(selector); elm.Size() > 0 {
		text := elm.Text()
		// Remove all commas in the float string.
		text = strings.Replace(text, ",", "", -1)

		if intRegExp.MatchString(text) {
			foundVal := intRegExp.FindString(text)

			if intVal, err := strconv.Atoi(foundVal); err == nil {
				return intVal
			}
		}
	}

	return 0
}

// Text gets the innerText of the specified Element.
func (sc *Scraper) Text(selector string) string {
	if elm := sc.Document.Find(selector); elm.Size() > 0 {
		return strings.TrimSpace(elm.Text())
	}

	return ""
}

// Html get the HTML of the specified Element.
func (sc *Scraper) Html(selector string) string {
	if elm := sc.Document.Find(selector); elm.Size() > 0 {
		if html, err := elm.Html(); err == nil {
			return html
		}
	}

	return ""
}

// Exists checks to see if the provided selector matches an Element that exists in the DOM.
func (sc *Scraper) Exists(selector string) bool {
	if elm := sc.Document.Find(selector); elm.Size() > 0 {
		return true
	}

	return false
}

// GetAttr gets the attribute value for the specified Element.
func (sc *Scraper) GetAttr(selector string, attrName string) string {
	if elm := sc.Document.Find(selector); elm.Size() > 0 {
		if attrVal, exists := elm.Attr(attrName); exists {
			return attrVal
		}
	}

	return ""
}
