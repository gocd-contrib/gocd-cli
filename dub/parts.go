package dub

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type Part interface {
	Build(*multipart.Writer) error
}

type FieldPart struct {
	param, value string
}

func (f *FieldPart) Build(w *multipart.Writer) error {
	return w.WriteField(f.param, f.value)
}

func NewFieldPart(param, value string) *FieldPart {
	return &FieldPart{param: param, value: value}
}

type FilePart struct {
	param, filepath string
}

func (f *FilePart) Build(w *multipart.Writer) (err error) {
	var file *os.File
	var part io.Writer

	if file, err = os.Open(f.filepath); err != nil {
		return
	}

	defer file.Close()

	if part, err = w.CreateFormFile(f.param, filepath.Base(f.filepath)); err != nil {
		return
	}

	_, err = io.Copy(part, file)
	return
}

func NewFilePart(param, path string) *FilePart {
	return &FilePart{param: param, filepath: path}
}

type StreamPart struct {
	param, filename string
	data            io.Reader
}

func (s *StreamPart) Build(w *multipart.Writer) (err error) {
	var part io.Writer

	if c, ok := s.data.(io.Closer); ok {
		defer c.Close()
	}

	if part, err = w.CreateFormFile(s.param, s.filename); err != nil {
		return
	}

	_, err = io.Copy(part, s.data)
	return
}

func NewStreamPart(param, filename string, data io.Reader) *StreamPart {
	return &StreamPart{param: param, filename: filename, data: data}
}
