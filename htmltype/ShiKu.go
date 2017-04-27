package htmltype

import (
	"fmt"
	"net/http"
	"poemcrawler/util"
	"strings"

	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
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

// 标题为h2标签，诗歌内容为h2标签后第二个标签内
// 例子页面 http://www.shiku.org/shiku/xs/mudan.htm        第二个标签为p
// 例子页面 http://www.shiku.org/shiku/xs/xuzhimo.htm      第二个标签为pre
func (t ShiKu) GetPoemsH2AndP() (poems []util.Poem) {
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
func (t ShiKu) GetPoemsPAndP() (poems []util.Poem) {
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

// 通过在标题前加分隔符后解析纯文本的方式解析页面
func (t ShiKu) GetPoems() (poems []util.Poem) {
	titles := make([]string, 0, 0)
	poet := t.GetPoet()
	poems = make([]util.Poem, 0, 0)
	sep := "=============="

	has999999 := false
	t.doc.Find("body").Find("a").Each(func(i int, s *goquery.Selection) {
		href, existHref := s.Attr("href")
		if existHref {
			if strings.HasPrefix(href, "#") {
				gbkTitle := s.Text()
				titleBytes := []byte(gbkTitle)
				title := strings.TrimSpace(util.GBK2Unicode(titleBytes))
				titles = append(titles, title)
			}
		}

		name, existName := s.Attr("name")

		if existName {
			if strings.Contains("0123456789", name[0:1]) {
				t.doc.Find("body").Find("a[name=\"" + name + "\"]").AppendHtml(sep)
			}
		}

		if name == "999999" {
			has999999 = true
		}
	})

	gbkText := t.doc.Text()
	TextBytes := []byte(gbkText)
	text := strings.TrimSpace(util.GBK2Unicode(TextBytes))
	//中国诗歌库 中华诗库 中国诗典 中国诗人 中国诗坛 首页
	text = strings.Replace(text, "中国诗歌库", "", -1)
	text = strings.Replace(text, "中华诗库", "", -1)
	text = strings.Replace(text, "中国诗典", "", -1)
	text = strings.Replace(text, "中国诗人", "", -1)
	text = strings.Replace(text, "中国诗坛", "", -1)
	text = strings.Replace(text, "首页", "", -1)
	text = strings.TrimSpace(text)

	textArr := strings.Split(text, sep)
	if strings.Contains(text, sep) {
		content := textArr[1:]
		if has999999 {
			content = textArr[1:len(textArr)-1]
		}

		fmt.Println("解析到的诗歌体数量为：", len(content))
		fmt.Println("解析到的诗歌标题数量为：", len(titles))
		count := len(content)
		//if has999999 {
		//	count = len(content) - 1
		//}

		for i := 0; i < count; i++ {
			whole := strings.TrimSpace(sep + content[i])
			title := titles[i]
			body := strings.Replace(whole, sep+title, "", -1)
			poem := util.Poem{
				Author: poet.Name,
				Source: t.uctx.URL().String(),
				Title:  title,
				Body:   body,
			}

			poems = append(poems, poem)
		}
	}

	return
}
