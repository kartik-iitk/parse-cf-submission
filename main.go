package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

type info struct {
	submissionID string
	author       string
	problem      string
	language     string
	verdict      string
	code         string
}

func scrapeInfo(doc *goquery.Document) *info {
	// Method to find selector for find: go to browser, use inspect element
	// feature, and hover over required portion of the html.
	// Next, see the required class-type of immediate or close-by parents and
	// try searching from there directly, you need not start from the class of
	// the very first parent! Classes and ids are what you should look out for
	// while scraping webpages.
	// How to test selector? Use the console, to write jQueries, and see if you
	// are getting the correct output. functions like alert, innerText,
	// specifying index simply by [...] after selector (...) all are helpful.
	// The $(..) command can only be used if
	// your website uses jQueries (which in most cases it does) else you have
	// to actually write jQueries!
	var data info

	// Find everything other than the code.
	doc.Find(".datatable div table tr").Eq(1).Each(
		// Find the class datatable (you don't need to start from .body)
		// and then find list of div tags, then find amongst them
		// the ones with table in them, in the new list find the tr tags and
		// Each() will iterate over all the rows, but Eq(1) will restrict
		// into the second row (provides index).
		func(index int, item *goquery.Selection) {
			// In the row, iterate over all columns (td tags).
			doc.Find("td").Each(
				func(index int, item *goquery.Selection) {
					value := item.Text()
					switch index {
					case 0:
						data.submissionID = strings.TrimSpace(value)
					case 1:
						auth := strings.Split(strings.TrimSpace(value),
							"\n") // Splits the two words.
						data.author = auth[len(auth)-1] // Get last element.
					case 2:
						// Remove all the spaces and "\n" from the string.
						data.problem = strings.TrimSpace(
							strings.ReplaceAll(
								strings.ReplaceAll(value, " ", ""),
								"\n", ""))
					case 3:
						data.language = strings.TrimSpace(value)
					case 4:
						data.verdict = strings.TrimSpace(value)
					}
				})
		})

	// Find the code. find should ideally only return 1 object, but to access
	// it we use index 0.
	codeNode := doc.Find(".prettyprint").Eq(0)
	data.code = codeNode.Text()

	return &data
}

func main() {
	res, err := http.Get("https://codeforces.com/contest/1696/submission/161733082")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close() // deferred closing statement.
	if res.StatusCode > 400 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document.
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	data := scrapeInfo(doc)
	fmt.Printf("%+v\n", *data)
}
