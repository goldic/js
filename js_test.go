package js

import (
	"fmt"
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

	err := tryCall(func() {
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

func TestArray(t *testing.T) {
	v := tryVal(Parse([]byte(`[null,1,2,3.3,"4",{"a":1}]`)))

	ii := v.Array().Ints()
	ff := v.Array().Nums()
	ss := v.Array().Strings()

	require(t, eq(ii, []int{0, 1, 2, 3, 4, 0}))
	require(t, eq(ff, []float64{0, 1, 2, 3.3, 4, 0}))
	require(t, eq(ss, []string{"", "1", "2", "3.3", "4", `{"a":1}`}))
}

func TestArray_Objects(t *testing.T) {
	v := tryVal(Parse([]byte(`[null,{},null,{"a":1}]`)))

	oo := v.Array().Objects()

	require(t, ToStr(oo) == `[null,{},null,{"a":1}]`)
}

func TestOr(t *testing.T) {

	res0 := Or(0, false, nil, "")
	res1 := Or(0, false, "", 666, "abc")

	require(t, res0 == nil)
	require(t, res1 == 666)
}

func TestAnd(t *testing.T) {

	res0 := And(0, false, "", 666, "abc")
	res1 := And(666, -0.123, true, "abc")

	require(t, res0 == false)
	require(t, res1 == true)
}

func require(t *testing.T, ok bool) {
	if !ok {
		t.Fail()
	}
}

func eq[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func tryCall(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	fn()
	return
}
