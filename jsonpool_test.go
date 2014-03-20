package bytepool

import (
	"reflect"
	"testing"
)

func TestJsonPoolEachItemIsOfASpecifiedSize(t *testing.T) {
	expected := 9
	p := NewJson(1, expected)
	item := p.Checkout()
	defer item.Close()
	if cap(item.bytes) != expected {
		t.Errorf("expecting array to have a capacity of %d, got %d", expected, cap(item.bytes))
	}
}

func TestJsonPoolDynamicallyCreatesAnItemWhenPoolIsEmpty(t *testing.T) {
	p := NewJson(1, 2)
	item1 := p.Checkout()
	item2 := p.Checkout()
	if cap(item2.bytes) != 2 {
		t.Error("Dynamically created item was not properly initialized")
	}
	if item2.pool != nil {
		t.Error("The dynamically created item should have a nil pool")
	}

	item1.Close()
	item2.Close()
	if p.Len() != 1 {
		t.Errorf("Expecting a pool lenght of 1, got %d", p.Len())
	}
	if p.Misses() != 1 {
		t.Errorf("Expecting a miss count of 1, got %d", p.Misses())
	}

}
func TestJsonPoolReleasesAnItemBackIntoThePool(t *testing.T) {
	p := NewJson(1, 20)
	item1 := p.Checkout()
	pointer := reflect.ValueOf(item1).Pointer()
	item1.Close()

	item2 := p.Checkout()
	defer item2.Close()
	if reflect.ValueOf(item2).Pointer() != pointer {
		t.Error("Pool returned an unexected item")
	}
}

func TestJsonPoolStatsTracksAndResetMisses(t *testing.T) {
	p := NewJson(1, 1)
	p.Checkout()
	p.Checkout()
	p.Checkout()

	misses := p.Stats()["misses"]
	if misses != 2 {
		t.Errorf("Expected 2 misses, got %d", misses)
	}

	//calling stats should reset this
	misses = p.Stats()["misses"]
	if misses != 0 {
		t.Errorf("Expected 0 misses, got %d", misses)
	}
}

func TestJsonPoolStatsTracksAndResetsMax(t *testing.T) {
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

	max := p.Stats()["max"]
	if max != 8 {
		t.Errorf("Expected 8 max, got %d", max)
	}

	//calling stats should reset this
	max = p.Stats()["max"]
	if max != 0 {
		t.Errorf("Expected 0 max, got %d", max)
	}
}
