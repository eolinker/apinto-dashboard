# 名称：apinto通用镜像
# 创建时间：2022-10-25
FROM centos:7.9.2009
MAINTAINER liujian

ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ARG APP
ARG VERSION
ARG DIR

RUN mkdir -p /${APP}

COPY ${DIR}/${APP}_${VERSION}_linux_amd64.tar.gz /${APP}

RUN cd /${APP} && tar -zxvf ${APP}_${VERSION}_linux_amd64.tar.gz && rm -rf ${APP}_${VERSION}_linux_amd64.tar.gz && ln -snF ${APP}_${VERSION} ${APP}

WORKDIR /$APP/$APP

CMD sh docker_run.sh