/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
/* eslint-disable no-undef */
import { Component, ElementRef, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { NavigationEnd, Router } from '@angular/router'
import { NzTreeNodeOptions, NzFormatEmitEvent, NzTreeNode } from 'ng-zorro-antd/tree'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { Subscription } from 'rxjs'
import { EoNgTreeDefaultComponent } from 'eo-ng-tree'
import { ApiManagementEditGroupComponent } from './edit-group/edit-group.component'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiGroup, EmptyHttpResponse } from 'projects/core/src/app/constant/type'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'

@Component({
  selector: 'eo-ng-api-management',
  templateUrl: './group.component.html',
  styles: [
    `
    eo-ng-tree-default span{
      white-space:nowrap;
      text-overflow:ellipsis;
      overflow:hidden;
    }

    .group-icon{
      padding:0px;
      margin-left:var(--LAYOUT_PADDING);
      .iconpark-icon{
        height:30px !important;
        width:30px !important;
        color: var(--background-color);
      }
    }

    .group-icon[disabled] .iconpark-icon{
      cursor:not-allowed;
    }
    `
  ]
})
export class ApiManagementComponent implements OnInit {
  @ViewChild('addGroupTpl', { read: TemplateRef, static: true }) addGroupTpl: TemplateRef<any> | string = '';
  @ViewChild('groupComponent') groupComponent!: ElementRef
  @ViewChild('eoNgTreeDefault') eoNgTreeDefault!: EoNgTreeDefaultComponent

  public nodesList:NzTreeNodeOptions[] = []
  public apiNodesMap:Map<string, any> = new Map()
  public apiNodesList:Array<any> = []
  groupUuid:string = '' // 供右侧list页面用
  queryName:string = '' // 支持搜索目录名称和api名称
  expandAll:boolean = false
  firstLevelMap:Set<string> = new Set()
  showAll:boolean = true
  searchValue:string = ''
  activatedNode?: NzTreeNode;
  editUuid:string = '' // 正在编辑的分组uuid

  autoTips: Record<string, Record<string, string>> = defaultAutoTips

  groupModal: NzModalRef |undefined
  editParentUuid:string = ''
  private subscription: Subscription = new Subscription()

  constructor (private message: EoNgFeedbackMessageService,
    private modalService:EoNgFeedbackModalService,
    private baseInfo:BaseInfoService,
    private api:ApiService,
    private router:Router) {
  }

  ngOnInit (): void {
    this.groupUuid = this.baseInfo.allParamsInfo.apiGroupId

    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.groupUuid = this.baseInfo.allParamsInfo.apiGroupId
      }
    })

    this.getMenuList()
  }

  onDestroy () {
    this.subscription.unsubscribe()
  }

  // 获取api列表
  getMenuList (flash?:boolean) {
    this.api.get('router/groups').subscribe((resp:{code:number, data:ApiGroup, msg:string}) => {
      if (resp.code === 0) {
        this.queryName = ''
        this.nodesList = this.nodesTransfer(resp.data.root.groups, true)
        setTimeout(() => {
          if (flash && !this.editParentUuid && !this.editUuid) {
            this.groupScrollToBottom()
          }
        })
      }
    })
  }

  // 遍历目录列表，转化成tree组件需要的参数格式
  // 第一级目录不可以创建api，当root为true时，标志该目录为第一级，并放入firstLevelMap
  nodesTransfer (data:any, root?:boolean): NzTreeNodeOptions[] {
    if (root) {
      this.firstLevelMap = new Set()
    }
    const res:NzTreeNodeOptions[] = []
    for (const index in data) {
      data[index].key = data[index].uuid
      data[index].title = data[index].name
      if (this.groupUuid && data[index].uuid === this.groupUuid) {
        data[index].selected = true
        this.showAll = false
      } else if ((this.editParentUuid && data[index].uuid === this.editParentUuid)) {
        this.showAll = false
      }
      if (data[index].children?.length > 0) {
        data[index].children = this.nodesTransfer(data[index].children)
        data[index].expanded = this.findExpandChildren(data[index].children)
      }
      res.push(data[index])
      if (root) {
        this.firstLevelMap.add(data[index].uuid)
      }
    }
    return res
  }

  findExpandChildren (data:NzTreeNodeOptions[]):boolean {
    for (const index in data) {
      if (data[index].selected || data[index].expanded) {
        return true
      }
    }
    return false
  }

  // 添加分组时的弹窗
  addGroupModal = (uuid?:any) => {
    let title:string = '添加分组'
    if (uuid !== 'root') {
      title = '添加子分组'
    }
    this.groupModal = this.modalService.create({
      nzTitle: title,
      nzContent: ApiManagementEditGroupComponent,
      nzWidth: MODAL_SMALL_SIZE,
      nzComponentParams: {
        uuid: uuid,
        type: 'add',
        closeModal: this.closeModal,
        showUuid: false
      },
      nzClosable: true,
      nzCancelText: '取消',
      nzOkText: '确定',
      nzOnOk: (component:ApiManagementEditGroupComponent) => {
        this.editParentUuid = uuid === 'root' ? '' : uuid || ''
        component.addGroup(uuid)
        return false
      }
    })
  }

  groupScrollToBottom () {
    try {
      this.groupComponent.nativeElement.scrollTop = this.groupComponent.nativeElement.scrollHeight
    } catch (err) {}
  }

  // 编辑分组的弹窗
  editGroupModal = (uuid:string, name?:string) => {
    this.groupModal = this.modalService.create({
      nzTitle: '编辑分组',
      nzContent: ApiManagementEditGroupComponent,
      nzWidth: MODAL_SMALL_SIZE,
      nzComponentParams: {
        groupName: name,
        uuid: uuid,
        type: 'edit',
        closeModal: this.closeModal,
        showUuid: true
      },
      nzClosable: true,
      nzOkText: '确定',
      nzCancelText: '取消',
      nzOnOk: (component:ApiManagementEditGroupComponent) => {
        this.editUuid = uuid
        component.editGroup(uuid)
        return false
      }
    })
  }

  closeModal = () => {
    this.groupModal?.close()
    this.getMenuList(true)
  }

  // 删除分组的弹窗
  deleteGroupModal = (name:string, uuid:string) => {
    this.groupModal = this.modalService.create({
      nzTitle: '删除',
      nzContent: `删除${name}后，该分组下的所有子分组将全部移入回收站，该操作无法撤销，确认删除？`,
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkText: '确定',
      nzOkDanger: true,
      nzCancelText: '取消',
      nzOnOk: () => {
        this.deleteGroup(uuid, name)
        return false
      }
    })
  }

  // 删除分组
  deleteGroup (groupUuid:string, name:string) {
    this.api.delete('group/api/' + groupUuid, { name: name }).subscribe((resp:EmptyHttpResponse) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '删除成功', { nzDuration: 1000 })
        this.closeModal()
      }
    })
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
    this.router.navigate(['/', 'router', 'api', 'group', 'list', data.node!.origin['uuid']])
    this.activatedNode = data.node!
  }

  // 右侧页面切换至所有api的列表页
  viewAllApis () {
    this.showAll = true
    if (this.groupUuid && this.eoNgTreeDefault?.getTreeNodeByKey(this.groupUuid)?.isSelected) {
      this.eoNgTreeDefault.getTreeNodeByKey(this.groupUuid)!.isSelected = false
    }
    if (this.activatedNode?.isSelected) {
    this.activatedNode!.isSelected = false
    }
    this.router.navigate(['/', 'router', 'api', 'group', 'list'])
  }

  // 右侧页面切换至新建API的页面
  addApi = (uuid:string) => {
    this.router.navigate(['/', 'router', 'api', 'create', uuid])
  }
}
