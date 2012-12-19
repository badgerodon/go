package main

import (
	"compress/flate"
	"io"
)

type (
	compressedRWC struct {
		reader io.ReadCloser
		writer *flate.Writer
	}
)

func NewCompressedReadWriteCloser(rwc io.ReadWriteCloser) io.ReadWriteCloser {
	w, _ := flate.NewWriter(rwc, 1)
	return &compressedRWC{
		flate.NewReader(rwc),
		w,
	}
}
// Reader
func (this *compressedRWC) Read(p []byte) (n int, err error) {
	return this.reader.Read(p)
}
// Writer
func (this *compressedRWC) Write(p []byte) (n int, err error) {
	n, err = this.writer.Write(p)
	this.writer.Flush()
	return n, err
}
// Closer
func (this *compressedRWC) Close() error {
	e1 := this.reader.Close()
	e2 := this.writer.Close()
	if e1 == nil {
		return e2
	}
	return e1
}