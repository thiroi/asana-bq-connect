package connector

import (
	"fmt"
	"net/http"
	"log"
	"bytes"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"time"
	"golang.org/x/net/context"
	"github.com/nlopes/slack"
)

func SlackSender(w http.ResponseWriter, ctx *http.Request) {
	fmt.Fprint(w, "Hello, world222!")
	//sendSlack(ctx)
	log.Println("slack done")
}

func sendNlope(ctx context.Context, message string, reporter Reporter){
	log.Println("Start sending with nlope")
	slack.SetHTTPClient(new(http.Client))
	api := slack.New(config.Slack.Token)
	params := slack.PostMessageParameters{}
	params.IconEmoji = reporter.Icon
	params.Username = reporter.Name
	message = reporter.Talk + "\n\n" + "======================================== \n" + message + "========================================"
	log.Println(message)
	channelID, timestamp, err := api.PostMessage("tec_growlio_dev", message, params)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	log.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
}

func BytesToShiftJIS(b []byte) (string, error) {
	return transformEncoding(bytes.NewReader(b), japanese.ShiftJIS.NewEncoder())
}

func transformEncoding( rawReader io.Reader, trans transform.Transformer) (string, error) {
	ret, err := ioutil.ReadAll(transform.NewReader(rawReader, trans))
	if err == nil {
		return string(ret), nil
	} else {
		return "", err
	}
}

func getNowString() (string){
	now := time.Now()
	nowUTC := now.UTC()
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	nowJST := nowUTC.In(jst)
	const layout2 = "2006-01-02 15:04"
	return nowJST.Format(layout2)
}