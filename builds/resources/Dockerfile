# 名称：apinto-dashboard镜像
# 创建时间：2022-4-27
FROM alpine:latest
MAINTAINER eolink

#声明端口
EXPOSE 8080

#定义数据卷
VOLUME /apinto-dashboard/data

#设置环境变量
ENV APINTO_CONFIG=/apinto-dashboard/config.yml APINTO_DATA_DIR=/var/lib/apinto/data

#安装curl和bash
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && apk update  && apk add curl && apk add bash
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
#解压dashboard程序压缩包
COPY ./apinto-dashboard.linux.x64.tar.gz /
RUN tar -zxvf apinto-dashboard.linux.x64.tar.gz && rm -rf apinto-dashboard.linux.x64.tar.gz

#复制程序默认配置文件
RUN cp /apinto-dashboard/config.yml.tmp ${APINTO_CONFIG}

WORKDIR /apinto-dashboard

#容器启动命令
CMD ./apinto-dashboard