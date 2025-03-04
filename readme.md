imyulizzz/chromegin-ubuntu
---
基于Go、Chromedp 的截图服务

[√]  支持表情包显示；
[TODO] 截图清晰度待提高

**构建命令**

```
docker build -t imyulizzz/chromegin-ubuntu:v0.0.6 . 
```

当前稳定的基础镜像为： **imyulizzz/chromegin-ubuntu:v0.0.6**

#### imyulizzz/chromegin-ubuntu:v0.0.6 Dockerfile

```
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

# 注: 使用 ubuntu:22.04 启动容器后， 手动执行了apt-get 相关命令，然后使用commit 构建的 imyulizzz/chromegin-ubuntu:v0.0.6 镜像
# 原因：安装 google-chrome-stable_current_amd64.deb 时，需要确认命令，这可能导致执行不成功。
```