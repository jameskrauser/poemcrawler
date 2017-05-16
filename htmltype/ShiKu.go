package htmltype

import (
	"net/http"
	"poemcrawler/util"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

type ShiKu struct {
	uctx          *gocrawl.URLContext
	res           *http.Response
	doc           *goquery.Document
	HasParseError bool
	Poet          util.Poet
}

const (
	ALinkTitleSep = "++++++++++++++++++++++++"
	BodyTitleSep  = "========================"
	BodyTitleSep1 = "111111111111111111111111"
)

func NewShiKu(uctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) *ShiKu {
	return &ShiKu{uctx: uctx, res: res, doc: doc}
}

// 标题为h2标签，诗歌内容为h2标签后第二个标签内
// 例子页面 http://www.shiku.org/shiku/xs/mudan.htm        第二个标签为p
// 例子页面 http://www.shiku.org/shiku/xs/xuzhimo.htm      第二个标签为pre
//func (t ShiKu) GetPoemsH2AndP() (poems []util.Poem) {
//	poet := t.GetPoet()
//	poems = make([]util.Poem, 0, 0)
//
//	t.doc.Find("body").Find("h2").Each(func(i int, s *goquery.Selection) {
//		gbkFullTitle := s.Text()
//		fullTitleBytes := []byte(gbkFullTitle)
//		fullTitle := strings.TrimSpace(util.GBK2Unicode(fullTitleBytes))
//		title := strings.TrimSpace(strings.Split(fullTitle, " ")[0])
//		subTitle := fullTitle[len(title):]
//
//		gbkContent := s.Next().Next().Text()
//		contentBytes := []byte(gbkContent)
//		content := util.GBK2Unicode(contentBytes)
//
//		poem := util.Poem{
//			Author:   poet.Name,
//			Source:   t.uctx.URL().String(),
//			Title:    title,
//			Subtitle: subTitle,
//			Body:     content,
//		}
//
//		poems = append(poems, poem)
//	})
//
//	return
//}

// 标题为p align="center" 标签， 诗歌内容为标题标签后第一个p标签内
// 例子页面 http://www.shiku.org/shiku/xs/guangweiran.htm
//func (t ShiKu) GetPoemsPAndP() (poems []util.Poem) {
//	poet := t.GetPoet()
//	poems = make([]util.Poem, 0, 0)
//
//	t.doc.Find("body").Find("p[align=\"center\"]").Each(func(i int, s *goquery.Selection) {
//		gbkFullTitle := s.Text()
//		fullTitleBytes := []byte(gbkFullTitle)
//		fullTitle := strings.TrimSpace(util.GBK2Unicode(fullTitleBytes))
//		title := strings.TrimSpace(strings.Split(fullTitle, " ")[0])
//		subTitle := fullTitle[len(title):]
//
//		gbkContent := s.Next().Text()
//		contentBytes := []byte(gbkContent)
//		content := util.GBK2Unicode(contentBytes)
//
//		// 起始作者信息介绍里，p align="center" 标签内内容为空，忽略
//		if title != "" {
//			poem := util.Poem{
//				Author:   poet.Name,
//				Source:   t.uctx.URL().String(),
//				Title:    title,
//				Subtitle: subTitle,
//				Body:     content,
//			}
//
//			poems = append(poems, poem)
//		}
//	})
//
//	return
//}
