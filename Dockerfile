#构建时
FROM golang:1.23 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o chromegin

# 运行时
FROM imyulizzz/chromegin-ubuntu:v0.0.6
WORKDIR /gobin
COPY --from=builder /app/chromegin .

# 设置环境变量
ENV PATH="/usr/bin/google-chrome-stable:${PATH}"

EXPOSE 6666

ENTRYPOINT ["/gobin/chromegin"]