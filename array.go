package js

import (
	"encoding/json"
	"strings"
)

type Array []any

// NewArray creates a new array.
func NewArray(value ...any) Array {
	return value
}

// ToArray converts a slice of any type to an Array.
func ToArray[T any](vv ...T) Array {
	aa := make([]any, len(vv))
	for i, v := range vv {
		aa[i] = v
	}
	return aa
}

// String converts the array to a string (JSON).
func (arr Array) String() string {
	return string(arr.Bytes())
}

// Bytes converts the array to bytes (JSON).
func (arr Array) Bytes() []byte {
	b, _ := json.Marshal(arr)
	return b
}

// Len returns the length of the array.
func (arr Array) Len() int {
	return len(arr)
}

// IsNull checks if the array is empty (null).
func (arr Array) IsNull() bool {
	return arr == nil
}

// Push adds an element to the array.
func (arr *Array) Push(v any) {
	*arr = append(*arr, v)
}

// Eq retrieves the value from the array by index.
func (arr Array) Eq(i int) (v Value) {
	if arr != nil && i >= 0 && i < len(arr) {
		return NewValue(arr[i])
	}
	return
}

// ForEach executes a function for each element in the array.
func (arr Array) ForEach(fn func(Value)) {
	for _, v := range arr {
		fn(NewValue(v))
	}
}

// ForEachObject executes a function for each object in the array.
func (arr Array) ForEachObject(fn func(obj Object)) {
	for _, v := range arr {
		fn(newObject(v))
	}
}

// FindObject searches for an object in the array by a function.
func (arr Array) FindObject(fn func(obj Object) bool) Object {
	for _, v := range arr.Objects() {
		if obj := newObject(v); obj != nil && fn(obj) {
			return obj
		}
	}
	return nil
}

// FindObjectBy searches for an object in the array by parameter and value.
func (arr Array) FindObjectBy(param string, val any) Object {
	for _, v := range arr {
		if obj := newObject(v); obj != nil && obj[param] == val {
			return obj
		}
	}
	return nil
}

// Filter filters the array by a function.
func (arr Array) Filter(fn func(v Value) bool) Array {
	vv := make(Array, 0, len(arr))
	for _, v := range arr {
		if fn(NewValue(v)) {
			vv = append(vv, v)
		}
	}
	return vv
}

// Map transforms the array using a function.
func (arr Array) Map(fn func(v Value) any) Array {
	vv := make(Array, 0, len(arr))
	for _, v := range arr {
		vv = append(vv, fn(NewValue(v)))
	}
	return vv
}

// Join concatenates string representations of the array elements with a separator.
func (arr Array) Join(sep string) string {
	return strings.Join(arr.Strings(), sep)
}

// Objects converts the array to a slice of objects.
func (arr Array) Objects() []Object {
	vv := make([]Object, len(arr))
	for i, v := range arr {
		vv[i] = newObject(v)
	}
	return vv
}

// Strings converts the array to a slice of strings.
func (arr Array) Strings() []string {
	vv := make([]string, len(arr))
	for i, v := range arr {
		vv[i] = ToStr(v)
	}
	return vv
}

// Ints converts the array to a slice of integers.
func (arr Array) Ints() []int {
	vv := make([]int, len(arr))
	for i, v := range arr {
		vv[i] = ToInt(v)
	}
	return vv
}

// Nums converts the array to a slice of floating-point numbers.
func (arr Array) Nums() []float64 {
	vv := make([]float64, len(arr))
	for i, v := range arr {
		vv[i] = ToNum(v)
	}
	return vv
}

// MarshalTo deserializes the array into another structure.
func (arr Array) MarshalTo(v any) error {
	return json.Unmarshal(arr.Bytes(), v)
}
