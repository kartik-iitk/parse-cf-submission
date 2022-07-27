package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
)

type Info struct {
	SubmissionID string
	Author       string
	Problem      string
	Language     string
	Verdict      string
	Code         string
}

func scrapeInfo(doc *goquery.Document) *Info {
	// Method to find selector for find: go to browser, use inspect element
	// feature, and hover over required portion of the html.
	// Next, see the required class-type of immediate or close-by parents and
	// try searching from there directly, you need not start from the class of
	// the very first parent! Classes and ids are what you should look out for
	// while scraping webpages.
	// How to test selector? Use the console, to write jQueries, and see if you
	// are getting the correct output. functions like alert, innerText,
	// specifying index simply by [...] after selector (...) all are helpful.
	// The $(..) command can only be used if your website uses jQueries (which
	// in most cases it does) else you have to actually write jQueries!
	var data Info

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
						data.SubmissionID = strings.TrimSpace(value)
					case 1:
						data.Author = doc.Find(".rated-user").Text()
					case 2:
						// Remove all the spaces and "\n" from the string.
						data.Problem = strings.ReplaceAll(
							strings.ReplaceAll(value, " ", ""),
							"\n", "")
					case 3:
						data.Language = strings.TrimSpace(value)
					case 4:
						data.Verdict = strings.TrimSpace(value)
					}
				})
		})

	// Find the code. find should ideally only return 1 object, but to access
	// it we use index 0.
	codeNode := doc.Find(".prettyprint").Eq(0)
	data.Code = codeNode.Text()

	return &data
}

func main() {
	res, err := http.Get("https://codeforces.com/contest/1706/submission/164749937")
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
