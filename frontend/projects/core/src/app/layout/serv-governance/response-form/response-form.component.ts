/* eslint-disable no-useless-constructor */
/* eslint-disable dot-notation */
/* eslint-disable camelcase */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-10-27 17:39:12
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-10-28 00:15:51
 * @FilePath: /projects/core/src/app/layout/serv-governance/response-form/response-form.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, Input, OnInit, Output, EventEmitter, SimpleChanges } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from '../../../constant/conf'
import { ApiService } from '../../../service/api.service'

@Component({
  selector: 'eo-ng-response-form',
  templateUrl: './response-form.component.html',
  styles: [
  ]
})
export class ResponseFormComponent implements OnInit {
@Input()
  get responseForm () {
    return this._responseForm as FormGroup
  }

set responseForm (val:FormGroup) {
  this._responseForm = val
}

@Output() responseFormChange:EventEmitter<any> = new EventEmitter()

@Input() disabled: boolean = false
@Input()
get responseHeaderList () {
  return this._responseHeaderList
}

set responseHeaderList (val) {
  this.responseHeaderListChange.emit(val)
  this._responseHeaderList = val
}

@Output() responseHeaderListChange:EventEmitter<any> = new EventEmitter()

private _responseHeaderList: Array<{
  key: string
  value: string
  [key: string]: any
}> = [{ key: '', value: '' }]

@Input() editPage:boolean = false

_responseForm:FormGroup = this.responseForm
contentTypeList: Array<{ label: string; value: string; [key: string]: any }> =
  []

_contentTypeList: Array<{
  label: string
  value: string
  [key: string]: any
}> = []

charsetList: Array<{ label: string; value: string; [key: string]: any }> = []
contentTypeMap: Map<string, string> = new Map()

responseHeaderTableBody: Array<any> = [
  {
    key: 'key',
    type: 'input',
    placeholder: '请输入Key',
    disabledFn: () => {
      return this.disabled
    }
  },
  {
    key: 'value',
    type: 'input',
    placeholder: '请输入Value',
    disabledFn: () => {
      return this.disabled
    }
  },
  {
    type: 'btn',
    showFn: (item: any) => {
      return item === this._responseHeaderList[0]
    },
    btns: [
      {
        title: '添加',
        action: 'add',
        disabledFn: () => {
          return this.disabled
        }
      }
    ]
  },
  {
    type: 'btn',
    showFn: (item: any) => {
      return item !== this._responseHeaderList[0]
    },
    btns: [
      {
        title: '添加',
        action: 'add',
        disabledFn: () => {
          return this.disabled
        }
      },
      {
        title: '减少',
        action: 'delete',
        disabledFn: () => {
          return this.disabled
        }
      }
    ]
  }
]

autoTips: Record<string, Record<string, string>> = defaultAutoTips

constructor (
  private message: EoNgFeedbackMessageService,
  private api: ApiService,
  private fb: UntypedFormBuilder) {
  this._responseForm = this.fb.group({
    status_code: [200, [Validators.required, Validators.pattern(/^[1-9]{1}\d{2}$/)]],
    content_type: ['application/json', [Validators.required]],
    charset: ['UTF-8', [Validators.required]],
    header: [],
    body: []
  })
}

ngOnInit (): void {
  this.getContentTypeList()
  this.getCharsetList()
}

ngOnChanges (changes:SimpleChanges): void {
  if (changes['disabled'] && this.disabled) {
    this.responseForm.disable()
  }
}

getContentTypeList () {
  this.api
    .get('strategy/content-type')
    .subscribe(
      (resp: {
        code: number
        data: { items: Array<{ content_type: string; body: string }> }
        msg: string
      }) => {
        if (resp.code === 0) {
          this.contentTypeMap = new Map()
          this._contentTypeList = []
          this.contentTypeList = []
          for (const index in resp.data.items) {
            this._contentTypeList.push({
              label: resp.data.items[index].content_type,
              value: resp.data.items[index].content_type
            })
            this.contentTypeList.push({
              label: resp.data.items[index].content_type,
              value: resp.data.items[index].content_type
            })
            this.contentTypeMap.set(
              resp.data.items[index].content_type,
              resp.data.items[index].body
            )
          }
          if (!this.editPage) {
            const body: string =
              this.contentTypeMap
                .get(this._responseForm.controls['content_type'].value)
                ?.toString() || ''
            this._responseForm.controls['body'].setValue(body)
          }
        } else {
          this.message.error(resp.msg || '获取数据失败!')
        }
      }
    )
}

getCharsetList () {
  this.api
    .get('strategy/charset')
    .subscribe(
      (resp: {
        code: number
        data: { items: Array<string> }
        msg: string
      }) => {
        if (resp.code === 0) {
          this.charsetList = []
          for (const index in resp.data.items) {
            this.charsetList.push({
              label: resp.data.items[index],
              value: resp.data.items[index]
            })
          }
        } else {
          this.message.error(resp.msg || '获取数据失败!')
        }
      }
    )
}

changeContentType (value: any) {
  this.contentTypeList = [...this._contentTypeList]
  if (value) {
    for (const index in this._contentTypeList) {
      if (this._contentTypeList[index].value === value) {
        return
      }
    }
    this.contentTypeList.unshift({ label: value, value: value })
  }
}
}
