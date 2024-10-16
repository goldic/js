package js

import (
	"testing"
)

func TestNewObject(t *testing.T) {

	obj0 := NewObject(nil)
	obj1 := NewObject(struct{ A, B int }{1, 2})
	obj2 := NewObject(map[string]int{"a": 1, "b": 2})
	obj3 := NewObject(map[int]int{1: 1, 2: 2})
	obj4 := NewObject(map[int]int{})
	obj5 := NewObject((*struct{})(nil)) // any nil pointer

	require(t, obj0 == nil && obj0.String() == `null`)
	require(t, obj1 != nil && obj1.String() == `{"A":1,"B":2}`)
	require(t, obj2 != nil && obj2.String() == `{"a":1,"b":2}`)
	require(t, obj3 != nil && obj3.String() == `{"1":1,"2":2}`)
	require(t, obj4 != nil && obj4.String() == `{}`)
	require(t, obj5 == nil && obj5.String() == `null`)
}

func TestNewObject_fail(t *testing.T) {

	err := call(func() {
		NewObject(123)
	})

	require(t, err != nil)
}

func TestParseObject(t *testing.T) {
	o, err := ParseObject([]byte(`{
		"i":	123,
		"num":	-12.3,
		"str":	"abc",
		"arr":	[1,2,3],
		"obj":	{"a":1},
		"objs":	[{"a":1},null,{"b":2}]
	}`))

	i := o.GetInt("i")
	num := o.GetNum("num")
	str := o.GetStr("str")
	arr := o.GetArr("arr")
	obj := o.GetObj("obj")
	objs := o.GetArr("objs").Objects()

	require(t, err == nil)
	require(t, i == 123)
	require(t, num == -12.3)
	require(t, str == "abc")
	require(t, arr.String() == `[1,2,3]`)
	require(t, obj.String() == `{"a":1}`)
	require(t, len(objs) == 3)
	require(t, objs[0].String() == `{"a":1}`)
	require(t, objs[1].String() == `null`)
	require(t, objs[2].String() == `{"b":2}`)
}

func TestParseObject_null(t *testing.T) {
	o, err := ParseObject(nil)

	require(t, err == nil)
	require(t, o == nil)
	require(t, ToStr(o) == "null")
}
