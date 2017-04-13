package dispatcher

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
	uctx *gocrawl.URLContext
	res  *http.Response
	doc  *goquery.Document
}

func NewDispatcher(uctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) *Dispatcher {
	return &Dispatcher{uctx: uctx, res: res, doc: doc}
}

func (d Dispatcher)Dispatch() {
	ps := strings.Split(d.uctx.URL().Path, "/")
	fn := ps[len(ps)-1] + "l"

	t:=strings.Split(d.doc.Url.Path, "/")[2]
	fmt.Println(t)

	path := strings.Split(d.uctx.URL().Path, "/")
	suffix := path[len(path)-1]

	fmt.Println(suffix)
	if strings.Contains(suffix, "index") {
		return
	}

	switch t {
	case "xs":
		c := htmltype.NewXianDaiShi(d.uctx, d.res, d.doc)
		poems := c.GetPoems()
		util.Save(fn, d.uctx.URL().String(), poems)
	case "gs":
		c := htmltype.NewGuDianShi(d.uctx, d.res, d.doc)
		poems := c.GetPoems()
		util.Save(fn, d.uctx.URL().String(), poems)
	case "ws":
		c := htmltype.NewGuoJiShi(d.uctx, d.res, d.doc)
		poems := c.GetPoems()
		util.Save(fn, d.uctx.URL().String(), poems)
	}
}
