/* eslint-disable dot-notation */
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { defaultAutoTips } from 'projects/core/src/app/constant/conf'
import { ApiService } from 'projects/core/src/app/service/api.service'
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
  dataType = ['string', 'boolean', 'number']
  jsonKey = ['eo:type', 'dependencies', 'skill', 'switch', 'label', 'ui:sort', 'properties', 'type'] // json schema自定义的关键字
  constructor (private fb: UntypedFormBuilder, private api:ApiService, private message:EoNgFeedbackMessageService) {
  }

  ngOnInit (): void {
    this.validateConfigForm = this.fb.group({
      name: [this.data?.name || '', [Validators.required]]
    })
    if (this.data) {
      this.code = JSON.stringify(this.handleJsonSchema2Json(JSON.parse(this.data.config)))
    }
    this.getPluginList()
  }

  handleJsonSchema2Json (data:any) {
    const obj:{[k:string]:any} = {}
    for (const key in data.properties) {
      const param = data.properties[key]
      const type = param.type
      if (!this.jsonKey.includes(key)) {
        obj[key] = {}
        if (this.dataType.includes(type)) {
          obj[key] = param.default || ''
        } else if (type === 'array') {
          const items = param.items
          if (this.dataType.includes(items.type)) {
            obj[key] = [items.default]
          } else {
            obj[key] = [this.handleJsonSchema2Json(items)]
          }
        } else if (type === 'object') {
          obj[key] = this.handleJsonSchema2Json(param)
        } else if (type === 'number') {
          obj[key] = 0
        } else if (type === 'boolean') {
          obj[key] = false
        }
      }
    }
    return obj
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
        }
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
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
      this.configList = [{ name: this.validateConfigForm.controls['name'].value, config: this.code, disable: this.data?.disable || false }, ...(this.configList?.length > 0 ? this.configList : [])]
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
    this.code = JSON.stringify(this.handleJsonSchema2Json(JSON.parse(this.pluginsList.filter((plugin) => {
      return plugin.name === this.validateConfigForm.controls['name'].value
    })[0].config)))
  }
}
