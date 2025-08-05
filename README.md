### fun-service
#### 依赖环境
```
go v1.24.5
```
#### 项目目录
```
|   .air.toml
|   .gitignore
|   go.mod
|   go.sum
|   main.go
|   README.md
|   tree.txt
|
+---config
|       app.yaml
|       config.go
|
+---docs
|       docs.go
|       swagger.json
|       swagger.yaml
|
+---internal
|   +---handler
|   |       user.go
|   |
|   +---middleware
|   |       logger.go
|   |
|   +---model
|   |       user.go
|   |
|   \---router
|           router.go
|
+---pkg
|   +---database
|   |       mysql.go
|   |
|   +---logger
|   |       logger.go
|   |
|   +---redis
|   |       redis.go
|   |
|   \---utils
|           utils.go
|
+---
```
#### 指令
安装依赖
```
go tidy
```
热启动项目
```
air 
```
生成swagger文档
```
swag init
```
格式化文档
```
go fmt ./...
```
