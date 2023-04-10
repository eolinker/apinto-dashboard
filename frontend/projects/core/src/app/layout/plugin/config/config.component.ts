/* eslint-disable dot-notation */
import { Component, OnInit } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { SelectOption } from 'eo-ng-select'
import { THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { defaultAutoTips } from '../../../constant/conf'
import { EmptyHttpResponse } from '../../../constant/type'
import { ApiService } from '../../../service/api.service'
import { EoNgMessageService } from '../../../service/eo-ng-message.service'
import { PluginInstallConfigTableHeadName, PluginInstallConfigTableBody } from '../types/conf'
import { PluginInstallConfigData, PluginInstallData } from '../types/types'

@Component({
  selector: 'eo-ng-plugin-config',
  templateUrl: './config.component.html',
  styles: [
  ]
})
export class PluginConfigComponent implements OnInit {
  pluginId:string = ''
  configTableHeadName:THEAD_TYPE[] = [...PluginInstallConfigTableHeadName]
  configTableBody:EO_TBODY_TYPE[] = [...PluginInstallConfigTableBody]
  name:string = ''
  navigation:string = ''
  apiGroup:string = ''
  server:string = ''

  showServer:boolean = false
  invisible:boolean = false
  showApiGroup:boolean = false

  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validateForm:FormGroup = new FormGroup({})

  headerList:Array<PluginInstallConfigData> = []
  queryList:Array<PluginInstallConfigData> = []
  initializeList:Array<PluginInstallConfigData> = []

  navigationList:SelectOption[] = []
  nzDisabled:boolean = false

  modalRef:NzModalRef|undefined
  constructor (private fb: UntypedFormBuilder, private api:ApiService, private message:EoNgMessageService) {
    this.validateForm = this.fb.group({
      navigation: ['', [Validators.required]],
      server: ['', [Validators.required]],
      apiGroup: ['', [Validators.required]]
    })
  }

  ngOnInit (): void {
    this.getMessage()
  }

  getMessage () {
    this.api.get('system/plugin/enable', { id: this.pluginId }).subscribe((resp:{code:number, data:PluginInstallData, msg:string}) => {
      if (resp.code === 0) {
        this.name = resp.data.module.name
        this.navigation = resp.data.module.navigation
        this.apiGroup = resp.data.module.apiGroup
        this.server = resp.data.module.server
        this.headerList = resp.data.module.header.map((header:PluginInstallConfigData) => {
          header.placeholder = header.placeholder || '请输入'
          return header
        })
        this.queryList = resp.data.module.query.map((query:PluginInstallConfigData) => {
          query.placeholder = query.placeholder || '请输入'
          return query
        })
        this.initializeList = resp.data.module.initialize.map((initItem:PluginInstallConfigData) => {
          initItem.placeholder = initItem.placeholder || '请输入'
          return initItem
        })
        this.showServer = resp.data.render.internet
        this.invisible = resp.data.render.invisible
        this.showApiGroup = resp.data.render.apiGroup
      }
    })
  }

  goToNavigation () {
    this.modalRef?.close()
    window.open('../../../navigation')
  }

  checkValid () {
    let valid:boolean = true
    if (this.invisible && this.validateForm.controls['navigation'].invalid) {
      valid = false
      this.validateForm.controls['navigation'].markAsDirty()
      this.validateForm.controls['navigation'].updateValueAndValidity({ onlySelf: true })
    }
    if (this.showApiGroup && this.validateForm.controls['apiGroup'].invalid) {
      valid = false
      this.validateForm.controls['apiGroup'].markAsDirty()
      this.validateForm.controls['apiGroup'].updateValueAndValidity({ onlySelf: true })
    }
    if (this.showServer && this.validateForm.controls['server'].invalid) {
      valid = false
      this.validateForm.controls['server'].markAsDirty()
      this.validateForm.controls['server'].updateValueAndValidity({ onlySelf: true })
    }

    return valid
  }

  enablePlugin () {
    if (this.checkValid()) {
      const data = {
        name: this.name,
        navigation: this.validateForm.controls['navigation'].value,
        apiGroup: this.validateForm.controls['apiGroup'].value,
        server: this.validateForm.controls['server'].value,
        header: this.headerList.map((header:PluginInstallConfigData) => {
          return { name: header.name, value: header.value }
        }),
        query: this.queryList.map((query:PluginInstallConfigData) => {
          return { name: query.name, value: query.value }
        }),
        initialize: this.initializeList.map((initItem:PluginInstallConfigData) => {
          return { name: initItem.name, value: initItem.value }
        })
      }

      this.api.post('system/plugin/enable', { id: this.pluginId }, data).subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '启用插件成功')
        }
      })
    }
  }

  // 不插入数据
  nzCheckAddRow = () => {
    return false
  }
}
