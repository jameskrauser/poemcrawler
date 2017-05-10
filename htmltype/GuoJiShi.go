package htmltype

import (
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

// 处理国际诗歌的类型
// 首页 http://www.shiku.org/shiku/ws/wg/index.htm
type GuoJiShi struct {
	Base *ShiKu
	uctx *gocrawl.URLContext
	res  *http.Response
	doc  *goquery.Document
}

func NewGuoJiShi(uctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) *GuoJiShi {
	return &GuoJiShi{
		Base: NewShiKu(uctx, res, doc),
		uctx: uctx,
		res:  res,
		doc:  doc,
	}
}
