package main

import (
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

type Book struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	// Promo  string `json:"promo"`
}

type books []Book

// var r informations

func main() {
	c := colly.NewCollector()

	books := books{}
	// Find and visit all links
	c.OnHTML(".views-row-unformatted", func(e *colly.HTMLElement) {
		e.ForEach("div[class=book-info]", func(_ int, el *colly.HTMLElement) {
			bI := el.DOM.Children()
			book := Book{
				Title: bI.Eq(0).Text(),
				// Author: bI.Eq(1).Text(),
				// Genre  :bI.Eq(2).Text(),
				// Promo: bI.Eq(4).Text(),
			}
			books = append(books, book)
		})

	})

	c.OnHTML("li.pager-current + li > a[href]", func(e *colly.HTMLElement) {
		log.Println("Next page link found:", e.Attr("href"))
		link := "https://www.bookreporter.com" + e.Attr("href")
		e.Request.Visit(link)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		// fmt.Println(r.Body)
	})
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	c.Visit("https://www.bookreporter.com/coming-soon")
	fmt.Println(books)
	// jsonStr, _ := json.Marshal(r)
	// fmt.Println("Done", string(jsonStr))
}
