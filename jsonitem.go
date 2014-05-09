package bytepool

import (
	"strconv"
	"strings"
	"time"
)

type JsonItem struct {
	*Item
	depth int
	pool  *JsonPool
}

func NewJsonItem(capacity int, pool *JsonPool) *JsonItem {
	return &JsonItem{
		pool: pool,
		Item: NewItem(capacity, nil),
	}
}

var JsonEncode = func(s string) string {
	return strings.Replace(strings.Replace(s, `"`, `\"`, -1), "\n", "\\n", -1)
}

func (item *JsonItem) WriteString(s string) int {
	return item.WriteSafeString(JsonEncode(s))
}

func (item *JsonItem) WriteInt(value int) int {
	n := item.Item.WriteString(strconv.Itoa(value))
	return item.delimit(n)
}

func (item *JsonItem) WriteBool(value bool) int {
	n := item.Item.WriteString(strconv.FormatBool(value))
	return item.delimit(n)
}

func (item *JsonItem) WriteSafeString(s string) int {
	return item.writeString(s, true)
}

func (item *JsonItem) Write(b []byte) int {
	n, _ := item.Item.Write(b)
	return item.delimit(n)
}

func (item *JsonItem) WriteKeyString(key, value string) int {
	return item.WriteKeySafeString(key, JsonEncode(value))
}

func (item *JsonItem) WriteKeySafeString(key, value string) int {
	return item.writeKeyValue(key, `"`+value+`"`)
}

func (item *JsonItem) WriteKeyInt(key string, value int) int {
	return item.writeKeyValue(key, strconv.Itoa(value))
}

func (item *JsonItem) WriteKeyBool(key string, value bool) int {
	return item.writeKeyValue(key, strconv.FormatBool(value))
}

func (item *JsonItem) WriteKeyTime(key string, value time.Time) int {
	return item.WriteKeySafeString(key, value.Format(time.RFC3339Nano))
}

func (item *JsonItem) WriteKeyArray(key string) int {
	n := item.writeString(key, false)
	if item.WriteByte(byte(':')) {
		n++
	}
	if item.BeginArray() {
		n++
	}
	return n
}

func (item *JsonItem) WriteKeyObject(key string) int {
	n := item.writeString(key, false)
	if item.WriteByte(byte(':')) {
		n++
	}
	if item.BeginObject() {
		n++
	}
	return n
}

func (item *JsonItem) writeKeyValue(key, value string) int {
	n := item.writeString(key, false)
	if item.WriteByte(byte(':')) {
		n++
	}
	n += item.Item.WriteString(value)
	return item.delimit(n)
}

func (item *JsonItem) writeString(s string, delimit bool) int {
	n := item.Item.WriteString(`"` + s + `"`)
	if delimit == false {
		return n
	}
	return item.delimit(n)
}

func (item *JsonItem) BeginArray() bool {
	item.depth++
	return item.WriteByte('[')
}

func (item *JsonItem) EndArray() bool {
	item.depth--
	item.TrimLastIf(',')
	item.WriteByte(']')
	item.delimit(0)
	return true
}

func (item *JsonItem) BeginObject() bool {
	item.depth++
	return item.WriteByte('{')
}

func (item *JsonItem) EndObject() bool {
	item.depth--
	item.TrimLastIf(',')
	item.WriteByte('}')
	item.delimit(0)
	return true
}

func (item *JsonItem) delimit(length int) int {
	if item.depth == 0 {
		return length
	}
	item.WriteByte(byte(','))
	return length + 1
}

func (item *JsonItem) Close() error {
	if item.pool != nil {
		item.pool.track(int64(item.length))
		item.read = 0
		item.depth = 0
		item.length = 0
		item.pool.list <- item
	}
	return nil
}
