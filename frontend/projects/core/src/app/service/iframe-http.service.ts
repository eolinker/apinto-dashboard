import { HttpClient } from '@angular/common/http'
import { Inject, Injectable } from '@angular/core'
import { Observable, Subject, Subscriber, Subscription, take } from 'rxjs'
import { API_URL, ApiService } from './api.service'
import { NavigationEnd, Router } from '@angular/router'
import { EoNgNavigationService } from './eo-ng-navigation.service'
import { BaseInfoService } from './base-info.service'
import { RouterService } from '../layout/api/router.service'

@Injectable({
  providedIn: 'root'
})
export class IframeHttpService {
  moduleName:string = ''
  subscription: Subscription = new Subscription()

  modalMaskEl: HTMLDivElement | undefined;
  private changeIframe: Subject<string> = new Subject<string>()

  constructor (private http:HttpClient,
    private api:ApiService,
    private router:Router,
    private navigation:EoNgNavigationService,
    private baseInfo:BaseInfoService,
    private routerService:RouterService,
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
        this.navigation.reqFlashBreadcrumb(breadcrumbOption, 'iframe')
        resolve(true)
      })
    },
    // user-center
    renewUserInfo: async () => {
      return new Promise((resolve) => {
        this.navigation.reqUpdateRightList()
        return this.navigation.repUpdateRightList().pipe(take(1)).subscribe(() => {
          resolve({
            userId: this.navigation.getUserId(),
            userRoleId: this.navigation.getUserRoleId(),
            userModuleAccess: this.navigation.accessMap.get(this.moduleName)
          })
        })
      })
    },
    getCurrentUserInfo: async () => {
      return new Promise((resolve) => {
        resolve({
          userId: this.navigation.getUserId(),
          userRoleId: this.navigation.getUserRoleId(),
          userModuleAccess: this.navigation.accessMap.get(this.moduleName)
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
    randomPsw: async () => {
      return new Promise((resolve) => {
        return this.api.get('random/password/id').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // monitor
    strategyRemote: (type:string, data?:any) => {
      return new Promise((resolve) => {
        return this.api.get(`strategy/filter-remote/${type}`, data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    clusterList: async () => {
      return new Promise((resolve) => {
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
        return this.api.get('common/enum/Service').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    appsSimpleList: () => {
      return new Promise((resolve) => {
        return this.api.get('application/enum').subscribe((resp:any) => {
          if (resp.code === 0) {
            resp.data.Application = resp.data.applications.map((app:{id:string, title:string}) => {
              return { name: app.id, title: app.title }
            })
          }
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
    },
    noticeChannels: () => {
      return new Promise((resolve) => {
        this.api.get('channels').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // apiSpace
    // 根据uuid判断api是否已存在
    checkApi: (uuid:string) => {
      return new Promise((resolve) => {
        this.api.get('router/check', { uuid: uuid }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 根据分组名称和uuid判断分组是否已存在
    checkGroupName: (type:string, data:any) => {
      return new Promise((resolve) => {
        this.api.put(`group/${type}/check`, data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 新建分组
    addGroup: (type:string, data:any) => {
      return new Promise((resolve) => {
        this.api.post(`group/${type}`, data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 新建实例（智能插件，比如上游是智能插件，新建上游=新建name为Service的示例
    addDynamic: (name:string, data:any) => {
      return new Promise((resolve) => {
        this.api.post(`dynamic/${name}`, data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 添加api
    addApi: (data:any) => {
      return new Promise((resolve) => {
        this.api.post('router', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 获取简易集群列表
    getSimpleClusters: () => {
      return new Promise((resolve) => {
        this.api.get('clusters/simple').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 将智能插件的实例发布到指定集群里
    dynamicOnlines: (name:string, uuid:string, cluster:Array<string>) => {
      return new Promise((resolve) => {
        this.api.put(`dynamic/${name}/online/${uuid}`, { cluster: cluster }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 发起发布管理弹窗，type支持api
    publishModal: (type:string, uuid:string) => {
      return this.publishModal(type, uuid)
    },
    // 发起批量发布管理弹窗，type支持api
    batchPublishResModal: (publishType:string, type:'online'|'offline', data:any, showLastStep?:boolean) => {
      return this.batchPublishResModal(publishType, type, data, showLastStep)
    },
    // remote插件存储数据
    storeKey: (key:string, data:any) => {
      return new Promise((resolve) => {
        this.api.put(`remote/${this.moduleName}/store/${key}`, data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // remote插件获得存储的数据
    getStore: (key:string) => {
      return new Promise((resolve) => {
        this.api.get(`remote/${this.moduleName}/store/${key}`).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    showModalMask: (style = {}) => {
      return new Promise((resolve) => {
        this.modalMaskEl?.remove()
        this.modalMaskEl = document.createElement('div')
        this.modalMaskEl.classList.add('ant-modal-mask', 'cdk-overlay-backdrop-showing')
        document.getElementById('iframePanel')?.parentElement!.appendChild(this.modalMaskEl)
        this.modalMaskEl.classList.add('cdk-overlay-backdrop')
        const iframeWrapper = document.querySelector<HTMLDivElement>('#iframePanel')
        if (iframeWrapper) {
          iframeWrapper.style.zIndex = '10000'
        }
        // this.modalMaskEl.onclick = () => {
        //   this.modalMaskEl!.style.display = 'none'
        //   this.modalMaskEl!.classList.remove('cdk-overlay-backdrop')
        //  this.modalMaskEl!.remove()
        // }
        Object.entries(style).forEach(([key, value]) => {
          // @ts-ignore
          this.modalMaskEl.style[key] = value
        })
        resolve(true)
      })
    },
    hideModalMask: () => {
      return new Promise((resolve) => {
        if (this.modalMaskEl) {
          this.modalMaskEl!.style.display = 'none'
          this.modalMaskEl!.classList.remove('cdk-overlay-backdrop')
         this.modalMaskEl!.remove()
        }
        const iframeWrapper = document.querySelector<HTMLDivElement>('#iframePanel')
        if (iframeWrapper) {
          iframeWrapper.style.zIndex = 'unset'
        }
        resolve(true)
      })
    }
  }

  // 该方法是为了打开需要传入header的iframe页面(remote插件），local插件则不需要
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

  publishModal (type:string, uuid:string) {
    switch (type) {
      case 'api':
        this.routerService.publishApiModal(uuid)
        break
      default:
        console.warn(`eo预警：无法调用发布管理方法，请检查传入的type=${type}是否正确`)
    }
  }

  batchPublishResModal (publishType:string, type:'online'|'offline', data:any, showLastStep?:boolean) {
    switch (publishType) {
      case 'api':
        this.routerService.batchPublishApiResModal(type, data, showLastStep
          ? () => {
              return new Promise((resolve) => {
                resolve({ data: { lastStep: true } })
                console.log({ data: { lastStep: true } })
              })
            }
          : undefined, () => {
          return new Promise((resolve) => {
            resolve({ data: { finishPublish: true } })
            console.log({ data: { finishPublish: true } })
          })
        })
        break
      default:
        console.warn(`eo预警：无法调用发布管理方法，请检查传入的type=${type}是否正确`)
    }
  }

  // 控制台手动触发iframe路由跳转，暂时用在面包屑中
  reqFlashIframe (url:string) {
    this.changeIframe.next(url)
  }

  repFlashIframe () {
    return this.changeIframe.asObservable()
  }
}
