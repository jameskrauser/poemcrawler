package htmltype

import (
	"github.com/PuerkitoBio/gocrawl"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"fmt"
	"poemcrawler/util"
)

// 处理国际诗歌的类型
// 首页 http://www.shiku.org/shiku/ws/wg/index.htm
type GuoJiShi struct {
	Base *ShiKu
	uctx *gocrawl.URLContext
	res  *http.Response
	doc  *goquery.Document
	Poet util.Poet
}

func NewGuoJiShi(uctx *gocrawl.URLContext, res *http.Response, doc *goquery.Document) *GuoJiShi {
	poet := GetPoet(uctx, doc)
	return &GuoJiShi{
		Base: NewShiKu(uctx, res, doc),
		uctx: uctx,
		res:  res,
		doc:  doc,
		Poet: poet,
	}
}

func GetFirstPoemTitleWithSep(doc *goquery.Document) string {
	titles := make([]string, 0, 0)
	var title string
	doc.Find("body").Find("a").Each(func(i int, s *goquery.Selection) {
		href, existHref := s.Attr("href")
		if existHref {
			if strings.HasPrefix(href, "#") {
				gbkTitle := s.Text()
				titleBytes := []byte(gbkTitle)
				title = strings.TrimSpace(util.GBK2Unicode(titleBytes))
				titles = append(titles, title+ALinkTitleSep)
				s.AppendHtml(ALinkTitleSep)
				return
			}
		}
	})
	if len(titles) > 0 {
		return titles[0]
	}

	return ""
}

func GetFirstPoemTitleWithSep1FromPoemBody(doc *goquery.Document) string {
	titles := make([]string, 0, 0)
	var title string
	doc.Find("body").Find("h2").Each(func(i int, s *goquery.Selection) {
		gbkTitle := s.Text()
		titleBytes := []byte(gbkTitle)
		title = strings.TrimSpace(util.GBK2Unicode(titleBytes))
		titles = append(titles, title+BodyTitleSep1)
		s.AppendHtml(BodyTitleSep1)
	})
	if len(titles) > 0 {
		return titles[0]
	}

	return ""
}

func GetPoet(uctx *gocrawl.URLContext, doc *goquery.Document) (poet util.Poet) {
	gbkStr := doc.Find("title").Text()
	bytes := []byte(gbkStr)
	title := strings.TrimSpace(util.GBK2Unicode(bytes))

	var name string
	arr := strings.Split(title, "::")

	for _, v := range arr {
		if strings.Contains(v, "诗选") || strings.Contains(v, "诗集") {
			name = strings.TrimSpace(strings.Split(v, "诗选")[0])
			name = strings.TrimSpace(strings.Split(name, "诗集")[0])
		}
	}

	if name == "" {
		gbkStr = doc.Find("body").Find("h1").Text()
		bytes = []byte(gbkStr)
		title = strings.TrimSpace(util.GBK2Unicode(bytes))
		name = strings.TrimSpace(strings.Split(title, "诗选")[0])
		name = strings.TrimSpace(strings.Split(name, "诗集")[0])
	}

	ft := GetFirstPoemTitleWithSep(doc)
	if ft == "" {
		// 处理页面上诗人简介下面的标题不是链接的情况
		// http://www.shiku.org/shiku/ws/wg/corneille.htm
		bft := GetFirstPoemTitleWithSep1FromPoemBody(doc)
		if bft == "" {
			poet = util.Poet{
				Name:   name,
				Intro:  "",
				Source: uctx.URL().String(),
			}

		} else {
			gbkStr = doc.Find("body").Text()
			bytes = []byte(gbkStr)
			text := strings.TrimSpace(util.GBK2Unicode(bytes))
			text = strings.Replace(text, name+"诗选", "", 1)
			text = strings.Replace(text, name+"诗选", "", 1)

			index := strings.Index(text, bft)
			if index > 0 {
				text = text[0:index]
			}

			intro := strings.TrimSpace(text)
			poet = util.Poet{
				Name:   name,
				Intro:  intro,
				Source: uctx.URL().String(),
			}
		}
	} else {
		gbkStr = doc.Find("body").Text()
		bytes = []byte(gbkStr)
		text := strings.TrimSpace(util.GBK2Unicode(bytes))
		text = strings.Replace(text, name+"诗选", "", 1)

		index := strings.Index(text, ft)
		if index > 0 {
			text = text[0:index]
		}

		intro := strings.TrimSpace(text)
		poet = util.Poet{
			Name:   name,
			Intro:  intro,
			Source: uctx.URL().String(),
		}

	}

	return poet
}

// 通过在标题前加分隔符后解析纯文本的方式解析页面
func (t GuoJiShi) GetPoems() (poems []util.Poem) {
	titles := make([]string, 0, 0)
	poet := t.Poet
	poems = make([]util.Poem, 0, 0)

	has999999 := false
	t.doc.Find("body").Find("a").Each(func(i int, s *goquery.Selection) {
		href, existHref := s.Attr("href")
		if existHref {
			if strings.Contains(href, "#") {
				gbkTitle := s.Text()
				titleBytes := []byte(gbkTitle)
				title := strings.TrimSpace(util.GBK2Unicode(titleBytes))
				title = strings.Replace(title, ALinkTitleSep, "", -1)
				titles = append(titles, title)
			}
		}

		name, existName := s.Attr("name")

		if existName {
			if strings.Contains("0123456789", name[0:1]) {
				t.doc.Find("body").Find("a[name=\"" + name + "\"]").AppendHtml(BodyTitleSep)
			}
		}

		if name == "999999" {
			has999999 = true
		}
	})

	gbkText := t.doc.Text()
	TextBytes := []byte(gbkText)
	text := strings.TrimSpace(util.GBK2Unicode(TextBytes))
	// 去掉页脚文字：中国诗歌库 中华诗库 中国诗典 中国诗人 中国诗坛 首页
	text = strings.Replace(text, "中国诗歌库", "", -1)
	text = strings.Replace(text, "中华诗库", "", -1)
	text = strings.Replace(text, "中国诗典", "", -1)
	text = strings.Replace(text, "中国诗人", "", -1)
	text = strings.Replace(text, "中国诗坛", "", -1)
	text = strings.Replace(text, "首页", "", -1)
	// 可能页面底部有以_uacct开头的js文字
	// 例如页面： http://www.shiku.org/shiku/ws/wg/corneille.htm
	index := strings.Index(text, "_uacct")
	if index > 0 {
		text = text[0:index]
	}

	text = strings.TrimSpace(text)
	textArr := strings.Split(text, BodyTitleSep)
	if strings.Contains(text, BodyTitleSep) {
		content := textArr[1:]
		if has999999 {
			content = textArr[1:len(textArr)-1]
		}

		fmt.Println("解析到的诗歌体数量为：", len(content))
		fmt.Println("解析到的诗歌标题数量为：", len(titles))

		// 标题非链接的情况，获取不到标题，例如： http://www.shiku.org/shiku/ws/wg/corneille.htm
		// 标题链接少于实际的诗歌体数量的情况，例如：http://www.shiku.org/shiku/ws/wg/mallarme.htm
		if len(titles) != len(content) {
			for _, whole := range content {
				var title string
				whole = strings.TrimLeft(whole, " ")
				if strings.Contains(whole, BodyTitleSep1){
					title = strings.Split(whole, BodyTitleSep1)[0]
				}else{
					title = strings.Split(whole, " ")[0]
				}

				str := strings.TrimSpace(whole)
				body := strings.TrimLeft(str, title)
				body = strings.Replace(body, BodyTitleSep, "", -1)
				body = strings.Replace(body, BodyTitleSep1, "", -1)

				poem := util.Poem{
					Author: poet.Name,
					Source: t.uctx.URL().String(),
					Title:  title,
					Body:   body,
				}

				poems = append(poems, poem)
			}
		} else {
			count := len(content)
			for i := 0; i < count; i++ {
				whole := strings.TrimSpace(BodyTitleSep + content[i])
				title := titles[i]
				body := strings.Replace(whole, BodyTitleSep+title, "", -1)
				body = strings.Replace(body, BodyTitleSep, "", -1)

				// 网页本身有错误，标题为空
				// 标题与内容混在一起：http://www.shiku.org/shiku/xs/hanzuorong.htm
				if title == "" {
					body = strings.Trim(body, " ")
					body = strings.Trim(body, "\n")
					arr := strings.Split(body, "\n")
					if len(arr) > 1 {
						title = strings.Replace(arr[0], " ", "", -1)
					}
				}

				poem := util.Poem{
					Author: poet.Name,
					Source: t.uctx.URL().String(),
					Title:  title,
					Body:   body,
				}

				poems = append(poems, poem)
			}
		}
	}

	return
}
