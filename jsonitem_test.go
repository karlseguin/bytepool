package bytepool

import (
	. "github.com/karlseguin/expect"
	"testing"
)

type JsonItemTests struct {}

func Test_JsonItems(t *testing.T) {
	Expectify(new(JsonItemTests), t)
}

func (j *JsonItemTests) WriteAnEncodedString() {
	item := NewJsonItem(100, nil)
	item.WriteString(`over "9000"`)
	Expect(item.String()).To.Equal(`"over \"9000\""`)
}

func (j *JsonItemTests) WriteAString() {
	item := NewJsonItem(100, nil)
	item.WriteSafeString(`over "9000"`)
	Expect(item.String()).To.Equal(`"over "9000""`)
}

func (j *JsonItemTests) JsonWritesAnEmptyArray() {
	item := NewJsonItem(100, nil)
	item.BeginArray()
	item.EndArray()
	Expect(item.String()).To.Equal("[]")
}

func (j *JsonItemTests) JsonWritesASingleValueArray() {
	item := NewJsonItem(100, nil)
	item.BeginArray()
	item.WriteInt(90)
	item.EndArray()
	Expect(item.String()).To.Equal("[90]")
}

func (j *JsonItemTests) JsonWritesAMultiValueArray() {
	item := NewJsonItem(100, nil)
	item.BeginArray()
	item.WriteInt(90)
	item.WriteBool(false)
	item.WriteString("abc")
	item.WriteBool(true)
	item.EndArray()
	Expect(item.String()).To.Equal(`[90,false,"abc",true]`)
}

func (j *JsonItemTests) JsonWritesAnEmptyObject() {
	item := NewJsonItem(100, nil)
	item.BeginObject()
	item.EndObject()
	Expect(item.String()).To.Equal("{}")
}

func (j *JsonItemTests) JsonWritesADelimitedByte() {
	item := NewJsonItem(100, nil)
	item.BeginArray()
	item.Write([]byte(`"abc"`))
	item.Write([]byte(`"123"`))
	item.EndArray()
	Expect(item.String()).To.Equal(`["abc","123"]`)
}

func (j *JsonItemTests) JsonASingleValueObject() {
	item := NewJsonItem(100, nil)
	item.BeginObject()
	item.WriteKeyString("over", "90\"00!")
	item.EndObject()
	Expect(item.String()).To.Equal(`{"over":"90\"00!"}`)
}

func (j *JsonItemTests) JsonAMultiValueObject() {
	item := NewJsonItem(100, nil)
	item.BeginObject()
	item.WriteKeySafeString("name", "goku")
	item.WriteKeyInt("power", 9000)
	item.WriteKeyBool("over", true)
	item.EndObject()
	Expect(item.String()).To.Equal(`{"name":"goku","power":9000,"over":true}`)
}

func (j *JsonItemTests) WriteMultipleKeyObjects() {
	item := NewJsonItem(100, nil)
	item.BeginObject()

	item.WriteKeyObject("name")
	item.WriteKeyString("en", "leto")
	item.EndObject()

	item.WriteKeyObject("desc")
	item.WriteKeyString("en", "worm")
	item.EndObject()

	item.EndObject()
	Expect(item.String()).To.Equal(`{"name":{"en":"leto"},"desc":{"en":"worm"}}`)
}

func (j *JsonItemTests) JsonNestedObjects() {
	item := NewJsonItem(100, nil)
	item.BeginArray()
	item.WriteInt(1)
	item.BeginObject()
	item.WriteKeyString("name", "goku")
	item.WriteKeyArray("levels")
	item.WriteInt(2)
	item.BeginObject()
	item.WriteKeyObject("over")
	item.WriteKeyString("9000", "!")
	item.EndObject()
	item.EndObject()
	item.EndArray()
	item.EndObject()
	item.EndArray()
	Expect(item.String()).To.Equal(`[1,{"name":"goku","levels":[2,{"over":{"9000":"!"}}]}]`)
}
