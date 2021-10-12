gin + gorm

golang 版本是 1.17

`go mod tidy `安装依赖
下载遇到网络问题，可以查看main.go首部的注释，在goland中可以直接执行

开发过程中使用`fresh`命令，热更新  
`fresh`命令的配置文件是`runner.conf`

`fresh` 使用 `go install github.com/pilu/fresh` 安装完成

别的版本golang可能需要使用 `go get github.com/pilu/fresh`

##### API 文档的生成  
    `swagger 命令 swag init`

### 部署：
    `build_linux.bat` 生成可执行文件，上传到服务器，同时上传`config.yaml`文件  
    使用supervisor启动程序，使用nginx反向代理程序。

#### development:  
    `git update-index --skip-worktree "config.yaml"` 暂时忽略配置文件  
    `git update-index --no-skip-worktree "config.yaml"`
