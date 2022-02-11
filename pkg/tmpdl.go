// Copyright 2021-2022 myl7
// SPDX-License-Identifier: Apache-2.0

package pkg

import (
	"bufio"
	"github.com/myl7/tgchan2tw/pkg/cfg"
	"io"
	"io/ioutil"
	"net/http"
)

func tmpDl(urls []string) ([]io.ReadCloser, string) {
	dir, err := ioutil.TempDir(cfg.Cfg.TmpDir, "dl")
	if err != nil {
		panic(err)
	}

	var files []io.ReadCloser
	for i := range urls {
		s := urls[i]
		res, err := http.Get(s)
		if err != nil {
			panic(err)
		}

		f, err := ioutil.TempFile(dir, "image")
		if err != nil {
			panic(err)
		}

		_, err = bufio.NewReader(res.Body).WriteTo(f)
		if err != nil {
			panic(err)
		}

		_, err = f.Seek(0, 0)
		if err != nil {
			panic(err)
		}

		files = append(files, f)
	}

	return files, dir
}
