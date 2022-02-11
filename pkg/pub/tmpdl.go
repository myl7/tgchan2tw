package pub

import (
	"bufio"
	"github.com/myl7/tgchan2tw/pkg/cfg"
	"io"
	"io/ioutil"
	"net/http"
)

func tmpDl(urls []string) ([]io.ReadCloser, string, error) {
	dir, err := ioutil.TempDir(cfg.TmpDir, "dl")
	if err != nil {
		return nil, "", err
	}

	var files []io.ReadCloser
	for i := range urls {
		s := urls[i]
		res, err := http.Get(s)
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

		files = append(files, f)
	}

	return files, dir, nil
}
