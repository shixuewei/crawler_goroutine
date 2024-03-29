package controller

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"learning/crawler_goroutine/frontend/view"

	"learning/crawler_goroutine/engine"
	"learning/crawler_goroutine/frontend/model"
	"reflect"

	"regexp"

	"gopkg.in/olivere/elastic.v5"
)

//搜索结果句柄
type SearchResultHandler struct {
	//展示页面（模板页面）
	view view.SearchResultView
	//客户连接端
	client *elastic.Client
}

//工厂函数（创建SearchResultHandler）
func CreateSearchResultHandler(template string) SearchResultHandler {
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		panic(err)
	}

	return SearchResultHandler{
		view:   view.CreateSearchResultView(template),
		client: client,
	}
}

// localhost:8888/search?q=男 已购房 已购车&from=20
//查询
func (h SearchResultHandler) ServeHTTP(
	w http.ResponseWriter, req *http.Request) {
	q := strings.TrimSpace(req.FormValue("q"))

	from, err := strconv.Atoi(
		req.FormValue("from"))
	if err != nil {
		from = 0
	}

	page, err := h.getSearchResult(q, from)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.view.Render(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

const pageSize = 10

func (h SearchResultHandler) getSearchResult(
	q string, from int) (model.SearchResult, error) {
	var result model.SearchResult
	result.Query = q

	resp, err := h.client.
		Search("dating_book").
		Query(elastic.NewQueryStringQuery(rewriteQueryString(q))).
		From(from).
		Do(context.Background())

	if err != nil {
		return result, err
	}

	result.Hits = resp.TotalHits()
	result.Start = from
	result.Items = resp.Each(reflect.TypeOf(engine.Item{}))
	if result.Start == 0 {
		result.PrevFrom = -1
	} else {
		result.PrevFrom =
			(result.Start - 1) /
				pageSize * pageSize
	}
	result.NextFrom =
		result.Start + len(result.Items)

	return result, nil
}

// Rewrites query string. Replaces field names
// like "Age" to "Payload.Age"
func rewriteQueryString(q string) string {
	re := regexp.MustCompile(`([A-Z][a-z]*):`)
	//其中$代表就是([A-Z][a-z]*)里边的值
	return re.ReplaceAllString(q, "Payload.$1:")
}
