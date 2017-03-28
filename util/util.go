package util

import (
	"fmt"

	"github.com/axgle/mahonia"
	"io/ioutil"
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

type Poem struct {
	Author string // 作者
	Intro  string // 作者简介
	Source string // 来源
	Title  string // 标题
	Body   string // 内容
}

func Save(fn string, p []Poem) {
	b := ""
	for _, v := range p {
		b += v.Author + "\r\n" + v.Intro + "\r\n" + v.Title + "\r\n\r\n" + v.Body + "\r\n" + v.Source + "\r\n\r\n"
	}
	bytes := []byte(b)

	e := ioutil.WriteFile("./"+fn, bytes, 0666) //写入文件(字节数组)
	if e != nil {
		fmt.Println(e)
	}
}
