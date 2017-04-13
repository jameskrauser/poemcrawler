package main

import (
	"fmt"
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"poemcrawler/htmltype"
	"poemcrawler/util"
	"strings"
)

type Dispatcher struct {
	Uctx *gocrawl.URLContext
	Res *http.Response
	Doc *goquery.Document
}

func NewDispatcher(uctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) *Dispatcher {
	return &Dispatcher{Uctx: uctx, Res: res, Doc: doc}
}

func (d Dispatcher)Dispatch() {
	ps := strings.Split(d.Uctx.URL().Path, "/")
	fn := ps[len(ps)-1] + "l"

	fmt.Println(d.Doc.Url.Fragment)
	fmt.Println(d.Doc.Url.Opaque)
	fmt.Println(d.Doc.Url.RawPath)
	fmt.Println(d.Doc.Url.RawQuery)
	fmt.Println(d.Doc.Url.Scheme)
	fmt.Println(d.Doc.Url.User)
	fmt.Println(d.Doc.Url)

	switch d.Doc.Url.Host {
	case "www.shiku.org":
		//c := htmltype.NewShiKu(ctx, res, doc)
		//poems := c.GetPoems()
		//util.Save(fn, ctx.URL().String(), poems)
	case "www.shigeku.com":
		c := htmltype.NewShiGeKu(d.Uctx, d.Res, d.Doc)
		poems := c.GetPoems()
		util.Save(fn, d.Uctx.URL().String(), poems)
	}
}
