import { HttpClient } from '@angular/common/http'
import { Inject, Injectable } from '@angular/core'
import { Observable, Subscriber, Subscription, take } from 'rxjs'
import { API_URL, ApiService } from './api.service'
import { NavigationEnd, Router } from '@angular/router'
import { EoNgNavigationService } from './eo-ng-navigation.service'
import { BaseInfoService } from './base-info.service'

@Injectable({
  providedIn: 'root'
})
export class IframeHttpService {
  moduleName:string = ''
  subscription: Subscription = new Subscription()

  constructor (private http:HttpClient,
    private api:ApiService,
    private router:Router,
    private navigation:EoNgNavigationService,
    private baseInfo:BaseInfoService,
    @Inject(API_URL) public urlPrefix:string) {
    this.moduleName = this.baseInfo.allParamsInfo.moduleName
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.moduleName = this.baseInfo.allParamsInfo.moduleName
      }
    })
  }

  // 所有对外提供的接口都放在这里
  apinto2PluginApi = {
    get: async (url:string, params?:{[key:string]:any}) => {
      return new Promise((resolve) => {
        return this.api.get(`module/${this.moduleName}/${url}`, params).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    post: async (url:string, body?: any, params?:{[key:string]:any}) => {
      return new Promise((resolve) => {
        return this.api.post(`module/${this.moduleName}/${url}`, body, params).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    put: async (url:string, body?: any, params?:{[key:string]:any}) => {
      return new Promise((resolve) => {
        return this.api.put(`module/${this.moduleName}/${url}`, body, params).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    delete: async (url:string, params?:{[key:string]:any}) => {
      return new Promise((resolve) => {
        return this.api.delete(`module/${this.moduleName}/${url}`, params).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    patch: async (url:string, body?:any, params?:{[key:string]:any}) => {
      return new Promise((resolve) => {
        return this.api.patch(`module/${this.moduleName}/${url}`, body, params).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    changeRouter: async (url:string) => {
      return new Promise((resolve) => {
        let newRouterArr:Array<string> = this.router.url.split('#')
        if (this.router.url.includes('#')) {
          newRouterArr.pop()
        }
        newRouterArr = newRouterArr.join('').split('/')
        newRouterArr[newRouterArr.length - 1] = `${newRouterArr[newRouterArr.length - 1]}#${url}`
        window.location.href = newRouterArr.join('/')
        resolve(true)
      })
    },
    freshMenu: async () => {
      return new Promise((resolve) => {
        this.navigation.reqFlashMenu()
        resolve(true)
      })
    },
    changeBreadcrumb: async (breadcrumbOption:Array<any>) => {
      return new Promise((resolve) => {
        for (const breadcrumb of breadcrumbOption) {
          if (breadcrumb.routerLink) {
            breadcrumb.routerLink = `${this.navigation.iframePrefix}/${this.moduleName}/${breadcrumb.routerLink}`
          }
        }
        this.navigation.reqFlashBreadcrumb(breadcrumbOption)
        resolve(true)
      })
    },
    renewUserInfo: async () => {
      return new Promise((resolve) => {
        this.navigation.reqUpdateRightList()
        return this.navigation.repUpdateRightList().pipe(take(1)).subscribe(() => {
          resolve({
            userId: this.navigation.getUserId(),
            userRoleId: this.navigation.getUserRoleId(),
            userModuleAccess: this.navigation.originAccessData[this.moduleName]
          })
        })
      })
    },
    getCurrentUserInfo: async () => {
      return new Promise((resolve) => {
        resolve({
          userId: this.navigation.getUserId(),
          userRoleId: this.navigation.getUserRoleId(),
          userModuleAccess: this.navigation.originAccessData[this.moduleName]
        })
      })
    },
    accessList: async () => {
      return new Promise((resolve) => {
        return this.api.get('access').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    roleProfile: async (roleId:string) => {
      return new Promise((resolve) => {
        return this.api.get('role', { id: roleId }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    addRole: async (data:any) => {
      return new Promise((resolve) => {
        return this.api.post('role', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    editRole: async (roleId:string, data:any) => {
      return new Promise((resolve) => {
        return this.api.put('role', data, { id: roleId }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    deleteRole: async (roleId:string) => {
      return new Promise((resolve) => {
        return this.api.delete('role', { id: roleId }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    getMyProfile: async () => {
      return new Promise((resolve) => {
        return this.api.get('my/profile').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    editMyProfile: async (data:any) => {
      return new Promise((resolve) => {
        return this.api.put('my/profile', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    otherUserProfile: async (userId:string) => {
      return new Promise((resolve) => {
        return this.api.get('user/profile', { id: userId }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    addUser: async (data:any) => {
      return new Promise((resolve) => {
        return this.api.post('user/profile', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    editUser: async (userId:string, data:any) => {
      return new Promise((resolve) => {
        return this.api.put('user/profile', data, { id: userId }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    changeUserPartInfo: async (userId:string, data:any) => {
      return new Promise((resolve) => {
        return this.api.patch('user/profile', data, { id: userId }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    rolesList: async () => {
      return new Promise((resolve) => {
        return this.api.get('role/options').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    usersList: async (data:any) => {
      return new Promise((resolve) => {
        return this.api.get('user/list', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    updateRole: async (data:any) => {
      return new Promise((resolve) => {
        return this.api.post('role/batch-update', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    rolesGroupList: async () => {
      return new Promise((resolve) => {
        return this.api.get('roles').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    deleteUsers: async (idsList:string[]) => {
      return new Promise((resolve) => {
        return this.api.post('user/delete', { ids: idsList }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    batchDeleteRoles: async (data:any) => {
      return new Promise((resolve) => {
        return this.api.post('role/batch-delete', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    randomPsw: async () => {
      return new Promise((resolve) => {
        return this.api.get('random/password/id').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    resetPsw: async (data:any) => {
      return new Promise((resolve) => {
        return this.api.post('user/password-reset', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    monitorPartition: async (uuid:string) => {
      return new Promise((resolve) => {
        return this.api.get('monitor/partition', { uuid: uuid }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    addMonitorPartition: (data:any) => {
      return new Promise((resolve) => {
        return this.api.post('monitor/partition', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    editMonitorPartition: (data:any) => {
      return new Promise((resolve) => {
        return this.api.put('monitor/partition', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    deleteMonitorPartition: async (uuid:string) => {
      return new Promise((resolve) => {
        this.api.delete('monitor/partition', { uuid: uuid }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    monitorPartitionList: async () => {
      return new Promise((resolve) => {
        this.api.get('monitor/partitions').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    monitorOverviewSummary: async (data:any) => {
      return new Promise((resolve) => {
        this.api.post('monitor/overview/summary', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    monitorOverviewInvoke: async (data:any) => {
      return new Promise((resolve) => {
        this.api.post('monitor/overview/invoke', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    monitorOverviewMessage: async (data:any) => {
      return new Promise((resolve) => {
        this.api.post('monitor/overview/message', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    monitorOverviewTop: async (data:any) => {
      return new Promise((resolve) => {
        this.api.post('monitor/overview/top', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    monitorApiTableList: async (data:any) => {
      return new Promise((resolve) => {
        this.api.post('monitor/api', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    monitorAppTableList: async (data:any) => {
      return new Promise((resolve) => {
        this.api.post('monitor/app', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    monitorInvokeDetail: async (type:string, data:any) => {
      return new Promise((resolve) => {
        this.api.post(`monitor/${type}/details`, data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    monitorDetail: async (currentType:string, detailType:string, data:any) => {
      return new Promise((resolve) => {
        this.api.post(`monitor/${currentType}/details/${detailType}`, data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    monitorDetailTrend: async (currentType:string, detailType:string, data:any) => {
      return new Promise((resolve) => {
        this.api.post(`monitor/${currentType}/details/${detailType}/trend`, data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    warnHistoryList: async (data:any) => {
      return new Promise((resolve) => {
        this.api.get('warn/history', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    warnStrategyData: (uuid:string) => {
      return new Promise((resolve) => {
        return this.api.get('warn/strategy', { uuid: uuid }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    addWarnStrategy: (data:any) => {
      return new Promise((resolve) => {
        return this.api.post('warn/strategy', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    editWarnStrategy: (data:any) => {
      return new Promise((resolve) => {
        return this.api.put('warn/strategy', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    deleteWarnStrategy: (uuid:string) => {
      return new Promise((resolve) => {
        return this.api.delete('warn/strategy', { uuid: uuid }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    disabledStrategy: (data:any) => {
      return new Promise((resolve) => {
        this.api.patch('warn/strategy', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    warnStrategyList: (data:any) => {
      return new Promise((resolve) => {
        return this.api.get('warn/strategys', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 策略remote选项
    strategyRemote: (type:string, data?:any) => {
      return new Promise((resolve) => {
        return this.api.get(`strategy/filter-remote/${type}`, data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    clusterList: async () => {
      return new Promise((resolve) => {
        console.log('clusterList')
        return this.api.get('cluster/enum').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    apiGroupList: () => {
      return new Promise((resolve) => {
        return this.api.get('router/groups').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    apiSimpleList: (data:any) => {
      return new Promise((resolve) => {
        return this.api.get('router/enum', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    serviceSimpleList: () => {
      return new Promise((resolve) => {
        return this.api.get('service/enum').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    appsSimpleList: () => {
      return new Promise((resolve) => {
        return this.api.get('application/enum').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    userList: () => {
      return new Promise((resolve) => {
        this.api.get('user/enum').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    }
  }

  openIframe (url:string, option?:{headers?:Array<{name:string, value:string}>}) {
    return new Observable((observer: Subscriber<any>) => {
      let objectUrl: string|null
      const header:{[k:string]:any} = {}
      if (option?.headers?.length) {
        for (const item of option.headers) {
          header[item.name] = item.value
        }
      }

      this.http
        .get(`${this.urlPrefix}${url}`, { ...header })
        .subscribe((m:any) => {
          objectUrl = URL.createObjectURL(new Blob([m.blob()], { type: 'application/json' }))
          observer.next(objectUrl)
        })
      return () => {
        if (objectUrl) {
          URL.revokeObjectURL(objectUrl)
          objectUrl = null
        }
      }
    })
  }
}
