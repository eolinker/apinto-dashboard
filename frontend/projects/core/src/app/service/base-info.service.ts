import { Injectable } from '@angular/core'
import {
  ActivationStart,
  GuardsCheckEnd,
  NavigationCancel,
  NavigationEnd,
  NavigationStart,
  ResolveStart,
  Router
} from '@angular/router'
import { assign } from 'lodash-es'
import { filter, of, switchMap, tap } from 'rxjs'

/**
 * 定义项目目前所有的路由参数，建议所有需要从该service取的参数值都在此定义
 * 另外参数名可以定的更具体一点，避免后续参数太多，出现重名问题
 */
export interface RouteParams {

  // API管理
  apiGroupId:string
  apiId:string
  pluginTemplateId:string

  // 应用管理
  appId:string

  // 基础设施
  clusterName:string
  pluginName:string // 插件管理中的插件名

  // 上游管理
  discoveryName:string
  serviceName:string

  // 服务治理
  strategyId:string

  // 监控告警
  // 分区id
  partitionId:string
  // 分区数据id
  monitorDataId:string
  // 策略id
  strategyUuid:string

  // 系统-外部应用id
  extAppId:string
  roleId:string

  // 插件管理
  pluginGroupId:string
  pluginId:string
  mdFileName:string
}

@Injectable({
  providedIn: 'root'
})
export class BaseInfoService {
  /**
   * 收集所有的路由params参数
   */
  private _allParams: RouteParams = {} as RouteParams

  // token = getCookieItem('Authorization') || window.localStorage.getItem('Authorization') || ''

  // updateToken (token: string) {
  //   this.token = token
  // }

  private routeMap = new Map<string, any[]>()

  tmpFn: any = null
  constructor (private router: Router) {
    this.init()
  }

  init () {
    this.onRouterChange()
  }

  onRouterChange () {
    let lastParamInfo!: RouteParams
    this.router.events
      .pipe(
        tap((e) => {
          if (e instanceof NavigationStart) {
            // 用 NavigationStart 会有问题，路由拦截之后，路由没有变化，但是 _allParams 被置空
            // 路由载入开始时重置参数
            lastParamInfo = Object.entries(this._allParams).length ? this._allParams : lastParamInfo || {}
            this._allParams = {} as RouteParams
            return
          }
          if (e instanceof GuardsCheckEnd && e.shouldActivate === false) {
            this._allParams = (lastParamInfo || {}) as RouteParams
          }
        }),
        filter((e: any) => e.snapshot),
        // filter((e: ActivationEnd) => e.snapshot.routeConfig?.path === 'inside/:hash'), // 匹配指定的 项目内页路由
        switchMap((e) => of(e.snapshot))
      )
      .subscribe(this.updateParamsInfo)
  }

  updateParamsInfo = (e: any) => {
    if (!e.routeConfig) {
      return
    }

    const routePath = e.routeConfig.path
    if (routePath && this.routeMap.has(routePath)) {
      // 清空原有的参数
      this.routeMap.get(routePath)?.forEach((key: keyof RouteParams) => delete this._allParams[key])
    }
    this.routeMap.set(routePath, Object.keys(assign({}, e.params || {}, e.queryParams || {})))

    this._allParams = { ...this._allParams, ...e.params, ...e.queryParams }
  }

  getSpaceKeyFromCache () {
    return window.localStorage.getItem('EO_SPACE_KEY') || ''
  }

  get spaceKey () {
    const hostName = window.location.hostname
    const oldSpaceKey = this.getSpaceKeyFromCache()
    const result = hostName.split('.')[0] || (this._allParams as any).spaceKey || oldSpaceKey
    if (oldSpaceKey !== result) {
      window.localStorage.setItem('EO_SPACE_KEY', result)
    }
    return result
  }

  // get projectHashKey () {
  //   return this._allParams.projectID
  // }

  get allParamsInfo () {
    return { ...this._allParams }
  }
}
