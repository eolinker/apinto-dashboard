# 名称：apinto镜像，携带了部署k8s集群所需要的脚本
# 创建时间：2022-10-25
FROM centos:7.9.2009
MAINTAINER apinto-dashboard

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

#声明端口
EXPOSE 8080

ARG Version

#定义数据卷
VOLUME /var/lib/apinto-dashboard

RUN mkdir "apinto-dashboard"

COPY apserver_v${Version}_linux_amd64.tar.gz /apinto-dashboard

RUN cd /apinto-dashboard && tar -zxvf apserver_v${Version}_linux_amd64.tar.gz && rm -rf apserver_v${Version}_linux_amd64.tar.gz && ln -sfT apserver_v${Version} apserver

WORKDIR /apinto-dashboard/apserver

#容器启动命令
CMD sh docker_run.sh
