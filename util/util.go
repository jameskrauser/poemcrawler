package util

import (
	"fmt"

	"github.com/axgle/mahonia"
	"io/ioutil"
	"strings"
)

func GBK2Unicode(data []byte) string {
	var dec mahonia.Decoder

	dec = mahonia.NewDecoder("gbk")

	_, ret, err := dec.Translate(data, true)
	if err != nil {
		fmt.Println(err)
	}

	return string(ret)
}

func SaveToFile(fn, url string, poet Poet, poems []Poem) {
	b := "网页地址：" + url + "\r\n\r\n诗人名字：" + poet.Name + "\r\n\r\n诗人简介：" + poet.Intro+"\r\n\r\n"
	for _, poem := range poems {
		b += "作者：" + poem.Author + "\r\n\r\n诗歌标题：" + poem.Title + "\r\n\r\n诗歌副标题：" +
			poem.Subtitle + "\r\n\r\n诗歌内容：" + poem.Body + "\r\n\r\n诗歌出处：" + poem.Source + "\r\n\r\n"
	}
	bytes := []byte(b)

	e := ioutil.WriteFile("./"+fn, bytes, 0666)
	if e != nil {
		fmt.Println(e)
	}
}

func CheckPoet(p Poet) bool {
	name := strings.Replace(p.Name, " ", "", -1)
	intro := strings.Replace(p.Intro, " ", "", -1)

	if name == "" || intro == "" {
		return true
	}
	return false
}

func CheckPoems(ps []Poem) bool {
	for _, p := range ps {
		author := strings.Replace(p.Author, " ", "", -1)
		title := strings.Replace(p.Title, " ", "", -1)
		body := strings.Replace(p.Body, " ", "", -1)

		if author == "" || title == "" || body == "" {
			return true
		}
	}

	return false
}

func TrimRightSpace(s string) string {
	if s == "" {
		return s
	}

	for {
		if strings.HasSuffix(s, " ") {
			s = strings.TrimSuffix(s, " ")
			s = strings.TrimSuffix(s, "\r\n")
			s = strings.TrimSuffix(s, "\r")
			s = strings.TrimSuffix(s, "\n")
		} else {
			break
		}
		if strings.HasSuffix(s, "\r") {
			s = strings.TrimSuffix(s, " ")
			s = strings.TrimSuffix(s, "\r\n")
			s = strings.TrimSuffix(s, "\r")
			s = strings.TrimSuffix(s, "\n")
		} else {
			break
		}
		if strings.HasSuffix(s, "\n") {
			s = strings.TrimSuffix(s, " ")
			s = strings.TrimSuffix(s, "\r\n")
			s = strings.TrimSuffix(s, "\r")
			s = strings.TrimSuffix(s, "\n")
		} else {
			break
		}
	}

	return s
}
