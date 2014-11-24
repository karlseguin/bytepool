package bytepool

import (
	"encoding/binary"
	"io"
)

type bytes interface {
	write(b []byte) (bytes, int, error)
	writeByte(b byte) (bytes, error)
	readNFrom(n int64, r io.Reader) (bytes, int64, error)

	Read(b []byte) (int, error)
	Bytes() []byte
	String() string
	Len() int
}

type Bytes struct {
	bytes
	pool    *Pool
	fixed   *fixed
	scratch []byte
	enc     binary.ByteOrder
}

func NewBytes(capacity int) *Bytes {
	return NewEndianBytes(capacity, binary.BigEndian)
}

func NewEndianBytes(capacity int, enc binary.ByteOrder) *Bytes {
	return newPooled(nil, capacity, enc)
}

func newPooled(pool *Pool, capacity int, enc binary.ByteOrder) *Bytes {
	b := &Bytes{
		enc:  enc,
		pool: pool,
		fixed: &fixed{
			capacity: capacity,
			bytes:    make([]byte, capacity),
		},
		scratch: make([]byte, 8),
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

func (b *Bytes) PutUint16(n uint16) {
	b.enc.PutUint16(b.scratch, n)
	b.bytes, _, _ = b.write(b.scratch[:2])
}

func (b *Bytes) PutUint32(n uint32) {
	b.enc.PutUint32(b.scratch, n)
	b.bytes, _, _ = b.write(b.scratch[:4])
}

func (b *Bytes) PutUint64(n uint64) {
	b.enc.PutUint64(b.scratch, n)
	b.bytes, _, _ = b.write(b.scratch[:8])
}

// Write a string
func (b *Bytes) WriteString(str string) (int, error) {
	return b.Write([]byte(str))
}

// Read from the io.Reader
func (b *Bytes) ReadFrom(r io.Reader) (n int64, err error) {
	return b.ReadNFrom(0, r)
}

// Read N bytes from the io.Reader
func (b *Bytes) ReadNFrom(n int64, r io.Reader) (m int64, err error) {
	b.bytes, m, err = b.readNFrom(n, r)
	return m, err
}

// Reset the object without releasing it
func (b *Bytes) Reset() {
	b.fixed.length, b.fixed.r = 0, 0
	b.bytes = b.fixed
}

// Release the item back into the pool
func (b *Bytes) Release() {
	if b.pool != nil {
		b.Reset()
		b.pool.list <- b
	}
}

// Alias for Release
func (b *Bytes) Close() error {
	b.Release()
	return nil
}
