/*
 * @Date: 2023-11-09 15:27:45
 * @LastEditors: maggieyyy
 * @LastEditTime: 2023-11-17 10:26:15
 * @FilePath: \apinto\projects\apinto-auth\src\app\service\base-info.service.ts
 */
import { Injectable } from '@angular/core'
import {
  GuardsCheckEnd,
  NavigationStart,
  Router
} from '@angular/router'
import { assign } from 'lodash-es'
import { filter, of, switchMap, tap } from 'rxjs'
import { ApiService } from './api.service'

/**
 * 定义项目目前所有的路由参数，建议所有需要从该service取的参数值都在此定义
 * 另外参数名可以定的更具体一点，避免后续参数太多，出现重名问题
 */
export interface RouteParams {
  roleId:string
  operateDisable?:string
}

@Injectable({
  providedIn: 'root'
})
export class BaseInfoService {
  /**
   * 收集所有的路由params参数
   */
  private _allParams: RouteParams = {} as RouteParams
  private routeMap = new Map<string, any[]>()
  userInfoUpdated:boolean = false
  userModuleAccess:string = ''
  userId:string = ''
  userRoleId:string = ''
  userProfile:{id:number, roleIds:string[], nickName:string, userName:string, [k:string]:any} | undefined
  constructor (private router: Router, private apiService: ApiService) {
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

  get allParamsInfo () {
    return { ...this._allParams }
  }
}
