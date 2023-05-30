/* eslint-disable camelcase */
/* eslint-disable no-useless-constructor */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { THEAD_TYPE } from 'eo-ng-table'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { UpstreamListTableBody, UpstreamListTableHeadName } from '../types/conf'

@Component({
  selector: 'eo-ng-upstream-list',
  templateUrl: './list.component.html',
  styles: [
  ]
})
export class UpstreamListComponent implements OnInit {
  constructor (private message: EoNgFeedbackMessageService,
                private modalService:EoNgFeedbackModalService,
                private api:ApiService,
                private router:Router,
                private navigationService:EoNgNavigationService) {
    this.navigationService.reqFlashBreadcrumb([{ title: '上游管理', routerLink: 'upstream/upstream' }])
  }

 nzDisabled:boolean = false
  upstreamName:string = ''
  upstreamNameForSear:string = ''
  upstreamsList : Array<any> = []
  upstreamsTableHeadName: THEAD_TYPE[] = [...UpstreamListTableHeadName]
  upstreamsTableBody: EO_TBODY_TYPE[] =[...UpstreamListTableBody]

  ngOnInit (): void {
    this.getUpstreamsList()
    this.upstreamsTableBody[5].btns[0].click = (item:any) => {
      this.router.navigate(['/', 'upstream', 'upstream', 'content', item.data.name])
    }

    this.upstreamsTableBody[5].btns[1].disabledFn = (data:any, item:any) => {
      return this.nzDisabled || !item.data.isDelete
    }

    this.upstreamsTableBody[5].btns[1].click = (item:any) => {
      this.delete(item.data)
    }
  }

  pageSize:number = 20
  pageNum :number= 1
  total:number = 0

  // 刷新列表，除了初始化组件以外，每次调用该方法时都需要出现刷新成功的消息提示，所以有个可选参数initFlag, 为true时不提示
  getUpstreamsList () {
    this.api.get('services', { name: this.upstreamNameForSear, pageSize: this.pageSize, pageNum: this.pageNum }).subscribe(resp => {
      if (resp.code === 0) {
        this.upstreamsList = resp.data.services
        this.total = resp.data.total
        this.upstreamName = this.upstreamNameForSear
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
      }
    })
  }

  upstreamTableClick = (item:any) => {
    this.router.navigate(['/', 'upstream', 'upstream', 'content', item.data.name])
  }
}
