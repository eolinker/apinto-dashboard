# 名称：apinto-dashboard镜像
# 创建时间：2022-7-1
FROM alpine:latest
MAINTAINER eolink

#设置环境变量
ENV DASHBOARD_CONFIG=/apinto-dashboard/config.yml

#声明端口
EXPOSE 8081

#数据库挂载目录
VOLUME /apinto-dashboard/data

#安装curl和bash
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && apk update  && apk add curl && apk add bash
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2


#Copy : Dockerfile  account.yml  apinto-dashboard  config.yml static  tpl  version
COPY ./apinto-dashboard-* /apinto-dashboard

WORKDIR /apinto-dashboard

#容器启动命令
CMD ./apinto-dashboard