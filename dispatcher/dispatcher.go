package dispatcher

import (
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"poemcrawler/htmltype"
	"poemcrawler/util"
	"poemcrawler/db"
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

	var poet util.Poet
	var poems []util.Poem

	switch t {
	case "xs":
		c := htmltype.NewXianDaiShi(d.uctx, d.res, d.doc)
		if len(ps) == 4 { // 诗集的情况，一个页面一首诗
			poems = c.GetOnePoemFromCollection()
			poet = c.GetPoet()
		} else {
			poems = c.Base.GetPoems()
			poet = c.GetPoet()
		}
	case "gs":
		c := htmltype.NewGuDianShi(d.uctx, d.res, d.doc)
		poems = c.Base.GetPoems()
		poet = c.Base.GetPoet()

	case "ws":
		c := htmltype.NewGuoJiShi(d.uctx, d.res, d.doc)
		poems = c.Base.GetPoems()
		poet = c.Base.GetPoet()
	}

	poetErr := util.CheckPoet(poet)
	poemErr := util.CheckPoems(poems)
	if poetErr || poemErr {
		ep := util.ErrorPage{Url: d.uctx.URL().String()}
		db.SaveErrorPage(ep)
		util.SaveToFile(fn, d.uctx.URL().String(), poet, poems)
	} else {
		db.SavePoet(poet)
		db.SavePoems(poems)
	}
}
