# Online-Exercise-System 在线联系系统
技术栈 Gin Gorm

## 安装mysql
```shell
docker pull mysql:5.7

#启动mysql
docker run --name mysql  -e MYSQL_ROOT_PASSWORD=123456  -p 3306:3306 -d mysql:5.7 
```


## 依赖
```shell
go get -u github.com/gin-gonic/gin
go get -u gorm.io/gorm
go install github.com/swaggo/swag/cmd/swag
go get github.com/swaggo/gin-swagger
go get github.com/swaggo/files
```
## swagger 
https://gitcode.net/mirrors/swaggo/gin-swagger 
```shell
swag init 
```
