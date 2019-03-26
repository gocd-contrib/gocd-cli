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

func NewWireTapPayload(wrapped MultipartPayload, onRead func([]byte) error) MultipartPayload {
	return &wiretapPayload{delegate: wrapped, onRead: onRead}
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
			defer func() { close(ec) }()
			defer p.pw.Close()

			if len(parts) > 0 {
				for _, p := range parts {
					if err := p.Build(w); err != nil {
						ec <- err

						return
					}
				}
			}

			if err := w.Close(); err != nil {
				ec <- err

				return
			}

			ec <- nil
		}(p.asmErrCh)

		p.ready = true
	})

	// Do not return the value received through the error channel
	// or this will block on w.Write(). Assembly errors will be
	// handled during Read(), which happens concurrently.
	return nil
}

// Implements a wrapping payload that is similar in concept to io.TeeReader.
// It wraps another MultipartPayload and has a function member, `onRead()`,
// that will receive the bytes read. Note that any error returned by onRead()
// will be reported as a `Read()` error. `onRead()` can modify the []byte data
// before it returns from `Read()`.
type wiretapPayload struct {
	delegate MultipartPayload
	onRead   func([]byte) error
}

// Exposes the `[]byte` read to an `onRead()` function, where it can be spied
// on or modified in-flight. Any error returned by onRead() will be reported
// as a `Read()` error.
func (p *wiretapPayload) Read(b []byte) (n int, err error) {
	if n, err = p.delegate.Read(b); err == nil {
		if p.onRead != nil && n > 0 {
			err = p.onRead(b[:n])
		}
	}
	return
}

func (p *wiretapPayload) Close() error {
	return p.delegate.Close()
}

func (p *wiretapPayload) Len() int {
	return p.delegate.Len()
}

func (p *wiretapPayload) Ready() bool {
	return p.delegate.Ready()
}

func (p *wiretapPayload) DoAssemble(w *multipart.Writer, parts []Part) error {
	return p.delegate.DoAssemble(w, parts)
}
