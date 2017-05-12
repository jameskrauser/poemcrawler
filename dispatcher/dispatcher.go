package dispatcher

import (
	"net/http"
	"poemcrawler/db"
	"poemcrawler/htmltype"
	"poemcrawler/util"
	"strings"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	//"fmt"
	"fmt"
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
	var isPoemCollection = false
	fmt.Println(ps)
	switch t {
	case "xs":
		c := htmltype.NewXianDaiShi(d.uctx, d.res, d.doc)
		if len(ps) == 4 {
			if ps[2] == "yeshibin" {
				// 诗集的情况，一个页面多首诗
				// 如：http://www.shiku.org/shiku/xs/yeshibin/yeshibin_ztz_1.htm
				poems = c.Base.GetPoems()
				poet = c.Base.GetPoet()
			} else {
				// 诗集的情况，一个页面一首诗
				// 如：http://www.shiku.org/shiku/xs/haizi/154.htm
				poems = c.GetPoemFromOnePageOfCollection()
				poet = c.GetPoetFromOnePageOfCollection()
			}
			isPoemCollection = true
		} else {
			poems = c.Base.GetPoems()
			poet = c.Base.GetPoet()
		}
		//case "gs":
		//	c := htmltype.NewGuDianShi(d.uctx, d.res, d.doc)
		//	poems = c.Base.GetPoems()
		//	poet = c.Base.GetPoet()
		//
		//case "ws":
		//	c := htmltype.NewGuoJiShi(d.uctx, d.res, d.doc)
		//	poems = c.Base.GetPoems()
		//	poet = c.Base.GetPoet()
	}

	err := util.CheckPoet(poet)
	if err != nil {
		ep := util.ErrorPage{Url: d.uctx.URL().String(), Message: err.Error()}
		db.SaveErrorPage(ep)
		util.SaveToFile(fn, d.uctx.URL().String(), poet, poems)
	} else {
		if !isPoemCollection {
			db.SavePoet(poet)
		}
	}

	err = util.CheckPoems(poems)
	if err != nil {
		ep := util.ErrorPage{Url: d.uctx.URL().String(), Message: err.Error()}
		db.SaveErrorPage(ep)
		util.SaveToFile(fn, d.uctx.URL().String(), poet, poems)
	} else {
		db.SavePoems(poems)
	}
}
