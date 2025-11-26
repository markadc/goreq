package goreq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

// ==================== 全局配置 ====================
var (
	Timeout = 30 * time.Second
	Proxy   string
	Headers = make(http.Header)
)

// ==================== 类型别名（写起来最爽）================
type P map[string]string // Params
type J map[string]any    // JSON body
type F map[string]string // Form body
type H map[string]string // Headers

// ==================== Response ====================
type Response struct {
	*http.Response
	body []byte
	err  error
}

func (r *Response) OK() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

func (r *Response) RaiseForStatus() {
	if !r.OK() {
		panic(fmt.Errorf("HTTP %d: %s", r.StatusCode, r.Text()))
	}
}

func (r *Response) Text() string       { return string(r.body) }
func (r *Response) Bytes() []byte      { return r.body }
func (r *Response) Json() gjson.Result { return gjson.ParseBytes(r.body) }

// 一行保存文件（支持超大文件、自动创建目录）
func (r *Response) Save(filepath string) error {
	if !r.OK() {
		return fmt.Errorf("bad status %d, cannot save", r.StatusCode)
	}
	if err := os.MkdirAll(path.Dir(filepath), 0755); err != nil {
		return err
	}
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, bytes.NewReader(r.body))
	return err
}

// ==================== Session ====================
type Session struct {
	client  *http.Client
	headers http.Header
}

func NewSession() *Session {
	jar, _ := cookiejar.New(nil)
	s := &Session{
		client:  &http.Client{Jar: jar, Timeout: Timeout},
		headers: Headers.Clone(),
	}
	if Proxy != "" {
		u, _ := url.Parse(Proxy)
		s.client.Transport = &http.Transport{Proxy: http.ProxyURL(u)}
	}
	return s
}

func (s *Session) SetHeader(k, v string) { s.headers.Set(k, v) }

func (s *Session) do(method, rawurl string, body any, extra ...any) *Response {
	resp, _ := s.request(method, rawurl, body, extra...)
	return resp
}

func (s *Session) request(method, rawurl string, body any, extra ...any) (*Response, error) {
	var params P
	var headers H

	for _, e := range extra {
		// 使用反射检查具体类型，区分 P 和 H
		t := reflect.TypeOf(e)
		if t != nil {
			switch t.String() {
			case "goreq.P":
				if v, ok := e.(P); ok {
					params = v
				}
			case "goreq.H":
				if v, ok := e.(H); ok {
					headers = v
				}
			}
		}
	}

	// 处理 query 参数
	if params != nil {
		q := url.Values{}
		for k, v := range params {
			q.Set(k, v)
		}
		sep := "?"
		if strings.Contains(rawurl, "?") {
			sep = "&"
		}
		rawurl += sep + q.Encode()
	}

	// 处理 body
	var reqBody io.Reader
	var contentType string

	if body != nil {
		switch v := body.(type) {
		case J:
			b, _ := json.Marshal(v)
			reqBody = bytes.NewReader(b)
			contentType = "application/json"
		case map[string]any:
			b, _ := json.Marshal(v)
			reqBody = bytes.NewReader(b)
			contentType = "application/json"
		case F:
			f := url.Values{}
			for k, vv := range v {
				f.Set(k, vv)
			}
			reqBody = strings.NewReader(f.Encode())
			contentType = "application/x-www-form-urlencoded"
		case map[string]string:
			f := url.Values{}
			for k, vv := range v {
				f.Set(k, vv)
			}
			reqBody = strings.NewReader(f.Encode())
			contentType = "application/x-www-form-urlencoded"
		case string:
			reqBody = strings.NewReader(v)
		case []byte:
			reqBody = bytes.NewReader(v)
		}
	}

	req, _ := http.NewRequest(method, rawurl, reqBody)
	req.Header = s.headers.Clone()
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return &Response{err: err}, err
	}
	defer resp.Body.Close()

	b, _ := io.ReadAll(resp.Body)
	return &Response{Response: resp, body: b}, nil
}

func (s *Session) Get(u string, extra ...any) *Response { return s.do("GET", u, nil, extra...) }
func (s *Session) Post(u string, data any, extra ...any) *Response {
	return s.do("POST", u, data, extra...)
}
func (s *Session) Put(u string, data any, extra ...any) *Response {
	return s.do("PUT", u, data, extra...)
}
func (s *Session) Delete(u string, data any, extra ...any) *Response {
	return s.do("DELETE", u, data, extra...)
}

// ==================== 全局快捷函数 ====================
func do(method, u string, body any, extra ...any) *Response {
	s := NewSession()
	s.headers = Headers.Clone()
	r, _ := s.request(method, u, body, extra...)
	return r
}

func Get(u string, extra ...any) *Response { return do("GET", u, nil, extra...) }
func Post(u string, data any, extra ...any) *Response {
	return do("POST", u, data, extra...)
}
func Put(u string, data any, extra ...any) *Response { return do("PUT", u, data, extra...) }
func Delete(u string, data any, extra ...any) *Response {
	return do("DELETE", u, data, extra...)
}

func SetHeader(k, v string) { Headers.Set(k, v) }
