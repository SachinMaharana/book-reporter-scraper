package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gocolly/colly"
)

//Book is a struct to hold book information..
type Book struct {
	Title   string   `json:"title"`
	Author  string   `json:"author"`
	Genre   []string `json:"genres"`
	Publish string   `json:"publish"`
	// Promo  string `json:"promo"`
}

type books []Book

func main() {
	c := colly.NewCollector()

	books := books{}
	c.OnHTML(".views-row-unformatted", func(e *colly.HTMLElement) {
		var genre []string
		e.ForEach("div[class=book-info]", func(_ int, el *colly.HTMLElement) {
			el.ForEach("span.genre", func(_ int, el *colly.HTMLElement) {
				genre = append(genre, el.Text)
			})
			bI := el.DOM.Children()
			title := bI.Eq(0).Text()
			author := bI.Eq(1).Text()
			publish := bI.Eq(3).Text()

			book := Book{
				Title:   title,
				Author:  author,
				Genre:   genre,
				Publish: publish,
			}
			books = append(books, book)
		})
	})

	c.OnHTML("li.pager-current + li > a[href]", func(e *colly.HTMLElement) {
		log.Println("Next page link found:", e.Text)
		link := "https://www.bookreporter.com" + e.Attr("href")
		e.Request.Visit(link)
	})

	c.OnScraped(func(_ *colly.Response) {
		log.Println("Scraping Done")
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})
	c.OnResponse(func(r *colly.Response) {
		log.Println("Visited", r.Request.URL)
	})

	c.Visit("https://www.bookreporter.com/coming-soon")
	jsonStr, _ := json.Marshal(books)
	fmt.Println("Done", string(jsonStr))
}
