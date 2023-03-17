/* eslint-disable camelcase */
/* eslint-disable no-useless-constructor */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'

@Component({
  selector: 'eo-ng-upstream-list',
  templateUrl: './list.component.html',
  styles: [
  ]
})
export class UpstreamListComponent implements OnInit {
  @ViewChild('deleteBtn', { read: TemplateRef, static: true }) deleteBtn: TemplateRef<any> | undefined
  @ViewChild('configTd', { read: TemplateRef, static: true }) configTd: TemplateRef<any> | undefined
  constructor (private message: EoNgFeedbackMessageService,
                private modalService:EoNgFeedbackModalService,
                private api:ApiService,
                private router:Router,
                private appConfigService:AppConfigService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '上游管理', routerLink: 'upstream/upstream' }])
  }

 nzDisabled:boolean = false
  upstreamName:string = ''
  upstreamNameForSear:string = ''
  upstreamsList : Array<any> = []
  upstreamsTableHeadName: Array<object> = [
    { title: '上游名称' },
    { title: '协议类型' },
    { title: '服务类型' },
    { title: '地址' },
    { title: '更新时间' },
    {
      title: '操作',
      right: true
    }
  ]

  upstreamsTableBody: Array<any> =[
    { key: 'name' },
    { key: 'scheme' },
    { key: 'service_type' },
    { key: 'config' },
    { key: 'update_time' },
    {
      type: 'btn',
      right: true,
      btns: [
        {
          title: '查看',
          click: (item:any) => {
            this.router.navigate(['/', 'upstream', 'upstream', 'content', item.data.name])
          }
        },
        {
          title: '删除',
          disabledFn: (data:any, item:any) => {
            return this.nzDisabled || !item.data.is_delete
          },
          click: (item:any) => {
            this.delete(item.data)
          }
        }
      ]
    }
  ]

  ngOnInit (): void {
    this.getUpstreamsList()
  }

  page_size:number = 20
  page_num :number= 1
  total:number = 0

  // 刷新列表，除了初始化组件以外，每次调用该方法时都需要出现刷新成功的消息提示，所以有个可选参数initFlag, 为true时不提示
  getUpstreamsList () {
    this.api.get('services', { name: this.upstreamNameForSear, page_size: this.page_size, page_num: this.page_num }).subscribe(resp => {
      if (resp.code === 0) {
        this.upstreamsList = resp.data.services
        this.total = resp.data.total
        this.upstreamName = this.upstreamNameForSear
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  addUpstream () {
    this.router.navigate(['/', 'upstream', 'upstream', 'create'])
  }

  delete (item:any, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzOkDanger: true,
      nzWidth: MODAL_SMALL_SIZE,
      nzOnOk: () => {
        this.deleteService(item)
      }
    })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  deleteService (item:any) {
    this.api.delete('service', { name: item.name }).subscribe(resp => {
      if (resp.code === 0) {
        this.getUpstreamsList()
        this.message.success(resp.msg || '删除成功!', { nzDuration: 1000 })
      } else {
        this.message.error(resp.msg || '删除失败!')
      }
    })
  }

  upstreamTableClick = (item:any) => {
    this.router.navigate(['/', 'upstream', 'upstream', 'content', item.data.name])
  }
}
