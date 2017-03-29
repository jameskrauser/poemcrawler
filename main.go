package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"

	"poemcrawler/htmltype"
	"poemcrawler/util"
	"strings"
)

type Ext struct {
	*gocrawl.DefaultExtender
}

func (e *Ext) Visit(ctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) (interface{}, bool) {
	fmt.Printf("Visit: %s\n", ctx.URL())
	ps := strings.Split(ctx.URL().Path, "/")
	fn := ps[len(ps)-1] + "l"

	ctx.URL().String()
	switch doc.Url.Host {
	case "www.shiku.org":
		c := htmltype.NewShiKu(ctx, res, doc)
		poems := c.GetPoems()
		util.Save(fn, poems)
	case "www.shigeku.com":
		c := htmltype.NewShiGeKu(ctx, res, doc)
		poems := c.GetPoems()
		util.Save(fn, poems)
	}

	return nil, true
}

func (e *Ext) Filter(ctx *gocrawl.URLContext, isVisited bool) bool {
	if isVisited {
		return false
	}
	//if ctx.URL().Host == "github.com" || ctx.URL().Host == "golang.org" || ctx.URL().Host == "0value.com" {
	//	return true
	//}
	if ctx.URL().Host == "www.shiku.org" || ctx.URL().Host == "www.shigeku.com" {
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
	opts.MaxVisits = 1//9999999999

	c := gocrawl.NewCrawlerWithOptions(opts)
	//c.Run("http://www.shiku.org/shiku/xs/xuzhimo.htm")
	c.Run("http://www.shigeku.com/xlib/xd/sgdq/ajian.htm")
	//c.Run("http://www.shigeku.com/xlib/xd/sgdq/caitianxin.htm")

	//c.Run("http://www.shiku.org/shiku/xs/index.htm")

}
