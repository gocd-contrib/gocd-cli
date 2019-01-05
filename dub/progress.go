package dub

import (
	"io"
	"net/http"
)

type Progress struct {
	Current     int64
	Total       int64
	RawRequest  *http.Request
	RawResponse *http.Response

	onUpdate []ProgressHandler
}

func newProgress(total int64, handlers []ProgressHandler) *Progress {
	return &Progress{Total: total, onUpdate: handlers}
}

type progressWriter struct {
	progress *Progress
}

func (p *progressWriter) Write(b []byte) (int, error) {
	n := len(b)
	p.progress.Current += int64(n)

	if err := p.doUpdate(); err != nil {
		return 0, err
	}

	return n, nil
}

func (p *progressWriter) doUpdate() error {
	for _, h := range p.progress.onUpdate {
		if err := h(p.progress); err != nil {
			return err
		}
	}
	return nil
}

func newProgressWriter(p *Progress) io.Writer {
	return &progressWriter{progress: p}
}
