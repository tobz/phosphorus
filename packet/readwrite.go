package packet

import "io"

type Reader struct {
	conn io.Reader
}

func NewReader(conn io.Reader) *Reader {
	return &Reader{conn}
}

func (r *Reader) Next() (*Inbound, error) {
	// ...
}

func (r *Reader) Stop() {}

type Writer struct {
	conn io.Writer
}

func NewWriter(conn io.Writer) *Writer {
	return &Writer{conn}
}

func (r *Writer) Write(p Packet) error {
	// ...
}

func (r *Writer) Stop() {}
