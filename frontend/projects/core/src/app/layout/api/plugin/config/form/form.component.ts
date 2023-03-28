/* eslint-disable dot-notation */
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgJsonService } from 'projects/core/src/app/service/eo-ng-json.service'
import { PluginTemplateConfigItem } from '../../../types/types'

@Component({
  selector: 'eo-ng-api-plugin-config-form',
  templateUrl: './form.component.html',
  styles: [
  ]
})
export class ApiPluginConfigFormComponent implements OnInit {
  @Input() validateConfigForm:FormGroup = new FormGroup({})
  @Input() editData?:PluginTemplateConfigItem
  @Input() configList:PluginTemplateConfigItem[] = []

  @Output() eoNgCloseDrawer: EventEmitter<string> = new EventEmitter()
  @Output() configListChange : EventEmitter<PluginTemplateConfigItem[]> = new EventEmitter()

  pluginsList:(SelectOption & {name:string, config:string})[] =[]
  configHeaderSet:Set<string> = new Set()
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  data:PluginTemplateConfigItem|undefined
  accessUrl:string = '' // 用来判断权限的url
  nzDisabled:boolean = false
  code:string = ''
  showValid:boolean = false
  constructor (private jsonService:EoNgJsonService, private fb: UntypedFormBuilder, private api:ApiService, private message:EoNgFeedbackMessageService) {
  }

  ngOnInit (): void {
    this.validateConfigForm = this.fb.group({
      name: [this.data?.name || '', [Validators.required]]
    })
    if (this.data) {
      try {
        this.code = JSON.stringify(this.jsonService.handleJsonSchema2Json(JSON.parse(this.data.config))) === '{}' ? this.data.config : this.data.config
      } catch {
        this.code = this.data.config
      }
    }
    this.getPluginList()
  }

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  getPluginList () {
    this.api.get('plugin/enum').subscribe((resp: {code:number, data:{plugins:Array<{config:string, name:string}>}, msg:string}) => {
      if (resp.code === 0) {
        this.pluginsList = resp.data.plugins.map((item:{config:string, name:string, label?:string, value?:string}) => {
          item.label = item.name
          item.value = item.name
          return item
        }) as (SelectOption & {name:string, config:string})[]
        if (!this.editData) {
          this.validateConfigForm.controls['name'].setValue(this.pluginsList[0].value)
          this.changePluginChange()
        }
      }
    })
  }

  codeChange () {
    this.showValid = !this.code
  }

  saveConfig () {
    if (!this.code) {
      this.showValid = true
      return
    }
    if (this.validateConfigForm.valid) {
      if (!this.data) {
        if (this.configHeaderSet.has(this.validateConfigForm.controls['name'].value)) {
          for (const index in this.configList) {
            if (this.configList[index].name === this.validateConfigForm.controls['name'].value && this.configList[index].config === this.code) {
              this.configList.splice(Number(index), 1)
              break
            }
          }
        }
      } else {
        for (const index in this.configList) {
          if (this.configList[index].name === this.editData!.name && this.configList[index].config === this.editData!.config && this.configList[index].disable === this.editData!.disable) {
            this.configList.splice(Number(index), 1)
            break
          }
        }
      }
      this.configList = [{ name: this.validateConfigForm.controls['name'].value, config: this.code, disable: this.data?.disable !== false }, ...(this.configList?.length > 0 ? this.configList : [])]
      this.closeDrawer()
    } else {
      Object.values(this.validateConfigForm.controls).forEach(control => {
        if (control.invalid) {
          control.markAsDirty()
          control.updateValueAndValidity({ onlySelf: true })
        }
      })
    }
  }

  closeDrawer () {
    this.eoNgCloseDrawer.emit('config')
  }

  changePluginChange () {
    this.validateConfigForm.patchValue({
      config: this.pluginsList.filter((plugin) => {
        return plugin.name === this.validateConfigForm.controls['name'].value
      })[0].config
    })
    this.code = JSON.stringify(this.jsonService.handleJsonSchema2Json(JSON.parse(this.pluginsList.filter((plugin) => {
      return plugin.name === this.validateConfigForm.controls['name'].value
    })[0].config)))
  }
}
