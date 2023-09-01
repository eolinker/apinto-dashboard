import { Component, OnInit } from '@angular/core'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgMessageService } from 'projects/core/src/app/service/eo-ng-message.service'
import { WebhookListData } from '../../types/type'
import { SystemWebhookConfigComponent } from '../config/config.component'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { webhooksTableHead, webhooksTableBody } from '../../types/conf'

@Component({
  selector: 'eo-ng-system-webhook-list',
  templateUrl: './list.component.html',
  styles: [
  ]
})
export class SystemWebhookListComponent implements OnInit {
  webhooksList:Array<WebhookListData> = []
  nzDisabled:boolean = false
  webhooksTableHead: THEAD_TYPE[] = [...webhooksTableHead]
  webhooksTableBody: TBODY_TYPE[] = [...webhooksTableBody]
  webhookId:string = ''
  modalRef:NzModalRef |undefined

  constructor (private message: EoNgMessageService,
    private modalService:EoNgFeedbackModalService,
     private api:ApiService,
     private navigationService:EoNgNavigationService) {
    this.navigationService.reqFlashBreadcrumb([{ title: 'Webhook管理', routerLink: 'system/webhook' }])
  }

  ngOnInit (): void {
    this.webhooksTableHead.push({ title: '操作', right: true })
    this.webhooksTableBody.push({
      type: 'btn',
      right: true,
      btns: [
        {
          title: '查看',
          click: (item:any) => {
            this.openWebhookModal(item.data.uuid)
          }
        },
        {
          title: '删除',
          disabledFn: (data:any, item:any) => {
            return !item.data.isDelete || this.nzDisabled
          },
          click: (item:any) => {
            this.delete(item.data)
          }
        }
      ]
    })

    this.getWebhooksList()
  }

  disabledEdit (editRight:boolean) {
    this.nzDisabled = editRight
  }

  webhookTableClick = (webhookItem:any) => {
    this.openWebhookModal(webhookItem.data.uuid)
  }

  getWebhooksList () {
    this.api.get('webhooks').subscribe((resp:{code:number, data:{webhooks:Array<WebhookListData>}, msg:string}) => {
      if (resp.code === 0) {
        this.webhooksList = resp.data.webhooks
      }
    })
  }

  openWebhookModal (id:string = '') {
    this.webhookId = id || ''
    this.modalRef = this.modalService.create({
      nzTitle: this.webhookId ? '编辑Webhook' : '新建Webhook',
      nzContent: SystemWebhookConfigComponent,
      nzComponentParams: {
        webhookId: this.webhookId,
        closeModal: this.closeModal
      },
      nzClosable: true,
      nzWidth: MODAL_NORMAL_SIZE,
      nzOkText: this.webhookId ? '保存' : '提交',
      nzOnOk: (component:SystemWebhookConfigComponent) => {
        component.saveWebhook()
        return false
      }
    })
  }

  closeModal = () => {
    this.getWebhooksList()
    this.modalRef?.close()
  }

  delete (webhookItem:any, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: `${webhookItem.title}一旦删除，数据将会丢失。`,
      nzClosable: true,
      nzWidth: MODAL_SMALL_SIZE,
      nzClassName: 'delete-modal',
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteDiscovery(webhookItem.uuid)
      }
    })
  }

  deleteDiscovery (id:string) {
    this.api.delete('webhook', { uuid: id }).subscribe(resp => {
      if (resp.code === 0) {
        this.getWebhooksList()
        this.message.success(resp.msg || '删除成功')
      }
    })
  }
}
