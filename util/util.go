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

func Save(fn, url string, p []Poem) {
	b := url
	for _, v := range p {
		b += v.Author + "\r\n" + v.Title + "\r\n" + v.Subtitle + "\r\n\r\n" + v.Body + "\r\n" + v.Source + "\r\n\r\n"
	}
	bytes := []byte(b)

	e := ioutil.WriteFile("./"+fn, bytes, 0666)
	if e != nil {
		fmt.Println(e)
	}
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
