package bytepool

import (
	stdbytes "bytes"
	"io"
)

type fixed struct {
	length   int
	capacity int
	bytes    []byte
}

func (f *fixed) Len() int {
	return f.length
}

func (f *fixed) Bytes() []byte {
	return f.bytes[:f.length]
}

func (f *fixed) String() string {
	return string(f.Bytes())
}

func (f *fixed) write(data []byte) (bytes, int, error) {
	if l := len(data); f.hasSpace(l) == false {
		buf := f.toBuffer()
		n, err := buf.Write(data)
		return buf, n, err
	}
	n := copy(f.bytes[f.length:], data)
	f.length += n
	return f, n, nil
}

func (f *fixed) writeByte(data byte) (bytes, error) {
	if f.length == f.capacity {
		buf := f.toBuffer()
		err := buf.WriteByte(data)
		return buf, err
	}
	f.bytes[f.length] = data
	f.length++
	return f, nil
}

func (f *fixed) readFrom(reader io.Reader) (bytes, int64, error) {
	var read int64
	for {
		if f.full() {
			buf := f.toBuffer()
			n, err := buf.ReadFrom(reader)
			return buf, read + n, err
		}
		r, err := reader.Read(f.bytes[f.length:])
		read += int64(r)
		f.length += r
		if err == io.EOF {
			return f, read, nil
		}
		if err != nil {
			return f, read, err
		}
	}
}

func (f *fixed) toBuffer() *buffer {
	buf := &buffer{stdbytes.NewBuffer(f.bytes)}
	buf.Truncate(f.length)
	return buf
}

func (f *fixed) hasSpace(toAdd int) bool {
	return f.length+toAdd <= f.capacity
}

func (f *fixed) full() bool {
	return f.length == f.capacity
}
