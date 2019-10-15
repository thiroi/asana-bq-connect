package connector

import (
	"cloud.google.com/go/bigquery"
	"golang.org/x/net/context"
	"strings"
	"log"
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
		log.Println("LET'S UPLOAD: ", bqStructs[i].tableName)
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
			log.Println( "ERROR: ", delErr)
		}
		log.Println("LET'S CREATE: ", bqStruct.tableName)
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
		log.Println("Failed to create client:", err)
	}

	var query string
	query = strings.Replace(basicSQL, "<project>", config.Bq.Project, -1)
	query = strings.Replace(query, "<data_set>", config.Bq.Dataset, -1)


	log.Println("QUERY SQL:" + query)
	// 引数で渡した文字列を元にQueryを生成
	q := client.Query(query)

	log.Println( "LAUNCH SQL")
	// 実行のためのqueryをサービスに送信してIteratorを通じて結果を返す
	// itはIterator
	it, readErr := q.Read(ctx)

	if readErr != nil {
		log.Println("Failed to Read Query:", readErr)
	}

	var countData CountData
	nextErr := it.Next(&countData)
	if nextErr != nil {
		log.Println("Failed to it.Next(&countData):", nextErr)
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
		log.Println("Failed to create client:", err)
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
		log.Println("Failed to Read Query:", readErr)
	}

	var countData CountData
	nextErr := it.Next(&countData)
	if nextErr != nil {
		log.Println("Failed to it.Next(&countData):", nextErr)
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
		log.Println("Failed to create client:", err)
	}

	var query = strings.Replace(basicSQL, "<project>", config.Bq.Project, -1)
	query = strings.Replace(query, "<data_set>", config.Bq.Dataset, -1)

	// 引数で渡した文字列を元にQueryを生成
	q := client.Query(query)

	job, runErr := q.Run(ctx)
	if runErr != nil {
		log.Println( "Failed to RUNNING:", runErr)
		log.Println("FAILED SQL:", query)
	}
	if job.LastStatus().Err() != nil {
		log.Println("FAILED! JOB STATUS IS NOT GOOD:", job.LastStatus().Err())
		log.Println("FAILED SQL:", query)
	}
}

func readQuery(ctx context.Context, basicSQL string)(*bigquery.RowIterator, error){
	// contextとprojectIDを元にBigQuery用のclientを生成
	client, err := bigquery.NewClient(ctx, config.Bq.Project)
	if err != nil {
		log.Println("Failed to create client:", err)
	}

	var query = strings.Replace(basicSQL, "<project>", config.Bq.Project, -1)
	query = strings.Replace(query, "<data_set>", config.Bq.Dataset, -1)

	// 引数で渡した文字列を元にQueryを生成
	q := client.Query(query)

	return q.Read(ctx)
}
