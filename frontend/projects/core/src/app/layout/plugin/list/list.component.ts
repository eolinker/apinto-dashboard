import { Component, OnInit, ViewChild } from '@angular/core'
import { ActivatedRoute, NavigationEnd, Router } from '@angular/router'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { RadioOption } from 'eo-ng-radio'
import { EoNgTreeDefaultComponent } from 'eo-ng-tree'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { NzTreeNode, NzTreeNodeOptions } from 'ng-zorro-antd/tree'
import { Subscription } from 'rxjs'
import { CardItem } from '../../../component/card-list/card-list.component'
import { MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { ApiService } from '../../../service/api.service'
import { EoNgNavigationService } from '../../../service/eo-ng-navigation.service'
import { BaseInfoService } from '../../../service/base-info.service'
import { EoNgMessageService } from '../../../service/eo-ng-message.service'
import { PluginCreateComponent } from '../create/create.component'
import { PluginListStatusItems } from '../types/conf'
import { PluginItem } from '../types/types'
import { EoNgPluginService } from '../eo-ng-plugin.service'

@Component({
  selector: 'eo-ng-plugin-list',
  templateUrl: './list.component.html',
  styles: [
    `:host ::ng-deep{
      display:block;
      height:100%;
      overflow:hidden;
    }`
  ]
})
export class PluginListComponent implements OnInit {
  @ViewChild('eoNgTreeDefault') eoNgTreeDefault!: EoNgTreeDefaultComponent
  nzDisabled:boolean = false
  radioOptions:RadioOption[] = [...PluginListStatusItems]
  radioValue:string|boolean = ''
  modalRef:NzModalRef | undefined
  showAll:boolean = true
  groupUuid:string = '' // 供右侧list页面用
  queryName:string = '' // 支持搜索目录名称和api名称
  activatedNode?: NzTreeNode;
  mdFileName:string = ''
  private subscription: Subscription = new Subscription()
  private subscription1: Subscription = new Subscription()
  public nodesList:NzTreeNodeOptions[] = []

  constructor (
    public api:ApiService,
    private modalService:EoNgFeedbackModalService,
    private appConfigService:EoNgNavigationService,
    private router:Router,
    private route: ActivatedRoute,
    private baseInfo:BaseInfoService,
    private message:EoNgMessageService,
    public service:EoNgPluginService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '企业插件' }])
  }

  ngOnInit (): void {
    this.groupUuid = this.baseInfo.allParamsInfo.pluginGroupId
    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.groupUuid = this.baseInfo.allParamsInfo.pluginGroupId
        this.service.groupUuid = this.groupUuid
        this.service.getPluginList()
        console.log('sub1')
      }
    })
    this.subscription1 = this.service.repFlashList().subscribe(() => {
      this.queryName = this.service.queryName
      this.service.getPluginList()
      console.log('sub2')
    })
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
    this.subscription1.unsubscribe()
    this.service.pluginList = [] // 清空pluginList，否则第二次进入该组件时会多请求一次所有图片
  }

  // 根据状态展示响应的插件（前端筛选）
  filterPluginList () {
    if (this.radioValue === '') {
      return this.service.pluginList
    } else {
      return this.service.pluginList.filter((plugin:PluginItem) => {
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
    this.router.navigate(['/', 'plugin', 'group', 'list', ''])
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
          this.service.getPluginList()
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
