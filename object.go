package js

import (
	"encoding/json"
	"io"
	"maps"
	"net/url"
	"sort"
	"time"
)

type Object map[string]any

// NewObject creates a new Object from any Go value.
func NewObject(v any) Object {
	switch val := v.(type) {
	case nil:
		return nil
	case Object:
		return val
	case map[string]any:
		return val
	case Value:
		return NewObject(val.val)
	}
	return MustParseObject(tryVal(json.Marshal(v)))
}

// String converts the entire object to a JSON string.
func (obj Object) String() string {
	if obj.Has("_string") {
		return obj.GetStr("_string")
	}
	return string(obj.Bytes())
}

// Bytes converts the object to bytes (JSON).
func (obj Object) Bytes() []byte {
	return tryVal(json.Marshal(map[string]any(obj)))
}

// IndentString converts the object to a formatted string (pretty JSON).
func (obj Object) IndentString() string {
	return string(tryVal(json.MarshalIndent(map[string]any(obj), "", "  ")))
}

// Len returns the number of key-value pairs in the object.
func (obj Object) Len() int {
	if obj == nil {
		return 0
	}
	return len(obj)
}

// Set adds or updates a key-value pair in the object.
func (obj Object) Set(name string, v any) Object {
	if obj == nil {
		obj = Object{}
	}
	obj[name] = v
	return obj
}

// Delete removes the specified key-value pair from the object.
func (obj Object) Delete(name string) Object {
	if obj == nil {
		obj = Object{}
	}
	delete(obj, name)
	return obj
}

// Clone creates a clone of the object.
func (obj Object) Clone() Object {
	return maps.Clone(obj)
}

// Extend merges the object with other objects.
func (obj Object) Extend(o ...Object) Object {
	if obj == nil {
		obj = Object{}
	}
	for _, vObj := range o {
		if vObj != nil {
			for k, v := range vObj {
				obj[k] = v
			}
		}
	}
	return obj
}

// Get retrieves the value for the given key as a Value.
func (obj Object) Get(name string) (v Value) {
	if obj != nil {
		v.val = obj[name]
	}
	return
}

// GetBool retrieves the boolean value by key.
func (obj Object) GetBool(name string) bool {
	return obj.Get(name).Bool()
}

// GetStr retrieves the string value by key.
func (obj Object) GetStr(name string) string {
	return obj.Get(name).String()
}

// GetNum retrieves the numeric value (float64) by key.
func (obj Object) GetNum(name string) float64 {
	return obj.Get(name).Float64()
}

// GetInt retrieves the integer value (int) by key.
func (obj Object) GetInt(name string) int {
	return obj.Get(name).Int()
}

// GetInt64 retrieves the 64-bit integer value (int64) by key.
func (obj Object) GetInt64(name string) int64 {
	return obj.Get(name).Int64()
}

// GetUint64 retrieves the unsigned 64-bit integer value (uint64) by key.
func (obj Object) GetUint64(name string) uint64 {
	return obj.Get(name).Uint64()
}

// GetTime retrieves the time value (time.Time) by key.
func (obj Object) GetTime(name string) time.Time {
	return obj.Get(name).Time()
}

// GetObj retrieves the object by key.
func (obj Object) GetObj(name string) Object {
	return obj.Get(name).Object()
}

// GetArr retrieves the array by key.
func (obj Object) GetArr(name string) Array {
	return obj.Get(name).Array()
}

// GetNoNil retrieves the first non-nil value by keys.
func (obj Object) GetNoNil(name ...string) (v Value) {
	if obj != nil {
		for _, n := range name {
			if v.val = obj[n]; v.val != nil {
				return
			}
		}
	}
	return
}

// Has checks if the object contains the given key.
func (obj Object) Has(key string) bool {
	if obj != nil {
		_, ok := obj[key]
		return ok
	}
	return false
}

// Keys returns all keys of the object (in sorted order).
func (obj Object) Keys() (keys []string) {
	if obj != nil {
		for key := range obj {
			keys = append(keys, key)
		}
		sort.Strings(keys)
	}
	return
}

// MarshalTo unmarshals the object into another structure.
func (obj Object) MarshalTo(v any) error {
	return json.Unmarshal(obj.Bytes(), v)
}

// Encode converts the object to bytes (JSON) for encoding.
// This is a synonym for Bytes().
func (obj Object) Encode() []byte {
	return obj.Bytes()
}

// Decode decodes the object from bytes (JSON).
func (obj *Object) Decode(data []byte) (err error) {
	return json.Unmarshal(data, &obj)
}

// URLValues converts the object to url.Values.
func (obj Object) URLValues() (values url.Values) {
	values = make(url.Values, len(obj))
	for k, v := range obj {
		if v != nil {
			values.Set(k, obj.GetStr(k))
		}
	}
	return
}

func ParseObject(data []byte) (obj Object, err error) {
	if len(data) > 0 {
		err = json.Unmarshal(data, &obj)
	}
	return
}

func MustParseObject(data []byte) Object {
	return tryVal(ParseObject(data))
}

func ReadObject(r io.Reader) (obj Object, err error) {
	defer catch(&err)
	return ParseObject(readAll(r))
}

// ObjectFromURLValues converts url.Values to an object.
func ObjectFromURLValues(val url.Values) Object {
	obj := make(Object, len(val))
	for k, vv := range val {
		obj[k] = vv[0]
	}
	return obj
}
