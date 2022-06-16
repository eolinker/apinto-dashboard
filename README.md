# Apinto Dashboard

[![Go Report Card](https://goreportcard.com/badge/github.com/eolinker/apinto-dashboard)](https://goreportcard.com/report/github.com/eolinker/apinto) [![Releases](https://img.shields.io/github/release/eolinker/apinto-dashboard/all.svg?style=flat-square)](https://github.com/eolinker/apinto/releases) [![LICENSE](https://img.shields.io/github/license/eolinker/Apinto-dashboard.svg?style=flat-square)](https://github.com/eolinker/apinto-dashboard/blob/main/LICENSE)![](https://shields.io/github/downloads/eolinker/apinto-dashboard/total)

### 简介

**Apinto Dashboard**是开源网关**Apinto**的可视化UI项目，极大简化了配置**Apinto**网关的流程操作，降低学习和使用成本。



### 部署

1. 启动**Apinto**开源网关，Apinto启动教程请[点击](https://github.com/eolinker/apinto/#get-start)
2. 下载并解压安装包

```
wget https://github.com/eolinker/apinto-dashboard/releases/download/v0.5.1/apinto-dashboard-v1.0.0.linux.x64.tar.gz && tar -zxvf apinto-dashboard-v1.0.0.linux.x64.tar.gz && cd apinto-dashboard
```

3. 编辑配置文件config.yml

```yaml
zone: zh_cn # 时区，根据时区获取当地语言的前端渲染页面，可选项：zh_cn｜ja_jp｜ en_us，当前版本仅支持zh_cn
default: upstream
apinto:		# Apinto openAPI地址列表，若有多个节点，可填写多个节点的openAPI地址
  - "http://127.0.0.1:9400"   
port: 8081    # dashboard监听端口
user_details:	# 用户账号获取渠道
  type: file	# 文件，当前版本只支持读取文件
  file: ./account.yml	# 文件名称
professions:    # 流程阶段，下面配置中的name和profession为dashboard在apinto的映射名称，下述配置内容将会在dashboard导航栏中展现
  - name: services    # dashboard模块：服务
    profession: service # apinto模块：服务
    i18n_name:    # 国际化语言名称
      zh_cn: 服务    # 中文描述
      en_us: services    # 英文描述
  - name: upstreams    # dashboard模块：上游
    profession: upstream    # apinto模块：上游
    i18n_name:
      zh_cn: 上游
      en_us: upstreams
  - name: discoveries    # dashboard模块：服务发现
    profession: discovery    # apinto模块：服务发现
    i18n_name:
      zh_cn: 服务发现
      en_us: discoveries
  - name: auths        # dashboard模块：鉴权
    profession: auth    # apinto模块：鉴权
    i18n_name:
      zh_cn: 鉴权
      en_us: auths
  - name: outputs        # dashboard模块：输出器
    profession: output    # apinto模块：输出器
    i18n_name:
      zh_cn: 输出
      en_us: outputs
```

用户账号、密码默认均为**admin**。如若需要修改账号密码信息，可编辑**account.yml**文件，语法遵从**yaml**语法，配置详细说明如下：

```yaml
account_list: # 账号列表
- user_name: admin	# 账号
  password: admin		# 密码
  info:							# 基本信息
    desc: admin用户		# 描述
```

4. 启动程序

（1） 在当前窗口运行，该方式启动的程序，当窗口关闭，进程也会关闭

```
./apinto-dashboard
```

（2）在后台运行

``` 
nohup ./apinto-dashboard > logs/stdout_apinto-dashboard_"$(date \"+%Y%m%d-%H%M%S\")".log 2>&1 &
```

5. 浏览器打开**Apinto Dashboard**地址，本示例在本地部署，因此ip为127.0.0.1，端口为8081

![image-20220616181447371](/Users/liujian/Library/Application Support/typora-user-images/image-20220616181447371.png)

6. 在浏览器中输入账号密码登录即可

至此，部署启用教程已结束，如需了解更多使用教程，请点击[更多](https://help.apinto.com/docs/apinto-dashboard)
