package dub

import (
	"bytes"
	"io"
	"mime/multipart"
	"testing"
)

func TestAllocPayload(t *testing.T) {
	as := asserts(t)

	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)
	w.SetBoundary("testbound")

	p := NewAllocPayload(b)

	// should not be ready to read yet before assembly
	as.not(p.Ready())
	d, err := p.Read(make([]byte, 1))
	as.err("EOF", err)
	as.eq(0, d)

	as.ok(p.DoAssemble(w, []Part{NewFieldPart("foo", "bar")}))
	as.is(p.Ready())
	as.is(p.Len() > 0)

	d, err = p.Read(make([]byte, 1))
	as.ok(err)
	as.eq(1, d)
}

func TestPipedPayload(t *testing.T) {
	as := asserts(t)
	pr, pw := io.Pipe()
	w := multipart.NewWriter(pw)
	w.SetBoundary("testbound")

	p := NewPipedPayload(pr, pw)

	// should not be ready to read yet before assembly
	as.not(p.Ready())
	d, err := p.Read(make([]byte, 1))
	as.err("Pipe is not ready to read()", err)
	as.eq(0, d)

	as.ok(p.DoAssemble(w, []Part{NewFieldPart("foo", "bar")}))
	as.is(p.Ready())
	as.eq(-1, p.Len()) // never knows the length of content

	d, err = p.Read(make([]byte, 1))
	as.ok(err)
	as.eq(1, d)
}

func TestWiretapPayload(t *testing.T) {
	as := asserts(t)

	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)
	w.SetBoundary("testbound")

	carbonCopy := &bytes.Buffer{}

	p := NewAllocPayload(b)
	wt := NewWireTapPayload(p, func(data []byte) error {
		_, err := carbonCopy.Write(data)
		return err
	})

	// All methods should pass through to the wrapped payload
	as.not(wt.Ready())
	d, err := wt.Read(make([]byte, 1))
	as.err("EOF", err) // errors pass-thru
	as.eq(0, d)

	as.ok(p.DoAssemble(w, []Part{NewFieldPart("foo", "bar")})) // assemble the inner payload
	as.is(wt.Ready())
	as.is(wt.Len() > 0)

	original := make([]byte, wt.Len())
	d, err = wt.Read(original)
	as.ok(err)

	as.eq(string(original), carbonCopy.String()) // successfully tapped!
	as.eq(len(original), d)
}
