# gin项目练习
## 项目结构
```
├── /cmd
│   └── main.go
├── /config
│   └── config.yml
├── /docs
├── /internal
│   ├── /api
│   ├── /app
│   ├── /common
│   ├── /controller
│   ├── /middleware
│   ├── /models
│   ├── /service
│   └── /utils
├── /pkg
├── /tests
├── go.mod
├── go.sum
└── README.md
```
### 职责 
- /cmd：项目的入口文件，包含 main.go。
- /config：项目的配置文件，包含 config.yml。
- /docs：项目的 Swagger 文档目录，包含生成的 Swagger 文档。
- /internal：项目的内部代码目录，包含项目的主要逻辑。
    - /api：项目的 API 目录，包含路由定义。
    - /app：项目的应用目录，包含连接第三方服务的代码。
    - /common：项目的公共目录，包含通用的常量、错误定义等。
    - /controller：项目的控制器目录，包含处理 HTTP 请求的函数。
    - /middleware：项目的中间件目录，包含请求处理的中间件函数。
    - /models：项目的模型目录，包含数据库模型的定义。
    - /service：项目的服务目录，包含业务逻辑的函数。
    - /utils：项目的工具目录，包含通用的工具函数。
- /pkg：存放第三方库，如第三方中间件、工具库等。
- /tests：存放项目的测试文件。
- go.mod：Go 模块文件，记录项目的依赖信息。
- go.sum：Go 模块文件，记录项目的依赖版本信息。
- README.md：项目的说明文档，包含项目介绍、初始化项目、数据库配置、启动程序等。


## 初始化项目
```shell
cd homework4
go mod init homework4  #已经有了go.mod文件不需要执行这个
#根据go文件 自动导入依赖
go mod tidy
```
## 数据库配置
> 修改config/config.yml文件 配置数据库连接信息

## 生成swagger文档
```shell
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd\main.go
```

## 启动程序
```shell
go run cmd\main.go 
```

## 访问Swagger文档

```
http://localhost:9527/api/v1/swagger/index.html
```
## 使用Swagger文档

### 不需要认证的接口
- 用户注册: `POST /api/v1/user/register`
- 用户登录: `POST /api/v1/user/login`
- 获取文章列表: `GET /api/v1/post/list`
- 获取评论列表: `GET /api/v1/comment/list`
- 健康检查: `GET /health`

### 需要认证的接口
这些接口需要在请求头中添加JWT Token:
- 创建文章: `POST /api/v1/post/create`
- 更新文章: `PUT /api/v1/post/update`
- 删除文章: `DELETE /api/v1/post/delete`
- 创建评论: `POST /api/v1/comment/create`

#### 添加Token的方式
1. 先调用登录接口获取Token
2. 在Swagger页面右上角点击 "Authorize" 按钮
3. 输入Token值(格式: `Bearer {token}`)
4. 点击 "Authorize" 确认

## 注意事项

1. 每次修改API接口后,需要重新生成Swagger文档
2. Swagger文档会自动生成在 `docs` 目录下
3. 如果修改了 `main.go` 中的总体注释,需要重新生成文档




