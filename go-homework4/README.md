# 个人博客系统后端（Go + Gin + GORM + SQLite）

## 一、项目简介

基于 Go 语言实现的个人博客后端，使用 Gin 作为 HTTP 框架、GORM 作为 ORM、SQLite 作为数据库。
功能覆盖作业要求的全部内容：

- 用户注册 / 登录（bcrypt 密码加密 + JWT 认证）
- 文章 CRUD（创建、列表、详情、更新、删除），更新 / 删除仅作者本人可操作
- 文章评论（发表评论、查看某篇文章的评论列表）
- 统一响应格式、统一错误处理、请求日志中间件

## 二、运行环境

| 项目 | 要求 |
| --- | --- |
| 操作系统 | Windows / macOS / Linux 均可 |
| Go 版本 | Go 1.26 及以上（见 `go.mod` 中的 `go 1.26.7`） |
| 数据库 | SQLite，使用纯 Go 驱动（`github.com/libtnb/sqlite`），**无需安装 SQLite，也无需 CGO** |
| 测试工具 | curl / Postman（任选其一） |

## 三、依赖安装

依赖已在 `go.mod` / `go.sum` 中声明，进入项目目录执行：

```bash
go mod tidy
```

主要依赖：

| 依赖 | 作用 |
| --- | --- |
| `github.com/gin-gonic/gin` | HTTP 框架：路由、参数绑定、中间件 |
| `gorm.io/gorm` | ORM：模型映射、增删改查、软删除、关联预加载 |
| `github.com/libtnb/sqlite` | GORM 的 SQLite 驱动（纯 Go 实现） |
| `github.com/golang-jwt/jwt/v5` | 签发 / 校验 JWT |
| `golang.org/x/crypto/bcrypt` | 密码单向哈希 |

## 四、启动方式

```bash
# 1. 进入项目目录
cd go-homework4

# 2. 下载依赖
go mod tidy

# 3. 启动服务
go run .
```

启动成功后输出：

```
服务已启动: http://localhost:8080
[GIN-debug] POST   /api/register ...
```

首次启动会自动创建 `blog.db` 并建表（`users` / `posts` / `comments`）。

### 可选环境变量

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| `SERVER_PORT` | `8080` | 服务监听端口 |
| `JWT_SECRET` | `dev-secret-change-me` | JWT 签名密钥（生产环境必须修改） |
| `DB_PATH` | `blog.db` | SQLite 数据库文件路径 |

PowerShell 示例：

```powershell
$env:SERVER_PORT="8081"
$env:JWT_SECRET="my-secret"
go run .
```

## 五、项目结构

```
go-homework4/
├── main.go              # 程序入口：读配置 → 连数据库 → 注册路由 → 启动
├── go.mod / go.sum      # 依赖声明
├── blog.db              # SQLite 数据库文件（首次启动自动生成）
├── config/
│   └── config.go        # 读取配置（端口 / JWT 密钥 / 数据库路径）
├── database/
│   └── database.go      # 连接数据库 + AutoMigrate 建表
├── models/
│   └── models.go        # 数据库模型 + 请求/响应结构
├── handlers/
│   ├── auth.go          # 注册、登录
│   ├── post.go          # 文章 CRUD
│   └── comment.go       # 评论创建、评论列表
├── middleware/
│   ├── auth.go          # JWT 鉴权中间件，把用户 ID 写入上下文
│   └── logger.go        # 请求日志中间件
└── utils/
    ├── jwt.go           # JWT 签发与解析
    └── response.go      # 统一响应格式 + 错误分诊
```

## 六、接口说明

### 6.1 统一响应格式

```json
{
  "code": 0,
  "message": "成功",
  "data": { }
}
```

- `code = 0` 表示业务成功；失败时 `code` 与 HTTP 状态码一致（401 / 403 / 404 / 409 / 500）
- 成功时 HTTP 状态码：查询/更新 `200`，新建 `201`

### 6.2 接口列表（共 9 个）

基础路径：`http://localhost:8080/api`

| # | 方法 | 路径 | 是否需要登录 | 说明 |
| --- | --- | --- | --- | --- |
| 1 | POST | `/register` | 否 | 用户注册 |
| 2 | POST | `/login` | 否 | 用户登录，返回 JWT |
| 3 | GET | `/posts` | 否 | 文章列表（分页） |
| 4 | GET | `/posts/:id` | 否 | 文章详情（含作者、评论） |
| 5 | GET | `/posts/:id/comments` | 否 | 某篇文章的评论列表 |
| 6 | POST | `/posts` | 是 | 创建文章 |
| 7 | PUT | `/posts/:id` | 是 | 更新文章（仅作者） |
| 8 | DELETE | `/posts/:id` | 是 | 删除文章（仅作者，软删除） |
| 9 | POST | `/posts/:id/comments` | 是 | 发表评论 |

需要登录的接口必须在请求头带上：

```
Authorization: Bearer <登录接口返回的 token>
```

### 6.3 请求体示例

注册 `POST /api/register`：

```json
{ "username": "tester01", "password": "123456", "email": "tester01@test.com" }
```

登录 `POST /api/login`：

```json
{ "username": "tester01", "password": "123456" }
```

创建 / 更新文章 `POST|PUT /api/posts`：

```json
{ "title": "我的第一篇文章", "content": "这是测试文章的内容。" }
```

发表评论 `POST /api/posts/:id/comments`：

```json
{ "content": "沙发！这是第一条评论。" }
```

文章列表分页参数：`GET /api/posts?page=1&page_size=10`（默认 `page=1`、`page_size=10`），
返回结构为 `{"list": [...], "page": 1, "page_size": 10, "total": N}`。

## 七、测试用例与结果

以下结果均为在本机 Windows + PowerShell 5 + `go run .`（端口 8080）下，使用 `curl.exe` 实际执行得到。
测试用户：`tester01`（作者）、`tester02`（他人）。

| # | 接口 | 前置条件 | 预期 | 实际状态码 | 结论 |
| --- | --- | --- | --- | --- | --- |
| 1 | POST `/api/register` | 用户名/邮箱未被占用 | 201 创建成功 | 201 | 通过 |
| 2 | POST `/api/register` | 注册第二个用户 tester02 | 201 创建成功 | 201 | 通过 |
| 3 | POST `/api/login` | 用户名密码正确 | 200 返回 token | 200 | 通过 |
| 4 | POST `/api/login` | 密码错误 | 401 | 401 | 通过 |
| 5 | GET `/api/posts?page=1&page_size=10` | 未登录 | 200 返回分页列表 | 200 | 通过 |
| 6 | POST `/api/posts` | 未带 token | 401 | 401 | 通过 |
| 7 | POST `/api/posts` | 带 tester01 的 token | 201 创建文章 | 201 | 通过 |
| 8 | GET `/api/posts/1` | 文章存在 | 200 返回详情（含作者） | 200 | 通过 |
| 9 | PUT `/api/posts/1` | 作者本人 | 200 返回更新后内容 | 200 | 通过 |
| 10 | PUT `/api/posts/1` | 非作者（tester02） | 403 无权限 | 403 | 通过 |
| 11 | POST `/api/posts/1/comments` | 带 tester02 的 token | 201 创建评论 | 201 | 通过 |
| 12 | GET `/api/posts/1/comments` | 文章存在 | 200 返回评论列表（含评论人） | 200 | 通过 |
| 13 | DELETE `/api/posts/1` | 作者本人 | 200 删除成功 | 200 | 通过 |
| 14 | GET `/api/posts/1` | 文章已被删除 | 404 | 404 | 通过 |
| 15 | GET `/api/posts/999999` | 文章不存在 | 404 | 404 | 通过 |

### 7.1 关键响应实录

**1. 注册成功（201）**

```json
{"code":0,"message":"创建成功","data":{"id":1,"username":"tester01","email":"tester01@test.com","created_at":"2026-10-30T12:12:33+08:00"}}
```

**3. 登录成功（200）**

```json
{"code":0,"message":"成功","data":{"token":"eyJhbGciOiJIUzI1NiIs...","user":{"id":1,"username":"tester01","email":"tester01@test.com"}}}
```

**4. 密码错误（401）**

```json
{"code":401,"message":"用户名或密码错误"}
```

**6. 未登录创建文章（401）**

```json
{"code":401,"message":"请先登录"}
```

**10. 非作者更新文章（403）**

```json
{"code":403,"message":"没有权限执行此操作"}
```

**13. 删除成功（200）**

```json
{"code":0,"message":"成功","data":{"deleted_id":1}}
```

**14. 删除后查询（404）**

```json
{"code":404,"message":"资源不存在"}
```

### 7.2 复现测试的命令（可选）

PowerShell 中需要注意：直接往 `curl.exe` 传 JSON 字符串会被吞掉双引号，因此把请求体先写入文件再 `--data-binary "@文件"`：

```powershell
# 1. 注册
'{"username":"tester01","password":"123456","email":"tester01@test.com"}' | Set-Content -Encoding UTF8 reg.json
curl.exe -s -i -X POST http://127.0.0.1:8080/api/register -H "Content-Type: application/json" --data-binary "@reg.json"

# 2. 登录并取出 token
'{"username":"tester01","password":"123456"}' | Set-Content -Encoding UTF8 login.json
$resp = curl.exe -s -X POST http://127.0.0.1:8080/api/login -H "Content-Type: application/json" --data-binary "@login.json"
$token = ($resp | ConvertFrom-Json).data.token

# 3. 带 token 创建文章
'{"title":"我的第一篇文章","content":"测试内容"}' | Set-Content -Encoding UTF8 post.json
curl.exe -s -X POST http://127.0.0.1:8080/api/posts -H "Content-Type: application/json" -H "Authorization: Bearer $token" --data-binary "@post.json"

# 4. 公开查询
curl.exe -s "http://127.0.0.1:8080/api/posts?page=1&page_size=10"
curl.exe -s http://127.0.0.1:8080/api/posts/1
curl.exe -s http://127.0.0.1:8080/api/posts/1/comments
```

> 注意端口与路径：若修改了 `SERVER_PORT`，请同步替换 URL 中的 `8080`；接口前缀统一为 `/api`。

## 八、实现要点

1. **密码安全**：注册时用 bcrypt 哈希存储，数据库中不保存明文；`User.Password` 标记 `json:"-"`，任何响应都不会带出密码。
2. **JWT 认证**：登录签发 HS256 token（含 `user_id`、`username`、`exp`）；鉴权中间件解析 `Authorization: Bearer <token>`，失败返回 401 并中断；解析成功把 `user_id` 写入上下文供 handler 使用。
3. **权限控制**：更新 / 删除前先按 ID 查文章，比对 `post.UserID` 与当前登录用户 ID，不一致返回 403。
4. **统一错误处理**：`utils.HandleError` 把内部错误映射为 404 / 403 / 409 / 500，避免把数据库细节暴露给调用方。
5. **分页**：`page`、`page_size` 做边界校验后使用 `Count` + `Offset/Limit`，并返回总条数。
6. **关联数据**：使用 `Preload` 预加载文章作者、评论人，避免 N+1 查询。
7. **软删除**：模型内嵌 `gorm.Model`，删除只写 `deleted_at`，之后查询自动过滤，因此删除后再查返回 404。
8. **日志**：自定义 `middleware.Logger()` 记录「状态码 + 方法 + 路径 + 耗时」，配合 `gin.Recovery()` 防止单个请求 panic 导致服务退出。
