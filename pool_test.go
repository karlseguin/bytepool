package bytepool

import (
	. "github.com/karlseguin/expect"
	"reflect"
	"strings"
	"testing"
)

type PoolTests struct{}

func Test_Pool(t *testing.T) {
	Expectify(new(PoolTests), t)
}

func (_ PoolTests) EachItemIsOfASpecifiedSize() {
	p := New(9, 1)
	bytes := p.Checkout()
	Expect(cap(bytes.fixed.bytes)).To.Equal(9)
}

func (_ PoolTests) HasTheSpecifiedNumberOfItems() {
	p := New(23, 3)
	Expect(cap(p.list)).To.Equal(3)
}

func (_ PoolTests) DynamicallyCreatesAnItemWhenPoolIsEmpty() {
	p := New(23, 1)
	bytes1 := p.Checkout()
	bytes2 := p.Checkout()
	Expect(cap(bytes2.fixed.bytes)).To.Equal(23)
	Expect(bytes2.pool).To.Equal(nil)
	bytes1.Release()
	bytes2.Release()
	Expect(len(p.list)).To.Equal(1)
	Expect(p.Depleted()).To.Equal(int64(1))
}

func (_ PoolTests) ReleasesAnItemBackIntoThePool() {
	p := New(20, 1)
	bytes := p.Checkout()
	pointer := reflect.ValueOf(bytes).Pointer()
	bytes.Release()

	if reflect.ValueOf(p.Checkout()).Pointer() != pointer {
		Fail("Pool returned an unexected item")
	}
}

func (_ PoolTests) StatsTracksAndResetMisses() {
	p := New(1, 1)
	p.Checkout()
	p.Checkout()
	p.Checkout()

	Expect(p.Stats()["depleted"]).To.Equal(int64(2))
	//calling stats should reset this
	Expect(p.Stats()["depleted"]).To.Equal(int64(0))
}

func (_ PoolTests) TracksExpansion() {
	p := New(2, 1)
	for i := 1; i < 6; i++ {
		bytes := p.Checkout()
		bytes.WriteString(strings.Repeat("!", i))
		bytes.Release()
	}
	Expect(p.Stats()["expanded"]).To.Equal(int64(3))
	//calling stats should reset this
	Expect(p.Stats()["expanded"]).To.Equal(int64(0))
}

// why? because it's hard
func (_ PoolTests) DoesNotTrackUnpooledExpansion() {
	p := New(2, 1)
	for i := 1; i < 6; i++ {
		bytes := p.Checkout()
		bytes.WriteString(strings.Repeat("!", i))
	}
	Expect(p.Stats()["expanded"]).To.Equal(int64(0))
}

func (_ PoolTests) Each() {
	p := New(10, 4)
	i := byte(0)
	p.Each(func(b *Bytes) {
		b.WriteByte(i)
		i++
	})
	b := p.Checkout()
	b.Position(1)
	Expect(b.Bytes()).To.Equal([]byte{0})

	b = p.Checkout()
	b.Position(1)
	Expect(b.Bytes()).To.Equal([]byte{1})

	b = p.Checkout()
	b.Position(1)
	Expect(b.Bytes()).To.Equal([]byte{2})

	b = p.Checkout()
	b.Position(1)
	Expect(b.Bytes()).To.Equal([]byte{3})
}
