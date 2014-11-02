package bytepool

import (
	"reflect"
	"testing"
	. "github.com/karlseguin/expect"
)

type PoolTests struct{}

func Test_Pool(t *testing.T) {
	Expectify(new(PoolTests), t)
}

func (_ *PoolTests) EachItemIsOfASpecifiedSize() {
	p := New(9, 1)
	bytes := p.Checkout()
	Expect(cap(bytes.fixed.bytes)).To.Equal(9)
}

func (_ *PoolTests) HasTheSpecifiedNumberOfItems() {
	p := New(23, 3)
	Expect(cap(p.list)).To.Equal(3)
}

func (pt *PoolTests) DynamicallyCreatesAnItemWhenPoolIsEmpty() {
	p := New(23, 1)
	bytes1 := p.Checkout()
	bytes2 := p.Checkout()
	Expect(cap(bytes2.fixed.bytes)).To.Equal(23)
	Expect(bytes2.pool).To.Equal(nil)
	bytes1.Release()
	bytes2.Release()
	Expect(len(p.list)).To.Equal(1)
	Expect(p.Misses()).To.Equal(int64(1))
}

func (_ *PoolTests) ReleasesAnItemBackIntoThePool() {
	p := New(20, 1)
	bytes := p.Checkout()
	pointer := reflect.ValueOf(bytes).Pointer()
	bytes.Release()

	if reflect.ValueOf(p.Checkout()).Pointer() != pointer {
		Fail("Pool returned an unexected item")
	}
}

func (pt *PoolTests) StatsTracksAndResetMisses() {
	p := New(1, 1)
	p.Checkout()
	p.Checkout()
	p.Checkout()

	Expect(p.Stats()["misses"]).To.Equal(int64(2))
	//calling stats should reset this
	Expect(p.Stats()["misses"]).To.Equal(int64(0))
}
