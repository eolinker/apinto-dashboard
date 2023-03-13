/* eslint-disable no-useless-constructor */
import { Injectable } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TabsOptions } from 'eo-ng-tabs'
import { Observable } from 'rxjs'
// import { APITabParams, TabsOptions } from '../constant/monitor-tabs'
import { MonitorAlarmComponent } from '../layout/monitor-alarm/monitor-alarm.component'
import { ApiService } from './api.service'

interface MonitorTabOptions extends TabsOptions{
  clusters?:Array<string>,
  uuid?:string,
  name:string
}
export class TabEvent {
  tab: MonitorTabOptions
  constructor (tab: MonitorTabOptions) {
    this.tab = tab
  }
}
export class InitEvent extends TabEvent {}
export class DestroyEvent extends TabEvent {}

export class History {
  queue: MonitorTabOptions[] = new Proxy([], {
    set (target: any[], prop: any, value, receiver) {
      if (prop !== 'length' && target.length >= 2) {
        target.shift()
        target.push(value)
        return true
      }
      if (prop === 'length' && value >= 2) {
        target[prop] = 2
        return true
      }
      target[prop] = value
      return true
    }
  })
}

@Injectable({
  providedIn: 'root'
})
export class EoNgMonitorTabsService {
  history = new History()
  events!: Observable<InitEvent | DestroyEvent>
  index: number = 0
  prevIndex:number = 0
  isLock: boolean = false
  list: MonitorTabOptions[] = []
  showList: MonitorTabOptions[] = []
  addTabFlag:boolean = false // 是否正在新建tab(仅在手动添加时有效)
  tabsMap:Map<string, MonitorTabOptions> = new Map()
  componentInstance!: MonitorAlarmComponent
  constructor (private api: ApiService, private message:EoNgFeedbackMessageService, private router:Router) {}

  getTabsData (toPartition?:string) {
    this.api.get('monitor/partitions').subscribe((resp:{code:number, data:{partitions:Array<{uuid:string, name:string, clusterNames:Array<string>}>}, msg:string}) => {
      if (resp.code === 0) {
        if (resp.data.partitions.length > 0) {
          this.getTabs(resp.data.partitions)
        }
        if (toPartition) {
          this.router.navigate(['/', 'monitor-alarm', 'area', 'total', toPartition])
        }
      } else {
        this.message.error(resp.msg || '获取分区列表失败，请重试！')
      }
    })
  }

  // 从后端获取的分区数据，根据后端数据uuid判断service存的tab是否保留
  getTabs (tabsList:any, currentPartitionId?:string) {
    this.list = []
    for (const tab of tabsList) {
      this.list.push({
        title: tab.name,
        name: tab.name,
        routerLink: this.tabsMap.get(tab.uuid)?.routerLink || ('/monitor-alarm/area/total/' + tab.uuid),
        queryParams: this.tabsMap.get(tab.uuid)?.queryParams || {},
        queryParamsHandling: '',
        clusters: tab.clusterNames,
        uuid: tab.uuid
      })
      if (!this.tabsMap.get(tab.uuid)) {
        this.tabsMap.set(tab.uuid, tab)
      }
    }
    if (currentPartitionId) {
      this.changeTab(currentPartitionId)
    }
    this.showList = [...this.list]
  }

  // 根据uuid获取集群列表
  getClusters (uuid:string):Array<string> {
    for (const index in this.list) {
      if (this.list[index].uuid === uuid) {
        return this.list[index].clusters || []
      }
    }
    return []
  }

  // 根据uuid分区获取名称
  getTabName (uuid:string):string {
    for (const index in this.list) {
      if (this.list[index].uuid === uuid) {
        return this.list[index].name
      }
    }
    return ''
  }

  /** 新增标签页 */
  newTab () {
    this.prevIndex = this.index
    this.list.push({
      title: '新建分区',
      name: '新建分区',
      routerLink: '/monitor-alarm/area/config'
    })
    this.index = this.list.length - 1
  }

  // 修改对应tab的路由
  changeRouter (uuid:string, router:string, params?:any) {
    for (const index in this.list) {
      if (this.list[index].uuid === uuid) {
        this.list[index].routerLink = router
        this.list[index].queryParams = { ...this.list[index].queryParams, ...params }
        this.showList[index].routerLink = router
        this.showList[index].queryParams = { ...this.list[index].queryParams, ...params }
      }
    }
  }

  //  删除uuid对应的路由，并定位到下一个tab
  deleteTab (uuid:string) {
    for (let i = 0; i < this.list.length; i++) {
      if (this.list[i].uuid === uuid) {
        this.list.splice(i, 1)
        this.showList.splice(i, 1)
        if (i >= this.list.length && i > 0) {
          this.index = this.prevIndex ? this.prevIndex : i - 1
          this.prevIndex = 0
        }
        break
      }
    }
  }

  // 根据分区id确认index
  changeTab (uuid:string) {
    for (let i = 0; i < this.list.length; i++) {
      if (this.list[i].uuid === uuid) {
        this.index = i
        break
      }
    }
  }

  // getComponentId () {
  //   return this.getCurrentTab()?.componentId || ''
  // }

  // /** 改变当前标签页索引值 */
  // changeIndex (index: number) {
  //   this.index = index
  // }

  // getAllTab () {
  //   return this.list
  // }

  // /** 获取当前标签页 */
  // getCurrentTab () {
  //   return this.list[this.index]
  // }

  // /** 通过索引获取Tab */
  // getTabByIndex (index: number) {
  //   return this.list[index]
  // }

  // /** 获取最后一个标签页 */
  // getLastTab () {
  //   return this.list[this.list.length - 1]
  // }

  // /** 编辑状态改变 */
  // /** 标记动态组件已加载 */
  // loaded (index: number, componentId?: string) {
  //   this.getTabByIndex(index).loaded = true
  //   this.getTabByIndex(index).componentId = componentId
  // }

  // /** 替换当前标签页 */
  // // replaceCurrentTab (tab: MonitorTabOptions) {
  // //   this.list[this.index] = tab
  // //   setTimeout(() => {
  // //     this.componentInstance.loadComponent()
  // //   }, 0)
  // // }

  // /** 更改当前标签页 */
  // changeCurrentTabInfo (params: MonitorTabOptions) {
  //   return (this.list[this.index] = params)
  // }

  // listener = (event: BeforeUnloadEvent) => {
  //   event.preventDefault()
  //   event.returnValue = ''
  // }

  // /** 更改当前标签页编辑状态 */
  // changeCurrentTabModified (status: boolean) {
  //   this.list[this.index].modified = status
  //   if (status) {
  //     window.addEventListener('beforeunload', this.listener)
  //   }
  //   if (!status && this.listener) {
  //     window.removeEventListener('beforeunload', this.listener)
  //   }
  //   return this.list[this.index]
  // }

  // /** 更改当前标签页名称 */
  // changeCurrentTabName (name: string) {
  //   this.list[this.index].name = name
  //   return this.list[this.index]
  // }

  // changeCurrentTabPath (path: string[]) {
  //   this.list[this.index].path = path
  //   return this.list[this.index]
  // }

  // /** 查找Tab根据TabParams中某个值 */
  // findTabByParams (key: keyof APITabParams, value: string) {
  //   return this.list.find((tab) => tab.params[key] + '' === value + '')
  // }

  // /** 根据ApiId查找标签页 */
  // findTabByApiId (apiId: string) {
  //   return this.findTabByParams('apiId', apiId)
  // }

  // /** 根据GroupId查找标签页 */
  // findTabByGroupId (apiId: string) {
  //   return this.findTabByParams('apiId', apiId)
  // }

  // /** 根据ProjectHash查找标签页 */
  // findTabByProjectHash (apiId: string) {
  //   return this.findTabByParams('apiId', apiId)
  // }

  // /**
  //  * 根据ApiId删除标签页
  //  * @param apiId 要删除的ApiId
  //  * @returns 找到删除目标则返回待删除tab，找不到则返回false
  //  */
  // delTabByApiId (apiId: string): false | MonitorTabOptions {
  //   const target = this.findTabByApiId(apiId)
  //   target && this.delTab(target)
  //   return target || false
  // }

  // /**
  //  * 根据GroupId删除标签页
  //  * @param groupId 要删除的groupId
  //  * @returns 找到删除目标则返回待删除tab，找不到则返回false
  //  */
  // delTabByGroupId (groupId: string): false | MonitorTabOptions {
  //   const target = this.findTabByApiId(groupId)
  //   target && this.delTab(target)
  //   return target || false
  // }

  // /**
  //  * 根据ProjectHash删除标签页
  //  * @param groupId 要删除的ProjectHash
  //  * @returns 找到删除目标则返回待删除tab，找不到则返回false
  //  */
  // delTabByProjectHash (projectHash: string): false | MonitorTabOptions {
  //   const target = this.findTabByApiId(projectHash)
  //   target && this.delTab(target)
  //   return target || false
  // }

  // /** 删除标签页 */
  // delTab (target: MonitorTabOptions) {
  //   this.list = this.list.filter((tab) => tab !== target)
  // }

  // /** 根据索引删除标签页 */
  // delTabByIndex (index?: number) {
  //   const prevIndex = index as number
  //   this.batchDelTabByIndex([index || 0])
  // }

  // /** 根据索引批量删除标签页 */
  // batchDelTabByIndex (idxs: number[]) {
  //   this.list = this.list.filter((tab, tabIndex) => {
  //     return !idxs.includes(tabIndex)
  //   })
  // }

  // /** 保留这些标签页 */
  // saveTab (tabs: MonitorTabOptions[]) {
  //   this.list = this.list.filter((tab) => {
  //     return tabs.includes(tab)
  //   })
  // }

  // /** 清空标签页列表 */
  // clearTabList () {
  //   this.list = []
  //   setTimeout(() => {
  //     // this.router.navigate(['/', 'home', 'api-studio', 'inside', this.baseService.projectHashKey, 'api', '-1'])
  //   }, 0)
  // }

  // simpleClearTabList () {
  //   this.list = []
  // }
}
