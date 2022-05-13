FROM golang:alpine

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
	GOPROXY="https://mirrors.aliyun.com/goproxy/"
	
WORKDIR /root/golang/gopath/src/github.com/daysleep666/short_chain

# 将代码复制到容器中
COPY . .

# 将我们的代码编译成二进制可执行文件  可执行文件名为 app
RUN go build -o app cmd/app/main.go

# 移动到用于存放生成的二进制文件的 /dist 目录
WORKDIR /dist

# 将二进制文件从 目录复制到这里
RUN cp /root/golang/gopath/src/github.com/daysleep666/short_chain/app .
# 配置
RUN mkdir ./config

RUN cp /root/golang/gopath/src/github.com/daysleep666/short_chain/config/app.toml ./config/
# 静态文件
RUN cp -r /root/golang/gopath/src/github.com/daysleep666/short_chain/static ./
# 声明服务端口
EXPOSE 1234
EXPOSE 8080

# 启动容器时运行的命令
CMD ["/dist/app"]

