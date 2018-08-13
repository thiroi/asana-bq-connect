package main

import (
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"log"
)

type CommonBqStruct struct {
	tableName string
	data      interface{}
}

type CommonBqTableDefintion struct {
	tableName string
	metadata bigquery.TableMetadata
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
		if (err != nil) {
			return err
		}
	}
	return nil
}

func initAsanaData(ctx context.Context) error {
	return deleteAndCreateBq(
		ctx,
		[]CommonBqStruct{
			{"project", Project{}},
			{"section", Section{}},
			{"task", Task{}},
			{"tag", Tag{}},
			{"user", User{}},
		})
}

func deleteAndCreateBq(ctx context.Context, bqStructs []CommonBqStruct)(error){
	client, err := bigquery.NewClient(ctx, config.Bq.Project)
	if err != nil {
		return err
	}
	defer client.Close()
	dataset := client.Dataset(config.Bq.Dataset)

	//各BqStructを元にアップロードしていく
	for i := 0; i < len(bqStructs); i++ {
		bqStruct := bqStructs[i]
		schema, schemaError := bigquery.InferSchema(bqStruct.data)
		if (schemaError != nil){
			return schemaError
		}
		table := dataset.Table(bqStruct.tableName)
		delErr := table.Delete(ctx)
		if delErr != nil {
			log.Println(delErr)
		}
		log.Println("LET'S CREATE")
		createErr := table.Create(ctx, &bigquery.TableMetadata{
			Name: bqStruct.tableName,
			Schema:schema,
			})
		if createErr != nil {
			return createErr
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
