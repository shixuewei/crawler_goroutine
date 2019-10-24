package engine

//定义一个爬取页面信息的函数类型
type ParserFunc func(content []byte, url string) ParseResult

//定义一个Parser接口，其包含两个方法，
//1、Parse 方法：返回为ParseResult，获取实际的爬取到的数据
//2、Serialize 方法（序列化方法）分布式时会用到。
type Parser interface {
	Parse(content []byte, url string) ParseResult
	Serialize() (name string, args interface{})
}

type Request struct {
	Url    string
	Parser Parser
}

type ParseResult struct {
	Requests []Request
	Items    []Item
}

//用来保存相应的数据
type Item struct {
	//相应的链接
	Url string
	//数据库类型（数据库的名字）
	Type string
	//每本书的id
	Id string
	//每本书的详细信息（里边的内容较多，所以采用接口变量来装）
	Payload interface{}
}

//定义一个未爬取到的页面的变量，同时实现Parser接口
type NilParser struct {
}

func (NilParser) Serialize() (name string, args interface{}) {
	return "NilParser", nil
}

func (NilParser) Parse(_ []byte, _ string) ParseResult {
	return ParseResult{}
}

//定义一个爬取页面的函数类型的变量，同时也实现Parse接口
type FuncParser struct {
	//具体的函数方法
	parser ParserFunc
	//函数的名字
	Name string
}

func (f *FuncParser) Parse(content []byte, url string) ParseResult {
	return f.parser(content, url)
}

func (f *FuncParser) Serialize() (name string, args interface{}) {
	return f.Name, nil
}

//工厂函数：创建一个FuncParser实例
func NewFuncParser(p ParserFunc, name string) *FuncParser {
	return &FuncParser{
		parser: p,
		Name:   name,
	}
}
