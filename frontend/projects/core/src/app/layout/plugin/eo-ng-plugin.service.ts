import { Injectable } from '@angular/core'
import { Subject, Subscription } from 'rxjs'
import { CardItem } from '../../component/card-list/card-list.component'
import { PluginItem, PluginGroupItem } from './types/types'
import { ApiService } from '../../service/api.service'
import { NzTreeNodeOptions } from 'ng-zorro-antd/tree'
import { BaseInfoService } from '../../service/base-info.service'

@Injectable({
  providedIn: 'root'
})
export class EoNgPluginService {
  private flashFlag: Subject<boolean> = new Subject<boolean>()
  public nodesList: NzTreeNodeOptions[] = []
  groupUuid: string = '' // 供右侧list页面用
  queryName:string = ''
  showAll:boolean = true
  radioValue:string|boolean = ''
  pluginList:(PluginItem & CardItem)[] = []
  totalNum:number = 0
  private subGetList: Subscription = new Subscription()
  constructor (private api:ApiService,
    private baseInfo:BaseInfoService) { }

  reqFlashGroup () {
    this.flashFlag.next(true)
  }

  repFlashGroup () {
    return this.flashFlag.asObservable()
  }

  reqFlashList () {
    this.flashFlag.next(true)
  }

  repFlashList () {
    return this.flashFlag.asObservable()
  }

  // 获取分组和插件列表
  getPluginList () {
    this.subGetList.unsubscribe()
    this.groupUuid = this.baseInfo.allParamsInfo.pluginGroupId
    this.subGetList = this.api.get('system/plugin/installed', { group: this.groupUuid || '', search: this.queryName || '' }).subscribe((resp:{code:number, data:{plugins:PluginItem[], groups:PluginGroupItem[]}, msg:string}) => {
      if (resp.code === 0) {
        this.nodesList = this.nodesTransfer(resp.data.groups)
        this.totalNum = 0
        for (const node of this.nodesList) {
          // eslint-disable-next-line dot-notation
          this.totalNum += node['count'] || 0
        }
        this.pluginList = resp.data.plugins.map((plugin:PluginItem) => {
          const newPlugin:PluginItem & CardItem = { ...plugin, title: plugin.cname, desc: plugin.resume, iconAddr: plugin.icon } as PluginItem & CardItem
          return newPlugin
        })
      }
    })
  }

  // 遍历目录列表，转化成tree组件需要的参数格式
  // 第一级目录不可以创建api，当root为true时，标志该目录为第一级，并放入firstLevelMap
  nodesTransfer (data:any): NzTreeNodeOptions[] {
    const res:NzTreeNodeOptions[] = []
    for (const index in data) {
      data[index].key = data[index].uuid
      data[index].title = `${data[index].name}（${data[index].count || 0}）`
      if (this.groupUuid && data[index].uuid === this.groupUuid) {
        data[index].selected = true
        this.showAll = false
      }
      if (data[index].children?.length > 0) {
        data[index].children = this.nodesTransfer(data[index].children)
      }
      res.push(data[index])
    }
    return res
  }
}
