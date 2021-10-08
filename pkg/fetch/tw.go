package fetch

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

func getTwHttpClient() *http.Client {
	config := oauth1.NewConfig(conf.TwConsumerKey, conf.TwConsumerSecret)
	token := oauth1.NewToken(conf.TwTokenKey, conf.TwTokenSecret)
	return config.Client(oauth1.NoContext, token)
}

func getTw() *twitter.Client {
	httpClient := getTwHttpClient()
	return twitter.NewClient(httpClient)
}

func downloadImages(imageUrls []string) ([]io.ReadCloser, string, error) {
	dir, err := ioutil.TempDir("/tmp", "tgchan2tw")
	if err != nil {
		return nil, "", err
	}

	var images []io.ReadCloser
	for i := range imageUrls {
		url := imageUrls[i]
		res, err := http.Get(url)
		if err != nil {
			return nil, "", err
		}

		f, err := ioutil.TempFile(dir, "image")
		if err != nil {
			return nil, "", err
		}

		_, err = bufio.NewReader(res.Body).WriteTo(f)
		if err != nil {
			return nil, "", err
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			return nil, "", err
		}

		images = append(images, f)
	}

	return images, dir, nil
}

func tweet(body string, imageUrls []string, replyTo int64) (i int64, e error) {
	images, dir, err := downloadImages(imageUrls)
	if err != nil {
		return 0, err
	}

	defer func(path string) {
		err := os.RemoveAll(path)
		if err != nil {
			e = err
		}
	}(dir)

	var mediaIds []int64
	for i := range images {
		image := images[i]
		id, err := tweetImage(image)
		if err != nil {
			return 0, err
		}

		mediaIds = append(mediaIds, id)
	}

	t, err := tweetText(body, mediaIds, replyTo)
	if err != nil {
		return 0, err
	}

	return t.ID, nil
}

func tweetText(body string, mediaIds []int64, replyTo int64) (*twitter.Tweet, error) {
	tw := getTw()
	t, _, err := tw.Statuses.Update(body, &twitter.StatusUpdateParams{
		InReplyToStatusID: replyTo,
		MediaIds:          mediaIds,
	})
	if err != nil {
		return nil, err
	}

	return t, nil
}

func tweetImage(image io.Reader) (int64, error) {
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
