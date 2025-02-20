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

# 设置环境变量
ENV PATH="/usr/bin/google-chrome-stable:${PATH}"

ARG GO_DIR=.

ARG BUILD_DIR=/gobin

WORKDIR $BUILD_DIR

COPY $GO_DIR/chromegin .

COPY $GO_DIR .

EXPOSE 6666


ENTRYPOINT  ["/gobin/chromegin"]