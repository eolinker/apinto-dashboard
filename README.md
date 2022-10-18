# Apinto Dashboard

[![Go Report Card](https://goreportcard.com/badge/github.com/eolinker/apinto-dashboard)](https://goreportcard.com/report/github.com/eolinker/apinto-dashboard) [![Releases](https://img.shields.io/github/release/eolinker/apinto-dashboard/all.svg?style=flat-square)](https://github.com/eolinker/apinto-dashboard/releases) [![LICENSE](https://img.shields.io/github/license/eolinker/Apinto-dashboard.svg?style=flat-square)](https://github.com/eolinker/apinto-dashboard/blob/main/LICENSE) ![](https://shields.io/github/downloads/eolinker/apinto-dashboard/total)

* **Apinto Dashboard**项目**main**分支与**Apinto**项目**main**分支同步更新

* 当前**Apinto Dashboard**最新版本为**v1.1.0-beta**，**Apinto**要求版本不低于**v0.6.4**

注意：main分支为开发主要分支，频繁更新可能导致使用不稳定，若需要使用稳定版本，请查看[release](https://github.com/eolinker/apinto-dashboard/releases)


### 什么是Apinto Dashboard

**Apinto Dashboard**是开源网关[**Apinto**](https://github.com/eolinker/apinto)的可视化UI项目。

此后，大家将告别繁琐复杂的命令行Curl命令，只需在**Dashboard**上轻轻一点，便可实现与开源网关**Apinto**的交互，极大地简化了Apinto的配置流程，降低了学习及使用成本。

**Apinto Dashboard**与**Apinto**交互流程如下图所示

![Apinto Dashboard与Apinto交互流程图](https://user-images.githubusercontent.com/14105999/175314303-4df9bfad-2abc-4e4a-9f24-30a8e3b64802.jpg)

### 编译

1. 进入**build/cmd**文件夹，执行编译脚本

```
cd builds/cmd && ./build.sh {版本号}
```

2. 编译后的文件存放在 **out/apinto-dashboard-{版本号}** 文件夹中

### 部署

1. 启动**Apinto**开源网关，Apinto启动教程请[点击](https://github.com/eolinker/apinto/#get-start)

2. 下载并解压安装包

```
wget https://github.com/eolinker/apinto-dashboard/releases/download/${version}/apinto-dashboard-${version}.linux.x64.tar.gz && tar -zxvf apinto-dashboard-${version}.linux.x64.tar.gz && cd apinto-dashboard
```

上述命令中的 **${version}** 为 **Apinto dashboard**的版本号，需要根据 **Apinto** 版本部署对应的 **Apinto Dashboard** 版本

下表为Apinto和Apinto Dashboard的版本联系

| Apinto版本   | Apinto Dashboard版本 |
| ------------ | -------------------- |
| 0.8.x        | v1.1.0-beta          |
| v0.6.x-0.7.x | v1.0.4-beta          |


下列示例命令以Apinto Dashboard v1.1.0-beta版本为例

```
wget https://github.com/eolinker/apinto-dashboard/releases/download/v1.1.0-beta/apinto-dashboard-v1.1.0-beta.linux.x64.tar.gz && tar -zxvf apinto-dashboard-v1.1.0-beta.linux.x64.tar.gz && cd apinto-dashboard
```


3. 编辑配置文件config.yml

```yaml
zone: zh_cn # 时区，根据时区获取当地语言的前端渲染页面，可选项：zh_cn｜ja_jp｜ en_us，当前版本仅支持zh_cn
default: monitor
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
      zh_cn: 上游服务   # 中文描述
      en_us: upstream services  # 英文描述
  - name: templates  # dashboard模块：插件模版
    profession: template # apinto模块：插件模版
    i18n_name:
      zh_cn: 模版
      en_us: template
  - name: discoveries    # dashboard模块：服务发现
    profession: discovery    # apinto模块：服务发现
    i18n_name:
      zh_cn: 服务发现
      en_us: discoveries
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
nohup ./apinto-dashboard > logs/stdout_apinto-dashboard_"$(date ‘+%Y%m%d-%H%M%S‘)".log 2>&1 &
```

5. 浏览器打开**Apinto Dashboard**地址，本示例在本地部署，因此ip为127.0.0.1，端口为8081

![image-20220616181447371](https://user-images.githubusercontent.com/14105999/174442723-1fe42ac5-012c-4f60-b1ec-e147d8d8ca9b.png)

6. 在浏览器中输入账号密码登录即可

至此，部署启用教程已结束，如需了解更多使用教程，请点击[更多](https://help.apinto.com/docs/apinto-dashboard)（教程文档正在赶工中）

### 联系我们

- **帮助文档**：[https://help.apinto.com](https://help.apinto.com/)

- **QQ群**: 725853895
- **Slack**：[加入我们](https://join.slack.com/t/slack-zer6755/shared_invite/zt-u7wzqp1u-aNA0XK9Bdb3kOpN03jRmYQ)
- **官网**：[https://www.apinto.com](https://www.apinto.com/)
- **论坛**：[https://community.apinto.com](https://community.apinto.com/)
- **微信群**：[![img](https://user-images.githubusercontent.com/25589530/149860447-5879437b-3cda-4833-aee3-69a2e538e85d.png)](https://user-images.githubusercontent.com/25589530/149860447-5879437b-3cda-4833-aee3-69a2e538e85d.png)

### 关于我们

EOLINK 是领先的 API 管理服务供应商，为全球超过3000家企业提供专业的 API 研发管理、API自动化测试、API监控、API网关等服务。是首家为ITSS（中国电子工业标准化技术协会）制定API研发管理行业规范的企业。

官方网站：[https://www.eolink.com](https://www.eolink.com/)

免费下载PC桌面端：https://www.eolink.com/pc/
