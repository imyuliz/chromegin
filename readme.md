

如使用以上服务没有完成截图的功能实现，则启动容器后，手动分布执行dockerfile 的内容，commit 成为新镜像

现在：稳定的镜像为： **imyulizzz/chromegin-ubuntu:v0.0.6**

二进制命令构建为： ** GOOS=linux GOARCH=amd64 go build . **
构建命令为： ** docker build -t imyulizzz/chromegin-ubuntu:v0.0.6 . **

```
# 使用 Ubuntu 22.04 作为基础镜像
FROM ubuntu:22.04

# 更新软件包列表并安装中文字体、Google Chrome 和支持表情包的字体
RUN apt-get update && \
    apt-get install -y \
    wget \
    git \
    fonts-noto-cjk \
    fonts-wqy-zenhei \
    fonts-noto-color-emoji && \
    wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb && \
    apt-get install -y ./google-chrome-stable_current_amd64.deb && \
    google-chrome --version && \
    rm -rf *.deb

# 设置环境
ENV PATH="/usr/bin/google-chrome-stable:${PATH}"
# 设置 Go 代理和环境变量
ENV GOPROXY=https://goproxy.io \
    PATH="${PATH}:/usr/local/go/bin"

# 定义使用的 Golang 版本
ARG GO_VERSION=1.16

# 安装 Golang
RUN wget "https://dl.google.com/go/go$GO_VERSION.linux-amd64.tar.gz" && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf "go$GO_VERSION.linux-amd64.tar.gz" && \
    rm -rf *.tar.gz && \
    go version && go env

# 设置工作目录
ARG GO_DIR=.
ARG BUILD_DIR=/gobin
WORKDIR $BUILD_DIR

# 复制 Go 模块文件并下载依赖
COPY $GO_DIR/go.mod .
COPY $GO_DIR/go.sum .
RUN go mod download

# 复制源代码并编译
COPY $GO_DIR .
RUN export GITHASH=$(git rev-parse --short HEAD) && \
    export BUILDAT=$(date) && \
    go build -ldflags "-w -s -X 'main.BuildAt=$BUILDAT' -X 'main.GitHash=$GITHASH'"

# 设置挂载点和暴露端口
VOLUME /pic
EXPOSE 6666

# 删除 Golang 文件
RUN rm -rf /usr/local/go /root/go /root/.cache /root/.config

# 设置容器启动命令
CMD ["/gobin/chromegin"]

```