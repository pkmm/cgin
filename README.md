使用gin框架  
orm使用gorm  
三方依赖使用gvendor

安装govendor 
`go get -u github.com/kardianos/govendor`  
govendor sync 同步所有使用依赖包

开发过程中是 fresh 命令 可以热更新代码的改动
fresh 配置文件是runner.conf