package utils

import (
	"bufio"
	"io"
	"os"
	"path"
	"strings"

	humanize "github.com/dustin/go-humanize"
	"github.com/gocd-contrib/gocd-cli/dub"
)

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
			return InspectError(err, `creating tempfile %q for downloading %q`, tmpfile, name)
		}

		defer file.Close() // ensure we close the file handle even if we abort early

		if err = res.OnProgress(downloadProgress).Consume(func(body io.Reader) error {
			w := bufio.NewWriter(file)
			_, err := io.Copy(w, body)

			if err != nil {
				return InspectError(err, `writing downloaded data to file %q`, file.Name())
			}

			return w.Flush()
		}); err != nil {
			return InspectError(err, `updating file download progress for %q`, url)
		}

		// explicitly close before rename or it may fail during rename
		if err = file.Close(); err != nil {
			return InspectError(err, `closing file %q`, file.Name())
		}

		Echofln("")

		return InspectError(os.Rename(tmpfile, filepath), `renaming tmplfile %q to %q`, tmpfile, filepath)
	})

	return filepath, err
}
