/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-07-28 22:12:29
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-20 23:17:19
 * @FilePath: /apinto/src/app/constant/app.config.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */

import { InjectionToken } from '@angular/core'

export const APP_CONFIG = new InjectionToken('app.config')
export const MODAL_NORMAL_SIZE: number = 900
export const MODAL_SMALL_SIZE: number = 600
export const MODAL_LARGE_SIZE: number = 1200

// apinto项目的目录参数,其中view和edit字段需要与后端数据一致,以便匹配(权限用)
export const AppConfig: any = {
  menuList: [
    {
      title: '上游服务',
      icon: 'connection-box',
      menuIndex: 0,
      router: 'upstream',
      id: 2,
      menu: true,
      level: 0,
      children: [
        {
          title: '上游管理',
          routerLink: 'upstream/upstream',
          menuIndex: 0,
          id: 201,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'service_view',
          edit: 'service_edit'
        },
        {
          title: '服务发现',
          routerLink: 'upstream/serv-discovery',
          menuIndex: 0,
          id: 202,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'discovery_view',
          edit: 'discovery_edit'
        }
      ]
    },
    {
      title: 'API管理',
      icon: 'APIjiekou-7mme3dcg',
      routerLink: 'router',
      matchRouter: true,
      matchRouterExact: false,
      menuIndex: 0,
      id: 4,
      menu: true,
      level: 1,
      children: [
        {
          title: 'API列表',
          routerLink: 'router/api',
          menuIndex: 0,
          id: 401,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'api_view',
          edit: 'api_edit'
        },
        {
          title: '插件模板',
          routerLink: 'router/plugin',
          menuIndex: 0,
          id: 402,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'plugin_template_view',
          edit: 'plugin_template_edit'
        }
      ]
    },
    {
      title: '应用管理',
      routerLink: 'application',
      matchRouter: true,
      matchRouterExact: false,
      menuIndex: 0,
      id: 3,
      menu: true,
      level: 1,
      view: 'application_view',
      edit: 'application_edit',
      icon: 'yingyong-7mmhj11e'
    },
    {
      title: '基础设施',
      menuIndex: 0,
      icon: 'file-cabinet',
      id: 1,
      router: 'deploy',
      menu: true,
      level: 0,
      children: [
        {
          title: '网关集群',
          routerLink: 'deploy/cluster',
          menuIndex: 0,
          id: 101,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'cluster_view',
          edit: 'cluster_edit'
        },
        {
          title: '环境变量',
          routerLink: 'deploy/env',
          menuIndex: 0,
          id: 102,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'variable_view',
          edit: 'variable_edit'
        },
        {
          title: '插件管理',
          routerLink: 'deploy/plugin',
          menuIndex: 0,
          id: 103,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'plugin_view',
          edit: 'plugin_edit'
        }
      ]
    },
    {
      title: '服务治理',
      menuIndex: 0,
      id: 5,
      menu: true,
      router: 'serv-governance',
      icon: 'network-tree',
      level: 0,
      children: [
        {
          title: '流量策略',
          routerLink: 'serv-governance/traffic',
          menuIndex: 0,
          id: 501,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'strategy_traffic_view',
          edit: 'strategy_traffic_edit'
        },
        {
          title: '熔断策略',
          routerLink: 'serv-governance/fuse',
          menuIndex: 0,
          id: 502,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'strategy_fuse_view',
          edit: 'strategy_fuse_edit'
        },
        {
          title: '访问策略',
          routerLink: 'serv-governance/visit',
          menuIndex: 0,
          id: 503,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'strategy_visit_view',
          edit: 'strategy_visit_edit'
        },
        {
          title: '缓存策略',
          routerLink: 'serv-governance/cache',
          menuIndex: 0,
          id: 504,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'strategy_cache_view',
          edit: 'strategy_cache_edit'
        },
        {
          title: '灰度策略',
          routerLink: 'serv-governance/grey',
          menuIndex: 0,
          id: 505,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'strategy_grey_view',
          edit: 'strategy_grey_edit'
        }
      ]
    },
    {
      title: '系统管理',
      menuIndex: 0,
      menu: true,
      router: 'system',
      icon: 'system',
      id: 6,
      level: 0,
      children: [
        {
          title: '用户角色',
          routerLink: 'system/role',
          menuIndex: 0,
          id: 601,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'user_role_view',
          edit: 'user_role_edit'
        },
        {
          title: '外部应用',
          routerLink: 'system/ext-app',
          menuIndex: 0,
          id: 602,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'ext_app_view',
          edit: 'ext_app_edit'
        },
        {
          title: '邮箱设置',
          routerLink: 'system/email',
          menuIndex: 0,
          id: 603,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'notice_email_view',
          edit: 'notice_email_edit'
        },
        {
          title: 'webhook管理',
          routerLink: 'system/webhook',
          menuIndex: 0,
          id: 604,
          level: 1,
          matchRouter: true,
          matchRouterExact: false,
          view: 'notice_webhook_view',
          edit: 'notice_webhook_edit'
        }
      ]
    },
    {
      title: '审计日志',
      routerLink: 'audit-log',
      matchRouter: true,
      matchRouterExact: false,
      menuIndex: 0,
      menu: true,
      icon: 'form-one',
      id: 7,
      level: 0,
      view: 'audit_log_view'
    },
    {
      title: '监控告警',
      routerLink: 'monitor-alarm',
      menu: true,
      matchRouter: true,
      matchRouterExact: false,
      icon: 'jiankongshexiangtou',
      id: 901,
      level: 0,
      view: 'mon_partition_view',
      edit: 'mon_partition_edit'
    },
    {
      title: '查看授权',
      routerLink: 'auth-info',
      menu: false,
      id: 8,
      level: 0,
      view: 'authorization_view',
      edit: 'authorization_edit'
    },
    {
      title: '企业插件',
      routerLink: 'plugin',
      menu: true,
      id: 10,
      level: 0,
      view: 'enterprise_plugin_view',
      edit: 'enterprise_plugin_edit'
    },
    {
      title: '导航管理',
      routerLink: 'navigation',
      menu: true,
      id: 11,
      level: 0,
      view: 'navigation_view',
      edit: 'navigation_edit'
    },
    {
      title: '拦截器管理',
      routerLink: 'interceptor',
      menu: true,
      id: 12,
      level: 0,
      view: 'interceptor_view',
      edit: 'interceptor_edit'
    }
  ]
}
