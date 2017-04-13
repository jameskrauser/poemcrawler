package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"strings"

)

type Ext struct {
	*gocrawl.DefaultExtender
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	fmt.Printf("Visit: %s\n", ctx.URL())

	//d := main.NewDispatcher(ctx, res, doc)
	//d.Dispatch()



	fmt.Println(doc.Url.Fragment)
	fmt.Println(doc.Url.Opaque)
	fmt.Println(doc.Url.RawPath)
	fmt.Println(doc.Url.RawQuery)
	fmt.Println(doc.Url.Scheme)
	fmt.Println(doc.Url.User)
	fmt.Println(doc.Url)


	return nil, true
}

func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}

	if ctx.URL().Host == "www.shiku.org" && (strings.Contains(ctx.URL().Path, ".htm") ||
		strings.Contains(ctx.URL().Path, ".html")) {
		return true
	}

	return false
}

func main() {
	ext := &Ext{&gocrawl.DefaultExtender{}}
	// Set custom options
	opts := gocrawl.NewOptions(ext)
	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = false
	opts.MaxVisits = 9999999999

	c := gocrawl.NewCrawlerWithOptions(opts)
	c.Run("http://www.shiku.org/shiku/index.htm")

}
