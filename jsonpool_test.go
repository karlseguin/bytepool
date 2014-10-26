package bytepool

import (
	"reflect"
	"testing"
	. "github.com/karlseguin/expect"
)

type JsonPoolTests struct{}

func Test_JsonPool(t *testing.T) {
	Expectify(new(JsonPoolTests), t)
}

func (j *JsonPoolTests) EachItemIsOfASpecifiedSize() {
	p := NewJson(1, 9)
	item := p.Checkout()
	defer item.Close()
	Expect(cap(item.bytes)).To.Equal(9)
}

func (j *JsonPoolTests) DynamicallyCreatesAnItemWhenPoolIsEmpty() {
	p := NewJson(1, 2)
	item1 := p.Checkout()
	item2 := p.Checkout()
	Expect(cap(item2.bytes)).To.Equal(2)
	Expect(item2.pool).To.Equal(nil)

	item1.Close()
	item2.Close()
	Expect(p.Len()).To.Equal(1)
	Expect(p.Misses()).To.Equal(int64(1))
}

func (j *JsonPoolTests) ReleasesAnItemBackIntoThePool() {
	p := NewJson(1, 20)
	item1 := p.Checkout()
	pointer := reflect.ValueOf(item1).Pointer()
	item1.Close()

	item2 := p.Checkout()
	defer item2.Close()

	if reflect.ValueOf(item2).Pointer() != pointer {
		Fail("Pool returned an unexected item")
	}
}

func (j *JsonPoolTests) StatsTracksAndResetMisses() {
	p := NewJson(1, 1)
	p.Checkout()
	p.Checkout()
	p.Checkout()
	Expect(p.Stats()["misses"]).To.Equal(int64(2))
	//calling stats should reset this
	Expect(p.Stats()["misses"]).To.Equal(int64(0))
}

func (j *JsonPoolTests) StatsTracksAndResetsMax() {
	p := NewJson(1, 20)
	item := p.Checkout()
	item.WriteString("abc")
	item.Close()

	item = p.Checkout()
	item.WriteString("abc123")
	item.Close()

	item = p.Checkout()
	item.WriteString("abc2")
	item.Close()

	Expect(p.Stats()["max"]).To.Equal(int64(8))
	//calling stats should reset this
	Expect(p.Stats()["max"]).To.Equal(int64(0))
}

func (j *JsonPoolTests) StatsTracksAndResetTaken() {
	p := NewJson(10, 1)
	p.Checkout()
	p.Checkout()
	p.Checkout()

	Expect(p.Stats()["taken"]).To.Equal(int64(3))
	//calling stats should reset this
	Expect(p.Stats()["taken"]).To.Equal(int64(0))
}
