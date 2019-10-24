package parser

import (
	"learning/crawler_goroutine/engine"
	"learning/crawler_goroutine/model"
	"regexp"
	"strconv"
)

var authorRe = regexp.MustCompile(
	`<li><span class="gray2">作者：</span><a.*>([^<]+)</a></li>`)
var bookSortRe = regexp.MustCompile(
	`<li><span class="gray2">分类：</span><a.*>([^<]+)</a></li>`)
var languageRe = regexp.MustCompile(
	`<li><span class="gray2">语言：</span><a.*>([^<]+)</a></li>`)
var countryRe = regexp.MustCompile(
	`<li><span class="gray2">国家：</span><a.*>([^<]+)</a></li>`)
var clickRe = regexp.MustCompile(
	`<li><span class="gray2">点击：</span><span class="gray pr20">([^<]+)</span></li>`)
var wordsNumRe = regexp.MustCompile(
	`<li><span class="gray2">字数：</span><span class="gray">([^<]+)</span></li>`)
var idUrlRe = regexp.MustCompile(
	`https://www.biikan.com/.*book-([0-9]+).shtml`)

func ParseBookMes(content []byte, url string, name string) engine.ParseResult {
	bookMes := model.BookMes{}
	bookMes.Name = name
	bookMes.Author = extractString(content, authorRe)
	bookMes.BookSort = extractString(content, bookSortRe)
	bookMes.Language = extractString(content, languageRe)
	bookMes.Country = extractString(content, countryRe)
	bookMes.Click, _ = strconv.Atoi(extractString(content, clickRe))
	bookMes.WordsNum, _ = strconv.Atoi(extractString(content, wordsNumRe))

	result := engine.ParseResult{
		Items: []engine.Item{
			{
				Url:     url,
				Type:    "book",
				Id:      extractString([]byte(url), idUrlRe),
				Payload: bookMes,
			},
		},
	}
	return result
}

func extractString(content []byte, re *regexp.Regexp) string {
	match := re.FindSubmatch(content)
	if len(match) >= 2 {
		return string(match[1])
	} else {
		return ""
	}
}

//此处及以后的代码都是为了将ParseBookMes(content []byte, url string, name string)
//这个函数进行序列化，以方便在网络上进行传输。按照之前函数的序列化过程，这个具有三个参数
//所以单独定义结构，将其序列化

//定义爬取书的结构体（内置书名）
type BookParser struct {
	bookName string
}

//实现Parse接口的Parse方法
func (p *BookParser) Parse(content []byte, url string) engine.ParseResult {
	return ParseBookMes(content, url, p.bookName)
}

//实现Parse接口的Serialize()方法
func (p *BookParser) Serialize() (name string, args interface{}) {
	return "BookParser", p.bookName
}

//使用工厂函数创建BookParser
func NewBookParser(name string) *BookParser {
	return &BookParser{
		bookName: name,
	}
}
