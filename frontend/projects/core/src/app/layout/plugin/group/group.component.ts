import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgTreeDefaultComponent } from 'eo-ng-tree'
import { NzFormatEmitEvent, NzTreeNode } from 'ng-zorro-antd/tree'
import { ApiService } from '../../../service/api.service'
import { EoNgNavigationService } from '../../../service/app-config.service'
import { BaseInfoService } from '../../../service/base-info.service'
import { EoNgPluginService } from '../eo-ng-plugin.service'

@Component({
  selector: 'eo-ng-plugin-group',
  templateUrl: './group.component.html',
  styles: [
  ]
})
export class GroupComponent implements OnInit {
  @ViewChild('addGroupRef', { read: TemplateRef, static: true }) addGroupRef:
    | TemplateRef<any>
    | string = ''

  @ViewChild('eoNgTreeDefault') eoNgTreeDefault!: EoNgTreeDefaultComponent

  activatedNode: NzTreeNode | null = null
  clusterName: string = ''
  clusterKey: string = ''
  queryName:string = ''
  constructor (
    private baseInfo:BaseInfoService,
    private api: ApiService,
    private router: Router,
    private navigationService:EoNgNavigationService,
    public service:EoNgPluginService
  ) {
    this.navigationService.reqFlashBreadcrumb([{ title: '企业插件' }])
  }

  ngOnInit (): void {
    this.service.getPluginList()
  }

  viewAllPlugins () {
    this.service.showAll = true
    if (this.service.groupUuid && this.eoNgTreeDefault?.getTreeNodeByKey(this.service.groupUuid)?.isSelected) {
      this.eoNgTreeDefault.getTreeNodeByKey(this.service.groupUuid)!.isSelected = false
    }
    if (this.activatedNode?.isSelected) {
    this.activatedNode!.isSelected = false
    }
    this.router.navigate(['/', 'module-plugin', 'group', 'list'])
  }

  getPluginList () {
    this.service.reqFlashList()
  }

  // 点击分组节点时，切换activatedNode
  // 当节点是目录时，右侧页面需要跳转至list页
  // 当节点是api时，右侧页面需要跳转至API编辑页
  activeNode (data: NzFormatEmitEvent): void {
    if (
      data.keys![0] &&
      this.service.groupUuid !== data.keys![0] &&
      this.eoNgTreeDefault?.getTreeNodeByKey(this.service.groupUuid)?.isSelected
    ) {
      // @ts-ignore
      this.eoNgTreeDefault.getTreeNodeByKey(this.service.groupUuid).isSelected = false
    }

    this.service.showAll = false
    data.node!.isExpanded = !data.node!.isExpanded
    // eslint-disable-next-line dot-notation
    if (data.node!.origin['uuid']) {
      // eslint-disable-next-line dot-notation
      this.router.navigate(['/', 'module-plugin', 'group', 'list', data.node!.origin['uuid']])
    } else {
      // eslint-disable-next-line dot-notation
      this.router.navigate(['/', 'module-plugin', 'group', 'list', ''])
    }
    this.activatedNode = data.node!
  }
}
