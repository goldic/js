package js

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func ParseFile(filename string) (v Value, err error) {
	err = UnmarshalFile(filename, &v.val)
	return
}

func UnmarshalFile(filename string, v any) (err error) {
	defer catch(&err)
	f := tryVal(os.Open(filename))
	defer f.Close()
	try(json.NewDecoder(f).Decode(v))
	return
}

func MarshalToFile(filename string, v any) (err error) {
	defer catch(&err)
	f := tryVal(os.Create(filename))
	defer f.Close()
	try(json.NewEncoder(f).Encode(v))
	return
}

func MarshalIndentToFile(filename string, v any) (err error) {
	defer catch(&err)
	f := tryVal(os.Create(filename))
	defer f.Close()
	e := json.NewEncoder(f)
	e.SetIndent("", "  ")
	try(e.Encode(v))
	return
}

func Write(w io.Writer, v any) (err error) {
	defer catch(&err)
	if rw, ok := w.(http.ResponseWriter); ok {
		rw.Header().Set("Content-Type", "application/json")
	}
	try(json.NewEncoder(w).Encode(v))
	return
}
