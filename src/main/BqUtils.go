package main

import (
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
)

type CommonBqStruct struct {
	tableName string
	data interface{}
}

func uploadBq(ctx context.Context, bqStructs []CommonBqStruct) error {
	client, err := bigquery.NewClient(ctx, config.Bq.Project)
	if err != nil {
		return err
	}
	defer client.Close()
	dataset := client.Dataset(config.Bq.Dataset)

	//各BqStructを元にアップロードしていく
	for i := 0; i < len(bqStructs); i++ {
		bqStruct := bqStructs[i]
		uploader := dataset.Table(bqStruct.tableName).Uploader()
		err := uploader.Put(ctx, bqStruct.data)
		if(err != nil){
			return err
		}
	}
	return nil
}

//type Sample struct {
//	Id int
//	Name string
//	Created time.Time
//}
//
//func putSample(ctx context.Context) error {
//	client, err := bigquery.NewClient(ctx, config.Bq.Project)
//	if err != nil {
//		return err
//	}
//	defer client.Close()
//
//	log.Printf("project", config.Bq.Project)
//	log.Printf("dataset", config.Bq.Dataset)
//
//	u := client.Dataset(config.Bq.Dataset).Table("sample").Uploader()
//
//	samples := []*Sample{
//		{Id: 10, Name: "Taro", Created: time.Now()},
//	}
//	err = u.Put(ctx, samples)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}