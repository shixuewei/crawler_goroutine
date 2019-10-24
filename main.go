package main

import (
	"learning/crawler_distributed/config"
	"learning/crawler_goroutine/book/parser"
	"learning/crawler_goroutine/engine"
	"learning/crawler_goroutine/save"
	"learning/crawler_goroutine/scheduler"
)

//并发版爬虫
func main() {
	itemChan, err := save.ItemSaver("dating_book")
	if err != nil {
		panic(err)
	}
	e := engine.ConcurrentEngine{
		Scheduler:        &scheduler.QueuedScheduler{},
		WorkerCount:      100,
		ItemChan:         itemChan,
		RequestProcessor: engine.Worker,
	}
	e.Run(engine.Request{
		Url:    "https://www.biikan.com",
		Parser: engine.NewFuncParser(parser.ParseBookSort, config.ParseBookSort),
	})
}
