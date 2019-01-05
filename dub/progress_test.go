package dub

import (
	"io"
	"net/http"
	"testing"
)

func TestProgressDuringUpload(t *testing.T) {
	as := asserts(t)

	c := testCl(func(req *http.Request) (*http.Response, error) {
		if err := readOneByteAtATime(req.Body); err != nil {
			return nil, err
		}
		return okResp(), nil
	})

	var result []int64

	c.Post("http://test").
		DataString("abc").
		OnProgress(func(p *Progress) error {
			result = append(result, p.Current)
			return nil
		}).Do(ignoreResponse)

	as.deepEqI64([]int64{1, 2, 3}, result)
}

func TestProgressDuringDownload(t *testing.T) {
	as := asserts(t)

	c := testCl(func(req *http.Request) (*http.Response, error) {
		return resp(200, "abc"), nil
	})

	var result []int64

	c.Get("http://test").
		Do(func(res *Response) error {
			return res.OnProgress(func(p *Progress) error {
				result = append(result, p.Current)
				return nil
			}).Consume(readOneByteAtATime)
		})

	as.deepEqI64([]int64{1, 2, 3}, result)
}

func readOneByteAtATime(data io.Reader) error {
	for {
		if _, err := data.Read(make([]byte, 1)); err != nil {
			if io.EOF == err {
				break
			} else {
				return err
			}
		}
	}
	return nil
}
