## 基本介绍
Webhook（又名钩子） 是一个 API 概念，是微服务 API 的使用范式之一，也被称为反向 API，即前端不主动发送请求，完全由后端推送。

简单来说，Webhook 就是一个接收 HTTP POST（或GET，PUT，DELETE）的URL，一个实现了 Webhook 的 API 提供商就是在当事件发生的时候会向这个配置好的 URL 发送一条信息，与请求-响应式不同，使用 Webhook 你可以实时接受到变化。

在配置前，使用者需要提供一个Webhook API，支持对接企业微信、飞书、钉钉等常用办公软件等消息通知API。
当Apinto触发特定事件时，如请求失败次数告警，监控程序将会按照配置的Webhook规则，推送事件信息到目标地址，实时发送告警信息。
## 功能特性
- 配置Webhook规则，包括请求地址、请求方式、请求体模版格式
- 当监控触发告警事件时，监控程序按照Webhhok规则发送请求到目标地址，实时发送告警信息。
- 通过Webhook的自定义请求体模版特性，对接微信、飞书、钉钉等常用办公软件的消息通知API，
## 更新日志
### V1.0（2023-4-30）
- 插件上线