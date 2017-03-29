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

func Save(fn string, p []Poem) {
	b := ""
	for _, v := range p {
		b += v.Author + "\r\n" + v.Title + "\r\n\r\n" + v.Body + "\r\n" + v.Source + "\r\n\r\n"
	}
	bytes := []byte(b)

	e := ioutil.WriteFile("./"+fn, bytes, 0666)
	if e != nil {
		fmt.Println(e)
	}
}
