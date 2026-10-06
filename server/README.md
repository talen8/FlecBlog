# Flec-Server

> 基于 Go + Gin + GORM 的现代化博客后端 API 服务

## 技术栈

- **语言**: [Go 1.27.1](https://golang.org)
- **框架**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io)
- **数据库**: PostgreSQL
- **认证**: JWT (JSON Web Tokens), OAuth2, Goth
- **API 文档**: Swagger
- **定时任务**: [Cron](https://github.com/robfig/cron)
- **其他**: User-Agent 解析, 飞书 SDK, 微信公众号

## 文件结构

```
server/
├── api/              # API 定义和 Swagger 注释
│   ├── middleware/   # 中间件 (认证、CORS、日志、限流、RBAC等)
│   ├── router/       # 路由配置
│   └── v1/           # API v1 版本接口
├── cmd/              # 应用入口
│   └── main.go
├── config/           # 配置管理
├── docs/             # Swagger 生成的文档
├── internal/         # 内部业务逻辑
│   ├── dto/          # 数据传输对象
│   ├── model/        # 数据模型
│   ├── repository/   # 数据访问层
│   └── service/      # 业务逻辑层
├── pkg/              # 可复用的包
├── templates/        # 模板文件
├── Dockerfile
└── go.mod
```

详细文档请查看 [项目主 README](../README.md)。
