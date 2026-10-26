# Bluebell

基于 Go 语言开发的高性能社区论坛 Web 后端服务，采用经典分层架构（Controller - Logic - DAO）解耦设计。

## 核心特性

- **分层架构**：清晰的路由控制、业务处理与数据访问层划分
- **高并发支持**：集成主流 Go Web 生态组件，支持高性能请求处理与日志归档
- **热重载开发**：内置 Air 工具配置，提升本地开发与调试效率

## 技术栈

- **Web 框架**：Gin
- **配置与日志**：Viper / Uber Zap + Lumberjack
- **认证与算法**：JWT / Snowflake（分布式 ID）
- **数据存储**：MySQL + Redis
- **热重载**：Air

## 目录结构

```text
├── controllers/    # API 控制器与参数校验
├── dao/            # 数据访问层（MySQL、Redis）
├── logic/          # 核心业务逻辑
├── middlewares/    # 中间件（认证、限流、日志等）
├── models/         # 数据模型与结构体定义
├── pkg/            # 公共功能组件（雪花算法、Token 工具等）
├── routes/         # 路由注册与分发
├── settings/       # 配置管理
├── logger/         # 高性能日志模块
├── templates/      # 模板资源
├── config.yaml     # 核心配置文件
├── .air.toml       # Air 热重载配置
└── main.go         # 服务入口
