package main

import (
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"strings"
)

type CommonBqStruct struct {
	tableName string
	data      interface{}
}

type CommonBqTableDefintion struct {
	tableName string
	metadata  bigquery.TableMetadata
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
		log.Infof(ctx, "LET'S UPLOAD: %v", bqStructs[i].tableName)
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

func deleteAndCreateBq(ctx context.Context, bqStructs []CommonBqStruct) (error) {
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
		if (schemaError != nil) {
			return schemaError
		}
		table := dataset.Table(bqStruct.tableName)
		delErr := table.Delete(ctx)
		if delErr != nil {
			log.Errorf(ctx, "ERROR: %v", delErr)
		}
		log.Infof(ctx, "LET'S CREATE: %v", bqStruct.tableName)
		createErr := table.Create(ctx, &bigquery.TableMetadata{
			Name:   bqStruct.tableName,
			Schema: schema,
		})
		if createErr != nil {
			return createErr
		}
	}
	return nil
}

const BASIC_DATA_CHECK_QUERY = "SELECT COUNT(1) FROM `<project>.<data_set>.<table>` WHERE name = '<nameFilter>'"

type CountData struct {
	Count int64      `json:"id,omitempty"`
}

func hasDataSimple(ctx context.Context, basicSQL string) bool {
	// contextとprojectIDを元にBigQuery用のclientを生成
	client, err := bigquery.NewClient(ctx, config.Bq.Project)

	if err != nil {
		log.Errorf(ctx, "Failed to create client:%v", err)
	}

	var query string
	query = strings.Replace(basicSQL, "<project>", config.Bq.Project, -1)
	query = strings.Replace(query, "<data_set>", config.Bq.Dataset, -1)


	log.Infof(ctx, "QUERY SQL:" + query)
	// 引数で渡した文字列を元にQueryを生成
	q := client.Query(query)

	log.Infof(ctx, "LAUNCH SQL")
	// 実行のためのqueryをサービスに送信してIteratorを通じて結果を返す
	// itはIterator
	it, readErr := q.Read(ctx)

	if readErr != nil {
		log.Errorf(ctx, "Failed to Read Query:%v", readErr)
	}

	var countData CountData
	nextErr := it.Next(&countData)
	if nextErr != nil {
		log.Errorf(ctx, "Failed to it.Next(&countData):%v", nextErr)
		return true
	}

	if countData.Count == 0 {
		return false
	} else {
		return true
	}
}

func hasData(ctx context.Context, tableName string, nameFilter string) bool {
	// contextとprojectIDを元にBigQuery用のclientを生成
	client, err := bigquery.NewClient(ctx, config.Bq.Project)

	if err != nil {
		log.Errorf(ctx, "Failed to create client:%v", err)
	}

	var query string
	query = strings.Replace(BASIC_DATA_CHECK_QUERY, "<project>", config.Bq.Project, -1)
	query = strings.Replace(query, "<data_set>", config.Bq.Dataset, -1)
	query = strings.Replace(query, "<table>", tableName, -1)
	query = strings.Replace(query, "<nameFilter>", nameFilter, -1)

	// 引数で渡した文字列を元にQueryを生成
	q := client.Query(query)

	// 実行のためのqueryをサービスに送信してIteratorを通じて結果を返す
	// itはIterator
	it, readErr := q.Read(ctx)

	if readErr != nil {
		log.Errorf(ctx, "Failed to Read Query:%v", readErr)
	}

	var countData CountData
	nextErr := it.Next(&countData)
	if nextErr != nil {
		log.Errorf(ctx, "Failed to it.Next(&countData):%v", nextErr)
		return true
	}
	if countData.Count == 0 {
		return false
	} else {
		return true
	}
}

func runQuery(ctx context.Context, basicSQL string){
	// contextとprojectIDを元にBigQuery用のclientを生成
	client, err := bigquery.NewClient(ctx, config.Bq.Project)
	if err != nil {
		log.Errorf(ctx, "Failed to create client:%v", err)
	}

	var query = strings.Replace(basicSQL, "<project>", config.Bq.Project, -1)
	query = strings.Replace(query, "<data_set>", config.Bq.Dataset, -1)

	// 引数で渡した文字列を元にQueryを生成
	q := client.Query(query)

	job, runErr := q.Run(ctx)
	if runErr != nil {
		log.Errorf(ctx, "Failed to RUNNING:%v", runErr)
		log.Errorf(ctx, "FAILED SQL:%v", query)
	}
	if job.LastStatus().Err() != nil {
		log.Errorf(ctx, "FAILED! JOB STATUS IS NOT GOOD:%v", job.LastStatus().Err())
		log.Errorf(ctx, "FAILED SQL:%v", query)
	}
}

func readQuery(ctx context.Context, basicSQL string)(*bigquery.RowIterator, error){
	// contextとprojectIDを元にBigQuery用のclientを生成
	client, err := bigquery.NewClient(ctx, config.Bq.Project)
	if err != nil {
		log.Errorf(ctx, "Failed to create client:%v", err)
	}

	var query = strings.Replace(basicSQL, "<project>", config.Bq.Project, -1)
	query = strings.Replace(query, "<data_set>", config.Bq.Dataset, -1)

	// 引数で渡した文字列を元にQueryを生成
	q := client.Query(query)

	return q.Read(ctx)
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
