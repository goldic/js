package js

import (
	"slices"
	"testing"
)

func TestArray_Objects(t *testing.T) {
	v := must(Parse([]byte(`[null,{},null,{"a":1},[],1]`)))

	oo := v.Array().Objects()

	require(t, ToStr(oo) == `[null,{},null,{"a":1},null,null]`)
}

func TestArray(t *testing.T) {
	v := must(Parse([]byte(`[null,1,2,3.3,"4",{"a":1}]`)))

	ii := v.Array().Ints()
	ff := v.Array().Nums()
	ss := v.Array().Strings()

	require(t, slices.Equal(ii, []int{0, 1, 2, 3, 4, 0}))
	require(t, slices.Equal(ff, []float64{0, 1, 2, 3.3, 4, 0}))
	require(t, slices.Equal(ss, []string{"", "1", "2", "3.3", "4", `{"a":1}`}))
}
