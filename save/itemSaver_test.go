package save

import (
	"learning/crawler_goroutine/engine"
	"learning/crawler_goroutine/model"
	"testing"

	"gopkg.in/olivere/elastic.v5"
)

func TestSave(t *testing.T) {
	expectedBook := engine.Item{
		Url:  "https://www.biikan.com/xinli/book-19239.shtml",
		Type: "book",
		Id:   "19239",
		Payload: model.BookMes{
			Name:     "岁朝清供",
			Author:   "汪曾祺",
			BookSort: "文学",
			Language: "中文",
			Country:  "中国当代",
			Click:    17,
			WordsNum: 125102,
		},
	}

	client, err := elastic.NewClient(
		elastic.SetSniff(false))

	if err != nil {
		panic(err)
	}
	const index = "dating_test"

	Save(client, expectedBook, index)
}
