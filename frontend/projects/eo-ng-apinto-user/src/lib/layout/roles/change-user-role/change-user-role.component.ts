import { Component, Inject, Input, OnInit } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { API_SERVICE_ADAPTER, ApiServiceAdapter } from '../../../constant/api-service-adapter'
import { APP_SERVICE_ADAPTER, AppServiceAdapter } from '../../../constant/app-service-adapter'
import { DataService } from '../../../service/data.service'

@Component({
  selector: 'eo-ng-change-user-role',
  templateUrl: './change-user-role.component.html',
  styles: ['']
})
export class ChangeUserRoleComponent implements OnInit {
  @Input() getRolesList?:(list?:any)=>void
  @Input() closeModal?:(value?:any)=>void
  addUserKeyword: string = '' // 添加用户列表的搜索关键词
  addUsersList: Array<any> = [] // 添加用户列表, 为当前角色以外的用户

  addUsersTableHeadName: Array<any> = [
    {
      type: 'checkbox',
      click: (item: any) => {
        this.changeAddUsersSet(item, 'all')
      },
      showFn: () => {
        return !this.nzDisabled
      }
    },
    {
      title: '账号'
    },
    {
      title: '名称'
    },
    {
      title: '邮箱'
    },
    {
      title: '角色'
    }
  ]

  addUsersTableBody: Array<any> = [
    {
      key: 'checked',
      type: 'checkbox',
      click: (item: any) => {
        this.changeAddUsersSet(item)
      }
    },
    {
      key: 'user_name'
    },
    {
      key: 'nick_name'
    },
    {
      key: 'email'
    },
    {
      key: 'role'
    }
  ]

  addUsersSet: Set<any> = new Set() // 被选择用来添加的用户
  nzDisabled:boolean = false // 权限控制
  usersList: Array<any> = [] // 用户列表
  role: string = '' // 当前列表的角色id(自定义角色的用户列表中使用)
  rolesMap: Map<string, string> = new Map() // 角色id与角色名的映射表

  constructor (
    private message: EoNgFeedbackMessageService,
    @Inject(API_SERVICE_ADAPTER) private apiService: ApiServiceAdapter,
    @Inject(APP_SERVICE_ADAPTER) public appService: AppServiceAdapter,
    private dataService: DataService
  ) {
  }

  ngOnInit (): void {
    this.getAddUsersList()
  }

  // 获取添加用户时的数据,需筛出与当前列表角色id相同的用户且不为当前用户的用户
  getAddUsersList () {
    this.apiService
      .get('user/list', { keyword: this.addUserKeyword || '' })
      .subscribe((resp: any) => {
        if (resp.code === 0) {
          for (const index in resp.data.users) {
            resp.data.users[index].role = resp.data.users[index].role_ids.length === 0
              ? '未分配'
              : this.rolesMap.get(resp.data.users[index].role_ids[0])
          }
          this.addUsersList = resp.data.users.filter((item: any) => {
            return item.role_ids[0] !== this.role && !item.operate_disable && item.id !== this.appService.getUserId()
          })
          this.getRolesList && this.getRolesList(this.addUsersList)
        }
      })
  }

  addUsersTableClick = (rowItem:any) => {
    this.changeAddUsersSet(rowItem.data)
    rowItem.checked = !rowItem.checked
    rowItem.data.checked = !rowItem.data.checked
  }

  // 向自定义角色添加用户
  changeUsersRole () {
    if (this.addUsersSet.size > 0) {
      this.apiService
        .post('role/batch-update', {
          ids: Array.from(this.addUsersSet),
          role_id: this.role
        })
        .subscribe((resp: any) => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '添加用户成功!', { nzDuration: 1000 })
            this.dataService.reqFlashGroup()
            this.closeModal && this.closeModal(true)
            return true
          } else {
            this.message.error(resp.msg || '添加用户失败!')
            return false
          }
        })
    }
  }

  changeAddUsersSet (item: any, type?:string) {
    if (type === 'all') {
      if (item) {
        for (const index in this.usersList) {
          this.addUsersSet.add(this.usersList[index].id)
        }
      } else {
        this.addUsersSet = new Set()
      }
    } else {
    // 被取消勾选
      if (item?.checked) {
        this.addUsersSet.delete(item.id)
      } else {
      // 被选中
        this.addUsersSet.add(item.id)
      }
    }
  }
}
