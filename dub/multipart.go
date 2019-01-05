package dub

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
)

type Multipart struct {
	Parts []Part
	MultipartPayload

	w *multipart.Writer
}

func (m *Multipart) Read(p []byte) (int, error) {
	if !m.MultipartPayload.Ready() {
		return 0, errors.New("Multipart stream is not ready to read(); call Multipart.Assemble() first")
	}

	return m.MultipartPayload.Read(p)
}

func (m *Multipart) Close() error {
	return m.MultipartPayload.Close()
}

func (m *Multipart) ContentType() string {
	return m.w.FormDataContentType()
}

func (m *Multipart) AddField(param, value string) *Multipart {
	m.Parts = append(m.Parts, NewFieldPart(param, value))
	return m
}

func (m *Multipart) AddFile(param, path string) *Multipart {
	m.Parts = append(m.Parts, NewFilePart(param, path))
	return m
}

func (m *Multipart) AddFileStream(param, filename string, data io.Reader) *Multipart {
	m.Parts = append(m.Parts, NewStreamPart(param, filename, data))
	return m
}

func (m *Multipart) Len() int {
	return m.MultipartPayload.Len()
}

func (m *Multipart) Assemble() error {
	if m.MultipartPayload.Ready() {
		return errors.New("This multipart stream has already been assembled")
	}

	return m.MultipartPayload.DoAssemble(m.w, m.Parts)
}

func NewPipedMultipart() *Multipart {
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)

	return &Multipart{
		w:                w,
		MultipartPayload: NewPipedPayload(pr, pw),
	}
}

func NewAllocMultipart() *Multipart {
	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)

	return &Multipart{
		w:                w,
		MultipartPayload: NewAllocPayload(buf),
	}
}
