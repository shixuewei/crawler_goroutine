package parser

import (
	"learning/crawler_goroutine/engine"
	"regexp"
)

const (
	bookRe = `<h1><a href="([^"]+)" target="_blank">([^<]+)</a></h1>`
)

func ParseBook(content []byte, _ string) engine.ParseResult {
	re := regexp.MustCompile(bookRe)
	matches := re.FindAllStringSubmatch(string(content), 20)
	result := engine.ParseResult{}
	for _, v := range matches {
		url := "https://www.biikan.com" + v[1]
		name := v[2]
		//result.Items = append(result.Items, "BookName:"+v[2])
		result.Requests = append(result.Requests, engine.Request{
			Url:    url,
			Parser: NewBookParser(name),
		})
	}
	return result
}
