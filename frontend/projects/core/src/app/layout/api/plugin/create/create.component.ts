/* eslint-disable dot-notation */
import { Component, Input, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { setFormValue } from 'projects/core/src/app/constant/form'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'
import { EoNgJsonService } from 'projects/core/src/app/service/eo-ng-json.service'
import { PluginTemplateConfigItem, PluginTemplateData } from '../../types/types'

@Component({
  selector: 'eo-ng-api-plugin-template-create',
  templateUrl: './create.component.html',
  styles: [
  ]
})
export class ApiPluginTemplateCreateComponent implements OnInit {
  @Input() uuid:string = ''
  @Input() editPage:boolean = false
  nzDisabled:boolean = false
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  modalRef:NzModalRef | undefined
  editData:any = null
  validateForm:FormGroup = new FormGroup({})
  configList: PluginTemplateConfigItem[] = []
  pluginConfigError:boolean = false
  startValid:boolean = false
  constructor (private message: EoNgFeedbackMessageService,
    private baseInfo:BaseInfoService,
    private api:ApiService,
    private appConfigService:AppConfigService,
    private fb: UntypedFormBuilder,
    private router: Router,
    private jsonService:EoNgJsonService
  ) {
    this.appConfigService.reqFlashBreadcrumb([
      { title: '插件模板', routerLink: 'router/plugin' },
      { title: '新建模板' }
    ])

    this.validateForm = this.fb.group({
      name: ['', [Validators.required]],
      desc: ['']
    })
  }

  ngOnInit (): void {
    if (this.baseInfo.allParamsInfo.pluginTemplateId) {
      this.appConfigService.reqFlashBreadcrumb([
        { title: '插件模板', routerLink: 'router/plugin' },
        { title: '编辑模板' }
      ])
    }
    if (this.editPage) {
      this.getPluginTemplateMessage()
    }
  }

  ngAfterViewInit () {
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  getPluginTemplateMessage () {
    this.api.get('plugin/template', { uuid: this.uuid }).subscribe((resp:{code:number, data:{template:PluginTemplateData}, msg:string}) => {
      if (resp.code === 0) {
        setFormValue(this.validateForm, resp.data.template)
        this.configList = resp.data.template.plugins.map((plugin) => {
          plugin.disable = !plugin.disable
          plugin.config = JSON.stringify(this.jsonService.handleJsonSchema2Json(JSON.parse(plugin.config))) === '{}' ? plugin.config : JSON.stringify(this.jsonService.handleJsonSchema2Json(JSON.parse(plugin.config)))
          return plugin
        })
      }
    })
  }

  // 返回列表页，当fromList为true时，该页面左侧有分组
  backToList () {
    this.router.navigate(['/', 'router', 'plugin'])
  }

  handlerConfigListChange () {
    this.pluginConfigError = this.startValid ? !this.configList || this.configList?.length === 0 : false
  }

  // 提交api数据
  savePluginTemplate () {
    this.startValid = true
    if (this.configList?.length === 0) {
      this.pluginConfigError = true
    }
    if (this.validateForm.valid && !this.pluginConfigError) {
      const pluginListApi: PluginTemplateConfigItem[] = [] // 提交接口时转换disable
      for (const plugin of this.configList) {
        pluginListApi.push({ ...plugin, disable: !plugin.disable, config: JSON.parse(plugin.config || 'null') })
      }
      if (this.editPage) {
        this.api.put('plugin/template', {
          uuid: this.uuid,
          name: this.validateForm.controls['name'].value,
          desc: this.validateForm.controls['desc'].value,
          plugins: pluginListApi
        }).subscribe(resp => {
          if (resp.code === 0) {
            this.backToList()
            this.message.success(resp.msg || '修改成功！', { nzDuration: 1000 })
          }
        })
      } else {
        this.api.post('plugin/template', {
          ...this.validateForm.value,
          uuid: this.uuid,
          name: this.validateForm.controls['name'].value,
          plugins: pluginListApi
        }).subscribe(resp => {
          if (resp.code === 0) {
            this.message.success(resp.msg || '添加成功！', { nzDuration: 1000 })
            this.backToList()
          }
        })
      }
    } else {
      Object.values(this.validateForm.controls).forEach((control) => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }
}
