package pub

import (
	"bufio"
	"bytes"
	"encoding/json"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/myl7/tgchan2tw/pkg/conf"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

type TweetMsg struct {
	ItemId    int
	Body      string
	ImageUrls []string
	ReplyTo   int64
}

func Tweet(msg TweetMsg) ([]int64, error) {
	images, tmpDir, err := tmpDl(msg.ImageUrls)
	if err != nil {
		return nil, err
	}

	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tmpDir)

	var mediaIds []int64
	for i := range images {
		image := images[i]
		id, err := tweetImage(image)
		if err != nil {
			return nil, err
		}

		mediaIds = append(mediaIds, id)
	}

	ts, err := tweetText(msg.Body, mediaIds, msg.ReplyTo)
	if err != nil {
		return nil, err
	}

	var tids []int64
	for i := range ts {
		tids = append(tids, ts[i].ID)
	}

	return tids, nil
}

func getTwHttpClient() *http.Client {
	config := oauth1.NewConfig(conf.TwConsumerKey, conf.TwConsumerSecret)
	token := oauth1.NewToken(conf.TwTokenKey, conf.TwTokenSecret)
	return config.Client(oauth1.NoContext, token)
}

func getTw() *twitter.Client {
	httpClient := getTwHttpClient()
	return twitter.NewClient(httpClient)
}

func tweetText(body string, mediaIds []int64, replyTo int64) ([]*twitter.Tweet, error) {
	tw := getTw()

	bodies, err := splitTweetBody(body)
	if err != nil {
		return nil, err
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
			return nil, err
		}

		ts = append(ts, t)
	}

	return ts, nil
}

func tweetImage(image io.ReadCloser) (int64, error) {
	defer func(image io.ReadCloser) {
		_ = image.Close()
	}(image)

	client := getTwHttpClient()
	url := "https://upload.twitter.com/1.1/media/upload.json"

	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := w.CreateFormFile("media", "")
	if err != nil {
		return 0, err
	}

	_, err = bufio.NewReader(image).WriteTo(f)
	if err != nil {
		return 0, err
	}

	err = w.Close()
	if err != nil {
		return 0, err
	}

	res, err := client.Post(url, w.FormDataContentType(), &b)
	if err != nil {
		return 0, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var r struct {
		MediaId int64 `json:"media_id"`
	}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return 0, err
	}

	return r.MediaId, nil
}
