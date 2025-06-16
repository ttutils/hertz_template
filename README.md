<h1 align="center">hertz_template</h1>

此项目是根据[@buyfakett](https://github.com/buyfakett)的hertz使用习惯的模板仓库

本仓库支持`sqlite3`、`mysql`、`postgres`数据库

默认不配置的情况下位`sqlite3`数据库

需要`postgres`或者`mysql`的话修改以下配置

```yaml
db:
  type: postgres           # mysql
  host: xxx
  port: 5432
  user: postgres
  password: xxx
  database: test
```

在启动的时候会去判断是否有用户，如果没有就会创建一个id为1的用户，这个用户为超级管理员

## 技术栈

- [Go](https://golang.org/)

- [hertz](https://github.com/cloudwego/hertz)

- [gorm](https://github.com/go-gorm/gorm)

## 项目目录

```tree
.
├── Dockerfile                  # Dockerfile
├── Dockerfile.template         # Dockerfile模板
├── Makefile                    # Makefile
├── biz                         # 业务代码
│     ├── dal                   # 数据库连接和操作
│     ├── dbmodel               # 数据库模型
│     ├── handler               # 服务逻辑
│     ├── model                 # 模型(proto生成的)
│     ├── mw                    # 中间件
│     └── router                # 路由
├── bootstrao                   # 启动代码
├── build.sh                    # 编译脚本
├── config                      # 配置文件
│     ├── config.yaml           # 配置文件(可以覆盖默认配置)
│     └── default.yaml          # 默认配置文件(服务端这里定义的默认配置)
├── docs                        # swagger文档
├── idl                         # idl文件(proto)
├── main.go                     # 启动文件
├── router.go                   # 自定义路由文件
├── static                      # 静态文件(前端编译结果，必须要index.html)
└── utils                       # 工具包
```

## 开发

### 启动

如果需要指定配置文件，可以使用以下命令

```bash
go run . --config=config/config.yaml
```

### 自动化

目前使用`github actions`自动化,开发环境每个`commit`会自动编译docker镜像,打v1.0.0的标签的时候会编译docker镜像和二进制文件到`release`下
