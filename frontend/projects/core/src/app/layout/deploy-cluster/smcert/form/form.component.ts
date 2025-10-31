/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { Buffer } from 'buffer'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'
import { DeploySmcertData } from '../../types/types'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'

@Component({
  selector: 'eo-ng-deploy-cluster-smcert-form',
  templateUrl: './form.component.html',
  styles: [
    `
    textarea.cert_textarea{
      width:446px !important;
      max-width:446px !important;
      height:86px;
    }
    label.ant-btn{
      font-size:14px;
      font-family:"Helvetica Neue", "Helvetica", "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "Noto Sans CJK SC", "WenQuanYi Micro Hei", "Arial", sans-serif;
      white-space: nowrap;
      text-align: center;
      border-radius: var(--border-radius);
      padding: 7px 8px !important;
      height: 32px;
      line-height: 14px;
      text-shadow: none;
      box-shadow: none;
    }

    label.ant-btn:not(.ant-btn-disabled):focus,
    label.ant-btn:not(.ant-btn-disabled):hover,
    label.ant-btn:not(.ant-btn-disabled):active
     {
      border: 1px solid var(--primary-color) !important;
      background-color: #fff !important;
      color: var(--primary-color) !important;
    }

    label.ant-btn span{
      font-size:14px;
    }

    label.ant-btn span:before{
      font-size:16px;
      margin-right:4px;
    }

    label.ant-btn.ant-btn-disabled,
    label.ant-btn.ant-btn-disabled:active,
    label.ant-btn.ant-btn-disabled:focus,
    label.ant-btn.ant-btn-disabled:hover {
      color: rgba(0, 0, 0, 0.25) !important;
      border-color: #d9d9d9 !important;
      background: #f5f5f5 !important;
      text-shadow: none !important;
      box-shadow: none !important;
      cursor: not-allowed;
    }
`
  ]
})
export class DeployClusterSmcertFormComponent implements OnInit {
  @Input() closeModal?:(value?:any)=>void
  @Input() editPage:boolean = false
  validateForm:FormGroup = new FormGroup({})
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  clusterName:string = ''
  smcertId:string = ''
  nzDisabled:boolean = false

  constructor (
    private message: EoNgFeedbackMessageService,
    public api:ApiService,
    private fb: UntypedFormBuilder) {
    this.validateForm = this.fb.group({
      sign_key: ['', [Validators.required]],
      sign_cert: ['', [Validators.required]],
      enc_key: ['', [Validators.required]],
      enc_cert: ['', [Validators.required]]
    })
  }

  ngOnInit (): void {
    if (this.editPage) {
      this.getCertMessage()
    }
  }

  getCertMessage () {
    this.api.get(`cluster/${this.clusterName}/gm_certificate/${this.smcertId}`)
      .subscribe((resp:{code:number, data:{certificate:DeploySmcertData}, msg:string}) => {
        if (resp.code === 0) {
          this.validateForm.patchValue({
            sign_key: this.decode(resp.data.certificate.signKey),
            sign_cert: this.decode(resp.data.certificate.signCert),
            enc_key: this.decode(resp.data.certificate.encKey),
            enc_cert: this.decode(resp.data.certificate.encCert)
          })
        }
      })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 保存
  save (usage:string):void {
    if (this.validateForm.valid) {
      switch (usage) {
        case 'addSmcert':
          this.api.post('cluster/' + this.clusterName + '/gm_certificate', {
            sign_key: this.encode(this.validateForm.controls['sign_key'].value),
            sign_cert: this.encode(this.validateForm.controls['sign_cert'].value),
            enc_key: this.encode(this.validateForm.controls['enc_key'].value),
            enc_cert: this.encode(this.validateForm.controls['enc_cert'].value)
          })
            .subscribe((resp:EmptyHttpResponse) => {
              if (resp.code === 0) {
                this.message.success(resp.msg || '新增成功')
                this.closeModal && this.closeModal()
              }
            })
          break
        case 'editSmcert':
          this.api.put('cluster/' + this.clusterName + '/gm_certificate/' + this.smcertId, {
            sign_key: this.encode(this.validateForm.controls['sign_key'].value),
            sign_cert: this.encode(this.validateForm.controls['sign_cert'].value),
            enc_key: this.encode(this.validateForm.controls['enc_key'].value),
            enc_cert: this.encode(this.validateForm.controls['enc_cert'].value)
          })
            .subscribe((resp:EmptyHttpResponse) => {
              if (resp.code === 0) {
                this.message.success(resp.msg || '修改成功')
                this.closeModal && this.closeModal()
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
  }

  // 读取文件内容
  readSingleFile (e: any, type:string) {
    if (e.target?.files[0]) {
      const reader = new FileReader()
      let content
      reader.onload = (e) => {
        content = e.target?.result
        if (content) {
          switch (type) {
            case 'sign_key':
              this.validateForm.controls['sign_key'].setValue(content?.toString() || '')
              break
            case 'sign_cert':
              this.validateForm.controls['sign_cert'].setValue(content?.toString() || '')
              break
            case 'enc_key':
              this.validateForm.controls['enc_key'].setValue(content?.toString() || '')
              break
            case 'enc_cert':
              this.validateForm.controls['enc_cert'].setValue(content?.toString() || '')
              break
            default:
              break
          }
        }
      }
      reader.readAsText(e.target?.files[0], 'utf-8')
    }
  }

  // 字符串转base64
  encode (str:string) {
    return str ? Buffer.from(str).toString('base64') : ''
  }

  decode (str:string) {
    return str ? Buffer.from(str, 'base64').toString('utf8') : ''
  }
}
