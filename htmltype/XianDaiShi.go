package htmltype

import (
	//"fmt"
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"poemcrawler/util"
	"strings"
	//"strconv"
	//"fmt"
	"strconv"
	"fmt"
)

// 处理现代诗歌的类型
// 页面样例 http://www.shiku.org/shiku/xs/xuzhimo.htm
type XianDaiShi struct {
	uctx *gocrawl.URLContext
	res  *http.Response
	doc  *goquery.Document
}

func NewXianDaiShi(uctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) *XianDaiShi {
	return &XianDaiShi{uctx: uctx, res: res, doc: doc}
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

func (t XianDaiShi) initPoems() (poems []util.Poem) {
	poet := t.GetPoet()
	poems = make([]util.Poem, 0)

	selector := "h2"
	t.doc.Find("body").Find("p[align=\"center\"]").Each(func(i int, s *goquery.Selection) {
		gbkFullTitle := s.Text()
		fullTitleBytes := []byte(gbkFullTitle)
		fullTitle := strings.TrimSpace(util.GBK2Unicode(fullTitleBytes))
		title := strings.TrimSpace(strings.Split(fullTitle, " ")[0])

		if i == 1 && title != "" {
			selector = "p[align=\"center\"]"
		}
	})

	t.doc.Find("body").Find(selector).Each(func(i int, s *goquery.Selection) {
		gbkFullTitle := s.Text()
		fullTitleBytes := []byte(gbkFullTitle)
		fullTitle := strings.TrimSpace(util.GBK2Unicode(fullTitleBytes))
		title := strings.TrimSpace(strings.Split(fullTitle, " ")[0])
		subTitle := fullTitle[len(title):]

		if title != "" {
			poem := util.Poem{
				Author:   poet.Name,
				Source:   t.uctx.URL().String(),
				Title:    title,
				Subtitle: subTitle,
			}
			poems = append(poems, poem)
		}
	})

	return
}

func (t XianDaiShi) getPoemsContent() string {
	ps := t.initPoems()

	lp := ps[len(ps)-1]

	gbkText := t.doc.Find("body").Text()
	textBytes := []byte(gbkText)
	text := util.GBK2Unicode(textBytes)

	lastTitle := lp.Title

	i := strings.Index(text, lastTitle)
	lastTitleLen := len(lastTitle)

	k := strings.Index(text, "中国诗歌库")
	poemsContent := text[i+lastTitleLen: k]

	return poemsContent
}

func (t XianDaiShi) GetPoemsByPureTextAnalysis() (poems []util.Poem) {
	poems = t.initPoems()
	content := t.getPoemsContent()
	for i := 0; i < len(poems); i++ {
		l := strings.Index(content, poems[i].Title) + len(poems[i].Title)
		content = content[l:]
		k := len(content)
		if i < (len(poems) - 1) {
			// 加空格是为了防止当前诗歌内容中包含有下一首的标题，导致得不到正确的索引号
			// 虽然这不是百分百可靠，但是出现的概率很低
			k = strings.Index(content, "  "+poems[i+1].Title)

			if k == 0 || k == -1 {
				k = strings.Index(content, poems[i+1].Title)
			}
		}
		content := content[0:k]
		poems[i].Body = content

		content = content[k:]
	}
	return
}

func (t XianDaiShi) isNeedParseFromPureText() bool {
	h, _ := t.doc.Find("body").Find("h2").Next().Next().Html()
	if h == "" {
		return true
	}
	return false
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

// 标题为p align="left" 标签，一个标签对应一卷，一卷内包含多首诗歌
// 例子页面 http://www.shiku.org/shiku/xs/zhengmin.htm
// 处理方案：在a标签后添加分隔符，然后获取所有文本，解析出诗歌标题和内容
func (t XianDaiShi) GetPoemsT3() (poems []util.Poem) {
	titles := make([]string, 0, 0)
	poet := t.GetPoet()
	poems = make([]util.Poem, 0, 0)
	sep := "=============="

	has999999 := false
	t.doc.Find("body").Find("a").Each(func(i int, s *goquery.Selection) {
		href, existHref := s.Attr("href")
		if existHref {
			fmt.Println(href)
			if strings.HasPrefix(href, "#") {
				//fmt.Println(strconv.Atoi(strings.TrimLeft(href, "#")))
				_, err := strconv.Atoi(strings.TrimLeft(href, "#"))

				if err == nil {
					gbkTitle := s.Text()
					titleBytes := []byte(gbkTitle)
					title := strings.TrimSpace(util.GBK2Unicode(titleBytes))
					titles = append(titles, title)
				}

			}
		}

		name, existName := s.Attr("name")
		if existName {
			t.doc.Find("body").Find("a[name=\"" + name + "\"]").AppendHtml(sep)
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
		fmt.Println(len(content))
		fmt.Println(len(titles))
		num := len(content)
		if has999999 {
			num = len(content) - 1
		}

		for i := 0; i < num; i++ {
			whole := strings.TrimSpace(sep + content[i])
			title := titles[i]
			body := strings.Replace(whole, sep+title, "", -1)
			poem := util.Poem{
				Author: poet.Name,
				Source: t.uctx.URL().String(),
				Title:  title,
				//Subtitle: subTitle,
				//Volume: volume,
				Body: body,
			}

			poems = append(poems, poem)

		}
	}

	return
}

func (t XianDaiShi) GetPoems() []util.Poem {
	//return t.GetPoemsPAndP()
	return t.GetPoemsT3()
	//if t.isNeedParseFromPureText() {
	//	fmt.Println("使用纯文本解析的方式获取诗歌，此方法有风险")
	//	return t.GetPoemsByPureTextAnalysis()
	//} else {
	//	fmt.Println("使用document目录树解析的方式获取诗歌")
	//	return t.GetPoemsH2AndP()
	//}
}
