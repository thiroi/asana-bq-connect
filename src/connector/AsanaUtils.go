
package connector

import (
	"net/http"
	"io/ioutil"
	"golang.org/x/net/context"
)

func loadAsana(ctx context.Context, url string)([]byte, error){
	//tokenとurlを元にGETする

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer " + config.Asana.Token)
	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil{
		return nil, err
	}
	defer res.Body.Close()

	//byt変換
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}