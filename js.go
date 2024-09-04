package js

import (
	"encoding/json"
	"fmt"
	"io"
	"runtime"
)

func Marshal(v any) ([]byte, error) {
	return json.Marshal(v)
}

func Encode(v any) string {
	return string(tryVal(json.Marshal(v)))
}

func IndentEncode(v any) string {
	return string(tryVal(json.MarshalIndent(v, "", "  ")))
}

func try(err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(2)
		panic(fmt.Errorf("%w\n\t%s:%d", err, file, line))
	}
}

func tryVal[T any](v T, err error) T {
	if err != nil {
		_, file, line, _ := runtime.Caller(2)
		panic(fmt.Errorf("%w\n\t%s:%d", err, file, line))
	}
	return v
}

func catch(err *error) {
	if r := recover(); r != nil && err != nil && *err == nil {
		if e, ok := r.(error); ok {
			*err = e
		} else {
			*err = fmt.Errorf("%v", r)
		}
	}
}

func readAll(r io.Reader) []byte {
	if c, ok := r.(io.ReadCloser); ok && c != nil {
		defer c.Close()
	}
	return tryVal(io.ReadAll(r))
}
