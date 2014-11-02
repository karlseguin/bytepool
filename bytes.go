package bytepool

type bytes interface {
	write(b []byte) (bytes, int, error)
	writeByte(b byte) (bytes, error)

	Bytes() []byte
	String() string
	Len() int
}

type Bytes struct {
	bytes
	pool *Pool
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
	b.bytes = b.fixed
	return b
}

func (b *Bytes) Write(data []byte) (n int, err error) {
	b.bytes, n, err = b.bytes.write(data)
	return n, err
}

func (b *Bytes) WriteByte(d byte) (err error) {
	b.bytes, err = b.bytes.writeByte(d)
	return err
}

func (b *Bytes) WriteString(str string) (int, error) {
	return b.Write([]byte(str))
}

func (b *Bytes) Release() {
	if b.pool != nil {
		b.fixed.length = 0
		b.bytes = b.fixed
		b.pool.list <- b
	}
}
