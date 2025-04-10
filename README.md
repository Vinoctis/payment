### 目录结构     
```
/payment
├── cmd
│   └── payment          # 服务入口（main.go）
├── internal
│   ├── handler          # HTTP/RPC请求处理层（路由、参数校验）
│   ├── service          # 业务逻辑层（支付核心逻辑）
│   ├── repository       # 数据访问层（数据库/缓存操作）
│   └── model            # 数据结构定义（DTO/DO）
├── pkg
│   ├── config           # 配置加载工具（支持yaml/json）
│   ├── middleware       # 通用中间件（鉴权、日志）
│   └── utils            # 通用工具函数
├── api
│   └── proto            # Protobuf协议定义（用于gRPC）
├── configs              # 配置文件（dev/prod环境）
├── scripts              # 部署/迁移脚本
├── test                 # 集成测试用例
├── go.mod               # Go模块依赖
└── README.md
```


