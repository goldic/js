package js

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"strings"
)

const http2Proto = "http2:"

var http2Transport = &http.Transport{ // use HTTP/2
	TLSNextProto: make(map[string]func(string, *tls.Conn) http.RoundTripper),
}

func Load(url string) (Value, error) {
	return Request(http.MethodGet, url, nil, nil)
}

func LoadObject(url string) (Object, error) {
	v, err := Load(url)
	return v.Object(), err
}

func PostData(url string, jsonPostData any) (Value, error) {
	return Request(http.MethodPost, url, nil, jsonPostData)
}

func Request(method, url string, headers Object, body any) (res Value, err error) {
	defer recoverErr(&err)

	client := http.DefaultClient
	if strings.HasPrefix(url, http2Proto) {
		c2 := *client // copy client
		c2.Transport = http2Transport
		client, url = &c2, "https:"+strings.TrimPrefix(url, http2Proto)
	}
	trace := strings.HasSuffix(url, "#trace")
	url = strings.TrimSuffix(url, "#trace")
	if method == "" {
		method = http.MethodGet
	}
	var contType = ""
	var reqBody io.Reader
	if !isNil(body) {
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
			reqBody, contType = bytes.NewBuffer(noerrVal(json.Marshal(body))), "application/json"
		}
	}
	req := noerrVal(http.NewRequest(method, url, reqBody))
	req.Header.Set("Accept", "*/*")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	if contType != "" {
		req.Header.Set("Content-Type", contType)
	}
	if headers != nil {
		for name := range headers {
			if arr := headers.Get(name).Array(); arr != nil {
				req.Header.Set(name, arr.Eq(0).String())
				for i, n := 1, arr.Len(); i < n; i++ {
					req.Header.Add(name, arr.Eq(i).String())
				}
			} else {
				req.Header.Set(name, headers.GetStr(name))
			}
		}
	}
	if trace {
		log.Printf("js> %s %s %s %s", req.Proto, method, url, IndentEncode(headers))
		if !isNil(reqBody) {
			log.Printf("js> http-Request-Body: %s", reqBody)
		}
	}
	resp := noerrVal(client.Do(req))
	defer resp.Body.Close()
	if trace {
		log.Printf("js> http-Response-StatusCode: %v `%v`", resp.StatusCode, resp.Status)
	}
	if resp.StatusCode != 200 {
		panic(fmt.Errorf("js> http-ERROR: StatusCode=%v `%s`", resp.StatusCode, resp.Status))
	}
	data := bytes.TrimSpace(readAll(resp.Body))
	if trace {
		log.Printf("js> http-Response:\n%s", string(data))
	}
	return Parse(data)
}
