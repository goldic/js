package js

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	url2 "net/url"
	"os"
	"path/filepath"
	"strings"
)

const http2Proto = "http2:"

var http2Transport = &http.Transport{ // use HTTP/2
	TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
}

func Load(url string) (Value, error) {
	return RequestValue(http.MethodGet, url, nil, nil)
}

func LoadObject(url string) (Object, error) {
	v, err := Load(url)
	return v.Object(), err
}

func PostData(url string, jsonPostData any) (Value, error) {
	return RequestValue(http.MethodPost, url, nil, jsonPostData)
}

func RequestValue(method, url string, headers Object, body any) (val Value, err error) {
	defer catch(&err)
	return Parse(must(Request(method, url, headers, body)))
}

func Request(method, url string, headers Object, body any) (data []byte, err error) {
	defer catch(&err)

	client := http.DefaultClient
	if strings.HasPrefix(url, http2Proto) {
		c2 := *client // copy client
		c2.Transport = http2Transport
		client, url = &c2, "https:"+strings.TrimPrefix(url, http2Proto)
	}
	trace := strings.HasSuffix(url, "#trace")
	url = strings.TrimSuffix(url, "#trace")
	var contType = ""
	var reqBody io.Reader
	if !isNil(body) {
		if method == "" {
			method = http.MethodPost
		}
		switch v := body.(type) {
		case url2.Values:
			reqBody, contType = bytes.NewBufferString(v.Encode()), "application/x-www-form-urlencoded"
		case io.Reader:
			reqBody = v
		case []byte:
			reqBody = bytes.NewBuffer(v)
		case string:
			reqBody = bytes.NewBufferString(v)
		default:
			reqBody, contType = bytes.NewBuffer(must(json.Marshal(body))), "application/json"
		}
	}
	if method == "" {
		method = http.MethodGet
	}
	req := must(http.NewRequest(method, url, reqBody))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept-Encoding", "gzip, deflate")
	if contType != "" {
		req.Header.Set("Content-Type", contType)
	}
	if headers != nil {
		for name := range headers {
			if v := headers.Get(name); v.IsArray() {
				req.Header[textproto.CanonicalMIMEHeaderKey(name)] = v.Array().Strings()
			} else {
				req.Header.Set(name, headers.GetStr(name))
			}
		}
	}
	if trace {
		log.Printf("js> %s %s %s %s", req.Proto, method, url, encHeader(req.Header))
		if !isNil(reqBody) {
			log.Printf("js> http-Request-Body: %s", reqBody)
		}
	}
	resp := must(client.Do(req))
	defer resp.Body.Close()
	if trace {
		log.Printf("js> http-Response: %v `%v` %s", resp.StatusCode, resp.Status, encHeader(resp.Header))
	}
	if resp.StatusCode != 200 {
		panic(fmt.Errorf("js> http-ERROR: StatusCode=%v `%s`", resp.StatusCode, resp.Status))
	}
	var respReader io.ReadCloser
	switch resp.Header.Get("Content-Encoding") {
	case "":
		respReader = resp.Body
	case "gzip":
		respReader = must(gzip.NewReader(resp.Body))
		defer respReader.Close()
	case "deflate":
		respReader = flate.NewReader(resp.Body)
		defer respReader.Close()
	//case "br":
	//	respReader = must(brotli.NewReader(resp.Body, nil))
	//	defer respReader.Close()
	//case "compress": ...
	//case "sdch": ...
	default:
		check(fmt.Errorf("js.Request: Unknown Content-Encoding `%s`", resp.Header.Get("Content-Encoding")))
	}
	data = readAll(respReader)
	if trace {
		log.Printf("js> http-Response:\n%s", string(data))
	}
	return
}

func PostMultipart(method, url string, headers Object, params Object, files map[string]string) (res Value, err error) {
	defer catch(&err)
	contType, body := makeMultipartBody(params, files)
	headers = headers.Set("Content-Type", contType)
	return RequestValue(method, url, headers, body)
}

func encHeader(h http.Header) (res string) {
	for name, vv := range h {
		for _, v := range vv {
			res += fmt.Sprintf("\n - %s: %s", name, v)
		}
	}
	return
}

func makeMultipartBody(params Object, files map[string]string) (contType string, body io.Reader) {
	b := bytes.NewBuffer(nil)
	w := multipart.NewWriter(b)
	if params != nil {
		for key, value := range params {
			part := must(w.CreateFormField(key))
			must(part.Write([]byte(ToStr(value))))
		}
	}
	if files != nil {
		for key, filePath := range files {
			file := must(os.Open(filePath))
			defer file.Close()
			part := must(w.CreateFormFile(key, filepath.Base(file.Name())))
			must(io.Copy(part, file))
		}
	}
	check(w.Close())
	return w.FormDataContentType(), b
}
