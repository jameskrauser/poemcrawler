package dispatcher

import (
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

func (d Dispatcher) Dispatch() {
	p := strings.TrimLeft(d.doc.Url.Path, "/")
	ps := strings.Split(p, "/")
	fn := ps[len(ps)-1] + "l"

	t := ps[1]

	suffix := ps[len(ps)-1]

	if strings.Contains(suffix, "index") {
		return
	}

	switch t {
	case "xs":
		c := htmltype.NewXianDaiShi(d.uctx, d.res, d.doc)
		if len(ps) == 4 { // 诗集的情况，一个页面一首诗
			poems := c.GetOnePoemFromCollection()
			util.Save(fn, d.uctx.URL().String(), poems)
		} else {
			poems := c.Base.GetPoems()
			util.Save(fn, d.uctx.URL().String(), poems)
		}
	case "gs":
	//c := htmltype.NewGuDianShi(d.uctx, d.res, d.doc)
	//poems := c.Base.GetPoems()
	//util.Save(fn, d.uctx.URL().String(), poems)
	case "ws":
		c := htmltype.NewGuoJiShi(d.uctx, d.res, d.doc)
		poems := c.Base.GetPoems()
		util.Save(fn, d.uctx.URL().String(), poems)
	}
}
