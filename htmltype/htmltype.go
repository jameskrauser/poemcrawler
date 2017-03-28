package htmltype

import (
	"poemcrawler/util"

	"github.com/PuerkitoBio/goquery"
	"strings"
)

type ShiKu struct {
	Source string // 来源
}

func NewShiKu(url string) *ShiKu {
	return &ShiKu{Source: url}
}

func (t ShiKu) Parse(doc *goquery.Document) (poems []util.Poem) {
	poems = make([]util.Poem, 0)

	gbkAuthor := doc.Find("body").Find("h1").Text()
	authorBytes := []byte(gbkAuthor)
	author := util.GBK2Unicode(authorBytes)
	author = strings.Replace(author, "诗选", "", -1)

	doc.Find("body").Find("h2").Each(func(i int, s *goquery.Selection) {
		gbkTitle := s.Text()
		titleBytes := []byte(gbkTitle)
		title := util.GBK2Unicode(titleBytes)

		gbkContent := s.Next().Next().Text()
		contentBytes := []byte(gbkContent)
		content := util.GBK2Unicode(contentBytes)

		poem := util.Poem{
			Author: author,
			Intro:  "",
			Source: t.Source,
			Title:  title,
			Body:   content,
		}

		poems = append(poems, poem)
	})

	return
}

type ShiGeKu struct {
	Source string // 来源
}

func NewShiGeKu(url string) *ShiGeKu {
	return &ShiGeKu{Source: url}
}

func (t ShiGeKu) Parse(doc *goquery.Document) (poems []util.Poem) {
	poems = make([]util.Poem, 0)

	gbkAuthor := doc.Find("body").Find("table").Next().Text()
	authorBytes := []byte(gbkAuthor)
	author := util.GBK2Unicode(authorBytes)
	author = strings.Replace(author, "诗选", "", -1)

	doc.Find("body").Find("p[align=\"center\"]").Each(func(i int, s *goquery.Selection) {
		gbkTitle := s.Text()
		var titleBytes = []byte(gbkTitle)
		title := util.GBK2Unicode(titleBytes)

		gbkContent := s.Next().Text()
		contentBytes := []byte(gbkContent)
		content := util.GBK2Unicode(contentBytes)
		content = strings.Replace(content, " ", "", -1)

		poem := util.Poem{
			Author: author,
			Intro:  "",
			Source: t.Source,
			Title:  title,
			Body:   content,
		}

		poems = append(poems, poem)
	})

	return
}
