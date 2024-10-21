
##  项目描述

该项目为账本APP后台同步服务，采用golang编写

## 编译

```shell
# 交叉编译windows
CGO_ENABLED=0  GOOS=windows  GOARCH=amd64  go build cmd/account.go

# 交叉编译Linux
CGO_ENABLED=0  GOOS=linux  GOARCH=amd64  go build -o account cmd/account.go
```