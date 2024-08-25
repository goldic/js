package js

import (
	"encoding/json"
	"strings"
)

type Array []any

func NewArray(value ...any) Array {
	return value
}

func ToArray[T any](vv ...T) Array {
	aa := make([]any, len(vv))
	for i, v := range vv {
		aa[i] = v
	}
	return aa
}

func (arr Array) String() string {
	return string(arr.Bytes())
}

func (arr Array) Bytes() []byte {
	b, _ := json.Marshal(arr)
	return b
}

func (arr Array) Len() int {
	return len(arr)
}

func (arr Array) IsNull() bool {
	return arr == nil
}

func (arr *Array) Append(v any) {
	*arr = append(*arr, v)
}

func (arr Array) Eq(i int) (v Value) {
	if arr != nil && i >= 0 && i < len(arr) {
		return NewValue(arr[i])
	}
	return
}

func (arr Array) ForEach(fn func(Value)) {
	for _, v := range arr {
		fn(NewValue(v))
	}
}

func (arr Array) ForEachObject(fn func(obj Object)) {
	for _, v := range arr {
		fn(NewObject(v))
	}
}

func (arr Array) FindObject(fn func(obj Object) bool) Object {
	for _, v := range arr.Objects() {
		if obj := NewObject(v); fn(obj) {
			return obj
		}
	}
	return nil
}

func (arr Array) FindObjectBy(param string, val any) Object {
	for _, v := range arr {
		if obj := NewObject(v); obj != nil && obj[param] == val {
			return obj
		}
	}
	return nil
}

func (arr Array) Filter(fn func(v Value) bool) Array {
	vv := make(Array, 0, len(arr))
	for _, v := range arr {
		if fn(NewValue(v)) {
			vv = append(vv, v)
		}
	}
	return vv
}

func (arr Array) Map(fn func(v Value) any) Array {
	vv := make(Array, 0, len(arr))
	for _, v := range arr {
		vv = append(vv, fn(NewValue(v)))
	}
	return vv
}

func (arr Array) Join(sep string) string {
	return strings.Join(arr.Strings(), sep)
}

func (arr Array) Objects() []Object {
	vv := make([]Object, len(arr))
	for i, v := range arr {
		vv[i] = NewObject(v)
	}
	return vv
}

func (arr Array) Strings() []string {
	vv := make([]string, len(arr))
	for i, v := range arr {
		vv[i] = ToStr(v)
	}
	return vv
}

func (arr Array) Ints() []int {
	vv := make([]int, len(arr))
	for i, v := range arr {
		vv[i] = ToInt(v)
	}
	return vv
}

func (arr Array) Nums() []float64 {
	vv := make([]float64, len(arr))
	for i, v := range arr {
		vv[i] = ToNum(v)
	}
	return vv
}

func (arr Array) Unmarshal(v any) error {
	return json.Unmarshal(arr.Bytes(), v)
}
