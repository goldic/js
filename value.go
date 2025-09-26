package js

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Value struct {
	val any
}

func NewValue(v any) Value {
	if val, ok := v.(Value); ok {
		return val
	}
	return Value{v}
}

func Parse(data []byte) (v Value, err error) {
	if len(data) == 0 {
		return
	}
	err = json.Unmarshal(bytes.TrimSpace(data), &v.val)
	return
}

func MustParse(data []byte) Value {
	return must(Parse(data))
}

func ReadValue(r io.Reader) (v Value, err error) {
	defer catch(&err)
	return Parse(readAll(r))
}

func (v Value) Value() any {
	return v.val
}

func (v Value) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.val)
}

func (v *Value) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &v.val)
}

func (v Value) MarshalTo(val any) error {
	return json.Unmarshal(v.Bytes(), val)
}

func (v Value) Bytes() (data []byte) {
	data, _ = json.Marshal(v.val)
	return
}

func (v Value) JSON() string {
	return string(v.Bytes())
}

func (v Value) Empty() bool {
	if isNil(v.val) {
		return true
	}
	switch v := v.val.(type) {
	case bool:
		return !v
	case string:
		return v == ""
	case int, int32, int64, uint, uint32, uint64, float32, float64:
		return v == 0
	case map[string]any:
		return len(v) == 0
	case []any:
		return len(v) == 0
	case []byte:
		return len(v) == 0
	}
	return false
}

func (v Value) Equal(b any) bool {
	return Cmp(v.val, b) == 0
}

func (v Value) Cmp(b any) int {
	return Cmp(v.val, b)
}

func (v Value) IsNull() bool {
	return v.val == nil
}

func (v Value) IsNum() bool {
	return IsNum(v.val)
}

func (v Value) IsObject() bool {
	return v.Object() != nil
}

func (v Value) IsArray() bool {
	kind := reflect.ValueOf(v.val).Kind()
	return kind == reflect.Slice || kind == reflect.Array
}

func (v Value) Array() Array {
	if !v.IsArray() {
		return nil
	}
	switch val := v.val.(type) {
	case []any:
		return val
	case Array:
		return val
	}
	var arr []any
	v.MarshalTo(&arr) // ignore error
	return arr
}

func (v Value) Object() Object {
	return newObject(v.val)
}

func (v Value) Objects() []Object {
	return v.Array().Objects()
}

func (v Value) Bool() bool {
	return v.Int64() != 0
}

func (v Value) String() string {
	if isNil(v.val) {
		return ""
	}
	switch val := v.val.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case int:
		return strconv.FormatInt(int64(val), 10)
	case uint:
		return strconv.FormatUint(uint64(val), 10)
	case int32:
		return strconv.FormatInt(int64(val), 10)
	case uint32:
		return strconv.FormatUint(uint64(val), 10)
	case int64:
		return strconv.FormatInt(int64(val), 10)
	case uint64:
		return strconv.FormatUint(uint64(val), 10)
	case float64:
		return strconv.FormatFloat(float64(val), 'f', -1, 64)
	case float32:
		return strconv.FormatFloat(float64(val), 'f', -1, 64)
	case fmt.Stringer:
		return val.String()
	case error:
		return val.Error()
	case io.Reader:
		data, _ := io.ReadAll(val)
		return string(data)
	}
	return v.JSON()
}

func (v Value) Int() int {
	return int(v.Int64())
}

func (v Value) Uint64() uint64 {
	return uint64(v.Int64())
}

func (v Value) Int64() (num int64) {
	if isNil(v.val) {
		return
	}
	switch val := v.val.(type) {
	case int:
		return int64(val)
	case uint:
		return int64(val)
	case int64:
		return int64(val)
	case uint64:
		return int64(val)
	case int32:
		return int64(val)
	case uint32:
		return int64(val)
	case float32:
		return int64(val)
	case float64:
		return int64(val)
	case bool:
		if val {
			return 1
		}
	case []byte:
		return NewValue(string(val)).Int64()
	case string:
		if strings.IndexByte(val, '.') >= 0 {
			f, _ := strconv.ParseFloat(val, 64)
			return int64(f)
		}
		num, _ = strconv.ParseInt(val, 0, 64)
	default:
		return NewValue(v.String()).Int64()
	}
	return
}

func (v Value) Num() float64 {
	return v.Float64()
}

func (v Value) Float64() (num float64) {
	if isNil(v.val) {
		return
	}
	switch val := v.val.(type) {
	case nil:
		return 0
	case int:
		return float64(val)
	case uint:
		return float64(val)
	case int32:
		return float64(val)
	case uint32:
		return float64(val)
	case int64:
		return float64(val)
	case uint64:
		return float64(val)
	case float32:
		return float64(val)
	case float64:
		return float64(val)
	case []byte:
		num, _ = strconv.ParseFloat(string(val), 64)
	case string:
		num, _ = strconv.ParseFloat(val, 64)
	default:
		num, _ = strconv.ParseFloat(v.String(), 64)
	}
	return
}

func (v Value) Time() time.Time {
	if IsNum(v) {
		return time.Unix(NewValue(v).Int64(), 0)
	}
	t, _ := ParseTime(ToStr(v))
	return t
}

func isNil(v any) bool {
	return v == nil || v == (*any)(nil)
}

func IsNum(v any) bool {
	switch v.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return true
	}
	return false
}

func IsInt(v any) bool {
	switch v.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return true
	}
	return false
}

func ToStr(v any) string {
	return NewValue(v).String()
}

func ToNum(v any) float64 {
	return NewValue(v).Float64()
}

func ToInt(v any) int {
	return NewValue(v).Int()
}

func ToUint64(v any) uint64 {
	return NewValue(v).Uint64()
}

func IsEmpty(v any) bool {
	return NewValue(v).Empty()
}

func Or(v ...any) any {
	for _, v := range v {
		if !IsEmpty(v) {
			return v
		}
	}
	return nil
}

func And(v ...any) bool {
	for _, v := range v {
		if IsEmpty(v) {
			return false
		}
	}
	return true
}

func Cmp(a, b any) int {
	if IsNum(a) && IsNum(b) {
		return _cmp(ToNum(a), ToNum(b))
	}
	return _cmp(ToStr(a), ToStr(b))
}

func _cmp[T int64 | float64 | string](a, b T) int {
	if a < b {
		return -1
	} else if a > b {
		return +1
	}
	return 0
}

func ParseTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	switch {
	case len(s) == 0:
		return time.Time{}, nil

	case strings.IndexByte(s, 'Z') > 0:
		return time.Parse(time.RFC3339, s)

	case strings.IndexByte(s, 'T') > 0:
		s = strings.Replace(s, "T", " ", 1)
	}
	//
	switch len(s) {
	case 0:
		return time.Time{}, nil
	case 6:
		return time.Parse("020106", s)
	case 8:
		return time.Parse("02.01.06", s)
	case 10:
		switch {
		case strings.IndexByte(s, '.') > 0:
			return time.Parse("02.01.2006", s)
		case strings.IndexByte(s, '-') > 0:
			return time.Parse("2006-01-02", s)
		case strings.IndexByte(s, '/') > 0:
			return time.Parse("2006/01/02", s)
		default: // unix timestamp
			ts, err := strconv.ParseInt(s, 10, 64)
			return time.Unix(ts, 0), err
		}

	case 16:
		switch {
		case strings.IndexByte(s, '-') > 0:
			return time.Parse("2006-01-02 15:04", s)
		case strings.IndexByte(s, '/') > 0:
			return time.Parse("2006/01/02 15:04", s)
		case strings.IndexByte(s, '.') > 0:
			return time.Parse("02.01.2006 15:04", s)
		}

	case 19:
		switch {
		case strings.IndexByte(s, '-') > 0:
			return time.Parse("2006-01-02 15:04:05", s)
		case strings.IndexByte(s, '/') > 0:
			return time.Parse("2006/01/02 15:04:05", s)
		case strings.IndexByte(s, '.') > 0:
			return time.Parse("02.01.2006 15:04:05", s)
		}

	case 29:
		switch {
		case strings.IndexByte(s, '-') > 0:
			return time.Parse("2006-01-02 15:04:05.999999999", s)
		case strings.IndexByte(s, '/') > 0:
			return time.Parse("2006/01/02 15:04:05.999999999", s)
		case strings.IndexByte(s, '.') > 0:
			return time.Parse("02.01.2006 15:04:05.999999999", s)
		}

	//case len(time.RFC1123):
	//	return time.Parse(time.RFC1123, s)
	case len(time.RFC1123Z):
		return time.Parse(time.RFC1123Z, s)
	case len(time.RFC3339):
		return time.Parse(time.RFC3339, s)

	case len(time.RFC3339Nano):
		return time.Parse(time.RFC3339Nano, s)
	}
	return time.Time{}, fmt.Errorf("unknown time-format: `%s`", s)
}
