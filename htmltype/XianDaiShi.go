package htmltype

import (
	//"fmt"
	"net/http"
	"poemcrawler/util"
	"strings"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
)

// 处理现代诗歌的类型
// 页面样例 http://www.shiku.org/shiku/xs/xuzhimo.htm
type XianDaiShi struct {
	Base *ShiKu
	uctx *gocrawl.URLContext
	res  *http.Response
	doc  *goquery.Document
}

func NewXianDaiShi(uctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) *XianDaiShi {
	return &XianDaiShi{
		Base: NewShiKu(uctx, res, doc),
		uctx: uctx,
		res:  res,
		doc:  doc,
	}
}

func (t XianDaiShi) GetPoet() util.Poet {
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

// 标题为h2标签，诗歌内容为h2标签后第二个标签内
// 例子页面 http://www.shiku.org/shiku/xs/mudan.htm        第二个标签为p
// 例子页面 http://www.shiku.org/shiku/xs/xuzhimo.htm      第二个标签为pre
func (t XianDaiShi) GetPoemsH2AndP() (poems []util.Poem) {
	poet := t.GetPoet()
	poems = make([]util.Poem, 0, 0)

	t.doc.Find("body").Find("h2").Each(func(i int, s *goquery.Selection) {
		gbkFullTitle := s.Text()
		fullTitleBytes := []byte(gbkFullTitle)
		fullTitle := strings.TrimSpace(util.GBK2Unicode(fullTitleBytes))
		title := strings.TrimSpace(strings.Split(fullTitle, " ")[0])
		subTitle := fullTitle[len(title):]

		gbkContent := s.Next().Next().Text()
		contentBytes := []byte(gbkContent)
		content := util.GBK2Unicode(contentBytes)

		poem := util.Poem{
			Author:   poet.Name,
			Source:   t.uctx.URL().String(),
			Title:    title,
			Subtitle: subTitle,
			Body:     content,
		}

		poems = append(poems, poem)
	})

	return
}

// 标题为p align="center" 标签， 诗歌内容为标题标签后第一个p标签内
// 例子页面 http://www.shiku.org/shiku/xs/guangweiran.htm
func (t XianDaiShi) GetPoemsPAndP() (poems []util.Poem) {
	poet := t.GetPoet()
	poems = make([]util.Poem, 0, 0)

	t.doc.Find("body").Find("p[align=\"center\"]").Each(func(i int, s *goquery.Selection) {
		gbkFullTitle := s.Text()
		fullTitleBytes := []byte(gbkFullTitle)
		fullTitle := strings.TrimSpace(util.GBK2Unicode(fullTitleBytes))
		title := strings.TrimSpace(strings.Split(fullTitle, " ")[0])
		subTitle := fullTitle[len(title):]

		gbkContent := s.Next().Text()
		contentBytes := []byte(gbkContent)
		content := util.GBK2Unicode(contentBytes)

		// 起始作者信息介绍里，p align="center" 标签内内容为空，忽略
		if title != "" {
			poem := util.Poem{
				Author:   poet.Name,
				Source:   t.uctx.URL().String(),
				Title:    title,
				Subtitle: subTitle,
				Body:     content,
			}

			poems = append(poems, poem)
		}
	})

	return
}

// 获取诗集中的单首诗歌，返回只有一首诗歌的诗歌数组
// 例子页面：http://www.shiku.org/shiku/xs/haizi/100.htm
func (t XianDaiShi) GetOnePoemFromCollection() (poems []util.Poem) {
	gbkTitle := t.doc.Find("body").Find("h1").Text()
	titleBytes := []byte(gbkTitle)
	title := strings.TrimSpace(util.GBK2Unicode(titleBytes))

	gbkAuthor := t.doc.Find("a").Eq(0).Text()
	authorBytes := []byte(gbkAuthor)
	author := strings.TrimSpace(util.GBK2Unicode(authorBytes))
	author = strings.Replace(author, "诗集", "", -1)

	gbkPoemBody := t.doc.Find("pre").Text()
	poemBodyBytes := []byte(gbkPoemBody)
	poemBody := strings.TrimSpace(util.GBK2Unicode(poemBodyBytes))

	poems = make([]util.Poem, 0, 0)
	if title == "诗人简介" {
		return
	}

	poem := util.Poem{
		Author: author,
		Source: t.uctx.URL().String(),
		Title:  title,
		Body:   poemBody,
	}

	poems = append(poems, poem)

	return
}
