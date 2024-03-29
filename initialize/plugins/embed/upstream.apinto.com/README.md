## 基本介绍
在系统架构中，我们可以通过河流理论来判断上下游，最好方法是想象一条河。
1. 下游的水肯定是来自上游的。
2. 如果有人破坏了河的下游部分，那将对上游没有影响。
3. 如果有人破坏了河的上游部分，这将影响下游，即不会得到任何水。

以订单服务调用为例

![](http://data.eolinker.com/course/aDDlSnc299d3dd7d3131812811d0b922d63f7ea5a7f60ab.png)
订单服务提供服务数据，是河流中的上游，软件A接受订单服务返回响应数据，作为河流中的下游。

为了保证上游服务的高可用及稳定性，一般来说，相同的上游服务会部署在多台机器/虚拟机中，通过增加机器/虚拟机的方式横向拓展，从而应对下游请求的流量高峰。

该插件将给Apinto网关配置上游服务信息，通过API绑定，实现完整转发流程。
## 功能特性
- 支持配置负载均衡算法：IP Hash、带权轮询，合理分配上游服务压力，保证上游服务稳定性。
- 支持引用环境变量，既统一了各环境各集群中上游服务的配置，又可以通过环境变量实现上游服务地址的差异化。
- 支持会话保持，开启会话保持后，相同Session的请求将会转发到相同的上游服务地址，确保流程的连贯性。

## 更新日志
### V1.0（2023-4-30）
- 插件上线