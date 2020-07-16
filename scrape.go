package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gocolly/colly"
)

type comment struct {
	Author string
	Time   string
	Text   string
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	// Instantiate default collector
	c := colly.NewCollector()

	// Find the latest page of the comic (and not other posts)
	c.OnHTML(".postonpage-1 .comment-link a", func(e *colly.HTMLElement) {
		// Navigate to page
		e.Request.Visit(e.Attr("href"))
	})

	// On every top level element
	c.OnHTML(".comment.depth-1 .comment-content", func(e *colly.HTMLElement) {
		// print comment text
		comm := comment{e.ChildText(".comment-author cite"), e.ChildText(".comment-time"), e.ChildText(".comment-text p")}

		commJson, err := json.MarshalIndent(comm, "", "  ")
		checkErr(err)
		fmt.Println(string(commJson))

		writeJson(commJson, os.O_APPEND)

		fmt.Println("Author: ", comm.Author)
		fmt.Println("Time: ", comm.Time)
		fmt.Println("Comment: ", comm.Text)
		fmt.Println()
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Print finished after everything has been scraped
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})

	// Start scraping on https://killsixbilliondemons.com
	c.Visit("https://killsixbilliondemons.com/")
}

func writeJson(data []byte, flag int) {
	f, err := os.OpenFile("output.json", flag, 0644)
	checkErr(err)
	defer f.Close()

	n, err := f.Write(data)
	checkErr(err)
	fmt.Printf("Wrote %d bytes to %s\n", n, f.Name())
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
