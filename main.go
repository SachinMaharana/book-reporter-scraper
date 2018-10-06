package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/gocolly/colly"
)

const (
	// Site to be Scraped
	rootSiteLink = "https://www.bookreporter.com"

	comingSoonLink = rootSiteLink + "/coming-soon"

	monthsSelector = "#sidebar-last > div:nth-child(1) > div > div > div > div > div > ul li"

	bookInfoSelector = "div[class=book-info]"

	eachBookSelector = ".views-row-unformatted"

	genreSelector = "span.genre"

	nextPageLink = "li.pager-current + li > a[href]"
)

//Book is a struct to hold book information..
type Book struct {
	Title   string   `json:"title"`
	Author  string   `json:"author"`
	Genre   []string `json:"genres"`
	Publish string   `json:"publish"`
}

// books is a collection of Book Type
type books []Book

type callback *colly.HTMLCallback

type data map[string]books

//createCollectors return the colectors. Using Naked Return..
func createCollectors() (monthsCollector, booksCollector *colly.Collector) {
	monthsCollector = colly.NewCollector()
	booksCollector = colly.NewCollector()
	return
}

func main() {
	books := books{}
	// d := make(data)
	monthsCollector, booksCollector := createCollectors()

	monthsCallback := func(e *colly.HTMLElement) {
		// d[e.ChildText("a")] = books{}
		monthLink := e.ChildAttr("a", "href")
		link := rootSiteLink + monthLink
		fmt.Println("Link Found: ", link)
		// visit that month and collect books info
		booksCollector.Visit(link)
	}

	// visit months avaiable
	monthsCollector.OnHTML(monthsSelector, monthsCallback)

	// bookCallback is called each time book selector is found
	bookCallback := func(e *colly.HTMLElement) {
		var genre []string
		e.ForEach(bookInfoSelector, func(_ int, el *colly.HTMLElement) {
			el.ForEach(genreSelector, func(_ int, el *colly.HTMLElement) {
				genre = append(genre, el.Text)
			})
			child := el.DOM.Children()

			book := Book{
				Title:   child.Eq(0).Text(),
				Author:  child.Eq(1).Text(),
				Genre:   genre,
				Publish: child.Eq(3).Text(),
			}
			books = append(books, book)
		})
	}

	// collect book from each page
	booksCollector.OnHTML(eachBookSelector, bookCallback)

	// visit next page to collect book
	booksCollector.OnHTML(nextPageLink, func(nextPage *colly.HTMLElement) {
		log.Println("Next page link found:", nextPage.Text)
		link := rootSiteLink + nextPage.Attr("href")
		nextPage.Request.Visit(link)
	})

	//events
	booksCollector.OnRequest(func(r *colly.Request) {
		log.Println("booksCollector : Visiting", r.URL)
	})

	booksCollector.OnResponse(func(r *colly.Response) {
		log.Println("booksCollector : Visited", r.Request.URL)
	})

	booksCollector.OnError(func(_ *colly.Response, err error) {
		log.Println("booksCollector : Something went wrong:", err)
	})

	monthsCollector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	monthsCollector.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	monthsCollector.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.Request.URL)
	})

	monthsCollector.Visit(comingSoonLink)
	// jsonStr, _ := json.Marshal(books)
	// fmt.Println("Done", string(jsonStr))
	file, err := os.Create("result.json")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	enc.Encode(books)
}
