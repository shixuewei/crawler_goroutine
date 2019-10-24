package parser

import (
	"learning/crawler_distributed/config"
	"learning/crawler_goroutine/engine"
	"regexp"
	"strconv"
)

const (
	bookPageRe = `<a href="https://www.biikan.com/([^/]+)/[0-9]+.shtml">([^<]+)</a>`
)

func ParseBookPage(content []byte, _ string) engine.ParseResult {
	re := regexp.MustCompile(bookPageRe)
	matches := re.FindStringSubmatch(string(content))
	result := engine.ParseResult{}
	sum, _ := strconv.Atoi(matches[2])
	for i := 1; i <= sum; i++ {
		result.Requests = append(result.Requests, engine.Request{
			Url:    "https://www.biikan.com" + "/" + matches[1] + "/" + strconv.Itoa(i) + ".shtml",
			Parser: engine.NewFuncParser(ParseBook, config.ParseBook),
		})
	}
	return result
}
