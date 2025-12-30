# 乡村振兴后端应用 (zxbe_demo)

基于Go语言和Gin框架开发的乡村振兴微信小程序后端API服务。

## 功能模块

- 📰 **资讯管理** - 乡村资讯的发布、查看、分类筛选
- 🏠 **农家乐管理** - 农家乐信息的展示、搜索、发布
- 📋 **政策公告** - 政策文件的发布、查看、分类管理
- 🏞️ **旅游景区** - 景区信息的管理、展示、评价
- 💼 **招聘信息** - 招聘岗位的发布、搜索、管理
- 🆘 **求助信息** - 求助信息的发布、响应、状态管理
- 👤 **用户管理** - 用户信息的维护和管理

## 技术栈

- **语言**: Go 1.17+
- **框架**: Gin
- **数据库**: MySQL
- **ORM**: GORM
- **配置管理**: godotenv
- **跨域处理**: gin-contrib/cors

## 项目结构

```
zxbe_demo/
├── cmd/                    # 命令行工具
│   └── seed.go            # 数据初始化脚本
├── config/                # 配置文件
│   └── database.go        # 数据库配置
├── controllers/           # 控制器
│   ├── news.go           # 资讯控制器
│   ├── farmhouse.go      # 农家乐控制器
│   ├── policy.go         # 政策控制器
│   ├── tourism.go        # 旅游控制器
│   ├── jobs.go           # 招聘控制器
│   ├── help.go           # 求助控制器
│   └── user.go           # 用户控制器
├── middleware/            # 中间件
│   └── cors.go           # 跨域中间件
├── models/               # 数据模型
│   ├── user.go          # 用户模型
│   ├── news.go          # 资讯模型
│   ├── farmhouse.go     # 农家乐模型
│   ├── policy.go        # 政策模型
│   ├── tourism.go       # 旅游模型
│   ├── job.go           # 招聘模型
│   └── help.go          # 求助模型
├── routes/               # 路由
│   └── routes.go        # 路由配置
├── utils/                # 工具函数
│   └── response.go      # 响应工具
├── .env                  # 环境变量配置
├── go.mod               # Go模块文件
├── main.go              # 主程序入口
└── README.md            # 项目说明
```

## 安装和运行

### 1. 环境要求

- Go 1.17+
- MySQL 5.7+

### 2. 克隆项目

```bash
cd zxbe_demo
```

### 3. 安装依赖

```bash
go mod tidy
```

### 4. 配置数据库

创建MySQL数据库：
```sql
CREATE DATABASE zx_demo CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

修改 `.env` 文件中的数据库配置：
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=zx_demo
```

### 5. 运行应用

```bash
# 启动服务器
go run main.go

# 或者先编译再运行
go build -o zxbe_demo main.go
./zxbe_demo
```

### 6. 初始化测试数据（可选）

```bash
go run cmd/seed.go
```

服务器将在 `http://localhost:8080` 启动。

## API接口文档

### 资讯相关接口

- `GET /api/news` - 获取资讯列表
- `GET /api/news/:id` - 获取资讯详情
- `GET /api/news/latest` - 获取最新资讯
- `GET /api/news/category/:category` - 根据分类获取资讯

### 农家乐相关接口

- `GET /api/farmhouse` - 获取农家乐列表
- `GET /api/farmhouse/:id` - 获取农家乐详情
- `POST /api/farmhouse` - 创建农家乐
- `PUT /api/farmhouse/:id` - 更新农家乐
- `DELETE /api/farmhouse/:id` - 删除农家乐

### 政策公告相关接口

- `GET /api/policy` - 获取政策列表
- `GET /api/policy/:id` - 获取政策详情
- `POST /api/policy` - 创建政策
- `PUT /api/policy/:id` - 更新政策
- `DELETE /api/policy/:id` - 删除政策

### 旅游景区相关接口

- `GET /api/tourism` - 获取景区列表
- `GET /api/tourism/:id` - 获取景区详情
- `POST /api/tourism` - 创建景区
- `PUT /api/tourism/:id` - 更新景区
- `DELETE /api/tourism/:id` - 删除景区

### 招聘信息相关接口

- `GET /api/jobs` - 获取招聘列表
- `GET /api/jobs/:id` - 获取招聘详情
- `POST /api/jobs` - 创建招聘
- `PUT /api/jobs/:id` - 更新招聘
- `DELETE /api/jobs/:id` - 删除招聘

### 求助信息相关接口

- `GET /api/help` - 获取求助列表
- `GET /api/help/:id` - 获取求助详情
- `POST /api/help` - 创建求助
- `PUT /api/help/:id` - 更新求助
- `DELETE /api/help/:id` - 删除求助

### 用户相关接口

- `GET /api/user/profile` - 获取用户信息
- `PUT /api/user/profile` - 更新用户信息

## 响应格式

所有API接口都返回统一的JSON格式：

```json
{
  "code": 200,
  "message": "success",
  "data": {}
}
```

- `code`: 状态码 (200: 成功, 400: 参数错误, 404: 未找到, 500: 服务器错误)
- `message`: 响应消息
- `data`: 响应数据

## 开发说明

1. 所有模型都包含软删除功能
2. 支持分页查询和关键词搜索
3. 自动处理CORS跨域问题
4. 统一的错误处理和响应格式
5. 数据库自动迁移

## 许可证

MIT License
