package model

import "encoding/json"

type BookMes struct {
	Name     string
	Author   string
	BookSort string
	Language string
	Country  string
	Click    int
	WordsNum int
}

func FromJsonObj(o interface{}) (BookMes, error) {
	var bookMes BookMes
	s, err := json.Marshal(o)
	if err != nil {
		return bookMes, err
	}

	err = json.Unmarshal(s, &bookMes)
	return bookMes, err
}
