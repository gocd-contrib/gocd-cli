package dub

import (
	"mime/multipart"
	"strings"
	"testing"
)

func TestFieldPart(t *testing.T) {
	as := asserts(t)
	p := NewFieldPart("foo", "bar")

	b := &strings.Builder{}
	w := multipart.NewWriter(b)
	w.SetBoundary("testbound")

	as.ok(p.Build(w))
	as.eq(reqfmt(
		`--testbound`,
		`Content-Disposition: form-data; name="foo"`,
		``,
		`bar`,
	), b.String())
}

func TestFilePart(t *testing.T) {
	as := asserts(t)
	p := NewFilePart("files[]", "testdata/file-01.txt")

	b := &strings.Builder{}
	w := multipart.NewWriter(b)
	w.SetBoundary("testbound")

	as.ok(p.Build(w))
	as.eq(reqfmt(
		`--testbound`,
		`Content-Disposition: form-data; name="files[]"; filename="file-01.txt"`,
		`Content-Type: application/octet-stream`,
		``,
		"hello! this is some content.\n",
	), b.String())
}

func TestStreamPart(t *testing.T) {
	as := asserts(t)
	p := NewStreamPart("files[]", "arbitrary-data.txt", strings.NewReader("streamed content"))

	b := &strings.Builder{}
	w := multipart.NewWriter(b)
	w.SetBoundary("testbound")

	as.ok(p.Build(w))
	as.eq(reqfmt(
		`--testbound`,
		`Content-Disposition: form-data; name="files[]"; filename="arbitrary-data.txt"`,
		`Content-Type: application/octet-stream`,
		``,
		`streamed content`,
	), b.String())
}
