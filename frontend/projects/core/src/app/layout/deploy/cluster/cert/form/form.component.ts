/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { Buffer } from 'buffer'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'

@Component({
  selector: 'eo-ng-deploy-cluster-cert-form',
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
export class DeployClusterCertFormComponent implements OnInit {
  @Input() closeModal?:(value?:any)=>void
  validateForm:FormGroup = new FormGroup({})
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  clusterName:string = ''
  certId:string = ''
  nzDisabled:boolean = false

  constructor (
    private message: EoNgFeedbackMessageService,
    public api:ApiService,
    private fb: UntypedFormBuilder,
    private appConfigService:AppConfigService) {
    this.validateForm = this.fb.group({
      key: ['', [Validators.required]],
      pem: ['', [Validators.required]]
    })
    this.appConfigService.reqFlashBreadcrumb([{ title: '网关集群', routerLink: 'deploy/cluster' }, { title: '证书管理' }])
  }

  ngOnInit (): void {
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  save (usage:string):void {
    if (this.validateForm.valid) {
      switch (usage) {
        case 'addCert':
          this.api.post('cluster/' + this.clusterName + '/certificate', { key: this.encode(this.validateForm.controls['key'].value), pem: this.encode(this.validateForm.controls['pem'].value) || '' })
            .subscribe((resp:EmptyHttpResponse) => {
              if (resp.code === 0) {
                this.closeModal && this.closeModal()
              }
            })
          break
        case 'editCert':
          this.api.put('cluster/' + this.clusterName + '/certificate/' + this.certId, { key: this.encode(this.validateForm.controls['key'].value), pem: this.encode(this.validateForm.controls['pem'].value) || '' })
            .subscribe((resp:EmptyHttpResponse) => {
              if (resp.code === 0) {
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
            case 'key':
              this.validateForm.controls['key'].setValue(content?.toString() || '')
              break
            case 'pem':
              this.validateForm.controls['pem'].setValue(content?.toString() || '')
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
    return Buffer.from(str).toString('base64')
  }
}
