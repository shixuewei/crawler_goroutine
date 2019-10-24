package save

import (
	"context"
	"errors"
	"learning/crawler_goroutine/engine"
	"log"

	"gopkg.in/olivere/elastic.v5"
)

//保存Item,使用到elasticsearch
//index为数据库的名
func ItemSaver(index string) (chan engine.Item, error) {
	//创建elasticsearch客户端
	//elastic.SetSniff(false)：不使用集群
	client, err := elastic.NewClient(
		elastic.SetSniff(false))

	if err != nil {
		return nil, err
	}
	//初始化输出
	out := make(chan engine.Item)
	//ItemSaver 代表一个协程
	go func() {
		itemCount := 0
		for {
			item := <-out
			log.Printf("Item Saver : got item #%d : %v ", itemCount, item)
			itemCount++

			err := Save(client, item, index)
			if err != nil {
				log.Printf("Item saver : error"+"saving item %v : %v", item, err)
			}
		}
	}()
	return out, nil
}

//保存Item数据到elastic search
func Save(client *elastic.Client, item engine.Item, index string) (err error) {

	if item.Type == "" {
		return errors.New("Must supply type")
	}
	//启用elastic search服务
	indexService := client.Index().Index(index).Type(item.Type).BodyJson(item)

	if item.Id != "" {
		indexService.Id(item.Id)
	}
	_, err = indexService.Do(context.Background())

	if err != nil {
		return err
	}
	return err
}
