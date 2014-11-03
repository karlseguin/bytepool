package bytepool

import (
	stdbytes "bytes"
	. "github.com/karlseguin/expect"
	"io"
	"testing"
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
