package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"

	"poemcrawler/dispatcher"
	"strings"
)

type Ext struct {
	*gocrawl.DefaultExtender
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	fmt.Printf("Visit: %s\n", ctx.URL())

	d := dispatcher.NewDispatcher(ctx, res, doc)
	d.Dispatch()

	return nil, true
}

func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}

	fmt.Println(ctx.URL().String())
	if strings.Contains(ctx.URL().String(), "http://www.shiku.org/shiku/ws") &&
		strings.Contains(ctx.URL().Path, ".htm") {
		return true
	}

	//if strings.Contains(ctx.URL().String(), "http://www.shiku.org/shiku/xs") &&
	//	strings.Contains(ctx.URL().Path, ".htm") {
	//	return true
	//}
	//if ctx.URL().Host == "www.shiku.org" && (strings.Contains(ctx.URL().Path, ".htm") ||
	//	strings.Contains(ctx.URL().Path, ".html")) {
	//	return true
	//}

	return false
}

func main() {
	ext := &Ext{&gocrawl.DefaultExtender{}}
	// Set custom options
	opts := gocrawl.NewOptions(ext)
	opts.CrawlDelay = 1 * time.Second
	opts.LogFlags = gocrawl.LogError
	opts.SameHostOnly = false
	opts.MaxVisits = 10000

	c := gocrawl.NewCrawlerWithOptions(opts)

	//c.Run("http://www.shiku.org/shiku/xs/index.htm")
	c.Run("http://www.shiku.org/shiku/ws/ww/catullus.htm")
	c.Run("http://www.shiku.org/shiku/ws/wg/fairburn.htm")
	//c.Run("http://www.shiku.org/shiku/ws/index.htm")
	//c.Run("http://www.shiku.org/shiku/ws/wg/dante/index.htm")
	//c.Run("http://www.shiku.org/shiku/ws/wg/corneille.htm")
	//c.Run("http://www.shiku.org/shiku/index.htm")

}
