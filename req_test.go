package goreq

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))
	defer server.Close()

	resp := Get(server.URL)
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
	if resp.Text() != "Hello, World!" {
		t.Errorf("Expected 'Hello, World!', got '%s'", resp.Text())
	}
}

func TestPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	resp := Post(server.URL, J{"name": "test"})
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
	if resp.Json().Get("status").String() != "ok" {
		t.Errorf("Expected status 'ok', got '%s'", resp.Json().Get("status").String())
	}
}

func TestPostForm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("Expected form content type, got %s", r.Header.Get("Content-Type"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp := Post(server.URL, F{"username": "admin", "password": "123456"})
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
}

func TestPut(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp := Put(server.URL, J{"data": "test"})
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp := Delete(server.URL, nil)
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
}

func TestQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("name") != "goreq" {
			t.Errorf("Expected query param name=goreq, got %s", r.URL.Query().Get("name"))
		}
		if r.URL.Query().Get("version") != "0.1" {
			t.Errorf("Expected query param version=0.1, got %s", r.URL.Query().Get("version"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp := Get(server.URL, P{"name": "goreq", "version": "0.1"})
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
}

func TestHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("User-Agent") != "TestAgent" {
			t.Errorf("Expected User-Agent header 'TestAgent', got '%s'", r.Header.Get("User-Agent"))
		}
		if r.Header.Get("X-Custom") != "CustomValue" {
			t.Errorf("Expected X-Custom header 'CustomValue', got '%s'", r.Header.Get("X-Custom"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp := Get(server.URL, H{"User-Agent": "TestAgent", "X-Custom": "CustomValue"})
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
}

func TestResponseOk(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp := Get(server.URL)
	if !resp.Ok() {
		t.Error("Expected Ok() to return true for 200 status")
	}

	server404 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server404.Close()

	resp404 := Get(server404.URL)
	if resp404.Ok() {
		t.Error("Expected Ok() to return false for 404 status")
	}
}

func TestResponseRaiseForStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer server.Close()

	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected RaiseForStatus to panic for 404 status")
		}
	}()

	resp := Get(server.URL)
	resp.RaiseForStatus()
}

func TestResponseBytes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test data"))
	}))
	defer server.Close()

	resp := Get(server.URL)
	bytes := resp.Bytes()
	if string(bytes) != "test data" {
		t.Errorf("Expected 'test data', got '%s'", string(bytes))
	}
}

func TestResponseJson(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"name":"goreq","version":"0.1","nested":{"key":"value"}}`))
	}))
	defer server.Close()

	resp := Get(server.URL)
	json := resp.Json()

	if json.Get("name").String() != "goreq" {
		t.Errorf("Expected name 'goreq', got '%s'", json.Get("name").String())
	}
	if json.Get("version").String() != "0.1" {
		t.Errorf("Expected version '0.1', got '%s'", json.Get("version").String())
	}
	if json.Get("nested.key").String() != "value" {
		t.Errorf("Expected nested.key 'value', got '%s'", json.Get("nested.key").String())
	}
}

func TestResponseSave(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("file content"))
	}))
	defer server.Close()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "subdir", "test.txt")

	resp := Get(server.URL)
	err := resp.Save(filePath)
	if err != nil {
		t.Fatalf("Failed to save file: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read saved file: %v", err)
	}

	if string(content) != "file content" {
		t.Errorf("Expected 'file content', got '%s'", string(content))
	}
}

func TestSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Session-Header") != "session-value" {
			t.Errorf("Expected session header, got '%s'", r.Header.Get("X-Session-Header"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	s := NewSession()
	s.SetHeader("X-Session-Header", "session-value")

	resp := s.Get(server.URL)
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
}

func TestSessionCookies(t *testing.T) {
	var cookieValue string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/set-cookie":
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: "test-session-id",
			})
			w.WriteHeader(http.StatusOK)
		case "/check-cookie":
			cookie, err := r.Cookie("session")
			if err != nil {
				t.Error("Expected cookie to be set")
			} else {
				cookieValue = cookie.Value
			}
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server.Close()

	s := NewSession()
	s.Get(server.URL + "/set-cookie")
	s.Get(server.URL + "/check-cookie")

	if cookieValue != "test-session-id" {
		t.Errorf("Expected cookie value 'test-session-id', got '%s'", cookieValue)
	}
}

func TestGlobalTimeout(t *testing.T) {
	originalTimeout := Timeout
	defer func() { Timeout = originalTimeout }()

	Timeout = 100 * time.Millisecond

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp := Get(server.URL)
	if resp.err == nil {
		t.Error("Expected timeout error")
	}
}

func TestGlobalHeaders(t *testing.T) {
	originalHeaders := Headers
	defer func() { Headers = originalHeaders }()

	Headers = make(http.Header)
	SetHeader("X-Global-Header", "global-value")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Global-Header") != "global-value" {
			t.Errorf("Expected global header, got '%s'", r.Header.Get("X-Global-Header"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp := Get(server.URL)
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
}

func TestParamsAndHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("key") != "value" {
			t.Errorf("Expected query param key=value, got %s", r.URL.Query().Get("key"))
		}
		if r.Header.Get("X-Test") != "test-value" {
			t.Errorf("Expected header X-Test=test-value, got %s", r.Header.Get("X-Test"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp := Get(server.URL, P{"key": "value"}, H{"X-Test": "test-value"})
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
}

func TestPostString(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp := Post(server.URL, "raw string data")
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
}

func TestPostBytes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	resp := Post(server.URL, []byte("raw bytes data"))
	if !resp.Ok() {
		t.Errorf("Expected OK status, got %d", resp.StatusCode)
	}
}
