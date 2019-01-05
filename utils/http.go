package utils

import (
	"bufio"
	"io"
	"os"
	"path"
	"strings"

	"net"
	"net/http"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gocd-contrib/gocd-cli/dub"
)

func Http() *http.Client {
	tx := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout: 5 * time.Second,
	}

	return &http.Client{
		Transport: tx,
	}
}

func downloadProgress(dp *dub.Progress) error {
	Echof("\r%s", strings.Repeat(" ", 35))
	if dp.Total > -1 {
		Echof("\r  Fetched %s/%s (%.1f%%) complete", humanize.Bytes(uint64(dp.Current)), humanize.Bytes(uint64(dp.Total)), float64(dp.Current)/float64(dp.Total)*float64(100))
	} else {
		Echof("\r  Fetched %s complete", humanize.Bytes(uint64(dp.Current)))
	}
	return nil
}

func Wget(url string, name string, destFolder string) (filepath string, err error) {
	tmpfile := path.Join(destFolder, "_"+name+".partialdownload")
	filepath = path.Join(destFolder, name)

	Echofln("Downloading %s", url)

	err = dub.New().Get(url).Do(func(res *dub.Response) (err error) {
		var file *os.File

		if file, err = os.Create(tmpfile); err != nil {
			return err
		}

		defer file.Close()

		if err = res.OnProgress(downloadProgress).Consume(func(body io.Reader) error {
			w := bufio.NewWriter(file)
			_, err := io.Copy(w, body)

			if err == nil {
				err = w.Flush()
			}

			return err
		}); err != nil {
			return err
		}

		Echofln("")

		return os.Rename(tmpfile, filepath)
	})

	return filepath, err
}
