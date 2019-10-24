package parser

import (
	"learning/crawler_distributed/config"
	"learning/crawler_goroutine/engine"
	"regexp"
)

const (
	bookSort = `<li ><a href="/([^"]+)">([^<]+)</a></li>`
)

func ParseBookSort(content []byte, _ string) engine.ParseResult {
	re := regexp.MustCompile(bookSort)
	matches := re.FindAllStringSubmatch(string(content), -1)
	result := engine.ParseResult{}
	for _, v := range matches {
		//result.Items = append(result.Items, "BookSort:"+v[2])
		result.Requests = append(result.Requests, engine.Request{
			Url:    "https://www.biikan.com" + "/" + v[1],
			Parser: engine.NewFuncParser(ParseBookPage, config.ParseBookPage),
		})
	}
	return result
}
