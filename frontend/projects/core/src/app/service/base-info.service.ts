import { Injectable } from '@angular/core'
import {
  GuardsCheckEnd,
  NavigationStart,
  Router
} from '@angular/router'
import { assign } from 'lodash-es'
import { filter, of, switchMap, tap } from 'rxjs'
import { environment } from '../../environments/environment'

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

  // 策略id
  strategyUuid:string

  // 系统-外部应用id
  extAppId:string
  roleId:string

  // 插件管理
  pluginGroupId:string
  pluginId:string
  mdFileName:string

  moduleName:string
}

@Injectable({
  providedIn: 'root'
})
export class BaseInfoService {
  private _version:string = '3.3.5'
  private _updateDate:string = '2023-12-15'
  private _powered:string = 'Powered by https://eolink.com'
  private _product:string = 'Apinto'
  private _showGuide:boolean = !environment.isBusiness
  /**
   * 收集所有的路由params参数
   */
  private _allParams: RouteParams = {} as RouteParams

  private routeMap = new Map<string, any[]>()

  tmpFn: any = null
  constructor (private router: Router) {
    this.init()
  }

  init () {
    this.onRouterChange()
  }

  onRouterChange () {
    let lastParamInfo!: RouteParams | {}
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

  get allParamsInfo () {
    return { ...this._allParams }
  }

  set version (newVersion:string) {
    this._version = newVersion || this._version
  }

  get version () {
    return this._version
  }

  set product (newProduct:string) {
    this._product = newProduct || this._product
  }

  get product () {
    return this._product
  }

  set powered (poweredStr:string) {
    this._powered = poweredStr ?? this._powered
  }

  get powered () {
    return this._powered
  }

  set showGuide (show:boolean) {
    this._showGuide = show ?? this._showGuide
  }

  get showGuide () {
    return this._showGuide
  }

  get updateDate () {
    return this._updateDate
  }

  set updateDate (time:string) {
    const nowDate = new Date(time)
    const year = nowDate.getFullYear()
    let month:string | number = nowDate.getMonth() + 1
    let day :string | number = nowDate.getDate()
    if (month < 10) month = '0' + month
    if (day < 10) day = '0' + day
    this._updateDate = time ? year + '-' + month + '-' + day : this._updateDate
  }
}
