package view

import (
	"html/template"
	"learning/crawler_goroutine/engine"
	"learning/crawler_goroutine/frontend/model"
	common "learning/crawler_goroutine/model"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	template := template.Must(template.ParseFiles("template.html"))

	out, err := os.Create("template.test.html")
	page := model.SearchResult{}
	page.Hits = 123
	item := engine.Item{
		Url:  "https://www.biikan.com/xinli/book-19239.shtml",
		Type: "book",
		Id:   "19239",
		Payload: common.BookMes{
			Name:     "岁朝清供",
			Author:   "汪曾祺",
			BookSort: "文学",
			Language: "中文",
			Country:  "中国当代",
			Click:    17,
			WordsNum: 125102,
		},
	}
	for i := 0; i < 10; i++ {
		page.Items = append(page.Items, item)
	}

	//err = template.Execute(os.Stdout, page)
	err = template.Execute(out, page)
	if err != nil {
		panic(err)
	}
}
