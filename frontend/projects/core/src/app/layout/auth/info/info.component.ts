/* eslint-disable no-useless-constructor */
import { Component, Input, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from '../../../service/api.service'
import { NzUploadChangeParam, NzUploadFile } from 'ng-zorro-antd/upload'
import { EoNgNavigationService } from '../../../service/eo-ng-navigation.service'

@Component({
  selector: 'eo-ng-auth-info',
  templateUrl: './info.component.html',
  styles: [
  ]
})
export class AuthInfoComponent implements OnInit {
  @Input() updateAuth:boolean = false
  macCode:string = ''
  canActive:boolean = false
  certCanActive:boolean = false
  authFile:NzUploadFile|undefined
  showActivationInfo:boolean = true
  fileList: NzUploadFile[] = [];
  free:boolean = true
  authInfo:{title:string, infos:Array<{key:string, value:string}>}
  = {
    title: '授权管理',
    infos: []
  }

  constructor (
    private message: EoNgFeedbackMessageService,
    private api:ApiService,
    private router:Router,
    private navigationService:EoNgNavigationService
  ) { }

  // 首次激活或用户通过输入url到达该组件，则updateAuth为false，需要检查是否激活过，激活则直接进入登录页
  // 更新授权则会传入macCode，无需获取
  ngOnInit (): void {
    this.getInfo()
    this.getAuthMaxCode()
  }

  getInfo () {
    this.api.authGet('activation/info')
      .subscribe((resp:{code:number, data:{infos:Array<{key:string, value:string}>, title:string}, msg:string}) => {
        if (resp.code === 0) {
          this.authInfo = resp.data
          this.navigationService.reqFlashBreadcrumb([{ title: '授权管理' }])
        }
      })
  }

  getAuthMaxCode () {
    this.api.authGet('mac')
      .subscribe((resp:{code:number, data:{info:{mac:string}}, msg:string}) => {
        if (resp.code === 0) {
          this.macCode = resp.data.info.mac
        }
      })
  }

  activeOrUpdate () {
    let url = 'activation' // 首次激活
    if (this.updateAuth) {
      url = 'reactivation'
    }
    const formData = new FormData()
    formData.append('authFile', this.authFile as any)
    this.api.authPostWithFile(url, formData)
      .subscribe((resp:{code:number, data:{infos:Array<{key:string, value:string}>, title:string}, msg:string}) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '激活成功！', { nzDuration: 1000 })
          if (!this.updateAuth) {
            this.showActivationInfo = true
            this.authInfo = resp.data
          } else {
            this.router.navigate(['/', 'auth-info'])
          }
        }
      })
  }

  goToLogin () {
    this.router.navigate(['/', 'login'])
  }

  copyCallback = () => {
    this.message.success('复制成功', { nzDuration: 1000 })
    if (!this.certCanActive) { this.certCanActive = true }
  };

  onChange (e: NzUploadChangeParam): void {
    if (e.file.status === 'done') {
      if (e.file.response.code === 0) {
        this.message.success(e.file.response.msg || '操作成功')
        this.showActivationInfo = true
        this.authInfo.infos = e.file.response.data.infos
      } else {
        this.message.error(e.file.response.msg || '操作失败')
      }
    } else if (e.file.status === 'error') {
      this.message.error(`${e.file.name}上传失败，请重试`)
    }
  }
}
