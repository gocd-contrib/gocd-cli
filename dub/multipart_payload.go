package dub

import (
	"errors"
	"io"
	"mime/multipart"
	"sync"
)

// Represents the constructed multipart form data payload as an io.ReadCloser.
// Implementations will be used to build a Multipart instance for a request body.
type MultipartPayload interface {
	io.ReadCloser
	Len() int
	Ready() bool
	DoAssemble(*multipart.Writer, []Part) error
}

func NewPipedPayload(pr *io.PipeReader, pw *io.PipeWriter) MultipartPayload {
	return &pipedPayload{pr: pr, pw: pw}
}

func NewAllocPayload(buffer io.Reader) MultipartPayload {
	return &allocPayload{r: buffer, ready: false}
}

// Implements a payload by allocating byte slice in memory for the entire payload.
// Payload assembly happens synchronously, so the final content length is calculable
// and will be set on the request. Limited to amount of memory allocable to the process.
type allocPayload struct {
	r     io.Reader
	once  sync.Once
	ready bool
}

func (p *allocPayload) Len() int {
	if b, ok := p.r.(lengther); ok {
		return b.Len()
	}
	return -1
}

func (p *allocPayload) Ready() bool {
	return p.ready
}

func (p *allocPayload) Read(b []byte) (int, error) {
	return p.r.Read(b)
}

func (p *allocPayload) Close() error {
	if c, ok := p.r.(io.Closer); ok {
		return c.Close()
	}
	return nil
}

func (p *allocPayload) DoAssemble(w *multipart.Writer, parts []Part) (err error) {
	p.once.Do(func() {
		if len(parts) > 0 {
			for _, p := range parts {
				if err = p.Build(w); err != nil {
					return
				}
			}
		}

		if err = w.Close(); err != nil {
			return
		}

		p.ready = true
	})
	return
}

// This payload type will assemble concurrently as it is read, so in theory is
// more efficient in that it does not require explicit allocation and pre-assembly.
// Data is read as it is built, by way of io.Pipe(). As this is a streaming payload,
// it does not have the memory limitations that allocPayload has, but is also unable
// to report content length.
type pipedPayload struct {
	pr       *io.PipeReader
	pw       *io.PipeWriter
	asmErrCh chan error
	once     sync.Once
	ready    bool
}

func (p *pipedPayload) Len() int {
	return -1
}

func (p *pipedPayload) Ready() bool {
	return p.ready
}

func (p *pipedPayload) Read(b []byte) (int, error) {
	if !p.ready { // prevent deadlock when write-end has yet to be written
		return 0, errors.New("Pipe is not ready to read()")
	}

	select {
	case err := <-p.asmErrCh:
		if err != nil {
			return 0, err
		}
		break
	default:
	}

	return p.pr.Read(b)
}

func (p *pipedPayload) Close() error {
	return p.pr.Close()
}

func (p *pipedPayload) DoAssemble(w *multipart.Writer, parts []Part) error {
	p.once.Do(func() {
		p.asmErrCh = make(chan error, 1)

		go func(ec chan<- error) {
			defer p.pw.Close()

			if len(parts) > 0 {
				for _, p := range parts {
					if err := p.Build(w); err != nil {
						ec <- err

						close(ec)
						return
					}
				}
			}

			if err := w.Close(); err != nil {
				ec <- err

				close(ec)
				return
			}

			ec <- nil
			close(ec)
		}(p.asmErrCh)

		p.ready = true
	})

	// Do not return the value received through the error channel
	// or this will block on w.Write(). Assembly errors will be
	// handled during Read(), which happens concurrently.
	return nil
}
