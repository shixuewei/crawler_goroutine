package engine

import (
	"learning/crawler_goroutine/fetcher"
	"log"
)

//简单版的爬虫引擎
type SimpleEngine struct {
}

func (e SimpleEngine) Run(seeds ...Request) {
	var requests []Request
	for _, v := range seeds {
		requests = append(requests, v)
	}
	//本质上和单机版是一样的
	for len(requests) > 0 {
		r := requests[0]
		requests = requests[1:]

		parseResult, err := Worker(r)
		if err != nil {
			continue
		}
		requests = append(requests, parseResult.Requests...)

		for _, v := range parseResult.Items {
			log.Printf("Got item %v", v)
		}

	}
}

//爬取页面的工作器
func Worker(r Request) (ParseResult, error) {
	//log.Println("Fetching %s", r.Url)

	body, err := fetcher.Fetch(r.Url)
	if err != nil {
		log.Printf("Fetcher : error"+"fetching url %s: %v", r.Url, err)
		return ParseResult{}, err
	}

	return r.Parser.Parse(body, r.Url), nil
}
