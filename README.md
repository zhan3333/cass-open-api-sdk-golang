# cass openapi Golang 使用示例

## 使用

1. 编辑 .env 文件

```shell script
cp .env.example .env
```

```.env
API_URL=
APPID=
# 用户公钥 (用于测试)
PUBLIC_KEY_STR=
# 用户私钥 (用于请求签名)
PRIVATE_KEY_STR=
# vzhuo 公钥 (用于验签)
VZHUO_PUBLIC_KEY_STR=
```

2. 安装 vendor

```shell script
go mod download
go mod tidy
```

2. 运行测试

`Cass/client_test.go -> TestOneBankPay`