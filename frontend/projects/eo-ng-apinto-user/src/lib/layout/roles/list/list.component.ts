/* eslint-disable dot-notation */
import {
  Component,
  EventEmitter,
  Inject,
  OnInit,
  Output,
  TemplateRef,
  ViewChild
} from '@angular/core'
import {
  EoNgFeedbackMessageService,
  EoNgFeedbackModalService
} from 'eo-ng-feedback'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { Subscription } from 'rxjs'
import { UserProfileComponent } from '../../../component/user-profile/user-profile.component'
import {
  API_SERVICE_ADAPTER,
  ApiServiceAdapter
} from '../../../constant/api-service-adapter'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from '../../../constant/app-config-adapter'
import { AppServiceAdapter, APP_SERVICE_ADAPTER } from '../../../constant/app-service-adapter'
import { BASEINFO_SERVICE_ADAPTER, BaseInfoServiceAdapter } from '../../../constant/base-info-service-adapter'
import { DataService } from '../../../service/data.service'
import { ChangeUserRoleComponent } from '../change-user-role/change-user-role.component'
import { RefreshPswComponent } from '../refresh-psw/refresh-psw.component'
@Component({
  selector: 'eo-ng-apinto-roles-list',
  templateUrl: './list.component.html',
  styles: [
    `
    `
  ]
})
export class RolesListComponent implements OnInit {
  @ViewChild('switchTpl', { read: TemplateRef, static: true }) switchTpl:
    | TemplateRef<any>
    | undefined

  @ViewChild('selectTpl', { read: TemplateRef, static: true }) selectTpl:
    | TemplateRef<any>
    | undefined

  @Output() flashList:EventEmitter<any> = new EventEmitter()
  usersList: Array<any> = [] // 用户列表
  usersSet: Set<any> = new Set() // 被选择用来删除或移除角色的用户
  modalRef: NzModalRef | undefined
  nzDisabled:boolean = false // 权限控制
  role: string = '' // 当前列表的角色id(自定义角色的用户列表中使用)
  keyword: string = '' // 当前列表的搜索关键词
  rolesMap: Map<string, string> = new Map() // 角色id与角色名的映射表
  rolesList: Array<any> = [] // 角色下拉选择器的选项数组
  listType: boolean = true // true是为所有用户列表
  disabledEditUser: boolean = false // 禁止编辑添加删除搜索用户,目前为超管列表
  srollX:string = 'auto' // 表格滚动范围
  usersTableHeadName: Array<any> = [
    {
      type: 'checkbox',
      showFn: () => {
        return !this.disabledEditUser && !this.nzDisabled
      },
      click: (item: any) => {
        this.changeUsersSet(item, 'all')
      }
    },
    { title: '账号' },
    { title: '名称' },
    {
      title: '状态',
      showFn: () => {
        return this.listType
      },
      width: 80
    },
    {
      title: '邮箱'
    },
    {
      title: '角色'
    },
    {
      title: '更新者'
    },
    {
      title: '更新时间'
    },
    {
      title: '操作',
      right: true
    }
  ]

  usersTableBody: Array<any> = [
    {
      key: 'checked',
      type: 'checkbox',
      showFn: (item: any) => {
        return !this.disabledEditUser && item.role !== '超级管理员' && !this.nzDisabled
      },
      click: (item: any) => {
        this.changeUsersSet(item)
      }
    },
    {
      showFn: (item: any) => {
        return !this.disabledEditUser && item.role === '超级管理员' && !this.nzDisabled
      }
    },
    { key: 'user_name' },
    { key: 'nick_name' },
    {
      key: 'status',
      showFn: (item: any) => {
        return this.listType && item.role !== '超级管理员'
      }
    },
    {
      key: '',
      showFn: (item: any) => {
        return this.listType && item.role === '超级管理员'
      }
    },
    { key: 'email' },
    {
      key: 'role',
      showFn: (item: any) => {
        return this.listType || (!this.listType && item.role === '超级管理员')
      }
    },
    {
      key: 'role',
      showFn: (item: any) => {
        return !this.listType && item.role !== '超级管理员'
      }
    },
    { key: 'operator' },
    { key: 'update_time' },
    {
      type: 'btn',
      right: true,
      showFn: (item: any) => {
        return (
          !item.operate_disable &&
          !this.disabledEditUser &&
          this.listType
        )
      },
      btns: [
        {
          title: '重置密码',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item: any) => {
            this.resetPswModal(item.data)
          }
        },
        {
          title: '编辑',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item: any) => {
            this.openModal('editUser', item.data)
          }
        },
        {
          title: '删除',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item: any) => {
            this.deleteUserModal(item.data)
          }
        }
      ]
    },
    {
      type: 'btn',
      right: true,
      showFn: (item: any) => {
        return (
          !this.disabledEditUser &&
          !this.listType &&
          !item.operate_disable
        )
      },
      btns: [
        {
          title: '移除',
          disabledFn: () => {
            return this.nzDisabled
          },
          click: (item: any) => {
            this.removeUserRoleModal(item.data)
          }
        }
      ]
    },
    {
      key: '',
      right: true,
      showFn: (item: any) => {
        return (
          item.operate_disable || this.disabledEditUser
        )
      }
    }
  ]

  addUsersSet: Set<any> = new Set() // 被选择用来添加的用户

  // eslint-disable-next-line no-useless-constructor
  constructor (
    private message: EoNgFeedbackMessageService,
    private modalService: EoNgFeedbackModalService,
    @Inject(BASEINFO_SERVICE_ADAPTER) private baseInfo: BaseInfoServiceAdapter,
    @Inject(API_SERVICE_ADAPTER) private apiService: ApiServiceAdapter,
    @Inject(APP_SERVICE_ADAPTER) public appService: AppServiceAdapter,
    private dataService: DataService
  ) {
    this.subscription = this.dataService.repFlashMenu().subscribe(() => {
      this.getUsersList()
      this.flash = true
    })
  }

  private subscription: Subscription = new Subscription()
  private subscription2: Subscription = new Subscription()
  flash:boolean = false
  ngOnInit (): void {
    this.role = this.baseInfo.allParamsInfo.roleId || ''
    this.disabledEditUser = this.baseInfo.allParamsInfo?.operate_disable === 'true'
    this.listType = !this.role
    this.getUsersList()
  }

  ngAfterViewInit (): void {
    this.usersTableBody[4].title = this.switchTpl
    this.usersTableBody[8].title = this.selectTpl
  }

  ngDoCheck () {
    if (Object.keys(this.baseInfo.allParamsInfo).length > 0 && this.baseInfo.allParamsInfo.roleId !== this.role) {
      this.role = this.baseInfo.allParamsInfo.roleId!
      this.disabledEditUser = this.baseInfo.allParamsInfo?.operate_disable === 'true'
      this.listType = !this.role

      this.getUsersList()
    }
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
    this.subscription2.unsubscribe()
  }

  getUsersList () {
    this.usersSet = new Set()
    this.apiService
      .get('user/list', { role: this.role || '', keyword: this.keyword || '' })
      .subscribe((resp: any) => {
        if (resp.code === 0) {
          this.usersList = resp.data.users
          if (this.rolesMap.size === 0 || this.flash) {
            this.getRolesList(this.usersList)
            this.flash = false
          } else {
            for (const index in this.usersList) {
              this.usersList[index].status_boolean =
                this.usersList[index].status === 2
              this.usersList[index].role =
                this.usersList[index].role_ids.length === 0
                  ? '未分配'
                  : this.rolesMap.get(this.usersList[index].role_ids[0])
            }
          }
        }
      })
  }

  // 获取角色id与title对应值, 传入list时,需要为该list的角色id与角色名匹配
  getRolesList = (list: Array<any>) => {
    this.apiService.get('role/options').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.rolesMap = new Map()
        this.rolesList = resp.data.roles.filter((item: any) => {
          this.rolesMap.set(item.id, item.title)
          return item.title !== '超级管理员'
        })
        for (const index in this.rolesList) {
          this.rolesList[index].label = this.rolesList[index].title
          this.rolesList[index].value = this.rolesList[index].id
        }
        this.rolesMap.set('', '未分配')
        if (this.role) {
          this.disabledEditUser = this.rolesMap.get(this.role) === '超级管理员'
        }
        if (list) {
          for (const index in list) {
            list[index].status_boolean = list[index].status === 2
            list[index].role =
              list[index].role_ids.length === 0
                ? '未分配'
                : this.rolesMap.get(list[index].role_ids[0])
          }
        }
      }
    })
  }

  usersTableClick = (item:any) => {
    if (!item.data.operate_disable) {
      this.openModal('editUser', item.data)
    }
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  changeUsersSet (item: any, type?:string) {
    if (type === 'all') {
      if (item) {
        for (const index in this.usersList) {
          this.usersSet.add(this.usersList[index].id)
        }
      } else {
        this.usersSet = new Set()
      }
    } else {
      if (item?.checked) {
        this.usersSet.delete(item.id)
      } else {
        this.usersSet.add(item.id)
      }
    }
  }

  changeUserStatus (e:Event, value: any) {
    e.stopPropagation()
    this.apiService
      .patch('user/profile', { status: value.status_boolean ? 1 : 2 }, { id: value.id || '' })
      .subscribe((resp: any) => {
        if (resp.code === 0) {
          value.status = value.status_boolean ? 1 : 2
          value.status_boolean = !value.status_boolean
          this.message.success(resp.msg || '修改用户状态成功!', { nzDuration: 1000 })
        } else {
          this.message.error(resp.msg || '修改用户状态失败!')
        }
      })
  }

  changeUserRole (value: any) {
    this.apiService
      .patch('user/profile', { role: value.role_ids }, { id: value.id || '' })
      .subscribe((resp: any) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '修改用户角色成功!', { nzDuration: 1000 })
          this.dataService.reqFlashGroup()
        } else {
          this.message.error(resp.msg || '修改用户角色失败!')
        }
        this.getUsersList()
      })
  }

  // 重置密码弹窗, 前端随机密码,用户确认重置后请求后端数据
  resetPswModal (item: any, e?:Event) {
    e?.stopPropagation()
    this.modalRef = this.modalService.create({
      nzTitle: '重置密码',
      nzContent: RefreshPswComponent,
      nzComponentParams: {
        userId: item.id,
        closeModal: this.closeDrawer
      },
      nzClosable: true,
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: (component:RefreshPswComponent) => {
        component.refreshPsw()
        return false
      }
    })
  }

  openModal (usage: string, data?: any, e?:Event) {
    e?.stopPropagation()
    if (usage !== 'addUserToRole') {
      this.modalRef = this.modalService.create({
        nzTitle: usage === 'addUser' ? '创建用户' : '编辑用户',
        nzWidth: MODAL_SMALL_SIZE,
        nzContent: UserProfileComponent,
        nzComponentParams: {
          type: usage,
          userId: data?.id,
          accessLink: 'system/role',
          nzDisabled: this.nzDisabled,
          closeModal: this.closeDrawer
        },
        nzOnOk: (component:UserProfileComponent) => {
          component.saveUserProfile()
          return false
        }
      })
    } else {
      this.modalRef = this.modalService.create({
        nzTitle: '添加用户',
        nzWidth: MODAL_NORMAL_SIZE,
        nzContent: ChangeUserRoleComponent,
        nzComponentParams: {
          role: this.role,
          rolesMap: this.rolesMap,
          getRolesList: this.getRolesList,
          closeModal: this.closeDrawer
        },
        nzFooter: [{
          label: '取消',
          onClick: () => {
            this.modalRef?.close()
          }
        },
        {
          label: '保存',
          type: 'primary',
          disabled: (component:ChangeUserRoleComponent|undefined) => {
            return this.nzDisabled || component?.addUsersSet.size === 0
          },
          onClick: (component:ChangeUserRoleComponent|undefined) => {
            component?.changeUsersRole()
          }
        }
        ]
      })
    }
  }

  deleteUserModal (value?: any, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzWidth: MODAL_SMALL_SIZE,
      nzClassName: 'delete-modal',
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteUsers(value?.id)
      }
    })
  }

  deleteUsers (ids?: any) {
    let idsList = []
    if (ids) {
      idsList = [ids]
    } else {
      idsList = Array.from(this.usersSet)
    }
    this.apiService
      .post('user/delete', { ids: idsList })
      .subscribe((resp: any) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '删除用户成功！', { nzDuration: 1000 })
          this.getUsersList()
          this.dataService.reqFlashGroup()
        } else {
          this.message.error(resp.msg || '删除用户失败！')
        }
      })
  }

  removeUserRoleModal (value?: any, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '移除',
      nzContent: '请确认是否移除指定用户的角色权限？',
      nzClosable: true,
      nzWidth: MODAL_SMALL_SIZE,
      nzClassName: 'delete-modal',
      nzOkDanger: true,
      nzOnOk: () => {
        this.removeUserRole(value?.id)
      }
    })
  }

  removeUserRole (ids?: Array<number>) {
    let idsList = []
    if (ids) {
      idsList = [ids]
    } else {
      idsList = Array.from(this.usersSet)
    }
    this.apiService
      .post('role/batch-delete', { ids: idsList, role_id: this.role })
      .subscribe((resp: any) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '移除用户角色成功', { nzDuration: 1000 })
          this.getUsersList()
          this.dataService.reqFlashGroup()
        } else {
          this.message.error(resp.msg || '移除用户角色失败！')
        }
      })
  }

  closeDrawer = (value: any) => {
    if (value) {
      this.dataService.reqFlashGroup()
      this.getUsersList()
    }
    this.modalRef?.close()
  }
}
