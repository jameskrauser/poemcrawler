package htmltype

import (
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"poemcrawler/util"
	"strings"
)

type ShiGeKu struct {
	uctx *gocrawl.URLContext
	res  *http.Response
	doc  *goquery.Document
}

func NewShiGeKu(uctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) *ShiGeKu {
	return &ShiGeKu{uctx: uctx, res: res, doc: doc}
}

func (t ShiGeKu) GetPoet() util.Poet {
	gbkAuthor := t.doc.Find("table").Next().Text()
	authorBytes := []byte(gbkAuthor)
	author := util.GBK2Unicode(authorBytes)
	author = strings.Replace(author, "诗选", "", -1)

	gbkIntro := t.doc.Find("body").Find("table").Next().Next().Text()
	introBytes := []byte(gbkIntro)
	intro := util.GBK2Unicode(introBytes)

	poet := util.Poet{
		Name:  author,
		Intro: intro,
	}

	return poet
}

func (t ShiGeKu) GetPoems() (poems []util.Poem) {
	poems = make([]util.Poem, 0)

	poet := t.GetPoet()

	t.doc.Find("body").Find("p[align=\"center\"]").Each(func(i int, s *goquery.Selection) {
		gbkTitle := s.Text()
		titleBytes := []byte(gbkTitle)
		title := util.GBK2Unicode(titleBytes)

		gbkContent := s.Next().Text()
		contentBytes := []byte(gbkContent)
		content := util.GBK2Unicode(contentBytes)
		//content = strings.Replace(content, " ", "", -1)

		poem := util.Poem{
			Author: poet.Name,
			Source: t.uctx.URL().String(),
			Title:  title,
			Body:   content,
		}

		poems = append(poems, poem)
	})

	// 第一个数据是诗人简介，所以只返回后面的数据
	if len(poems) > 0 {
		return poems[1:]
	}

	return poems
}
