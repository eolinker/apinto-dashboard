import { HttpClient } from '@angular/common/http'
import { Inject, Injectable } from '@angular/core'
import { Observable, Subject, Subscriber, Subscription, take } from 'rxjs'
import { API_URL, ApiService } from './api.service'
import { NavigationEnd, Router } from '@angular/router'
import { EoNgNavigationService } from './eo-ng-navigation.service'
import { BaseInfoService } from './base-info.service'
import { RouterService } from '../layout/api/router.service'
import { EditableEnvTableService } from '../component/editable-env-table/editable-env-table.service'
import { EoNgApplicationService } from '../layout/application/application.service'
import { EoIntelligentPluginService } from '../component/intelligent-plugin/intelligent-plugin.service'
import { ServiceGovernanceService } from '../layout/serv-governance/service-governance.service'

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
    private applicationService:EoNgApplicationService,
    private envTableService:EditableEnvTableService,
    private intelPluginService:EoIntelligentPluginService,
    private servGovernanceService:ServiceGovernanceService,
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
        window.history.replaceState(null, '', newRouterArr.join('/'))
        resolve(true)
      })
    },
    freshMenu: async () => {
      return new Promise((resolve) => {
        this.navigation.getMenuList().subscribe((resp:any) => {
          resolve(resp)
        })
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
        this.navigation.getMenuList().subscribe(() => {
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
    apiSimpleList: async (data:any) => {
      return new Promise((resolve) => {
        return this.api.get('router/enum', data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    serviceSimpleList: async () => {
      return new Promise((resolve) => {
        return this.api.get('common/enum/Service').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    appsSimpleList: async () => {
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
    userList: async () => {
      return new Promise((resolve) => {
        this.api.get('user/enum').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    noticeChannels: async () => {
      return new Promise((resolve) => {
        this.api.get('channels').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // apiSpace
    // 根据uuid判断api是否已存在
    checkApi: async (uuid:string) => {
      return new Promise((resolve) => {
        this.api.get('router/check', { uuid: uuid }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 根据分组名称和uuid判断分组是否已存在
    checkGroupName: async (type:string, data:any) => {
      return new Promise((resolve) => {
        this.api.put(`group/${type}/check`, data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 新建分组
    addGroup: async (type:string, body:any, option?:{query?:any, showModal?:boolean}) => {
      if (option?.showModal && type === 'api') {
        // 打开新建弹窗，仅支持api分组
        return this.routerService.addOrEditGroupModal('add', body.uuid)
      } else {
        return new Promise((resolve) => {
          this.api.post(`group/${type}`, body).subscribe((resp:any) => {
            resolve(resp)
          })
        })
      }
    },
    getGroup: async (type:string, query?:any) => {
      return new Promise((resolve) => {
        switch (type) {
          case 'api':
            // 有相同作用已对外使用的方法apiGroupList，如有改动需要同步
            this.api.get('router/groups', query).subscribe((resp:any) => {
              resolve(resp)
            })
            break
          default:
            resolve({ code: -1, msg: '暂不支持此类型，请检查输入' })
        }
      })
    },
    // 新建实例（智能插件，比如上游是智能插件，新建上游=新建name为Service的示例
    addDynamic: async (name:string, data:any) => {
      return this.addDynamic(name, data)
    },
    // 添加api
    addApi: async (data:any) => {
      return this.addApi(data)
    },
    // 获取简易集群列表
    getSimpleClusters: async () => {
      return new Promise((resolve) => {
        this.api.get('clusters/simple').subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 将智能插件的实例发布到指定集群里
    dynamicOnlines: async (name:string, uuid:string, cluster:Array<string>) => {
      return new Promise((resolve) => {
        this.api.put(`dynamic/${name}/online/${uuid}`, { cluster: cluster }).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // 发起批量发布管理弹窗，type支持api
    batchPublishModal: async (type:string, publishType:'online'|'offline', data:any) => {
      return this.batchPublishModal(type, publishType, data)
    },
    // 发起批量发布管理弹窗，type支持api
    batchPublishResModal: async (type:string, publishType:'online'|'offline', data:any, showLastStep?:boolean) => {
      return this.batchPublishResModal(type, publishType, data, showLastStep)
    },
    // remote插件存储数据
    storeKey: async (key:string, data:any) => {
      return new Promise((resolve) => {
        this.api.put(`remote/${this.moduleName}/store/${key}`, data).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    // remote插件获得存储的数据
    getStore: async (key:string) => {
      return new Promise((resolve) => {
        this.api.get(`remote/${this.moduleName}/store/${key}`).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    showModalMask: async (style?:{[k:string]:any}, closable?:boolean) => {
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
        if (style && Object.entries(style).length > 0) {
          Object.entries(style).forEach(([key, value]) => {
            // @ts-ignore
            this.modalMaskEl.style[key] = value
          })
        }
        if (closable) {
          this.modalMaskEl.onclick = () => {
            this.modalMaskEl!.style.display = 'none'
            this.modalMaskEl!.classList.remove('cdk-overlay-backdrop')
           this.modalMaskEl!.remove()
           resolve({ data: { closeMask: true } })
          }
        } else {
          this.modalMaskEl.onclick = null
          resolve(true)
        }
      })
    },
    hideModalMask: async () => {
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
    },
    addEntity: async (type:string, body:any, option?:{query?:any, showModal?:boolean}) => {
      switch (type) {
        case 'api':
          return this.addApi(body)
        case 'app-auth':
          return new Promise((resolve) => {
            this.api.post('application/auth', body, option?.query || {}).subscribe((resp:any) => {
              resolve(resp)
            })
          })
        default: // 智能应用
          return this.addDynamic(type, body)
      }
    },
    // 此处当entity为智能插件时，query是string类型，真实含义为uuid
    getEntity: async (type:string, query:any = {}) => {
      return new Promise((resolve) => {
        let url:string = ''
        let dynamic:boolean = false
        switch (type) {
          case 'api':
            url = 'router'
            break
          case 'app':
            url = 'application'
            break
          case 'app-auth':
            url = 'application/auth'
            break
          default: // 智能应用
            url = `dynamic/${type}/info/${query}`
            dynamic = true
        }

        this.api.get(url, dynamic ? {} : query).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    editEntity: async (type:string, body:any, option?:{uuid?:string, query?:any, showModal?:boolean}) => {
      return new Promise((resolve) => {
        let url:string = ''
        switch (type) {
          case 'api':
            url = 'router'
            break
          case 'app':
            url = 'application'
            break
          case 'app-auth':
            url = 'application/auth'
            break
          default: // 智能应用
            if (option?.uuid) {
              url = `dynamic/${type}/config/${option!.uuid}`
            } else {
              resolve({ code: -1, msg: '未检测到对应的uuid，请检查输入' })
            }
        }

        this.api.put(url, body).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    getEntities: async (type:string, query:any) => {
      return new Promise((resolve) => {
        let url:string = ''
        switch (type) {
          case 'api':
            url = 'routers'
            break
          case 'app-auth':
            url = 'application/auths'
            break
          default:
            resolve({ code: -1, msg: '暂不支持对输入的类型进行列表数据查询，请检查输入' })
        }

        this.api.get(url, query).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    publishEntity: async (type:string, data:{uuid:string, clusters:Array<string>, name?:string}, option?:{publishType?:'online'|'offline', showModal?:boolean}) => {
      return new Promise((resolve) => {
        if (option?.showModal) {
          switch (type) {
            case 'api':
              this.routerService.publishApiModal(data.uuid, undefined, (resp:any) => {
                resolve(resp)
              })
              break
            case 'app':
              this.applicationService.publishAppModal({ name: data.name!, id: data.uuid }, undefined, (resp:any) => { resolve(resp) })
              break
            default:
              this.intelPluginService.publishPluginModal(type, { name: data.name || '', id: data.uuid }, undefined, (resp:any) => { resolve(resp) })
          }
        }
        if (!option?.publishType) {
          resolve({ msg: '未检测到publishType，请检查输入' })
          return
        }
        let url:string = ''
        const query:{[k:string]:any} = {}
        let dynamic:boolean = false
        switch (type) {
          case 'api':
            url = `routers/${option!.publishType}`
            query['clusterNames'] = data.clusters
            break
          case 'app':
            url = `application/${option!.publishType}`
            query['clusterNames'] = data.clusters
            break
          default:
            url = `dynamic/${type}/${option!.publishType}/${data.uuid}`
            dynamic = true
        }

        this.api.put(url, dynamic ? { clusterName: data.clusters } : { cluster: data.clusters }, query).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    publishEntities: async (type:string, body:any, option:{query?:any, publishType:'online'|'offline', showModal?:boolean, showResModal?:boolean, showLastStep?:boolean}) => {
      if (option.showResModal) {
        return this.batchPublishResModal(type, option.publishType, body, option.showLastStep)
      }
      if (option.showModal) {
        return this.batchPublishModal(type, option.publishType, body)
      }
      return new Promise((resolve) => {
        let url:string = ''
        switch (type) {
          case 'api':
            url = `routers/batch-${option.publishType}`
            break
          default:
            resolve({ code: -1, msg: '暂不支持对输入类型的实体进行批量发布，请检查输入' })
        }

        this.api.post(url, body, option.query).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    publishStrategies: async (type:string, body:any, option:{query:any, showModal?:boolean}) => {
      return new Promise((resolve) => {
        if (option.showModal) {
          this.servGovernanceService.publishStrategyModal(type, option.query.cluster_name, undefined, (resp:any) => { resolve(resp) })
        }
        this.api.post(`strategy/${type}/publish`, body, option.query).subscribe((resp:any) => {
          resolve(resp)
        })
      })
    },
    chooseEnvVar: async () => {
      return this.envTableService.openModal()
    },
    getOptions: async (type:string, query?:any) => {
      return new Promise((resolve) => {
        let url:string = ''
        switch (type) {
          case 'cluster-env':
            url = 'cluster/enum'
            break
          case 'cluster':
            url = 'clusters/simple'
            break
          case 'service':
            url = 'common/enum/Service'
            break
          case 'app':
            url = 'application/enum'
            break
          case 'user':
            url = 'user/enum'
            break
          case 'notices':
            url = 'channels'
            break
          default:
            resolve({ msg: '查询类型有误，请检查输入' })
            return
        }
        this.api.get(url, query).subscribe((resp:any) => {
          resolve(resp)
        })
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

  batchPublishModal (type:string, publishType:'online'|'offline', data:any) {
    return new Promise((resolve) => {
      switch (type) {
        case 'api':
          this.routerService.batchPublishApiModal(publishType, data, (resp:any) => { resolve(resp) })
          break
        default:
          console.warn(`eo预警：无法调用发布管理方法，请检查传入的type=${type}是否正确`)
          resolve({ msg: `无法调用发布管理方法，请检查传入的type=${type}是否正确` })
      }
    })
  }

  batchPublishResModal (type:string, publishType:'online'|'offline', data:any, showLastStep?:boolean) {
    return new Promise((resolve) => {
      switch (type) {
        case 'api':
          this.routerService.batchPublishApiResModal(publishType, data, showLastStep
            ? () => {
                resolve({ data: { lastStep: true } })
              }
            : undefined, (resp:any) => {
            resolve(resp)
          })
          break
        default:
          console.warn(`eo预警：无法调用发布管理方法，请检查传入的type=${type}是否正确`)
          resolve({ msg: `无法调用发布管理方法，请检查传入的type=${type}是否正确` })
      }
    })
  }

  addDynamic (name:string, data:any) {
    return new Promise((resolve) => {
      this.api.post(`dynamic/${name}`, data).subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  addApi (data:any) {
    return new Promise((resolve) => {
      this.api.post('router', data).subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  // 控制台手动触发iframe路由跳转，暂时用在面包屑中
  reqFlashIframe (url:string) {
    this.changeIframe.next(url)
  }

  repFlashIframe () {
    return this.changeIframe.asObservable()
  }
}
