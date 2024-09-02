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
	return tryVal(ParseObject(tryVal(json.Marshal(v))))
}

func (obj Object) String() string {
	if obj.Has("_string") {
		return obj.GetStr("_string")
	}
	return string(obj.Bytes())
}

func (obj Object) Bytes() []byte {
	return tryVal(json.Marshal(map[string]any(obj)))
}

func (obj Object) IndentString() string {
	return string(tryVal(json.MarshalIndent(map[string]any(obj), "", "  ")))
}

func (obj Object) Copy() Object {
	if obj == nil {
		return nil
	}
	c := make(Object, len(obj))
	for k, v := range obj {
		c[k] = v
	}
	return c
}

func (obj Object) Set(name string, v any) Object {
	if obj == nil {
		obj = Object{}
	}
	obj[name] = v
	return obj
}

func (obj Object) Delete(name string) Object {
	if obj == nil {
		obj = Object{}
	}
	delete(obj, name)
	return obj
}

func (obj Object) Clone() Object {
	return maps.Clone(obj)
}

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

func (obj Object) Get(name string) (v Value) {
	if obj != nil {
		v.val = obj[name]
	}
	return
}

func (obj Object) GetBool(name string) bool {
	return obj.Get(name).Bool()
}

func (obj Object) GetStr(name string) string {
	return obj.Get(name).String()
}

func (obj Object) GetNum(name string) float64 {
	return obj.Get(name).Float64()
}

func (obj Object) GetInt(name string) int {
	return obj.Get(name).Int()
}

func (obj Object) GetInt64(name string) int64 {
	return obj.Get(name).Int64()
}

func (obj Object) GetUint64(name string) uint64 {
	return obj.Get(name).Uint64()
}

func (obj Object) GetTime(name string) time.Time {
	return obj.Get(name).Time()
}

func (obj Object) GetObj(name string) Object {
	return obj.Get(name).Object()
}

func (obj Object) GetArr(name string) Array {
	return obj.Get(name).Array()
}

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

func (obj Object) Has(key string) bool {
	if obj != nil {
		_, ok := obj[key]
		return ok
	}
	return false
}

func (obj Object) Keys() (keys []string) {
	if obj != nil {
		for key := range obj {
			keys = append(keys, key)
		}
		sort.Strings(keys)
	}
	return
}

func (obj Object) MarshalTo(v any) error {
	return json.Unmarshal(obj.Bytes(), v)
}

func (obj Object) Encode() []byte {
	return obj.Bytes()
}

func (obj *Object) Decode(data []byte) (err error) {
	return json.Unmarshal(data, &obj)
}

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
	if obj == nil {
		obj = Object{}
	}
	return
}

func MustParseObject(data []byte) Object {
	return tryVal(ParseObject(data))
}

func ReadObject(r io.Reader) (_ Object, err error) {
	defer catch(&err)
	return ParseObject(readAll(r))
}
