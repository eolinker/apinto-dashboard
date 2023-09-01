import { Injectable } from '@angular/core'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { ListComponent } from './list/list.component'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE } from '../../constant/app.config'
import { ServiceGovernancePublishComponent } from './publish/publish.component'

@Injectable({
  providedIn: 'root'
})
export class ServiceGovernanceService {
  modalRef:NzModalRef|undefined

  constructor (private modalService:EoNgFeedbackModalService) { }

  publishStrategyModal (strategyType:string, clusterName:string, component?:ListComponent, returnToSdk?:Function) {
    this.modalRef = this.modalService.create({
      nzTitle: '发布策略',
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: ServiceGovernancePublishComponent,
      nzComponentParams: {
        strategyType: strategyType,
        clusterName: clusterName,
        closeModal: () => {
          this.modalRef?.close()
          component?.getStrategiesList()
        },
        returnToSdk
      },
      nzOkDisabled: !!component?.nzDisabled,
      nzOnCancel: () => {
        returnToSdk && returnToSdk({ data: { closeModal: true } })
      },
      nzOnOk: (component:ServiceGovernancePublishComponent) => {
        component.publish()
        return false
      }
    })
  }
}
