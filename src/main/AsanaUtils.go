
package main

import (
	"net/http"
	"io/ioutil"
	"google.golang.org/appengine/urlfetch"
	"golang.org/x/net/context"
)

func loadAsana(ctx context.Context, url string)([]byte, error){
	//tokenとurlを元にGETする
	client := urlfetch.Client(ctx)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer " + config.Asana.Token)
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