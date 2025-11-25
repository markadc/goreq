# goreq

一个简洁优雅的 Go HTTP 请求库，灵感来自 Python requests。

## 特性

- 🚀 **简洁的 API** - 类似 Python requests 的使用体验
- 🎯 **类型别名** - `P`、`J`、`F`、`H` 让代码更简洁
- 🔄 **Session 支持** - 自动管理 Cookie
- 📦 **自动序列化** - 自动处理 JSON 和 Form 数据
- 🎨 **链式调用** - 支持优雅的链式操作
- 💾 **文件下载** - 一行代码保存响应到文件

## 安装

```bash
go get github.com/markadc/goreq
```

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "github.com/markadc/goreq"
)

func main() {
    // GET 请求
    resp := goreq.Get("https://api.github.com/users/markadc")
    fmt.Println(resp.Text())
    
    // POST JSON
    resp = goreq.Post("https://httpbin.org/post", goreq.J{
        "name": "goreq",
        "type": "http-client",
    })
    fmt.Println(resp.Json().Get("json.name").String())
    
    // POST Form
    resp = goreq.Post("https://httpbin.org/post", goreq.F{
        "username": "admin",
        "password": "123456",
    })
}
```

### 带参数和请求头

```go
// 使用类型别名让代码更简洁
resp := goreq.Get("https://api.github.com/search/repositories",
    goreq.P{"q": "golang", "sort": "stars"},  // 查询参数
    goreq.H{"User-Agent": "goreq/1.0"},       // 请求头
)
```

### Session 使用

```go
// 创建 Session（自动管理 Cookie）
s := goreq.NewSession()
s.SetHeader("User-Agent", "MyApp/1.0")

// 登录
s.Post("https://example.com/login", goreq.F{
    "username": "admin",
    "password": "123456",
})

// 后续请求会自动带上 Cookie
resp := s.Get("https://example.com/profile")
```

### 文件下载

```go
resp := goreq.Get("https://example.com/file.zip")
if err := resp.Save("/path/to/save/file.zip"); err != nil {
    panic(err)
}
```

### 错误处理

```go
resp := goreq.Get("https://httpbin.org/status/404")

// 检查状态码
if !resp.Ok() {
    fmt.Println("请求失败:", resp.StatusCode)
}

// 或者直接抛出异常
resp.RaiseForStatus()  // 状态码非 2xx 时 panic
```

### 全局配置

```go
import "time"

// 设置全局超时
goreq.Timeout = 10 * time.Second

// 设置全局代理
goreq.Proxy = "http://127.0.0.1:7890"

// 设置全局请求头
goreq.SetHeader("User-Agent", "MyApp/1.0")
```

## API 文档

### 类型别名

- `P` - 查询参数 (`map[string]string`)
- `J` - JSON 请求体 (`map[string]any`)
- `F` - Form 请求体 (`map[string]string`)
- `H` - 请求头 (`map[string]string`)

### 全局函数

- `Get(url, ...extra)` - GET 请求
- `Post(url, body, ...extra)` - POST 请求
- `Put(url, body, ...extra)` - PUT 请求
- `Delete(url, body, ...extra)` - DELETE 请求
- `SetHeader(key, value)` - 设置全局请求头

### Session 方法

- `NewSession()` - 创建新的 Session
- `SetHeader(key, value)` - 设置 Session 请求头
- `Get/Post/Put/Delete` - 与全局函数相同

### Response 方法

- `Ok()` - 检查状态码是否为 2xx
- `RaiseForStatus()` - 状态码非 2xx 时 panic
- `Text()` - 返回响应文本
- `Bytes()` - 返回响应字节
- `Json()` - 返回 gjson.Result 对象
- `Save(filepath)` - 保存响应到文件

## 依赖

- [github.com/tidwall/gjson](https://github.com/tidwall/gjson) - JSON 解析

## License

MIT
