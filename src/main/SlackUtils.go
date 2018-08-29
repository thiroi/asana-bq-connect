package main

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
"github.com/nlopes/slack"
"google.golang.org/appengine/urlfetch"
	"golang.org/x/net/context"
)

func SlackSender(w http.ResponseWriter, ctx *http.Request) {
	fmt.Fprint(w, "Hello, world222!")
	//sendSlack(ctx)
	log.Println("slack done")
}



func sendNlope(ctx context.Context, message string){
	log.Println("Start sending with nlope")
	slack.SetHTTPClient(urlfetch.Client(ctx))
	api := slack.New(config.Slack.Token)
	params := slack.PostMessageParameters{}
	//attachment := slack.Attachment{
	//	Pretext: "",
	//	Text:    message,
	//	// Uncomment the following part to send a field too
	//	/*
	//		Fields: []slack.AttachmentField{
	//			slack.AttachmentField{
	//				Title: "a",
	//				Value: "no",
	//			},
	//		},
	//	*/
	//}
	//params.Attachments = []slack.Attachment{attachment}
	params.IconEmoji = ":thx:"
	params.Username = "ザキヤマ"
	channelID, timestamp, err := api.PostMessage("times_hiroi", "========== \n" + message + "==========", params)
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