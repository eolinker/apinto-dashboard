/*
 * @Date: 2023-12-14 17:14:28
 * @LastEditors: maggieyyy
 * @LastEditTime: 2023-12-18 15:54:01
 * @FilePath: \apinto\projects\apinto-auth\src\app\layout\auth-info\auth-info.component.ts
 */
import { Component, Input, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { NzUploadFile, NzUploadChangeParam } from 'ng-zorro-antd/upload'
import { ApiService } from '../../service/api.service'
import { ModuleFederationService } from '../../service/module-federation.service'

@Component({
  selector: 'auth-info',
  templateUrl: './auth-info.component.html',
  styles: [
    `
    :host ::ng-deep{
        ol {
          list-style: none;
          padding-inline-start: 0;
        }

        .ant-alert-with-description .ant-alert-icon{
          height:24px;
        }

        .ant-alert-with-description .ant-alert-message{
          font-weight:bold;
          color:unset;
          font-size:14px;
          margin-bottom:0px;
        }
    }
    `]
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
  authInfo:{title:string, infos:Array<{key:string, value:string}>, status:'normal'|'waring'|'freeze', refer:string}
  = {
    title: '授权管理',
    infos: [],
    status: 'normal',
    refer: '节点数超出授权，请尽快更新授权或减少节点到授权数量内'
  }

  version:string = this.mfe.providerFromCore?.dashboardVersion?.version
  updateDate:string = this.mfe.providerFromCore?.dashboardVersion?.updateDate
  productName:string = this.mfe.providerFromCore?.dashboardVersion?.productName

  constructor (
    private message: EoNgFeedbackMessageService,
    private api:ApiService,
    private router:Router,
    private mfe:ModuleFederationService
    // private navigationService:EoNgNavigationService
  ) { }

  // 首次激活或用户通过输入url到达该组件，则updateAuth为false，需要检查是否激活过，激活则直接进入登录页
  // 更新授权则会传入macCode，无需获取
  ngOnInit (): void {
    this.getInfo()
    this.getAuthMaxCode()
    this.mfe.providerFromCore.breadcrumb = ([{ title: '授权管理' }])
  }

  getInfo () {
    this.api.authGet('activation/info')
      .subscribe((resp:{code:number, data:{infos:Array<{key:string, value:string}>, title:string, status:'normal'|'waring'|'freeze', refer:string}, msg:string}) => {
        if (resp.code === 0) {
          this.authInfo = resp.data
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
      .subscribe((resp:{code:number, data:{infos:Array<{key:string, value:string}>, title:string, status:'normal'|'waring'|'freeze', refer:string}, msg:string}) => {
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
        this.authInfo = e.file.response.data
      } else {
        this.message.error(e.file.response.msg || '操作失败')
      }
    } else if (e.file.status === 'error') {
      this.message.error(`${e.file.name}上传失败，请重试`)
    }
  }
}
