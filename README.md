# Alert WebHook Service
[TOC]

![image-20200521173449327](Alert%20WebHook%20Service-readme.assets/image-20200521173449327.png)

## Build

```bash
# 查看支持的环境
go tool dist list 
# 在 macos 下编译 linux 的包
GOOS=darwin GOARCH=linux/amd64 go build 
```
## Docker 
### Dockerfile1-internet
```dockerfile
FROM golang
MAINTAINER dwj_wz@163.com
ADD sms.go /app/
#RUN go build /app/sms.go
WORKDIR /app

RUN go get -d -v ./...
RUN go install -v ./...
EXPOSE 8080

ENTRYPOINT /app/sms
```
### Dockerfile
```dockerfile
FROM      alpine:latest
MAINTAINER dwj_wz@163.com

ENV GIN_MODE=release

WORKDIR /app
COPY xsgaSms /app
COPY config.yml /app
COPY sms.tmpl /app

EXPOSE 8080
ENTRYPOINT ["/app/xsgaSms"]
```
### Docker build
如果有外网环境`google`,执行 Dockerfile1-internet

或者执行 Dockerfile 文件
```bash
docker build . -t sms:latest
```
### Docker run 
```bash
docker run -d -p 8080:8080 -v /tmp/config.yml:/app/config.yml -v /tmp/sms.tmpl:/app/sms.tmpl sms:latest
```
## config.yml

```yml
# 短信接口配置
user: daiwj
pwd: xxx
url: http://127.0.0.1:8080/sms2
ext: xxx
priority: xxx
tempPath: sms.tmpl
userid: 10
# 模板文件路径
TempPath: xxx
# 告警默认手机号码,可在 alert 文件中指定
phoneNumberFromAlert: xxxx
```

## 消息模板

```template
{{range .}}
STATUS: {{.Status }}
Labels:
{{ range $key, $value := .Labels }}
{{ $key }}: {{ $value }}
{{end}}
Annotations:
{{ range $key, $value := .Annotations }}
{{ $key }}: {{ $value }}
{{end}}
=========
{{end}}

```

```bash
测试的消息样式:
STATUS: ok
Labels:
test1: test1
warnPhone: 123
Annotations:
=========
STATUS: not ok
Labels:
test2: test2
warnPhone: ssss
Annotations:
=========      
```



## Test

`https://juejin.im/post/5d6d462ef265da03e5234f57#heading-7`
```bash
# 执行全部测试
go test 
# 执行其中某个函数的测试
go test -run "TestGenApiData" -v 

# -v 参数可以看到日志输出,否则仅看到测试失败和成功
```


