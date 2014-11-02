package bytepool

import (
	stdbytes "bytes"
	"io"
)

type buffer struct {
	*stdbytes.Buffer
}

func (b *buffer) write(data []byte) (bytes, int, error) {
	n, err := b.Write(data)
	return b, n, err
}

func (b *buffer) writeByte(data byte) (bytes, error) {
	err := b.WriteByte(data)
	return b, err
}

func (b *buffer) readFrom(r io.Reader) (bytes, int64, error) {
	n, err := b.ReadFrom(r)
	return b, n, err
}
