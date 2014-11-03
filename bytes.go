package bytepool

import (
	"io"
)

type bytes interface {
	write(b []byte) (bytes, int, error)
	writeByte(b byte) (bytes, error)
	readFrom(r io.Reader) (bytes, int64, error)

	Read(b []byte) (int, error)
	Bytes() []byte
	String() string
	Len() int
}

type Bytes struct {
	bytes
	pool  *Pool
	fixed *fixed
}

func NewBytes(capacity int) *Bytes {
	return newPooled(nil, capacity)
}

func newPooled(pool *Pool, capacity int) *Bytes {
	b := &Bytes{
		pool: pool,
		fixed: &fixed{
			capacity: capacity,
			bytes:    make([]byte, capacity),
		},
	}
	if pool != nil {
		b.fixed.onExpand = pool.onExpand
	}
	b.bytes = b.fixed
	return b
}

// Write the bytes
func (b *Bytes) Write(data []byte) (n int, err error) {
	b.bytes, n, err = b.write(data)
	return n, err
}

// Write a byte
func (b *Bytes) WriteByte(d byte) (err error) {
	b.bytes, err = b.writeByte(d)
	return err
}

// Write a string
func (b *Bytes) WriteString(str string) (int, error) {
	return b.Write([]byte(str))
}

// Read from the io.Reader
func (b *Bytes) ReadFrom(r io.Reader) (n int64, err error) {
	b.bytes, n, err = b.readFrom(r)
	return n, err
}

// Release the item back into the pool
func (b *Bytes) Release() {
	if b.pool != nil {
		b.fixed.length = 0
		b.fixed.r = 0
		b.bytes = b.fixed
		b.pool.list <- b
	}
}

// Alias for Release
func (b *Bytes) Close() error {
	b.Release()
	return nil
}
