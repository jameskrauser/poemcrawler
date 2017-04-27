package htmltype

import (
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"net/http"
)

// 处理古代诗歌的类型
// 页面样例 http://www.shiku.org/shiku/gs/beichao.htm
type GuDianShi struct {
	Base *ShiKu
	uctx *gocrawl.URLContext
	res  *http.Response
	doc  *goquery.Document
}

func NewGuDianShi(uctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) *GuDianShi {
	return &GuDianShi{
		Base: NewShiKu(uctx, res, doc),
		uctx: uctx,
		res:  res,
		doc:  doc,
	}
}
