package dub

import (
	"io/ioutil"
	"strings"
	"testing"
)

func testMultipart() *Multipart {
	// might be easier to debug NewAllocMultipart() if anything goes wrong
	// but could have used NewPipedMultipart() as it might be used more
	m := NewAllocMultipart()
	m.w.SetBoundary("testbound")
	return m
}

func TestAddField(t *testing.T) {
	as := asserts(t)
	m := testMultipart()

	m.AddField("foo", "bar")
	as.eq(1, len(m.Parts))

	_, ok := m.Parts[0].(*FieldPart)
	as.is(ok)
	as.ok(m.Assemble())
	as.is(m.Ready())

	b, err := ioutil.ReadAll(m)
	as.ok(err)
	as.eq(reqfmt(
		`--testbound`,
		`Content-Disposition: form-data; name="foo"`,
		``,
		`bar`,
		`--testbound--`,
		``,
	), string(b))
}

func TestAddFile(t *testing.T) {
	as := asserts(t)
	m := testMultipart()

	m.AddFile("files[]", "testdata/file-01.txt")
	as.eq(1, len(m.Parts))

	_, ok := m.Parts[0].(*FilePart)
	as.is(ok)
	as.ok(m.Assemble())
	as.is(m.Ready())

	b, err := ioutil.ReadAll(m)
	as.ok(err)
	as.eq(reqfmt(
		`--testbound`,
		`Content-Disposition: form-data; name="files[]"; filename="file-01.txt"`,
		`Content-Type: application/octet-stream`,
		``,
		"hello! this is some content.\n",
		`--testbound--`,
		``,
	), string(b))
}

func TestAddFileStream(t *testing.T) {
	as := asserts(t)
	m := testMultipart()

	m.AddFileStream("files[]", "arbitrary-data.txt", strings.NewReader("streamed content"))
	as.eq(1, len(m.Parts))

	_, ok := m.Parts[0].(*StreamPart)
	as.is(ok)
	as.ok(m.Assemble())
	as.is(m.Ready())

	b, err := ioutil.ReadAll(m)
	as.ok(err)
	as.eq(reqfmt(
		`--testbound`,
		`Content-Disposition: form-data; name="files[]"; filename="arbitrary-data.txt"`,
		`Content-Type: application/octet-stream`,
		``,
		`streamed content`,
		`--testbound--`,
		``,
	), string(b))
}

func TestPipedMultipartComposition(t *testing.T) {
	as := asserts(t)
	m := NewPipedMultipart()
	m.w.SetBoundary("testbound")

	m.AddField("foo", "bar").
		AddFile("files[]", "testdata/file-01.txt").
		AddFileStream("files[]", "arbitrary-data.txt", strings.NewReader("streamed content"))

	as.eq(3, len(m.Parts))
	as.ok(m.Assemble())

	b, err := ioutil.ReadAll(m)
	as.ok(err)
	as.eq(reqfmt(
		`--testbound`,
		`Content-Disposition: form-data; name="foo"`,
		``,
		`bar`,
		`--testbound`,
		`Content-Disposition: form-data; name="files[]"; filename="file-01.txt"`,
		`Content-Type: application/octet-stream`,
		``,
		"hello! this is some content.\n",
		`--testbound`,
		`Content-Disposition: form-data; name="files[]"; filename="arbitrary-data.txt"`,
		`Content-Type: application/octet-stream`,
		``,
		`streamed content`,
		`--testbound--`,
		``,
	), string(b))
}

func TestAllocMultipartComposition(t *testing.T) {
	as := asserts(t)
	m := NewAllocMultipart() // be explicit here in case we decide to change testMultipart()
	m.w.SetBoundary("testbound")

	m.AddField("foo", "bar").
		AddFile("files[]", "testdata/file-01.txt").
		AddFileStream("files[]", "arbitrary-data.txt", strings.NewReader("streamed content"))

	as.eq(3, len(m.Parts))
	as.ok(m.Assemble())

	b, err := ioutil.ReadAll(m)
	as.ok(err)
	as.eq(reqfmt(
		`--testbound`,
		`Content-Disposition: form-data; name="foo"`,
		``,
		`bar`,
		`--testbound`,
		`Content-Disposition: form-data; name="files[]"; filename="file-01.txt"`,
		`Content-Type: application/octet-stream`,
		``,
		"hello! this is some content.\n",
		`--testbound`,
		`Content-Disposition: form-data; name="files[]"; filename="arbitrary-data.txt"`,
		`Content-Type: application/octet-stream`,
		``,
		`streamed content`,
		`--testbound--`,
		``,
	), string(b))
}

func TestAllocMultipartCalculatesContentLength(t *testing.T) {
	as := asserts(t)
	m := NewAllocMultipart() // be explicit here in case we decide to change testMultipart()
	m.w.SetBoundary("testbound")

	as.ok(m.AddField("foo", "bar").Assemble())
	as.eq(79, m.Len()) // show that it calculates before reading

	b, err := ioutil.ReadAll(m)
	as.ok(err)
	as.eq(79, len(b))
	as.eq(0, m.Len()) // once consumed, length drops to zero
}

func TestPipedMultipartReportsUnknownContentLength(t *testing.T) {
	as := asserts(t)
	m := NewPipedMultipart()
	m.w.SetBoundary("testbound")

	as.ok(m.AddField("foo", "bar").Assemble())
	as.eq(-1, m.Len())

	b, err := ioutil.ReadAll(m)
	as.ok(err)
	as.eq(79, len(b))
	as.eq(-1, m.Len()) // even after reading, does not calculate length
}

func TestMultipartOnlyReadsAfterAssembly(t *testing.T) {
	as := asserts(t)
	m := testMultipart()

	m.AddField("foo", "bar")
	as.not(m.Ready())

	_, err := m.Read(make([]byte, 1))
	as.err("Multipart stream is not ready to read(); call Multipart.Assemble() first", err)
	as.ok(m.Assemble())

	b, e := m.Read(make([]byte, 1))
	as.ok(e)
	as.eq(1, b)

	as.err("This multipart stream has already been assembled", m.Assemble())
}
