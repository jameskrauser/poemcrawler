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
	if strings.Contains(ctx.URL().String(), "http://www.shiku.org/shiku/xs") &&
		strings.Contains(ctx.URL().Path, ".htm") {
		return true
	}
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
	opts.MaxVisits = 1000000

	c := gocrawl.NewCrawlerWithOptions(opts)
	//c.Run("http://www.shiku.org/shiku/xs/bianzhilin.htm")
	//c.Run("http://www.shiku.org/shiku/xs/xuzhimo.htm")
	//c.Run("http://www.shiku.org/shiku/xs/mudan.htm")
	//c.Run("http://www.shiku.org/shiku/xs/guangweiran.htm")
	//c.Run("http://www.shiku.org/shiku/xs/zhengmin.htm")
	//c.Run("http://www.shiku.org/shiku/xs/yeshibin.htm")
	//c.Run("http://www.shiku.org/shiku/xs/shiwei.htm")
	//c.Run("http://www.shiku.org/shiku/ws/ww/homer.htm")

	// 标题不是链接: http://www.shiku.org/shiku/ws/wg/corneille.htm
	//c.Run("http://www.shiku.org/shiku/ws/wg/corneille.htm")
	//c.Run("http://www.shiku.org/shiku/xs/yeshibin.htm")

	//c.Run("http://www.shiku.org/shiku/ws/wg/mallarme.htm")
	//c.Run("http://www.shiku.org/shiku/xs/haizi/154.htm")

	//c.Run("http://www.shiku.org/shiku/ws/wg/tyutchev/000.htm")
	//c.Run("http://www.shiku.org/shiku/ws/wg/index.htm")
	//c.Run("http://www.shiku.org/shiku/xs/hanzuorong.htm")
	//c.Run("http://www.shiku.org/shiku/xs/beidao/160.htm")
	//c.Run("http://www.shiku.org/shiku/xs/guomoruo/guomr08.htm")
	//c.Run("http://www.shiku.org/shiku/xs/yeshibin/yeshibin_ztz.htm")
	//c.Run("http://www.shiku.org/shiku/xs/yeshibin/yeshibin_ztz_1.htm")
	//c.Run("http://www.shiku.org/shiku/xs/yeshibin.htm")
	//c.Run("http://www.shiku.org/shiku/xs/shenyinmo.htm")
	c.Run("http://www.shiku.org/shiku/xs/index.htm")
	//c.Run("http://www.shiku.org/shiku/index.htm")

}
