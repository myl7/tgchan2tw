// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package tw

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/myl7/tgchan2tw/pkg/cfg"
	"github.com/myl7/tgchan2tw/pkg/mdl"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
)

type TweetMsg struct {
	ItemId    int
	Body      string
	ImageUrls []string
	ReplyTo   int64
}

func Tweet(msg *mdl.Msg, images []io.ReadCloser, replyTo int64) []int64 {
	var mediaIds []int64
	for i := range images {
		image := images[i]
		id := tweetImage(image)

		mediaIds = append(mediaIds, id)
	}

	log.Printf("uploaded %d tw images\n", len(images))

	ts := tweetText(msg.Body, mediaIds, replyTo)

	var tids []int64
	for i := range ts {
		tids = append(tids, ts[i].ID)
	}
	return tids
}

func getTwHttpClient() *http.Client {
	config := oauth1.NewConfig(cfg.Cfg.TwConsumerKey, cfg.Cfg.TwConsumerSecret)
	token := oauth1.NewToken(cfg.Cfg.TwTokenKey, cfg.Cfg.TwTokenSecret)
	return config.Client(oauth1.NoContext, token)
}

func getTw() *twitter.Client {
	httpClient := getTwHttpClient()
	return twitter.NewClient(httpClient)
}

func tweetText(body string, mediaIds []int64, replyTo int64) []*twitter.Tweet {
	tw := getTw()

	bodies, err := splitTweetBody(body)
	if err != nil {
		panic(err)
	}

	var ts []*twitter.Tweet
	var t *twitter.Tweet
	for i := range bodies {
		params := &twitter.StatusUpdateParams{
			InReplyToStatusID: replyTo,
			MediaIds:          mediaIds,
		}
		if t != nil {
			params.InReplyToStatusID = t.ID
		}
		if len(mediaIds) > 0 {
			mediaIds = nil
		}

		t, _, err = tw.Statuses.Update(bodies[i], params)
		if err != nil {
			panic(err)
		}

		ts = append(ts, t)
	}

	return ts
}

func tweetImage(image io.ReadCloser) int64 {
	defer func(image io.ReadCloser) {
		_ = image.Close()
	}(image)

	client := getTwHttpClient()
	url := "https://upload.twitter.com/1.1/media/upload.json"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := w.CreateFormFile("media", "")
	if err != nil {
		panic(err)
	}

	_, err = bufio.NewReader(image).WriteTo(f)
	if err != nil {
		panic(err)
	}

	err = w.Close()
	if err != nil {
		panic(err)
	}

	res, err := client.Post(url, w.FormDataContentType(), &b)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	var r struct {
		MediaId int64 `json:"media_id"`
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(err)
	}

	return r.MediaId
}
