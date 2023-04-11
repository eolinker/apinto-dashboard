import { Component, OnInit, ViewChild } from '@angular/core'
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { RadioOption } from 'eo-ng-radio'
import { EoNgTreeDefaultComponent } from 'eo-ng-tree'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { NzTreeNode, NzTreeNodeOptions, NzFormatEmitEvent } from 'ng-zorro-antd/tree'
import { Subscription } from 'rxjs'
import { CardItem } from '../../../component/card-list/card-list.component'
import { MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { ApiService } from '../../../service/api.service'
import { EoNgNavigationService } from '../../../service/eo-ng-navigation.service'
import { BaseInfoService } from '../../../service/base-info.service'
import { EoNgMessageService } from '../../../service/eo-ng-message.service'
import { PluginCreateComponent } from '../create/create.component'
import { PluginListStatusItems } from '../types/conf'
import { PluginItem, PluginGroupItem } from '../types/types'

@Component({
  selector: 'eo-ng-plugin-list',
  templateUrl: './list.component.html',
  styles: [
  ]
})
export class PluginListComponent implements OnInit {
  @ViewChild('eoNgTreeDefault') eoNgTreeDefault!: EoNgTreeDefaultComponent
  nzDisabled:boolean = false
  radioOptions:RadioOption[] = [...PluginListStatusItems]
  radioValue:string|boolean = ''
  pluginList:(PluginItem & CardItem)[] = []
  modalRef:NzModalRef | undefined
  showAll:boolean = true
  groupUuid:string = '' // 供右侧list页面用
  queryName:string = '' // 支持搜索目录名称和api名称
  activatedNode?: NzTreeNode;
  mdFileName:string = ''
  private subscription: Subscription = new Subscription()
  public nodesList:NzTreeNodeOptions[] = []

  constructor (
    public api:ApiService,
    private modalService:EoNgFeedbackModalService,
    private appConfigService:EoNgNavigationService,
    private router:Router,
    private route: ActivatedRoute,
    private baseInfo:BaseInfoService,
    private message:EoNgMessageService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '企业插件' }])
  }

  ngOnInit (): void {
    this.groupUuid = this.baseInfo.allParamsInfo.pluginGroupId
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.groupUuid = this.baseInfo.allParamsInfo.pluginGroupId
        this.getPluginList()
      }
    })

    this.getPluginList()
  }

  onDestroy () {
    this.subscription.unsubscribe()
  }

  // 获取分组和插件列表
  getPluginList () {
    this.api.get('system/plugin/installed', { group: this.groupUuid, search: this.queryName }).subscribe((resp:{code:number, data:{plugins:PluginItem[], groups:PluginGroupItem[]}, msg:string}) => {
      if (resp.code === 0) {
        this.nodesList = this.nodesTransfer(resp.data.groups)
        this.pluginList = resp.data.plugins.map((plugin:PluginItem) => {
          const newPlugin:PluginItem & CardItem = { ...plugin, title: plugin.cname, desc: plugin.resume, iconAddr: plugin.icon } as PluginItem & CardItem
          return newPlugin
        })
      }
    })
  }

  // 根据状态展示响应的插件（前端筛选）
  filterPluginList () {
    if (this.radioValue === '') {
      return this.pluginList
    } else {
      return this.pluginList.filter((plugin:PluginItem) => {
        return plugin.enable === this.radioValue
      })
    }
  }

  // 右侧页面切换至所有插件的列表页
  viewAllPlugins () {
    this.showAll = true
    if (this.groupUuid && this.eoNgTreeDefault?.getTreeNodeByKey(this.groupUuid)?.isSelected) {
    this.eoNgTreeDefault.getTreeNodeByKey(this.groupUuid)!.isSelected = false
    }
    if (this.activatedNode?.isSelected) {
  this.activatedNode!.isSelected = false
    }
    this.router.navigate(['/', 'plugin', 'list', ''])
  }

  // 遍历目录列表，转化成tree组件需要的参数格式
  // 第一级目录不可以创建api，当root为true时，标志该目录为第一级，并放入firstLevelMap
  nodesTransfer (data:any): NzTreeNodeOptions[] {
    const res:NzTreeNodeOptions[] = []
    for (const index in data) {
      data[index].key = data[index].uuid
      data[index].title = data[index].name
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

  // 点击分组节点时，切换activatedNode
  // 当节点是目录时，右侧页面需要跳转至list页
  // 当节点是api时，右侧页面需要跳转至API编辑页
  activeNode (data: NzFormatEmitEvent): void {
    if (
      data.keys![0] &&
      this.groupUuid !== data.keys![0] &&
      this.eoNgTreeDefault?.getTreeNodeByKey(this.groupUuid)?.isSelected
    ) {
      // @ts-ignore
      this.eoNgTreeDefault.getTreeNodeByKey(this.groupUuid).isSelected = false
    }

    this.showAll = false
    data.node!.isExpanded = !data.node!.isExpanded
    // eslint-disable-next-line dot-notation
    if (data.node!.origin['uuid']) {
      // eslint-disable-next-line dot-notation
      this.router.navigate(['/', 'plugin', 'list', data.node!.origin['uuid']])
    } else {
      // eslint-disable-next-line dot-notation
      this.router.navigate(['/', 'plugin', 'list', ''])
    }
    this.activatedNode = data.node!
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  installPlugin () {
    this.modalRef = this.modalService.create({
      nzTitle: '安装插件',
      nzWidth: MODAL_SMALL_SIZE,
      nzContent: PluginCreateComponent,
      nzOkDisabled: this.nzDisabled,
      nzOnOk: (component:PluginCreateComponent) => {
        if (component.submit()) {
          this.getPluginList()
          return true
        } else {
          return false
        }
      }
    })
  }

  handerCardClick (card:CardItem) {
    // eslint-disable-next-line dot-notation
    this.router.navigate(['../../message', card['id'] || 'test', ''], { relativeTo: this.route })
  }
}
