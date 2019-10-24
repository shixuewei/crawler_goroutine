package view

import (
	"html/template"
	"io"

	"learning/crawler_goroutine/frontend/model"
)

//结果展示模板
type SearchResultView struct {
	template *template.Template
}

//工厂函数（创建SearchResultView）
func CreateSearchResultView(filename string) SearchResultView {
	return SearchResultView{
		template: template.Must(template.ParseFiles(filename)),
	}
}

//将数据写出到前端页面中
func (s SearchResultView) Render(
	w io.Writer, data model.SearchResult) error {
	return s.template.Execute(w, data)
}
