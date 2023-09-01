import { Injectable } from '@angular/core'
import { DynamicComponentComponent } from '../dynamic-component/dynamic-component.component'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { EditableEnvTableComponent } from './editable-env-table.component'
import { property } from 'lodash-es'
import { MODAL_NORMAL_SIZE } from '../../constant/app.config'
import { resolve } from 'path'

@Injectable({
  providedIn: 'root'
})
export class EditableEnvTableService {
  modalRef:NzModalRef|undefined
  constructor (private modalService:EoNgFeedbackModalService) { }

  openModal (component?:DynamicComponentComponent) {
    return new Promise((resolve) => {
      this.modalRef = this.modalService.create({
        nzTitle: '添加环境变量',
        nzContent: EditableEnvTableComponent,
        nzClosable: true,
        nzWidth: MODAL_NORMAL_SIZE,
        nzFooter: null,
        nzComponentParams: {
          ...component?.propertyWaitForChoose,
          chooseEnv: ($event:Event) => {
            component?.chooseEnv($event)
            resolve($event)
            this.modalRef?.close()
          }
        },
        nzOnCancel: () => {
          resolve({ closeModal: true })
        }

      })
      if (component) {
        component.subscription = this.modalRef.afterClose.subscribe(() => {
          component.envNameForSear = ''
          component.propertyWaitForChoose = null
        })
      }
    })
  }
}
