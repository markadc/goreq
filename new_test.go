package goreq

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestNewGet 测试 GET 请求功能
func TestNewGet(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（短模式）")
	}

	fmt.Println("🚀 开始测试 GET 请求...")

	// 测试基本 GET 请求
	resp := Get("https://httpbin.org/get")
	fmt.Printf("📡 请求状态码: %d\n", resp.StatusCode)

	if !resp.Ok() {
		t.Fatalf("❌ 期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	json := resp.Json()
	url := json.Get("url").String()
	if !strings.Contains(url, "httpbin.org/get") {
		t.Errorf("❌ 期望 URL 包含 'httpbin.org/get'，实际得到 '%s'", url)
	} else {
		fmt.Println("✅ GET 请求 URL 验证通过")
	}

	fmt.Println("🎉 GET 请求测试完成！")
}

// TestNewGetWithParams 测试带参数的 GET 请求
func TestNewGetWithParams(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（短模式）")
	}

	fmt.Println("🚀 开始测试带参数的 GET 请求...")

	// 测试带查询参数的 GET 请求
	params := P{
		"name":    "goreq",
		"version": "0.1",
		"author":  "markadc",
	}

	resp := Get("https://httpbin.org/get", params)
	fmt.Printf("📡 请求状态码: %d\n", resp.StatusCode)

	if !resp.Ok() {
		t.Fatalf("❌ 期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	json := resp.Json()

	// 验证查询参数
	if json.Get("args.name").String() != "goreq" {
		t.Errorf("❌ 期望 args.name='goreq'，实际得到 '%s'", json.Get("args.name").String())
	} else {
		fmt.Println("✅ 查询参数 name 验证通过")
	}

	if json.Get("args.version").String() != "0.1" {
		t.Errorf("❌ 期望 args.version='0.1'，实际得到 '%s'", json.Get("args.version").String())
	} else {
		fmt.Println("✅ 查询参数 version 验证通过")
	}

	if json.Get("args.author").String() != "markadc" {
		t.Errorf("❌ 期望 args.author='markadc'，实际得到 '%s'", json.Get("args.author").String())
	} else {
		fmt.Println("✅ 查询参数 author 验证通过")
	}

	fmt.Println("🎉 带参数的 GET 请求测试完成！")
}

// TestNewPost 测试 POST 请求功能
func TestNewPost(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（短模式）")
	}

	fmt.Println("🚀 开始测试 POST 请求...")

	// 测试 POST JSON 数据
	data := J{
		"name":        "goreq",
		"version":     "0.1.0",
		"author":      "markadc",
		"description": "简单易用的 Go HTTP 客户端",
	}
	fmt.Printf("📤 发送 JSON 数据: %+v\n", data)

	resp := Post("https://httpbin.org/post", data)
	fmt.Printf("📡 请求状态码: %d\n", resp.StatusCode)

	if !resp.Ok() {
		t.Fatalf("❌ 期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	json := resp.Json()
	contentType := json.Get("headers.Content-Type").String()
	fmt.Printf("📋 Content-Type: %s\n", contentType)

	// 验证 Content-Type
	if contentType != "application/json" {
		t.Errorf("❌ 期望 Content-Type='application/json'，实际得到 '%s'", contentType)
	} else {
		fmt.Println("✅ Content-Type 验证通过")
	}

	// 验证 JSON 数据
	if json.Get("json.name").String() != "goreq" {
		t.Errorf("❌ 期望 json.name='goreq'，实际得到 '%s'", json.Get("json.name").String())
	} else {
		fmt.Println("✅ JSON 字段 name 验证通过")
	}

	if json.Get("json.version").String() != "0.1.0" {
		t.Errorf("❌ 期望 json.version='0.1.0'，实际得到 '%s'", json.Get("json.version").String())
	} else {
		fmt.Println("✅ JSON 字段 version 验证通过")
	}

	if json.Get("json.author").String() != "markadc" {
		t.Errorf("❌ 期望 json.author='markadc'，实际得到 '%s'", json.Get("json.author").String())
	} else {
		fmt.Println("✅ JSON 字段 author 验证通过")
	}

	fmt.Println("🎉 POST JSON 请求测试完成！")
}

// TestNewPostForm 测试 POST 表单数据
func TestNewPostForm(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（短模式）")
	}

	fmt.Println("🚀 开始测试 POST 表单请求...")

	// 测试 POST 表单数据
	formData := F{
		"username": "admin",
		"password": "secret123",
		"email":    "admin@example.com",
		"role":     "administrator",
	}
	fmt.Printf("📤 发送表单数据: %+v\n", formData)

	resp := Post("https://httpbin.org/post", formData)
	fmt.Printf("📡 请求状态码: %d\n", resp.StatusCode)

	if !resp.Ok() {
		t.Fatalf("❌ 期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	json := resp.Json()
	contentType := json.Get("headers.Content-Type").String()
	fmt.Printf("📋 Content-Type: %s\n", contentType)

	// 验证 Content-Type
	if contentType != "application/x-www-form-urlencoded" {
		t.Errorf("❌ 期望 Content-Type='application/x-www-form-urlencoded'，实际得到 '%s'", contentType)
	} else {
		fmt.Println("✅ Content-Type 验证通过")
	}

	// 验证表单数据
	if json.Get("form.username").String() != "admin" {
		t.Errorf("❌ 期望 form.username='admin'，实际得到 '%s'", json.Get("form.username").String())
	} else {
		fmt.Println("✅ 表单字段 username 验证通过")
	}

	if json.Get("form.password").String() != "secret123" {
		t.Errorf("❌ 期望 form.password='secret123'，实际得到 '%s'", json.Get("form.password").String())
	} else {
		fmt.Println("✅ 表单字段 password 验证通过")
	}

	if json.Get("form.email").String() != "admin@example.com" {
		t.Errorf("❌ 期望 form.email='admin@example.com'，实际得到 '%s'", json.Get("form.email").String())
	} else {
		fmt.Println("✅ 表单字段 email 验证通过")
	}

	fmt.Println("🎉 POST 表单请求测试完成！")
}

// TestNewSave 测试文件保存功能
func TestNewSave(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（短模式）")
	}

	fmt.Println("🚀 开始测试文件保存功能...")

	// 创建临时目录
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test_download.json")
	fmt.Printf("📁 临时文件路径: %s\n", filePath)

	// 下载一个 JSON 文件
	resp := Get("https://httpbin.org/json")
	fmt.Printf("📡 请求状态码: %d\n", resp.StatusCode)

	if !resp.Ok() {
		t.Fatalf("❌ 期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	// 保存文件
	err := resp.Save(filePath)
	if err != nil {
		t.Fatalf("❌ 保存文件失败: %v", err)
	}
	fmt.Println("✅ 文件保存成功")

	// 验证文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("❌ 文件不存在: %s", filePath)
	} else {
		fmt.Println("✅ 文件存在验证通过")
	}

	// 读取文件内容并验证
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("❌ 读取文件失败: %v", err)
	}

	// 验证文件内容不为空
	if len(content) == 0 {
		t.Error("❌ 文件内容为空")
	} else {
		fmt.Printf("✅ 文件内容长度: %d 字节\n", len(content))
	}

	// 验证文件内容是有效的 JSON
	contentStr := string(content)
	if !strings.Contains(contentStr, "{") || !strings.Contains(contentStr, "}") {
		t.Error("❌ 文件内容不是有效的 JSON 格式")
	} else {
		fmt.Println("✅ 文件内容格式验证通过")
	}

	fmt.Println("🎉 文件保存功能测试完成！")
}

// TestNewSaveWithDirectory 测试自动创建目录的保存功能
func TestNewSaveWithDirectory(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（短模式）")
	}

	fmt.Println("🚀 开始测试自动创建目录的保存功能...")

	// 创建临时目录
	tempDir := t.TempDir()
	// 创建多层目录路径
	filePath := filepath.Join(tempDir, "downloads", "data", "test_file.txt")
	fmt.Printf("📁 多层目录文件路径: %s\n", filePath)

	// 下载文本内容
	resp := Get("https://httpbin.org/robots.txt")
	fmt.Printf("📡 请求状态码: %d\n", resp.StatusCode)

	if !resp.Ok() {
		t.Fatalf("❌ 期望状态码为 2xx，实际得到 %d", resp.StatusCode)
	}

	// 保存文件（应该自动创建目录）
	err := resp.Save(filePath)
	if err != nil {
		t.Fatalf("❌ 保存文件失败: %v", err)
	}
	fmt.Println("✅ 文件保存成功（自动创建目录）")

	// 验证目录是否被创建
	dirPath := filepath.Dir(filePath)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		t.Errorf("❌ 目录未被创建: %s", dirPath)
	} else {
		fmt.Println("✅ 目录自动创建验证通过")
	}

	// 验证文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("❌ 文件不存在: %s", filePath)
	} else {
		fmt.Println("✅ 文件存在验证通过")
	}

	// 读取并验证文件内容
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("❌ 读取文件失败: %v", err)
	}

	if len(content) == 0 {
		t.Error("❌ 文件内容为空")
	} else {
		fmt.Printf("✅ 文件内容长度: %d 字节\n", len(content))
		fmt.Printf("📄 文件内容预览: %s\n", string(content)[:min(100, len(content))])
	}

	fmt.Println("🎉 自动创建目录的保存功能测试完成！")
}

// min 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
