/* eslint-disable camelcase */
/* eslint-disable no-useless-constructor */
import { Component, Input, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { NzUploadFile } from 'ng-zorro-antd/upload'
import { ApiService } from '../../../service/api.service'

@Component({
  selector: 'eo-ng-auth-activation',
  templateUrl: './activation.component.html',
  styles: [
    `
    .activation-extra{
      margin-top:1px;
      background-color:var(--BAR_BG_COLOR);
      border:1px solid var(--bar-bg-color);
      box-shadow:none;
      text-shadow:none;
    }
    `
  ]
})
export class AuthActivationComponent implements OnInit {
  @Input() updateAuth:boolean = false
  macCode:string = ''
  canActive:boolean = false
  certCanActive:boolean = false
  authFile:NzUploadFile|undefined
  showActivationInfo:boolean = false
  authInfo:{title:string, infos:Array<{key:string, value:string}>}
  = { title: '标准版授权', infos: [] }

  fileList: NzUploadFile[] = [];

  constructor (
    private message: EoNgFeedbackMessageService,
    private api:ApiService,
    private router:Router
  ) { }

  // 首次激活或用户通过输入url到达该组件，则updateAuth为false，需要检查是否激活过，激活则直接进入登录页
  // 更新授权则会传入macCode，无需获取
  ngOnInit (): void {
    this.getAuthMaxCode()
  }

  getAuthMaxCode () {
    this.api.authGet('mac').subscribe((resp:{code:number, data:{info:{mac:string}}, msg:string}) => {
      if (resp.code === 0) {
        this.macCode = resp.data.info.mac
      } else {
        this.message.error(resp.msg || '获取机器码失败，请重试！')
      }
    })
  }

  beforeUpload = (file: NzUploadFile): boolean => {
    this.fileList = []
    this.fileList = this.fileList.concat(file)
    this.authFile = file
    this.canActive = !!this.authFile
    return false
  }

  removeFile () {
    this.fileList = []
    this.authFile = undefined
    this.canActive = !!this.authFile
    return true
  }

  activeOrUpdate () {
    let url = 'activation' // 首次激活
    if (this.updateAuth) {
      url = 'reactivation'
    }
    const formData = new FormData()
    formData.append('auth_file', this.authFile as any)
    this.api.authPostWithFile(url, formData).subscribe((resp:{code:number, data:{infos:Array<{key:string, value:string}>, title:string}, msg:string}) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '激活成功！', { nzDuration: 1000 })
        if (!this.updateAuth) {
          this.showActivationInfo = true
          this.authInfo = resp.data
        } else {
          this.router.navigate(['/', 'auth-info'])
        }
      } else {
        this.message.error(resp.msg || '激活失败，请重试！')
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
}
