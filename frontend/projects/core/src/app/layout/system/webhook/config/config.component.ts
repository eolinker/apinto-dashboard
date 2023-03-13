/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { UntypedFormBuilder, FormGroup, Validators } from '@angular/forms'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { EoNgMyValidators } from 'projects/core/src/app/constant/eo-ng-validator'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { responseHeaderTableBody } from 'projects/core/src/app/constant/table.conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgMessageService } from 'projects/core/src/app/service/eo-ng-message.service'
import { contentTypesList, methodsList, noticeTypesList } from '../../types/conf'
import { WebhookData } from '../../types/type'

@Component({
  selector: 'eo-ng-system-webhook-config',
  templateUrl: './config.component.html',
  styles: [
  ]
})
export class SystemWebhookConfigComponent implements OnInit {
  @Input() webhookId:string = ''
  @Input() closeModal?:(value?:any)=>void

  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  nzDisabled:boolean = false
  validateForm:FormGroup = new FormGroup({})
  methodsList:SelectOption[] = methodsList
  contentTypesList:SelectOption[] = contentTypesList
  noticeTypesList:SelectOption[] = noticeTypesList
  responseHeaderList: Array<{
    key: string
    value: string
    [key: string]: any
  }> = [{ key: '', value: '' }]

  responseHeaderTableBody: Array<any> = [...responseHeaderTableBody]

  constructor (
    private message: EoNgMessageService,
     private api:ApiService,
     private fb: UntypedFormBuilder) {
    this.validateForm = this.fb.group({
      title: ['', [Validators.required, EoNgMyValidators.maxByteLength(32), Validators.pattern('^[\u4E00-\u9FA5A-Za-z]+$')]],
      desc: [''],
      url: ['', [Validators.required]],
      method: ['POST', [Validators.required]],
      contentType: ['JSON', [Validators.required]],
      noticeType: ['single', [Validators.required]],
      userSeparator: [''],
      header: [''],
      template: ['{}']
    })
  }

  ngOnInit (): void {
    for (const resBody of this.responseHeaderTableBody) {
      resBody.disabledFn = () => {
        return this.nzDisabled
      }
    }

    this.responseHeaderTableBody.push(
      {
        type: 'btn',
        showFn: (item: any) => {
          return item === this.responseHeaderList[0]
        },
        btns: [
          {
            title: '添加',
            action: 'add',
            disabledFn: () => {
              return this.nzDisabled
            }
          }
        ]
      })

    this.responseHeaderTableBody.push(
      {
        type: 'btn',
        showFn: (item: any) => {
          return item !== this.responseHeaderList[0]
        },
        btns: [
          {
            title: '添加',
            action: 'add',
            disabledFn: () => {
              return this.nzDisabled
            }
          },
          {
            title: '减少',
            action: 'delete',
            disabledFn: () => {
              return this.nzDisabled
            }
          }
        ]
      })

    if (this.webhookId) {
      this.getWebhookMessage(this.webhookId)
    }
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  getWebhookMessage (id:string) {
    this.api.get('warn/webhook', { uuid: id }).subscribe((resp:{code:number, data:{webhook:WebhookData}, msg:string}) => {
      if (resp.code === 0) {
        setFormValue(this.validateForm, resp.data.webhook)
        this.validateForm.controls['template'].setValue(resp.data.webhook.template || '{}')
        this.responseHeaderList = this.transferToList(resp.data.webhook.header)
      }
    })
  }

  transferToList (rawData:any):Array<{key:string, value:string}> {
    const res:Array<{key:string, value:string}> = []
    const keys:Array<string> = Object.keys(rawData)
    if (keys?.length > 0) {
      for (const key of keys) {
        res.push({ key: key, value: rawData[key] })
      }
      return res
    }
    return [{ key: '', value: '' }]
  }

  saveWebhook ():boolean {
    if (this.validateForm.valid) {
      const data:WebhookData = { ...this.validateForm.value, header: this.transferToMap(this.responseHeaderList) }
      if (data.noticeType === 'single') { delete data.userSeparator }
      if (!this.webhookId) {
        this.api.post('warn/webhook', data).subscribe((resp:{code:number, data:{}, msg:string}) => {
          if (resp.code === 0) {
            this.closeModal && this.closeModal()
            this.message.success(resp.msg || '新建Webhook成功！')
          }
        })
      } else {
        this.api.put('warn/webhook', { ...data, uuid: this.webhookId }).subscribe((resp:{code:number, data:{}, msg:string}) => {
          if (resp.code === 0) {
            this.closeModal && this.closeModal()

            this.message.success(resp.msg || '修改Webhook成功！')
          }
        })
      }
    } else {
      Object.values(this.validateForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
    return false
  }

  transferToMap (header:Array<{key:string, value:string}>):{[key:string]:string} {
    const res:{[key:string]:string} = {}
    for (const kv of header) {
      if (kv.key && kv.value) { res[kv.key] = kv.value }
    }
    return res
  }
}
