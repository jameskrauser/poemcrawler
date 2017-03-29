package htmltype

import (
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"poemcrawler/util"
	"strings"
)

type ShiKu struct {
	uctx *gocrawl.URLContext
	res  *http.Response
	doc  *goquery.Document
}

func NewShiKu(uctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) *ShiKu {
	return &ShiKu{uctx: uctx, res: res, doc: doc}
}

func (t ShiKu) GetPoet() util.Poet {
	gbkAuthor := t.doc.Find("body").Find("h1").Text()
	authorBytes := []byte(gbkAuthor)
	author := util.GBK2Unicode(authorBytes)
	author = strings.Replace(author, "诗选", "", -1)

	gbkIntro := t.doc.Find("body").Find("h1").Next().Next().Text()
	introBytes := []byte(gbkIntro)
	intro := util.GBK2Unicode(introBytes)

	poet := util.Poet{
		Name:  author,
		Intro: intro,
	}

	return poet
}

func (t ShiKu) GetPoems() (poems []util.Poem) {
	poems = make([]util.Poem, 0)

	poet := t.GetPoet()

	t.doc.Find("body").Find("h2").Each(func(i int, s *goquery.Selection) {
		gbkTitle := s.Text()
		titleBytes := []byte(gbkTitle)
		title := util.GBK2Unicode(titleBytes)

		gbkContent := s.Next().Next().Text()
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

	return
}
