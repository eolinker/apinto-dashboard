/* eslint-disable dot-notation */
/* eslint-disable no-unused-expressions */
/* eslint-disable camelcase */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { BaseInfoService } from '../../../service/base-info.service'
import { ApplicationAuthenticationFormComponent } from './form/form.component'
import { differenceInCalendarDays } from 'date-fns'

@Component({
  selector: 'eo-ng-application-authentication',
  templateUrl: './authentication.component.html',
  styles: [
    `

    .authDrawer .label-large{
      width:126px;
      height: 32px;
      line-height: 32px;
      padding: 0;
      text-align:right;
    }

    .small-row{
      width:210px;
      margin-right:20px;
    }

    .dynamic-control{
      margin-bottom:0 !important;
    }
    `
  ]
})
export class ApplicationAuthenticationComponent implements OnInit {
  @ViewChild('authContentTpl', { read: TemplateRef, static: true }) authContentTpl: TemplateRef<any> | undefined
  appId:string = ''
  nzDisabled:boolean = false

  authenticationTableHeadName:Array<object> = [

    { title: '鉴权类型' },
    { title: '参数位置' },
    { title: '参数名' },
    { title: '参数信息' },
    { title: '过期时间' },
    { title: '隐藏鉴权信息' },
    { title: '更新者' },
    { title: '更新时间' },
    {
      title: '操作',
      right: true
    }
  ]

  authenticationTableBody:Array<any> = [
    { key: 'driver' },
    { key: 'param_position' },
    { key: 'param_name' },
    { key: 'param_info' },
    {
      key: 'expire_time_string',
      styleFn: (item:any) => {
        if (item.expire_time && differenceInCalendarDays(item.expire_time * 1000, new Date()) < 0) {
          return 'color:red'
        }
        return ''
      }
    },
    {
      key: 'is_transparent'
    },
    { key: 'operator' },
    { key: 'update_time' },
    {
      type: 'btn',
      right: true,
      btns: [{
        title: '查看',
        click: (item:any) => {
          this.openDrawer(item.data.uuid)
        }
      },
      {
        title: '删除',
        disabledFn: () => {
          return this.nzDisabled
        },
        click: (item:any) => {
          this.delete(item.data)
        }
      }
      ]
    }
  ]

  authenticationList:Array<{uuid:string, info:string, driver:string, expire_time_string:string, expire_time:number, operator:string, update_time:string, rule_info:string, is_transparent:boolean}>=[
  ]

  modalRef:NzModalRef | undefined

  constructor (
               private message: EoNgFeedbackMessageService,
               public api:ApiService,
               private baseInfo:BaseInfoService,
               private modalService:EoNgFeedbackModalService,
               private router:Router,
               private appConfigService:AppConfigService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '应用管理', routerLink: 'application' }, { title: '鉴权管理' }])
  }

  ngOnInit (): void {
    this.appId = this.baseInfo.allParamsInfo.appId
    if (!this.appId) {
      this.router.navigate(['/', 'application'])
    }
    this.getAuthsData()
  }

  authTableClick = (item:any) => {
    this.openDrawer(item.data.uuid)
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  openDrawer (authId?:string) {
    this.modalRef = this.modalService.create({
      nzTitle: '配置鉴权信息',
      nzWidth: MODAL_SMALL_SIZE,
      nzContent: ApplicationAuthenticationFormComponent,
      nzComponentParams: { authId: authId, appId: this.appId, closeModal: this.closeModal },
      nzOkDisabled: this.nzDisabled,
      nzOnOk: (component:ApplicationAuthenticationFormComponent) => {
        component.saveAuth()
        return false
      }
    })
  }

  delete (item:any, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.api.delete('application/auth', { uuid: item.uuid }).subscribe(resp => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '删除成功!', { nzDuration: 1000 })
            this.getAuthsData()
          } else {
            this.message.error(resp.msg || '删除失败!')
          }
        })
      }
    })
  }

  closeModal= (freshList?:boolean) => {
    if (freshList) {
      this.getAuthsData()
    }
    this.modalRef?.close()
  }

  // 获取鉴权列表
  getAuthsData () {
    this.api.get('application/auths', { app_id: this.appId }).subscribe(resp => {
      if (resp.code === 0) {
        for (const index in resp.data.auths) {
          resp.data.auths[index].driver = this.getAuthDriver(resp.data.auths[index].driver)
          resp.data.auths[index].is_transparent = resp.data.auths[index].is_transparent ? '是' : '否'
          resp.data.auths[index].expire_time_string = resp.data.auths[index].expire_time === 0 ? '永不过期' : this.getDateInList(resp.data.auths[index].expire_time)
        }
        this.authenticationList = resp.data.auths
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  // 将鉴权列表中的driver转换大小写
  getAuthDriver (driver:string):string {
    switch (driver) {
      case 'basic':
        return 'Basic'
      case 'apikey':
        return 'ApiKey'
      case 'aksk':
        return 'AkSk'
      case 'jwt':
        return 'Jwt'
      default:
        return driver
    }
  }

  // 将后端传的时间戳转成日期
  getDateInList (time:number):string {
    try {
      if (Number(time) || Number(time) === 0) {
        const date = new Date(Number(time) * 1000)
        const month = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1
        const day = date.getDate() < 10 ? '0' + date.getDate() : date.getDate()
        const hour = date.getHours() < 10 ? '0' + date.getHours() : date.getHours()
        const min = date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()
        const sec = date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds()
        return `${date.getFullYear()}-${month}-${day} ${hour}:${min}:${sec}`
      } else {
        return '日期数据格式有误'
      }
    } catch (error:any) {
      return '日期数据格式有误'
    }
  }
}
