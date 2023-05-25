import { Component, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_SMALL_SIZE } from '../../../constant/app.config'
import { ApiService } from '../../../service/api.service'
import { EoNgNavigationService } from '../../../service/app-config.service'
import { BaseInfoService } from '../../../service/base-info.service'
import { extraTableHeadName, extraTableBody } from '../types/conf'
import { ApplicationExtraFormComponent } from './form/form.component'
import { EmptyHttpResponse } from '../../../constant/type'
import { EoNgApplicationService } from '../application.service'
import { ApplicationParamData } from '../types/types'

@Component({
  selector: 'eo-ng-application-extra',
  templateUrl: './extra.component.html',
  styles: [
  ]
})
export class ApplicationExtraComponent {
  @ViewChild('conflictTpl', { read: TemplateRef, static: true }) conflictTpl: TemplateRef<any> | undefined
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
    this.extraTableBody[3].title = this.conflictTpl
    this.extraTableBody[4].btns[0].click = (item:any) => { this.openDrawer(item.data) }
    this.extraTableBody[4].btns[1].disabledFn = () => { return this.nzDisabled }
    this.extraTableBody[4].btns[1].click = (item:any) => { this.delete(item.data) }
  }

  authTableClick = (item:{data:ApplicationParamData}) => {
    this.openDrawer(item.data)
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  openDrawer (data?:ApplicationParamData) {
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
        const param = [...this.service.appData!.params as ApplicationParamData[]]
        return new Promise((resolve, reject) => {
          if (component.validateParamForm.valid) {
            if (data) {
              for (const index in param) {
                if (param[index].eoKey === component.data!.eoKey) {
                  param.splice(Number(index), 1)
                  break
                }
              }
            }
            param.unshift(component.validateParamForm.value as ApplicationParamData)
            this.api.put('application', { ...this.service.appData, params: param })
              .subscribe((resp:EmptyHttpResponse) => {
                if (resp.code === 0) {
                  resolve()
                  this.message.success(resp.msg || '操作成功!')
                  this.service.getApplicationData(this.appId)
                } else {
                  reject(new Error())
                }
              })
          } else {
            Object.values(component.validateParamForm.controls).forEach(control => {
              if (control.invalid) {
                control.markAsDirty()
                control.updateValueAndValidity({ onlySelf: true })
              }
            })
          }
        })
      }
    })
  }

  delete (item:ApplicationParamData, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        const newParams = this.service.appData!.params!.filter((p:ApplicationParamData) => {
          return !(item.key === p.key && item.conflict === p.conflict && item.value === p.value && item.position === p.position)
        })
        this.api.put('application', { ...this.service.appData, params: newParams })
          .subscribe((resp:EmptyHttpResponse) => {
            if (resp.code === 0) {
              this.message.success(resp.msg || '操作成功!')
              this.service.getApplicationData(this.appId)
            }
          })
      }
    })
  }
}
