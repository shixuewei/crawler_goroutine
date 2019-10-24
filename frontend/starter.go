package main

import (
	"learning/crawler_goroutine/frontend/controller"
	"net/http"
)

func main() {
	// 搜索首页
	http.Handle("/",
		http.FileServer(http.Dir("F:/develop/go_project/src/learning/crawler_goroutine/frontend/view")))
	// 搜索结果页面
	http.Handle("/search",
		controller.CreateSearchResultHandler("F:/develop/go_project/src/learning/crawler_goroutine/frontend/view/template.html"))
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		panic(err)
	}
}
