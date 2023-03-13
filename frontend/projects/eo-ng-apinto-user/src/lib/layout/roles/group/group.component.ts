/* eslint-disable dot-notation */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-09-22 23:02:01
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-23 00:16:31
 * @FilePath: /apinto/projects/eo-ng-apinto-user/src/lib/layout/roles/group/group.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/* eslint-disable no-useless-constructor */
import { Component, Inject, OnInit, ViewChild } from '@angular/core'
import { NavigationEnd, Router } from '@angular/router'
import { NzTreeNode, NzTreeNodeOptions } from 'ng-zorro-antd/tree'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { API_SERVICE_ADAPTER, ApiServiceAdapter } from '../../../constant/api-service-adapter'
import { DataService } from '../../../service/data.service'
import { Subscription } from 'rxjs'
import { EoNgTreeDefaultComponent } from 'eo-ng-tree'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from '../../../constant/app-config-adapter'
import { BASEINFO_SERVICE_ADAPTER, BaseInfoServiceAdapter, RoleProfileComponent } from '../../../../public-api'

@Component({
  selector: 'eo-ng-apinto-roles-group',
  templateUrl: './group.component.html',
  styleUrls: ['./group.component.scss']
})
export class RolesGroupComponent implements OnInit {
  @ViewChild('eoNgTreeDefault') eoNgTreeDefault!: EoNgTreeDefaultComponent

  public nodesList:NzTreeNodeOptions[] = []
  public apiNodesMap:Map<string, any> = new Map()
  public apiNodesList:Array<any> = []
  nzDisabled:boolean = false // 权限控制
  customRolesMenus: Array<{ title: string, click:any, showFn?:any }> = []
  usersAmount:number = 1
  drawerRef:NzModalRef | undefined
  activatedNode?: NzTreeNode;
  selectedNodeKey:string = ''
  showAll:boolean = true
  roleIdFromUrl:string = '' // 从url获取的roleId
  private subscription: Subscription
  constructor (private message: EoNgFeedbackMessageService,
    private modalService:EoNgFeedbackModalService,
    @Inject(API_SERVICE_ADAPTER) private apiService: ApiServiceAdapter,
    @Inject(BASEINFO_SERVICE_ADAPTER) private baseInfo: BaseInfoServiceAdapter,
    private router:Router,
    private flashGroupService: DataService
  ) {
    this.subscription = this.flashGroupService.repFlashGroup().subscribe(() => {
      this.getMenuList()
    })
  }

  ngOnInit (): void {
    this.roleIdFromUrl = this.baseInfo.allParamsInfo.roleId
    this.getMenuList()
    this.customRolesMenus = [
      {
        title: '编辑',
        click: this.editCustomRoleDrawer
      },
      {
        title: '删除',
        click: this.deleteCustomRoleModal
      }
    ]
    if (this.baseInfo.allParamsInfo.roleId) {
      this.showAll = false
    }

    this.subscription = this.router.events.subscribe((event) => {
      if (event instanceof NavigationEnd) {
        this.roleIdFromUrl = this.baseInfo.allParamsInfo.roleId!
        this.getGroupItemSelected()
      }
    })
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  // 获取api列表
  getMenuList () {
    this.apiService.get('roles').subscribe((resp:any) => {
      if (resp.code === 0) {
        this.nodesList = resp.data.roles
        this.usersAmount = resp.data.total
        for (const index in this.nodesList) {
          this.nodesList[index].key = this.nodesList[index]['id']
        }
        this.getGroupItemSelected()
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  // 查看所有用户列表
  viewAllUsers () {
    this.showAll = true
    if (this.activatedNode?.isSelected) {
      this.activatedNode.isSelected = false
    }
    if (
      this.eoNgTreeDefault?.getTreeNodeByKey(this.selectedNodeKey)?.isSelected
    ) {
      // @ts-ignore
      this.eoNgTreeDefault.getTreeNodeByKey(this.selectedNodeKey).isSelected = false
    }
    this.router.navigate(['/', 'system', 'role'])
  }

  openDrawer (usage:string, data?:any) {
    this.drawerRef = this.modalService.create({
      nzTitle: usage === 'addRole' ? '新建角色' : '编辑角色',
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: RoleProfileComponent,
      nzComponentParams: {
        type: usage,
        roleId: data?.id,
        accessLink: 'system/role',
        closeModal: this.closeModal
      },
      nzOkText: usage === 'addRole' ? '保存' : '提交',
      nzOnOk: (component:RoleProfileComponent) => {
        component.saveRoleProfile()
        return false
      }
    })
    this.drawerRef.afterClose.subscribe(() => {
      this.getMenuList()
      this.flashGroupService.reqFlashMenu()
    })
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  closeModal = () => {
    this.drawerRef?.close()
  }

  editCustomRoleDrawer = (value:any) => {
    this.openDrawer('editRole', value)
  }

  deleteCustomRoleModal = (value:any) => {
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '如果该角色关联了用户，删除后关联的用户将失去该角色权限并会成为未分配角色用户，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteCustomRole(value)
      }
    })
  }

  deleteCustomRole (data:any) {
    this.apiService.delete('role', { id: data.id }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '删除角色成功!', { nzDuration: 1000 })
        this.getMenuList()
        this.flashGroupService.reqFlashMenu()
      } else {
        this.message.error(resp.msg || '删除角色失败!')
      }
    })
  }

  // 点击分组,查看不同角色的用户列表
  activeNode (data:any) {
    if (
      data.keys[0] &&
      this.selectedNodeKey !== data.keys[0] &&
      this.eoNgTreeDefault?.getTreeNodeByKey(this.selectedNodeKey)?.isSelected
    ) {
      // @ts-ignore
      this.eoNgTreeDefault.getTreeNodeByKey(this.selectedNodeKey).isSelected = false
    }
    this.selectedNodeKey = data.keys[0]
    this.showAll = false
    this.activatedNode = data.node!
    this.router.navigate(['/', 'system', 'role', data.node.origin.id], { queryParams: { operate_disable: data.node.origin.operate_disable } })
  }

  getGroupItemSelected () {
    if (!this.nodesList.length) {
      this.showAll = true
      return
    }
    let findTheSelected:boolean = false
    for (const index in this.nodesList) {
      if (this.nodesList[index].key === this.roleIdFromUrl) {
        this.nodesList[index].selected = true
        this.nodesList[index]['isSelected'] = true
        this.showAll = false
        findTheSelected = true
        this.selectedNodeKey = this.roleIdFromUrl
      } else {
        this.nodesList[index].selected = false
        this.nodesList[index]['isSelected'] = false
      }
    }
    if (!this.roleIdFromUrl) {
      this.showAll = true
      this.activatedNode = undefined
      this.selectedNodeKey = ''
    } else if (!findTheSelected) {
      this.router.navigate(['/', 'system', 'role'])
    }
  }
}
