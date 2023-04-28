/* eslint-disable camelcase */
/* eslint-disable no-useless-constructor */
import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { ExternalAppListTableBody, ExternalAppListTableHeadName } from '../../types/conf'

@Component({
  selector: 'eo-ng-external-app-list',
  templateUrl: './list.component.html',
  styles: [
  ]
})
export class ExternalAppListComponent implements OnInit {
  @ViewChild('switchTpl', { read: TemplateRef, static: true }) switchTpl: TemplateRef<any> | undefined
  @ViewChild('copyTokenTpl', { read: TemplateRef, static: true }) copyTokenTpl: TemplateRef<any> | undefined
  nzDisabled:boolean = false
  appsTableHeadName:THEAD_TYPE[] = [...ExternalAppListTableHeadName]
  appsTableBody:TBODY_TYPE[] = [...ExternalAppListTableBody]
  appsList:Array<any> = []

  constructor (private router:Router, private message: EoNgFeedbackMessageService, private modalService:EoNgFeedbackModalService, private api:ApiService,
    private navigationService: EoNgNavigationService) {
    this.navigationService.reqFlashBreadcrumb([
      { title: '外部应用' }
    ])
  }

  ngOnInit (): void {
    this.getAppsData()
  }

  ngAfterViewInit ():void {
    this.appsTableBody[4].title = this.switchTpl
    this.appsTableBody[7].btns[0].disabledFn = () => { return this.nzDisabled }
    this.appsTableBody[7].btns[0].click = (item:any) => {
      this.updateToken(item.data)
    }
    this.appsTableBody[7].btns[1].type = this.copyTokenTpl
    this.appsTableBody[7].btns[2].click = (item:any) => {
      this.getAppMessage(item.data)
    }
    this.appsTableBody[7].btns[3].disabledFn = () => { return this.nzDisabled }
    this.appsTableBody[7].btns[3].click = (item:any) => {
      this.deleteAppModal(item.data)
    }
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  getAppsData ():void {
    this.api.get('external-apps').subscribe((resp:{code:number, data:{apps:Array<{id:string, name:string, token:string, tags:string, status:number, operator:string, updateTime:string}>}, msg:string}) => {
      if (resp.code === 0) {
        this.appsList = resp.data.apps
        for (const index in this.appsList) {
          this.appsList[index].statusBoolean = this.appsList[index].status === 2 // 禁用
        }
      }
    })
  }

  getAppMessage (item:any) {
    this.router.navigate(['/', 'system', 'ext-app', 'message', item.id])
  }

  addApp () {
    this.router.navigate(['/', 'system', 'ext-app', 'create'])
  }

  updateToken (item:any) {
    this.api.put('external-app/token', null, { id: item.id }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '刷新鉴权Token成功！')
        this.getAppsData()
      }
    })
  }

  appsTableClick = (item:any) => {
    this.getAppMessage(item.data)
  }

  disabledApp (item:any, e:Event) {
    e.stopPropagation()
    let url = ''
    if (!item.statusBoolean) {
      url = 'external-app/disable'
    } else {
      url = 'external-app/enable'
    }

    this.api.put(url, null, { id: item.id }).subscribe((resp:{code:number, data:{}, msg:string}) => {
      if (resp.code === 0) {
        item.statusBoolean = !item.statusBoolean
        this.message.success(resp.msg || ((url.split('/')[1] === 'disable' ? '禁用' : '启用') + '成功！'))
        this.getAppsData()
      }
    })
  }

  // 删除api弹窗
  deleteAppModal (item:any) {
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteApp(item)
      }
    })
  }

  // 删除单个api
  deleteApp = (items:any) => {
    this.api.delete('external-app', { id: items.id }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '删除成功!')
        this.getAppsData()
      }
    })
  }

  callback = () => {
    this.message.success('复制成功')
  };
}
