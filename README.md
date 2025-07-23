#### fun-service
```
go v1.24.5
```
#### 项目目录
```
fun-service/              # 项目根目录
├── cmd/                   # 程序入口（可包含多个子命令）
│   └── api/               # API 服务入口
│       └── main.go        # 主程序入口文件（初始化 Gin 引擎、路由等）
├── config/                # 配置文件相关
│   ├── config.go          # 配置结构体定义、解析逻辑
│   ├── app.yaml           # 应用配置（如端口、环境）
│   └── database.yaml      # 数据库配置（如 MySQL、Redis）
├── internal/              # 项目内部代码（不对外暴露）
│   ├── api/               # API 层（路由、控制器）
│   │   ├── middleware/    # 中间件（如认证、日志、跨域）
│   │   ├── handler/       # 处理器（处理 HTTP 请求，调用服务层）
│   │   └── router/        # 路由定义（注册路由与处理器的映射）
│   ├── service/           # 服务层（业务逻辑处理）
│   ├── repository/        # 数据访问层（与数据库交互）
│   └── model/             # 数据模型（结构体定义，对应数据库表）
├── pkg/                   # 公共库（可复用、对外暴露的工具）
│   ├── logger/            # 日志工具
│   ├── validator/         # 参数校验工具
│   ├── database/          # 数据库连接工具（如 MySQL、Redis 客户端）
│   └── utils/             # 通用工具函数（如加密、时间处理）
├── api/                   # 接口定义（如 OpenAPI/Swagger 文档）
│   └── swagger/           # Swagger 文档生成文件
├── web/                   # 静态资源（如前端页面、CSS、JS）
│   ├── static/            # 静态文件（图片、样式表等）
│   └── template/          # 模板文件（若使用 Gin 模板渲染）
├── storage/               # 存储目录（日志、上传文件等）
│   ├── logs/              # 日志文件
│   └── uploads/           # 上传的文件
├── test/                  # 测试代码
│   ├── api/               # API 层测试
│   └── service/           # 服务层测试
├── go.mod                 # Go 模块依赖文件
├── go.sum                 # 依赖校验文件
├── .env                   # 环境变量（本地开发用，不提交到代码库）
├── .gitignore             # Git 忽略文件
└── README.md              # 项目说明文档
```