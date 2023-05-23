import { Component } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { ApiService } from '../../../service/api.service'
import { EoNgNavigationService } from '../../../service/app-config.service'
import { BaseInfoService } from '../../../service/base-info.service'
import { AuthListData, ExtraListData } from '../types/types'
import { extraTableHeadName, extraTableBody } from '../types/conf'
import { ApplicationExtraFormComponent } from './form/form.component'
import { EmptyHttpResponse } from '../../../constant/type'
import { EoNgApplicationService } from '../application.service'

@Component({
  selector: 'eo-ng-application-extra',
  templateUrl: './extra.component.html',
  styles: [
  ]
})
export class ApplicationExtraComponent {
  appId:string = ''
  nzDisabled:boolean = false
  extraTableHeadName:THEAD_TYPE[] = [...extraTableHeadName]
  extraTableBody:TBODY_TYPE[] = [...extraTableBody]
  modalRef:NzModalRef | undefined

  constructor (
               private message: EoNgFeedbackMessageService,
               public api:ApiService,
               private baseInfo:BaseInfoService,
               private modalService:EoNgFeedbackModalService,
               private router:Router,
               private navigationService:EoNgNavigationService,
               public service:EoNgApplicationService) {
    this.navigationService.reqFlashBreadcrumb([{ title: '应用管理', routerLink: 'application' }, { title: '额外参数' }])
  }

  ngOnInit (): void {
    this.appId = this.baseInfo.allParamsInfo.appId
    this.initTable()
    if (!this.appId) {
      this.router.navigate(['/', 'application'])
    }
  }

  initTable () {
    this.extraTableBody[4].btns[0].click = (item:any) => { this.openDrawer(item.data) }
    this.extraTableBody[4].btns[1].disabledFn = () => { return this.nzDisabled }
    this.extraTableBody[4].btns[1].click = (item:any) => { this.delete(item.data) }
  }

  authTableClick = (item:{data:ExtraListData}) => {
    this.openDrawer(item.data)
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  openDrawer (data?:ExtraListData) {
    this.modalRef = this.modalService.create({
      nzTitle: `${data ? '编辑' : '添加'}参数`,
      nzWidth: MODAL_SMALL_SIZE,
      nzContent: ApplicationExtraFormComponent,
      nzComponentParams: {
        extraList: this.service.appData?.params,
        data: data,
        nzDisabled: this.nzDisabled
      },
      nzOkDisabled: this.nzDisabled,
      nzOnOk: (component: ApplicationExtraFormComponent) => {
        (component as ApplicationExtraFormComponent).saveParam()
        return new Promise((resolve, reject) => {
          this.api.put('application', { ...this.service.appData, params: component.extraList })
            .subscribe((resp:EmptyHttpResponse) => {
              if (resp.code === 0) {
                resolve()
                this.message.success(resp.msg || '操作成功!')
                this.service.getApplicationData(this.appId)
              } else {
                reject(new Error())
              }
            })
        })
      }
    })
  }

  delete (item:AuthListData, e?:Event) {
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
            this.message.success(resp.msg || '删除成功!')
            this.service.getApplicationData(this.appId)
          }
        })
      }
    })
  }
}
