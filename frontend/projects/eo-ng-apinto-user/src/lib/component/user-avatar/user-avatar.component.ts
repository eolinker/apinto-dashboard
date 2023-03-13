/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-09-21 22:19:44
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-09-21 22:41:25
 * @FilePath: /apinto/projects/eo-ng-apinto-user/src/lib/component/user-avatar/user-avatar.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, Inject, OnInit } from '@angular/core'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { API_SERVICE_ADAPTER, ApiServiceAdapter } from '../../constant/api-service-adapter'
import { Router } from '@angular/router'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from '../../constant/app-config-adapter'
import { UserProfileComponent } from '../user-profile/user-profile.component'
import { ResetPswComponent } from '../reset-psw/reset-psw.component'
import { APP_SERVICE_ADAPTER, AppServiceAdapter } from '../../constant/app-service-adapter'

@Component({
  selector: 'eo-ng-apinto-user-avatar',
  templateUrl: './user-avatar.component.html',
  styleUrls: ['./user-avatar.component.scss']
})
export class UserAvatarComponent implements OnInit {
  userMenu: Array<any> = []
  userNickName: string = ''
  userName: string = ''
  drawerRef:NzModalRef | undefined
  // eslint-disable-next-line no-useless-constructor
  constructor (private message: EoNgFeedbackMessageService,
                private modalService:EoNgFeedbackModalService,
                private router: Router,
                @Inject(API_SERVICE_ADAPTER) private apiService: ApiServiceAdapter,
                @Inject(APP_SERVICE_ADAPTER) private appService: AppServiceAdapter
  ) {
  }

  ngOnInit (): void {
    this.userMenu = [
      {
        title: '用户设置',
        click: this.userSetting
      },
      {
        title: '修改密码',
        click: this.changeUserPsw
      },
      {
        title: '退出登录',
        click: this.logout
      }
    ]
    this.getCurrentUserProfile()
  }

  getCurrentUserProfile () {
    this.apiService.get('my/profile').subscribe((resp:any) => {
      if (resp.code === 0) {
        this.userNickName = resp.data.profile.nick_name
        this.userName = resp.data.profile.user_name
        this.appService.setUserRoleId(resp.data.profile.role_ids[0])
        this.appService.setUserId(resp.data.profile.id)
      } else {
        this.message.error(resp.msg || '获取用户信息失败!')
      }
    })
  }

  userSetting = () => {
    this.openDrawer('editCurrentUser')
  }

  changeUserPsw = () => {
    this.openDrawer('changePsw')
  }

  openDrawer (usage:string) {
    switch (usage) {
      case 'editCurrentUser':
        this.drawerRef = this.modalService.create({
          nzTitle: '用户设置',
          nzWidth: MODAL_SMALL_SIZE,
          nzContent: UserProfileComponent,
          nzComponentParams: {
            type: usage,
            closeModal: this.closeModal
          },
          nzOnOk: (component:UserProfileComponent) => {
            component.saveUserProfile()
            return false
          }
        })
        break
      case 'changePsw':
        this.drawerRef = this.modalService.create({
          nzTitle: '修改密码',
          nzWidth: MODAL_SMALL_SIZE,
          nzContent: ResetPswComponent,
          nzComponentParams: { type: usage, userName: this.userName, closeModal: this.closeModal },
          nzOnOk: (component:ResetPswComponent) => {
            component.resetPsw()
            return false
          }
        })
        break
    }
  }

  logout = () => {
    this.apiService.logout().subscribe((resp:any) => {
      if (resp.code === 0) {
        this.router.navigate(['/', 'login'])
      } else {
        this.message.error(resp.msg || '退出登录失败!')
      }
    })
  }

  closeModal =() => {
    this.drawerRef?.close()
  }
}
