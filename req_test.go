package goreq

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestAll 综合测试所有功能
func TestAll(t *testing.T) {
	// ==================== GET 请求测试 ====================
	fmt.Println("🚀 开始测试 GET 请求...")

	resp := Get("https://httpbin.org/get")
	fmt.Printf("📡 GET 请求状态码: %d\n", resp.StatusCode)

	if !resp.OK() {
		t.Fatalf("❌ GET 请求失败，期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	json := resp.Json()
	url := json.Get("url").String()
	if !strings.Contains(url, "httpbin.org/get") {
		t.Errorf("❌ GET 请求 URL 验证失败，期望包含 'httpbin.org/get'，实际得到 '%s'", url)
	} else {
		fmt.Println("✅ GET 请求测试通过")
	}

	// ==================== GET 带参数请求测试 ====================
	fmt.Println("🚀 开始测试带参数的 GET 请求...")

	params := P{
		"name":    "goreq",
		"version": "0.1",
		"author":  "markadc",
	}

	resp = Get("https://httpbin.org/get", params)
	fmt.Printf("📡 带参数 GET 请求状态码: %d\n", resp.StatusCode)

	if !resp.OK() {
		t.Fatalf("❌ 带参数 GET 请求失败，期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	json = resp.Json()
	if json.Get("args.name").String() != "goreq" {
		t.Errorf("❌ 查询参数 name 验证失败，期望 'goreq'，实际得到 '%s'", json.Get("args.name").String())
	}
	if json.Get("args.version").String() != "0.1" {
		t.Errorf("❌ 查询参数 version 验证失败，期望 '0.1'，实际得到 '%s'", json.Get("args.version").String())
	}
	if json.Get("args.author").String() != "markadc" {
		t.Errorf("❌ 查询参数 author 验证失败，期望 'markadc'，实际得到 '%s'", json.Get("args.author").String())
	}
	fmt.Println("✅ 带参数 GET 请求测试通过")

	// ==================== POST JSON 请求测试 ====================
	fmt.Println("🚀 开始测试 POST JSON 请求...")

	data := J{
		"name":        "goreq",
		"version":     "0.1.0",
		"author":      "markadc",
		"description": "简单易用的 Go HTTP 客户端",
	}
	fmt.Printf("📤 发送 JSON 数据: %+v\n", data)

	resp = Post("https://httpbin.org/post", data)
	fmt.Printf("📡 POST JSON 请求状态码: %d\n", resp.StatusCode)

	if !resp.OK() {
		t.Fatalf("❌ POST JSON 请求失败，期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	json = resp.Json()
	contentType := json.Get("headers.Content-Type").String()
	if contentType != "application/json" {
		t.Errorf("❌ Content-Type 验证失败，期望 'application/json'，实际得到 '%s'", contentType)
	}
	if json.Get("json.name").String() != "goreq" {
		t.Errorf("❌ JSON 字段 name 验证失败，期望 'goreq'，实际得到 '%s'", json.Get("json.name").String())
	}
	fmt.Println("✅ POST JSON 请求测试通过")

	// ==================== POST 表单请求测试 ====================
	fmt.Println("🚀 开始测试 POST 表单请求...")

	formData := F{
		"username": "admin",
		"password": "secret123",
		"email":    "admin@example.com",
		"role":     "administrator",
	}
	fmt.Printf("📤 发送表单数据: %+v\n", formData)

	resp = Post("https://httpbin.org/post", formData)
	fmt.Printf("📡 POST 表单请求状态码: %d\n", resp.StatusCode)

	if !resp.OK() {
		t.Fatalf("❌ POST 表单请求失败，期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	json = resp.Json()
	contentType = json.Get("headers.Content-Type").String()
	if contentType != "application/x-www-form-urlencoded" {
		t.Errorf("❌ 表单 Content-Type 验证失败，期望 'application/x-www-form-urlencoded'，实际得到 '%s'", contentType)
	}
	if json.Get("form.username").String() != "admin" {
		t.Errorf("❌ 表单字段 username 验证失败，期望 'admin'，实际得到 '%s'", json.Get("form.username").String())
	}
	fmt.Println("✅ POST 表单请求测试通过")

	// ==================== 文件保存测试 ====================
	fmt.Println("🚀 开始测试文件保存功能...")

	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_download.json")
	fmt.Printf("📁 临时文件路径: %s\n", filePath)

	resp = Get("https://httpbin.org/json")
	fmt.Printf("📡 文件下载请求状态码: %d\n", resp.StatusCode)

	if !resp.OK() {
		t.Fatalf("❌ 文件下载请求失败，期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	err := resp.Save(filePath)
	if err != nil {
		t.Fatalf("❌ 保存文件失败: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("❌ 文件不存在: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("❌ 读取文件失败: %v", err)
	}

	if len(content) == 0 {
		t.Error("❌ 文件内容为空")
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "{") || !strings.Contains(contentStr, "}") {
		t.Error("❌ 文件内容不是有效的 JSON 格式")
	}

	// 清理临时文件
	err = os.Remove(filePath)
	if err != nil {
		fmt.Printf("⚠️ 清理临时文件失败: %v\n", err)
	} else {
		fmt.Println("🧹 临时文件清理完成")
	}
	fmt.Println("✅ 文件保存功能测试通过")

	// ==================== 自动创建目录保存测试 ====================
	fmt.Println("🚀 开始测试自动创建目录的保存功能...")

	tempDir = t.TempDir()
	filePath = filepath.Join(tempDir, "downloads", "data", "test_file.txt")
	fmt.Printf("📁 多层目录文件路径: %s\n", filePath)

	resp = Get("https://httpbin.org/robots.txt")
	fmt.Printf("📡 文本下载请求状态码: %d\n", resp.StatusCode)

	if !resp.OK() {
		t.Fatalf("❌ 文本下载请求失败，期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	err = resp.Save(filePath)
	if err != nil {
		t.Fatalf("❌ 保存文件失败: %v", err)
	}

	dirPath := filepath.Dir(filePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Errorf("❌ 目录未被创建: %s", dirPath)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("❌ 文件不存在: %s", filePath)
	}

	content, err = os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("❌ 读取文件失败: %v", err)
	}

	if len(content) == 0 {
		t.Error("❌ 文件内容为空")
	} else {
		fmt.Printf("✅ 文件内容长度: %d 字节\n", len(content))
	}

	// 清理临时目录
	err = os.RemoveAll(filepath.Dir(filePath))
	if err != nil {
		fmt.Printf("⚠️ 清理临时目录失败: %v\n", err)
	} else {
		fmt.Println("🧹 临时目录清理完成")
	}
	fmt.Println("✅ 自动创建目录保存功能测试通过")

	// ==================== 单元测试部分 ====================
	fmt.Println("🚀 开始单元测试...")

	// HTTP 方法测试
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response"))
	}))
	defer server.Close()

	// GET 方法测试
	resp = Get(server.URL)
	if !resp.OK() {
		t.Errorf("❌ 单元测试 GET 失败，期望状态码为 OK，实际得到 %d", resp.StatusCode)
	}

	// POST 方法测试
	resp = Post(server.URL, J{"test": "data"})
	if !resp.OK() {
		t.Errorf("❌ 单元测试 POST 失败，期望状态码为 OK，实际得到 %d", resp.StatusCode)
	}

	// PUT 方法测试
	resp = Put(server.URL, J{"update": "data"})
	if !resp.OK() {
		t.Errorf("❌ 单元测试 PUT 失败，期望状态码为 OK，实际得到 %d", resp.StatusCode)
	}

	// DELETE 方法测试
	resp = Delete(server.URL, nil)
	if !resp.OK() {
		t.Errorf("❌ 单元测试 DELETE 失败，期望状态码为 OK，实际得到 %d", resp.StatusCode)
	}
	fmt.Println("✅ HTTP 方法单元测试通过")

	// ==================== 参数和头部测试 ====================
	fmt.Println("🚀 开始参数和头部测试...")

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证查询参数
		if r.URL.Query().Get("key") != "value" {
			t.Errorf("❌ 查询参数验证失败，期望 key=value，实际得到 %s", r.URL.Query().Get("key"))
		}
		// 验证请求头
		if r.Header.Get("X-Test") != "test-value" {
			t.Errorf("❌ 请求头验证失败，期望 X-Test=test-value，实际得到 %s", r.Header.Get("X-Test"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server2.Close()

	resp = Get(server2.URL, P{"key": "value"}, H{"X-Test": "test-value"})
	if !resp.OK() {
		t.Errorf("❌ 参数和头部测试失败，期望状态码为 OK，实际得到 %d", resp.StatusCode)
	}
	fmt.Println("✅ 参数和头部测试通过")

	// ==================== 会话测试 ====================
	fmt.Println("🚀 开始会话测试...")

	server3 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Session-Header") != "session-value" {
			t.Errorf("❌ 会话头部验证失败，期望 'session-value'，实际得到 '%s'", r.Header.Get("X-Session-Header"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server3.Close()

	s := NewSession()
	s.SetHeader("X-Session-Header", "session-value")
	resp = s.Get(server3.URL)
	if !resp.OK() {
		t.Errorf("❌ 会话测试失败，期望状态码为 OK，实际得到 %d", resp.StatusCode)
	}
	fmt.Println("✅ 会话测试通过")

	// ==================== Cookie 测试 ====================
	fmt.Println("🚀 开始 Cookie 测试...")

	var cookieValue string
	server4 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
				t.Error("❌ 期望 Cookie 被设置")
			} else {
				cookieValue = cookie.Value
			}
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer server4.Close()

	s = NewSession()
	s.Get(server4.URL + "/set-cookie")
	s.Get(server4.URL + "/check-cookie")

	if cookieValue != "test-session-id" {
		t.Errorf("❌ Cookie 测试失败，期望 'test-session-id'，实际得到 '%s'", cookieValue)
	}
	fmt.Println("✅ Cookie 测试通过")

	// ==================== 超时测试 ====================
	fmt.Println("🚀 开始超时测试...")

	originalTimeout := Timeout
	defer func() { Timeout = originalTimeout }()
	Timeout = 100 * time.Millisecond

	server5 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server5.Close()

	resp = Get(server5.URL)
	if resp.err == nil {
		t.Error("❌ 期望超时错误，但没有发生")
	} else {
		fmt.Println("✅ 超时测试通过")
	}

	// ==================== 全局头部测试 ====================
	fmt.Println("🚀 开始全局头部测试...")

	originalHeaders := Headers
	defer func() { Headers = originalHeaders }()
	Headers = make(http.Header)
	SetHeader("X-Global-Header", "global-value")

	server6 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Global-Header") != "global-value" {
			t.Errorf("❌ 全局头部验证失败，期望 'global-value'，实际得到 '%s'", r.Header.Get("X-Global-Header"))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server6.Close()

	resp = Get(server6.URL)
	if !resp.OK() {
		t.Errorf("❌ 全局头部测试失败，期望状态码为 OK，实际得到 %d", resp.StatusCode)
	}
	fmt.Println("✅ 全局头部测试通过")

	// ==================== JSON 响应解析测试 ====================
	fmt.Println("🚀 开始 JSON 响应解析测试...")

	server7 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"name":"goreq","version":"0.1","nested":{"key":"value"}}`))
	}))
	defer server7.Close()

	resp = Get(server7.URL)
	json = resp.Json()

	if json.Get("name").String() != "goreq" {
		t.Errorf("❌ JSON 解析测试失败，期望 name='goreq'，实际得到 '%s'", json.Get("name").String())
	}
	if json.Get("version").String() != "0.1" {
		t.Errorf("❌ JSON 解析测试失败，期望 version='0.1'，实际得到 '%s'", json.Get("version").String())
	}
	if json.Get("nested.key").String() != "value" {
		t.Errorf("❌ JSON 解析测试失败，期望 nested.key='value'，实际得到 '%s'", json.Get("nested.key").String())
	}
	fmt.Println("✅ JSON 响应解析测试通过")

	// ==================== 异常状态处理测试 ====================
	fmt.Println("🚀 开始异常状态处理测试...")

	server8 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server8.Close()

	resp = Get(server8.URL)

	// 测试 panic 是否被正确触发
	defer func() {
		if r := recover(); r == nil {
			t.Error("❌ 期望 RaiseForStatus 触发 panic，但没有")
		} else {
			fmt.Println("✅ 异常状态处理测试通过")
		}
	}()

	resp.RaiseForStatus()

	fmt.Println("🎉 所有测试完成！")
}
