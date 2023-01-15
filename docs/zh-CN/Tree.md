# 项目目录结构

本文档随开发随时更新。

```text   
.
├── controller
├── docs
│   └── zh-CN
├── middleware
├── models
├── public
├── router
├── service
└── test
```

- ./controller:三层结构中靠近前端的一层，负责处理请求
- ./service:三层结构中的中间层，主要实现业务逻辑
- ./public:三层结构的最底层，包含所谓的实体类，与数据库交互
- ./router:包含路由映射
- ./middleware:中间件
- ./test:测试
- ./docs/*:文档