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

func (pt *PoolTests) EachItemIsOfASpecifiedSize() {
	p := New(1, 9)
	item := p.Checkout()
	defer item.Close()
	Expect(cap(item.bytes)).To.Equal(9)
}

func (pt *PoolTests) HasTheSpecifiedNumberOfItems() {
	p := New(3, 4)
	Expect(cap(p.list)).To.Equal(3)
}

func (pt *PoolTests) HasTheSpecifiedNumberOfItemsAfterACulling() {
	p := New(5, 4)
	p.SetCount(2)
	Expect(cap(p.list)).To.Equal(2)
}

func (pt *PoolTests) HasTheSpecifiedNumberOfItemsAfterAddition() {
	p := New(5, 4)
	p.SetCount(10)
	Expect(cap(p.list)).To.Equal(10)
}

func (pt *PoolTests) EachItemIsOfASpecifiedSizeAfterResize() {
	p := New(1, 5)
	p.SetCapacity(11)
	item := p.Checkout()
	defer item.Close()
	Expect(cap(item.bytes)).To.Equal(11)
}

func (pt *PoolTests) DynamicallyCreatesAnItemWhenPoolIsEmpty() {
	p := New(1, 2)
	item1 := p.Checkout()
	item2 := p.Checkout()
	Expect(cap(item2.bytes)).To.Equal(2)
	Expect(item2.pool).To.Equal(nil)
	item1.Close()
	item2.Close()
	Expect(p.Len()).To.Equal(1)
	Expect(p.Misses()).To.Equal(int64(1))
}

func (pt *PoolTests) ReleasesAnItemBackIntoThePool() {
	p := New(1, 20)
	item1 := p.Checkout()
	pointer := reflect.ValueOf(item1).Pointer()
	item1.Close()

	item2 := p.Checkout()
	defer item2.Close()
	if reflect.ValueOf(item2).Pointer() != pointer {
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

func (pt *PoolTests) StatsTracksAndResetsMax() {
	p := New(1, 20)
	item := p.Checkout()
	item.WriteString("abc")
	item.Close()

	item = p.Checkout()
	item.WriteString("abc123")
	item.Close()

	item = p.Checkout()
	item.WriteString("abc2")
	item.Close()

	Expect(p.Stats()["max"]).To.Equal(int64(6))
	//calling stats should reset this
	Expect(p.Stats()["max"]).To.Equal(int64(0))
}

func (pt *PoolTests) StatsTracksAndResetTaken() {
	p := New(10, 1)
	p.Checkout()
	p.Checkout()
	p.Checkout()

	Expect(p.Stats()["taken"]).To.Equal(int64(3))
	//calling stats should reset this
	Expect(p.Stats()["taken"]).To.Equal(int64(0))
}
