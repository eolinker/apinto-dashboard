# Apinto Dashboard


![](http://data.eolinker.com/course/eaC48Js3400ffd03c21e36b3eea434dce22d7877a3194f6.png)

Apinto Dashboad主版本应与 Apinto主版本一起使用。
最新发布的Apinto Dashboad版本是2.0.0，可与Apinto0.12.4以上版本兼容。

**什么是Apinto Dashboad V2.0.0**

Apinto Dashboard 是基于开源网关 Apinto并符合企业级API网关需求场景的可视化控制台项目。
通过Dashboard 可以管理集群、上游、应用以及API等模块，并以集群维度管理各个模块的生命周期。
具有优秀的用户操作体验，配置流程简短，上手难度低。内置了丰富的插件，用户可根据业务需求动态灵活地配置策略。

Apinto Dashboard与Apinto交互流程如下图所示：
![](http://data.eolinker.com/course/QSgDqEKb311ec59fb0052436e0d3bdbacdc2b984b71cc25.png)

集群管理：管理各个环境的集群，给集群配置证书、配置并发布该集群下的环境变量、监控并管理集群下各个网关节点、配置管理等。

上游服务：上游管理和服务发现。服务发现支持consul、eureka、nacos注册中心；上游管理是管理所有提供API调用的后端系统，都需要上线到指定的集群才生效；

API管理：支持业务域分组，管理所有后端系统提供的API及其生命周期，根据业务上下线到相应的集群。

应用管理：管理所有调用方，配置请求网关的鉴权，以及支持转发后端的额外参数鉴权，上下线到指定集群生效。

服务治理：针对不同集群配置并上线限流、访问、熔断、灰度、缓存等策略，保障网关集群以及后端系统稳定工作。

网关插件：即将开放，管理Apinto插件，Apinto内置几十个插件，同时支持自定义添加插件。

企业插件：即将提供并支持自定义业务型企业插件，供用户安装使用，业务型企业插件如：用户角色、监控告警、日志、API文档、开放平台、安全防护、数据分析、调用链、mock、在线调测、安全测试、国密、多协议……

系统管理：配置邮箱，配置告警模板等。

## 迭代计划
![](http://data.eolinker.com/course/gydll750fcfc7874b12137c49566f71a586dc093887aa93.png)
如果您是个人开发者，可基于API网关相关的业务场景开发有价值的网关插件或企业级插件，并且愿意分享给Apinto，您将会成为Apinto的杰出贡献者或得到一定的收益。
如果您是企业，可基于Apinto网关开发企业级插件，成为Apinto的合作伙伴。

### 部署

* 直接部署：[部署教程](https://help.apinto.com/docs/dashboard/quick/arrange.html)
* [快速入门教程](https://help.apinto.com/docs/dashboard-v2/quick/quick_start.html)

- ### **联系我们**


* **帮助文档**：[https://help.apinto.com](https://help.apinto.com/docs)

- **QQ群**: 725853895

- **Slack**：[加入我们](https://join.slack.com/t/slack-zer6755/shared_invite/zt-u7wzqp1u-aNA0XK9Bdb3kOpN03jRmYQ)

- **官网**：[https://www.apinto.com](https://www.apinto.com/)
- **论坛**：[https://community.apinto.com](https://community.apinto.com/)
- **微信群**：<img src="https://user-images.githubusercontent.com/25589530/149860447-5879437b-3cda-4833-aee3-69a2e538e85d.png" style="width:150px" />

### 关于我们

EOLINK 是领先的 API 管理服务供应商，为全球超过3000家企业提供专业的 API 研发管理、API自动化测试、API监控、API网关等服务。是首家为ITSS（中国电子工业标准化技术协会）制定API研发管理行业规范的企业。

官方网站：[https://www.eolink.com](https://www.eolink.com "EOLINK官方网站")

免费下载PC桌面端：[https://www.eolink.com/pc/](https://www.eolink.com/pc/ "免费下载PC客户端")
