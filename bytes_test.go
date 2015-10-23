package bytepool

import (
	stdbytes "bytes"
	"encoding/binary"
	"io"
	"testing"

	. "github.com/karlseguin/expect"
)

type BytesTest struct{}

func Test_Bytes(t *testing.T) {
	Expectify(new(BytesTest), t)
}

func (_ BytesTest) WriteByte() {
	bytes := NewBytes(1)
	bytes.WriteByte('!')
	Expect(bytes.String()).To.Equal("!")
	bytes.WriteByte('?')
	Expect(bytes.String()).To.Equal("!?")
}

func (_ BytesTest) WriteWithinCapacity() {
	bytes := NewBytes(16)
	bytes.Write([]byte("it's over 9000!"))
	Expect(bytes.String()).To.Equal("it's over 9000!")
	Expect(bytes.Len()).To.Equal(15)
}

func (_ BytesTest) WriteAtCapacity() {
	bytes := NewBytes(16)
	bytes.Write([]byte("it's over 9000!!"))
	Expect(bytes.String()).To.Equal("it's over 9000!!")
	Expect(bytes.Len()).To.Equal(16)
}

func (_ BytesTest) WriteOverCapacity1() {
	bytes := NewBytes(16)
	bytes.Write([]byte("it's over 9000!!!"))
	Expect(bytes.String()).To.Equal("it's over 9000!!!")
	Expect(bytes.Len()).To.Equal(17)
}

func (_ BytesTest) WriteOverCapacity2() {
	bytes := NewBytes(16)
	bytes.Write([]byte("it's over 9000"))
	bytes.Write([]byte("!!!"))
	Expect(bytes.String()).To.Equal("it's over 9000!!!")
	Expect(bytes.Len()).To.Equal(17)
}

func (_ BytesTest) WriteOverCapacity3() {
	bytes := NewBytes(15)
	bytes.Write([]byte("it's over 9000"))
	bytes.Write([]byte("!!"))
	bytes.WriteString("!")
	bytes.WriteByte('.')
	Expect(bytes.String()).To.Equal("it's over 9000!!!.")
	Expect(bytes.Len()).To.Equal(18)
}

func (_ BytesTest) ReleasesWhenNoOverflow() {
	bytes := New(20, 1).Checkout()
	bytes.Write([]byte("it's over 9000!"))
	bytes.Release()
	Expect(bytes.String()).To.Equal("")
	Expect(bytes.Len()).To.Equal(0)
	Expect(cap(bytes.bytes.(*fixed).bytes)).To.Equal(20)
}

func (_ BytesTest) ReleasesWhenOverflow() {
	bytes := New(10, 1).Checkout()
	bytes.Write([]byte("it's over 9000!"))
	bytes.Release()
	Expect(bytes.String()).To.Equal("")
	Expect(bytes.Len()).To.Equal(0)
	Expect(cap(bytes.bytes.(*fixed).bytes)).To.Equal(10)
}

func (_ BytesTest) ReadFrom() {
	bytes := NewBytes(10)
	bytes.ReadFrom(stdbytes.NewBufferString("hello"))
	Expect(bytes.String()).To.Equal("hello")
	bytes.ReadFrom(stdbytes.NewBufferString("world"))
	Expect(bytes.String()).To.Equal("helloworld")
	bytes.ReadFrom(stdbytes.NewBufferString("how"))
	Expect(bytes.String()).To.Equal("helloworldhow")
	bytes.ReadFrom(stdbytes.NewBufferString("goes"))
	Expect(bytes.String()).To.Equal("helloworldhowgoes")
}

func (_ BytesTest) WritesTo() {
	bytes := NewBytes(10)
	bytes.WriteString("over 9000")
	buffer := new(stdbytes.Buffer)
	bytes.WriteTo(buffer)
	Expect(bytes.Len()).To.Equal(0)
	Expect(bytes.fixed.r).To.Equal(0)
	Expect(buffer.String()).To.Equal("over 9000")
}

func (_ BytesTest) ReadNFrom() {
	bytes := NewBytes(10)
	bytes.ReadNFrom(4, stdbytes.NewBufferString("hello"))
	Expect(bytes.String()).To.Equal("hell")
	bytes.ReadNFrom(4, stdbytes.NewBufferString("world"))
	Expect(bytes.String()).To.Equal("hellworl")
	bytes.ReadNFrom(6, stdbytes.NewBufferString("thisisfun"))
	Expect(bytes.String()).To.Equal("hellworlthisis")
	bytes.ReadNFrom(2, stdbytes.NewBufferString("go"))
	Expect(bytes.String()).To.Equal("hellworlthisisgo")
}

func (_ BytesTest) ReadNFromExact() {
	bytes := NewBytes(3)
	bytes.ReadNFrom(3, stdbytes.NewBufferString("hello"))
	Expect(bytes.String()).To.Equal("hel")
}

func (_ BytesTest) FullRead() {
	bytes := NewBytes(10)
	bytes.WriteString("hello!")
	data := make([]byte, 10)
	n, err := bytes.Read(data)
	Expect(n, err).To.Equal(6, io.EOF)
	Expect(string(data[:n])).To.Equal("hello!")
}

func (_ BytesTest) Partial() {
	bytes := NewBytes(10)
	bytes.WriteString("hello!")
	data := make([]byte, 4)
	n, err := bytes.Read(data)
	Expect(n, err).To.Equal(4, nil)
	Expect(string(data)).To.Equal("hell")
}

func (_ BytesTest) Reset() {
	bytes := NewBytes(10)
	bytes.WriteString("hello!")
	bytes.Reset()
	bytes.WriteString("spice")
	Expect(bytes.String()).To.Equal("spice")
}

func (_ BytesTest) ResetFromExpansion() {
	bytes := NewBytes(2)
	bytes.WriteString("hello!")
	_, is := bytes.bytes.(*buffer)
	Expect(is).To.Equal(true)
	bytes.Reset()
	_, is = bytes.bytes.(*fixed)
	Expect(is).To.Equal(true)
}

func (_ BytesTest) WriteBigEndian() {
	p := New(10, 1)
	b := p.Checkout()
	b.WriteUint64(2933)
	Expect(b.Bytes()).To.Equal([]byte{0, 0, 0, 0, 0, 0, 11, 117})
	b.WriteUint32(8484848)
	Expect(b.Bytes()).To.Equal([]byte{0, 0, 0, 0, 0, 0, 11, 117, 0, 129, 119, 240})
}

func (_ BytesTest) WriteLittleEndian() {
	p := NewEndian(10, 1, binary.LittleEndian)
	b := p.Checkout()
	b.WriteUint64(2933)
	Expect(b.Bytes()).To.Equal([]byte{117, 11, 0, 0, 0, 0, 0, 0})
	b.WriteUint32(8484848)
	Expect(b.Bytes()).To.Equal([]byte{117, 11, 0, 0, 0, 0, 0, 0, 240, 119, 129, 0})
}

func (_ BytesTest) ReadBigEndian() {
	p := New(12, 1)
	b := p.Checkout()
	b.WriteUint64(2933)
	b.WriteUint32(10)
	Expect(b.ReadUint64()).To.Equal(uint64(2933), nil)
	Expect(b.ReadUint32()).To.Equal(uint32(10), nil)
	b.WriteUint16(1234)
	Expect(b.ReadUint16()).To.Equal(uint16(1234), nil)
	b.WriteUint64(94994949)
	Expect(b.ReadUint64()).To.Equal(uint64(94994949), nil)
}

func (_ BytesTest) ReadIntsEOF() {
	p := New(4, 1)
	b := p.Checkout()
	Expect(b.ReadUint64()).To.Equal(uint64(0), io.EOF)
	Expect(b.ReadUint32()).To.Equal(uint32(0), io.EOF)
	Expect(b.ReadUint16()).To.Equal(uint16(0), io.EOF)
}

func (_ BytesTest) ReadIntsEOFBuffer() {
	p := New(4, 1)
	b := p.Checkout()
	b.WriteUint64(23)
	b.ReadUint64()
	Expect(b.ReadUint64()).To.Equal(uint64(0), io.EOF)
	Expect(b.ReadUint32()).To.Equal(uint32(0), io.EOF)
	Expect(b.ReadUint16()).To.Equal(uint16(0), io.EOF)
}

func (_ BytesTest) ReadAndWriteByte() {
	p := New(3, 1)
	b := p.Checkout()
	b.WriteByte(9)
	b.WriteByte(4)
	Expect(b.ReadByte()).To.Equal(byte(9), nil)
	Expect(b.ReadByte()).To.Equal(byte(4), nil)
	Expect(b.ReadByte()).To.Equal(byte(0), io.EOF)
}

func (_ BytesTest) ReadAndWriteByteForBuffer() {
	p := New(2, 1)
	b := p.Checkout()
	b.WriteByte(9)
	b.WriteByte(4)
	b.WriteByte(39)
	Expect(b.ReadByte()).To.Equal(byte(9), nil)
	Expect(b.ReadByte()).To.Equal(byte(4), nil)
	Expect(b.ReadByte()).To.Equal(byte(39), nil)
	Expect(b.ReadByte()).To.Equal(byte(0), io.EOF)
}

func (_ BytesTest) PositionFixed() {
	bytes := NewBytes(10)
	bytes.Position(6)
	bytes.WriteString("abc")
	Expect(bytes.Bytes()).To.Equal([]byte{0, 0, 0, 0, 0, 0, 97, 98, 99})
	bytes.Position(2)
	bytes.WriteByte(4)
	Expect(bytes.Bytes()).To.Equal([]byte{0, 0, 4})
}

func (_ BytesTest) NegativePositionFixed() {
	bytes := NewBytes(10)
	bytes.WriteString("abc")
	bytes.Position(2)
	bytes.WriteString("bd")
	Expect(bytes.Bytes()).To.Equal([]byte{97, 98, 98, 100})
	bytes.Position(1)
	bytes.WriteByte(4)
	Expect(bytes.Bytes()).To.Equal([]byte{97, 4})
}

func (_ BytesTest) PositionBuffer() {
	bytes := NewBytes(2)
	bytes.Position(4)
	Expect(bytes.Bytes()).To.Equal([]byte{0, 0, 0, 0})
	bytes.WriteString("12")
	Expect(bytes.Bytes()).To.Equal([]byte{0, 0, 0, 0, 49, 50})
	bytes.Position(10)
	Expect(bytes.Bytes()).To.Equal([]byte{0, 0, 0, 0, 49, 50, 0, 0, 0, 0})
	bytes.Position(5)
	Expect(bytes.Bytes()).To.Equal([]byte{0, 0, 0, 0, 49})
}

func (_ BytesTest) NegativePositionBuffer() {
	bytes := NewBytes(2)
	bytes.WriteString("abc")
	bytes.Position(2)
	bytes.WriteString("bd")
	Expect(bytes.Bytes()).To.Equal([]byte{97, 98, 98, 100})
	bytes.Position(1)
	bytes.WriteByte(4)
	Expect(bytes.Bytes()).To.Equal([]byte{97, 4})
}

func (_ BytesTest) CustomOnExpand() {
	expanded := 0
	bytes := NewBytes(7)
	bytes.SetOnExpand(func() { expanded++ })
	bytes.WriteString("hello")
	Expect(expanded).To.Equal(0)
	bytes.WriteString("world")
	Expect(expanded).To.Equal(1)
	bytes.WriteString("world")
	Expect(expanded).To.Equal(1)
}
