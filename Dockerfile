#构建时
FROM golang:1.23 as builder
WORKDIR /usr/src/app
ENV GOPROXY=https://goproxy.cn
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o chromegin

# 运行时(使用 Ubuntu 22.04 作为基础镜像)
FROM imyulizzz/chromegin-ubuntu:v0.0.6 as server
ARG GO_DIR=.
ARG BUILD_DIR=/gobin
WORKDIR $BUILD_DIR
COPY --from=builder /usr/src/app/chromegin $GO_DIR/chromegin
COPY --from=builder /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

# 设置环境变量
ENV PATH="/usr/bin/google-chrome-stable:${PATH}"

EXPOSE 6666


ENTRYPOINT  ["/gobin/chromegin"]