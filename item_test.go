package bytepool

import (
	"bytes"
	"io"
	"testing"
	. "github.com/karlseguin/expect"
)

type ItemTests struct {}

func Test_Items(t *testing.T) {
	Expectify(new(ItemTests), t)
}

func (i *ItemTests) CanWriteAString() {
	item := NewItem(10, nil)
	item.WriteString("over ")
	item.WriteString("9000")
	Expect(item.String()).To.Equal("over 9000")
}

func (i *ItemTests) CanWriteAByteArray() {
	item := NewItem(60, nil)
	item.Write([]byte("the "))
	item.Write([]byte("spice "))
	item.Write([]byte("must "))
	item.Write([]byte("flow"))
	Expect(item.Bytes()).To.Equal([]byte("the spice must flow"))
}

func (i *ItemTests) WriteAByte() {
	item := NewItem(60, nil)
	item.Write([]byte("the "))
	item.WriteByte(byte('s'))
	item.WriteByte(byte('p'))
	Expect(item.Bytes()).To.Equal([]byte("the sp"))
}

func (i *ItemTests) DoesNotWriteAByteWhenFull() {
	item := NewItem(5, nil)
	item.Write([]byte("the "))
	item.WriteByte(byte('s'))
	item.WriteByte(byte('p'))
	Expect(item.Bytes()).To.Equal([]byte("the s"))
}

func (i *ItemTests) HAndlesReadingAnExactSize() {
	item := NewItem(5, nil)
	buffer := bytes.NewBufferString("12345")
	item.ReadFrom(buffer)

	Expect(item.String()).To.Equal("12345")
}

func (i *ItemTests) CanWriteFromAReader() {
	item := NewItem(60, nil)
	n, _ := item.ReadFrom(bytes.NewBuffer([]byte("I am in a reader")))
	Expect(item.Bytes()).To.Equal([]byte("I am in a reader"))
	Expect(int(n)).To.Equal(len("I am in a reader"))
}

func (i *ItemTests) WriteFromMultipleSources() {
	item := NewItem(100, nil)
	item.Write([]byte("start"))
	n, _ := item.ReadFrom(bytes.NewBuffer([]byte("I am in a reader")))
	item.WriteString("end")

	Expect(item.Bytes()).To.Equal([]byte("startI am in a readerend"))
	Expect(int(n)).To.Equal(len([]byte("I am in a reader")))
}

func (i *ItemTests) CanSetThePosition() {
	item := NewItem(100, nil)
	item.WriteString("hello world")
	item.Position(5)
	item.WriteString(".")
	Expect(item.String()).To.Equal("hello.")
}

func (i *ItemTests) CloseResetsTheLengthWhenAttachedToApool() {
	pool := New(1, 100)
	item := pool.Checkout()
	item.WriteString("hello world")
	item.Close()
	Expect(item.Len()).To.Equal(0)
	item.WriteString("hello")
	Expect(item.String()).To.Equal("hello")
}

func (i *ItemTests) CannotSetThePositionToANegativeValue() {
	item := NewItem(25, nil)
	item.WriteString("hello world")
	item.Position(-10)
	item.WriteString(".")
	Expect(item.String()).To.Equal("hello world.")
}

func (i *ItemTests) CannotSetThePositionBeyondTheLength() {
	item := NewItem(25, nil)
	item.WriteString("hello world")
	item.Position(30)
	item.WriteString(".")
	Expect(item.String()).To.Equal("hello world.")
}

func (i *ItemTests) TrimLastIfTrimsOnMatch() {
	item := NewItem(25, nil)
	item.WriteString("hello world.")
	r := item.TrimLastIf(byte('.'))
	Expect(r).To.Equal(true)
	Expect(item.String()).To.Equal("hello world")
}

func (i *ItemTests) TrimLastIfTrimsOnNoMatch() {
	item := NewItem(25, nil)
	item.WriteString("hello world.")
	r := item.TrimLastIf(byte(','))
	Expect(r).To.Equal(false)
	Expect(item.String()).To.Equal("hello world.")
}

func (i *ItemTests) TruncatesTheContentToTheLength() {
	item := NewItem(4, nil)
	item.WriteString("hello")
	Expect(item.String()).To.Equal("hell")
	item.WriteString("world")
	Expect(item.String()).To.Equal("hell")
}

func (i *ItemTests) CanReadIntoVariousSizedByteArray() {
	for size, expected := range map[int]string{3: "hel", 5: "hello", 7: "hello\x00\x00"} {
		item := NewItem(5, nil)
		item.WriteString("hello")
		target := make([]byte, size)
		item.Read(target)
		Expect(string(target)).To.Equal(expected)
	}
}

func (i *ItemTests) ReadDoesNotAutomaticallyRewind() {
	item := NewItem(5, nil)
	item.WriteString("hello")
	b := make([]byte, 5)

	n, err := item.Read(b[0:2])
	Expect(n, err).To.Equal(2, nil)
	Expect(string(b[0:2])).To.Equal("he")

	n, err = item.Read(b[2:])
	Expect(n, err).To.Equal(3, io.EOF)
	Expect(string(b[0:5])).To.Equal("hello")


	n, err = item.Read(b)
	Expect(n, err).To.Equal(0, io.EOF)
	Expect(string(b[0:5])).To.Equal("hello")
}

func (i *ItemTests) CloneDetachesTheObject() {
	item := NewItem(10, nil)
	item.WriteString("over")
	actual := item.Clone()
	item.Raw()[0] = '!'
	Expect(actual[0]).To.Equal(byte('o'))
}

func (i *ItemTests) ReturnsTheAvailableSpace() {
	item := NewItem(10, nil)
	Expect(item.Space()).To.Equal(10)
	item.WriteString("hello")
	Expect(item.Space()).To.Equal(5)
	item.WriteString("world")
	Expect(item.Space()).To.Equal(0)
}
