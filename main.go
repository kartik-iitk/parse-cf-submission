package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	CfjailFolderLocation    = "cfjail"
	CfjailSubmissionsFolder = "submissions"
	CfjailContestsFolder    = "contests"
	CfjailContestPrefix     = "contest-"
	CfjailProblemPrefix     = "problem-"
)

type SubmissionData struct {
	SubmissionID string
	Author       string
	Contest      string
	PrblmIndx    string
	PrblmRev     string
	Language     string
	Verdict      string
	CodePath     string
}

func checkError(err error) {
	// To avoid unnecessary duplication of code.
	if err != nil {
		log.Fatal(err)
	}
}

func loadHTMLFile(path string) []byte {
	// Load the HTML document. No http implementation to enable offline testing!
	file, err := os.ReadFile(path)
	checkError(err)
	return file
}

func createGoqueryDoc(html []byte) *goquery.Document {
	// Generate goquery document from the html file.
	// Important: The HTML file should be properly formatted (and NOT for eg,
	// such that body tag and head tag are on the same line) as parsing would
	// fail without generating an error.
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
	checkError(err)
	return doc
}

func extractUserInfo(doc *goquery.Document, sub *SubmissionData) {
	// Find everything other than the code.

	// How to find selector for find? go to browser, use inspect element
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
	doc.Find(".datatable div table tr").Eq(1).Each(
		// Find the class datatable (you don't need to start from .body)
		// and then find list of div tags, then find amongst them
		// the ones with table in them, in the new list find the tr tags and
		// Each() will iterate over all the rows, but Eq(1) will restrict
		// into the second row (provides index).
		func(index int, item *goquery.Selection) {
			// In the row, iterate over all columns (td tags).
			doc.Find("td").EachWithBreak(
				func(index int, item *goquery.Selection) bool {
					switch index {
					case 0:
						sub.SubmissionID = strings.TrimSpace(item.Text())
					case 1:
						sub.Author = doc.Find(".rated-user").Text()
					case 2:
						// Remove all the spaces and "\n" from the string.
						temp := strings.ReplaceAll(
							strings.ReplaceAll(item.Text(), " ", ""),
							"\n", "")
						tempData := strings.Split(temp, "-")
						sub.Contest = tempData[0][0 : len(tempData[0])-1]
						sub.PrblmIndx = tempData[0][len(tempData[0])-1:]
						sub.PrblmRev = tempData[1]
					case 3:
						sub.Language = strings.TrimSpace(item.Text())
					case 4:
						sub.Verdict = strings.TrimSpace(item.Text())
					default:
						return false // Break out of Each()
					}
					return true
				})
		})
}

func createCodeFolder(sub *SubmissionData) string {
	// Generate folder to store the code.
	dir := filepath.Join(CfjailFolderLocation,
		CfjailSubmissionsFolder,
		CfjailContestsFolder,
		fmt.Sprintf("%s%s", CfjailContestPrefix, sub.Contest),
		fmt.Sprintf("%s%s", CfjailProblemPrefix, sub.PrblmIndx))
	err := os.MkdirAll(dir, os.ModePerm) // Make directories if they don't exist.
	checkError(err)
	return dir
}

func getLanguageExtension(sub *SubmissionData) string {
	// Get language extension.
	if strings.Contains(sub.Language, "C++") {
		return "cpp"
	} else if strings.Contains(sub.Language, "Java") {
		return "java"
	} else if strings.Contains(sub.Language, "Py") {
		return "py"
	} else {
		fmt.Println("Unknown Language! Saving code as a text document.")
		return "txt"
	}
}

func createCodeFile(doc *goquery.Document, sub *SubmissionData) {
	dir := createCodeFolder(sub)
	langext := getLanguageExtension(sub)

	// Update CodePath.
	sub.CodePath = filepath.Join(dir,
		fmt.Sprintf("%s.%s", sub.SubmissionID, langext))

	// Find the code. find() should ideally only return 1 object, but to access
	// it we use index 0.
	code := doc.Find(".prettyprint").Eq(0).Text()

	// Write the code to file with permission 0644 which stands for:
	// 1. The file's owner can read and write (6)
	// 2. Users in the same group as the file's owner can read (first 4)
	// 3. All users can read (second 4)
	os.WriteFile(sub.CodePath, []byte(code), 0644)
}

func Scrape(path string) *SubmissionData {
	var sub SubmissionData
	doc := createGoqueryDoc(loadHTMLFile(path))
	extractUserInfo(doc, &sub)
	createCodeFile(doc, &sub)
	return &sub
}

func main() {
	// FromSlash() helps in maintaining os-independency.
	path := filepath.FromSlash(
		"./html-files/contest_1706_submission_165909579.html")
	submission1 := Scrape(path)
	fmt.Printf("%+v", *submission1) // Prints better than fmt.Println()!
}
