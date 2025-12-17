package js

import (
	"encoding/json"
	"sort"
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

// Unshift adds an element to the beginning of the array.
func (arr *Array) Unshift(v any) {
	*arr = append([]any{v}, *arr...)
}

// Eq retrieves the value from the array by index.
func (arr Array) Eq(i int) Value {
	if n := len(arr); n > 0 {
		if i < 0 {
			i += n
		}
		if i >= 0 && i < n {
			return NewValue(arr[i])
		}
	}
	return Value{}
}

// First retrieves the first value from the array.
func (arr Array) First() Value {
	return arr.Eq(0)
}

// Last retrieves the last value from the array.
func (arr Array) Last() Value {
	if n := len(arr); n > 0 {
		return NewValue(arr[n-1])
	}
	return Value{}
}

// ForEach executes a function for each element in the array.
func (arr Array) ForEach(fn func(Value, int)) {
	for i, v := range arr {
		fn(NewValue(v), i)
	}
}

// ForEachObject executes a function for each object in the array.
func (arr Array) ForEachObject(fn func(obj Object, index int)) {
	for i, v := range arr {
		fn(newObject(v), i)
	}
}

// IndexOf searches for an element in the array by a value.
func (arr Array) IndexOf(value any) int {
	val := NewValue(value)
	for i, v := range arr {
		if val.Equal(v) {
			return i
		}
	}
	return -1
}

// IndexOfFn searches for an element in the array by a function.
func (arr Array) IndexOfFn(fn func(Value) bool) int {
	for i, v := range arr {
		if fn(NewValue(v)) {
			return i
		}
	}
	return -1
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
func (arr Array) Map(fn func(v Value, index int) any) Array {
	vv := make(Array, 0, len(arr))
	for i, v := range arr {
		vv = append(vv, fn(NewValue(v), i))
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

// SortBy sorts the array of objects by a specified parameter name.
func (arr Array) SortBy(paramName string) {
	sort.Slice(arr, func(i, j int) bool {
		a := arr.Eq(i).Object().get(paramName)
		b := arr.Eq(j).Object().get(paramName)
		return Cmp(a, b) < 0
	})
}

// Sort sorts the array using a custom comparison function.
func (arr Array) Sort(less func(a, b Value) bool) {
	if less == nil {
		less = func(a, b Value) bool { return Cmp(a.val, b.val) < 0 }
	}
	sort.Slice(arr, func(i, j int) bool {
		return less(arr.Eq(i), arr.Eq(j))
	})
}

// Reverse reverses the order of elements in the array.
func (arr Array) Reverse() {
	for i, j := 0, len(arr)-1; i < j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}
