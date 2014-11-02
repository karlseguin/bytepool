package bytepool

import (
	stdbytes "bytes"
)

type fixed struct {
	length int
	capacity int
	bytes []byte
}

func (f *fixed) write(data []byte) (bytes, int, error) {
	if l := len(data); f.hasSpace(l) == false {
		buf := &buffer{stdbytes.NewBuffer(f.bytes)}
		buf.Truncate(f.length)
		n, err := buf.Write(data)
		return buf, n, err
	}
	n := copy(f.bytes[f.length:], data)
	f.length += n
	return f, n, nil
}

func (f *fixed) writeByte(data byte) (bytes, error) {
	if f.hasSpace(1) == false {
		buf := &buffer{stdbytes.NewBuffer(f.bytes)}
		buf.Truncate(f.length)
		err := buf.WriteByte(data)
		return buf, err
	}
	f.bytes[f.length] = data
	f.length++
	return f, nil
}

func (f *fixed) hasSpace(toAdd int) bool {
	return f.length + toAdd <= f.capacity
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
