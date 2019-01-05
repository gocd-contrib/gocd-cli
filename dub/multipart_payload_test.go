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
