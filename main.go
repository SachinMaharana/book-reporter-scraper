package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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

	path = "result.json"

	layout = "Jan 2, 2006 at 3:04pm (MST)"
)

type data struct {
	LastUpdated string `json:"lastUpdated"`
	Now         string `json:"now"`
	BooksData   books  `json:"booksData"`
}

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

// type data map[string]books

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
		if i, err := strconv.Atoi(nextPage.Text); err == nil {
			fmt.Println("Collected From Page: ", i-1)
		}
		log.Println("Next page link found:", nextPage.Text)
		link := rootSiteLink + nextPage.Attr("href")
		nextPage.Request.Visit(link)
	})

	//events
	booksCollector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	booksCollector.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.Request.URL)
	})

	booksCollector.OnError(func(_ *colly.Response, err error) {
		log.Println("booksCollector : Something went wrong:", err)
		os.Exit(0)
	})

	monthsCollector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	monthsCollector.OnError(func(_ *colly.Response, err error) {
		log.Println("monthsCollector: Something went wrong:", err)
		os.Exit(0)
	})
	monthsCollector.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.Request.URL)
	})

	monthsCollector.Visit(comingSoonLink)
	monthsCollector.Wait()
	createDump(books)

}

func createAndEncode(d data) {
	f, err := os.Create(path)
	if err != nil {
		fmt.Println("Error Creating File")
		os.Exit(0)
	}
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	encoder.Encode(d)
	defer f.Close()
}

func updateData(b books) data {
	file, _ := os.Open(path)
	defer file.Close()
	decoder := json.NewDecoder(file)
	data := data{}
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	now := time.Now().UTC().Format(layout)
	data.LastUpdated, data.Now = data.Now, now
	data.BooksData = b
	return data
}

func createDump(b books) {
	//check if file exist
	_, err := os.Stat(path)

	if !os.IsNotExist(err) {
		fmt.Println("Exist")
		updatedData := updateData(b)
		createAndEncode(updatedData)
	} else {
		// file not exist
		newData := data{
			LastUpdated: time.Now().UTC().Format(layout),
			Now:         time.Now().UTC().Format(layout),
			BooksData:   b,
		}
		fmt.Println("Done Creating file", path)
		createAndEncode(newData)
	}

}
