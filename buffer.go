package bytepool

import (
	stdbytes "bytes"
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
