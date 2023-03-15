# 名称：apinto镜像，携带了部署k8s集群所需要的脚本
# 创建时间：2022-10-25
FROM centos:7.9.2009
MAINTAINER apinto-dashboard

#声明端口
EXPOSE 8080


#定义数据卷
VOLUME /var/lib/apinto-dashboard

#设置环境变量
ENV APP=apinto-dashboard


RUN mkdir "apinto-dashboard"

RUN ls


COPY apinto-dashboard.tar.gz /apinto-dashboard

WORKDIR /apinto-dashboard

RUN pwd

RUN tar -zxvf apinto-dashboard.tar.gz && rm -rf apinto-dashboard.tar.gz

WORKDIR /apinto-dashboard/apinto-dashboard

RUN pwd

RUN ls

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

#修改脚本权限以及复制程序默认配置文件
RUN chmod 777 docker_run.sh


#容器启动命令
CMD sh docker_run.sh
