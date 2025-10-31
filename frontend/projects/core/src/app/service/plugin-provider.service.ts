import { Injectable } from '@angular/core'
import { EoNgNavigationService } from './eo-ng-navigation.service'
import { EoNgBreadcrumbOptions } from 'eo-ng-breadcrumb'
import { Router } from '@angular/router'
import { ApiService } from './api.service'
import { RouterService } from '../layout/api/router.service'
import { EoNgApplicationService } from '../layout/application/application.service'
import { EditableEnvTableService } from '../component/editable-env-table/editable-env-table.service'
import { EoIntelligentPluginService } from '../component/intelligent-plugin/intelligent-plugin.service'
import { ServiceGovernanceService } from '../layout/serv-governance/service-governance.service'
import { BaseInfoService } from './base-info.service'

@Injectable({
  providedIn: 'root'
})
export class PluginProviderService {
  constructor (
    private navigationService:EoNgNavigationService,
    private router:Router,
    private api:ApiService,
    private routerService:RouterService,
    private applicationService:EoNgApplicationService,
    private envTableService:EditableEnvTableService,
    private intelPluginService:EoIntelligentPluginService,
    private servGovernanceService:ServiceGovernanceService,
    private baseInfo:BaseInfoService) {
  }

  getModuleAccess (module:string) {
    return this.navigationService.getUserModuleAccess(module)
  }

  get breadcrumb () {
    return this.navigationService.getLatestBreadcrumb()
  }

  set breadcrumb (breadcrumb: EoNgBreadcrumbOptions[]) {
    this.navigationService.reqFlashBreadcrumb(breadcrumb)
  }

  get userInfo () {
    return this.navigationService.userInfo
  }

  set userInfo (userInfo:any) {
    this.navigationService.userInfo = userInfo
  }

  get mainPage () {
    return this.navigationService.getPageRoute
  }

  get dashboardVersion () {
    return {
      version: this.baseInfo.version,
      updateDate: this.baseInfo.updateDate,
      productName: this.baseInfo.product
    }
  }

  goToMainPage = () => {
    this.router.navigate([this.mainPage])
  }

  renewMenu () {
    this.navigationService.getMenuList().subscribe(() => {
    })
  }

  setRouterConfig = (isRoot:boolean, pluginRouterConfig:any, routerConfig:any = this.router.config) => {
    if (isRoot) {
      routerConfig.unshift(pluginRouterConfig)
    } else {
      const basicRouter = routerConfig.find((item:any) => item?.data?.type === 'basicLayout')
      basicRouter.children.unshift(pluginRouterConfig)
    }
  }

  // 获取作为策略选项的列表数据
  strategyRemote = (type:string, data?:any) => {
    return new Promise((resolve) => {
      return this.api.get(`strategy/filter-remote/${type}`, data).subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  clusterList= async () => {
    return new Promise((resolve) => {
      return this.api.get('cluster/enum').subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  apiGroupList = () => {
    return new Promise((resolve) => {
      return this.api.get('router/groups').subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  apiSimpleList= async (data:any) => {
    return new Promise((resolve) => {
      return this.api.get('router/enum', data).subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  serviceSimpleList= async () => {
    return new Promise((resolve) => {
      return this.api.get('common/enum/Service').subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  appsSimpleList= async () => {
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
  }

  userList= async () => {
    return new Promise((resolve) => {
      this.api.get('user/enum').subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  noticeChannels= async () => {
    return new Promise((resolve) => {
      this.api.get('channels').subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  // apiSpace
  // 根据uuid判断api是否已存在
  checkApi = async (uuid:string) => {
    return new Promise((resolve) => {
      this.api.get('router/check', { uuid: uuid }).subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  // 根据分组名称和uuid判断分组是否已存在
  checkGroupName = async (type:string, data:any) => {
    return new Promise((resolve) => {
      this.api.put(`group/${type}/check`, data).subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  // 新建分组
  addGroup = async (type:string, body:any, option?:{query?:any, showModal?:boolean}) => {
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
  }

  getGroup = async (type:string, query?:any) => {
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
  }

  addEntity = async (type:string, body:any, option?:{query?:any, showModal?:boolean}) => {
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
  }

  // 此处当entity为智能插件时，query是string类型，真实含义为uuid
  getEntity = async (type:string, query:any = {}) => {
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
  }

  editEntity = async (type:string, body:any, option?:{uuid?:string, query?:any, showModal?:boolean}) => {
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
  }

  getEntities = async (type:string, query:any) => {
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
  }

  publishEntity = async (type:string, data:{uuid:string, clusters:Array<string>, name?:string}, option?:{publishType?:'online'|'offline', showModal?:boolean}) => {
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
  }

  publishEntities = async (type:string, body:any, option:{query?:any, publishType:'online'|'offline', showModal?:boolean, showResModal?:boolean, showLastStep?:boolean}) => {
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
  }

  publishStrategies = async (type:string, body:any, option:{query:any, showModal?:boolean}) => {
    return new Promise((resolve) => {
      if (option.showModal) {
        this.servGovernanceService.publishStrategyModal(type, option.query.cluster_name, undefined, (resp:any) => { resolve(resp) })
      }
      this.api.post(`strategy/${type}/publish`, body, option.query).subscribe((resp:any) => {
        resolve(resp)
      })
    })
  }

  chooseEnvVar = async () => {
    return this.envTableService.openModal()
  }

  getOptions = async (type:string, query?:any) => {
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
}
