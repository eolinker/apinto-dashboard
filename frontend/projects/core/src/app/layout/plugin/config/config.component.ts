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
import { PluginInstallConfigData } from '../types/types'
import { EoNgNavigationService } from '../../../service/eo-ng-navigation.service'

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
  server:string = ''

  showServer:boolean = false
  nameConfilct:boolean = false
  autoTips: Record<string, Record<string, string>> = defaultAutoTips
  validateForm:FormGroup = new FormGroup({})

  headerList:Array<PluginInstallConfigData> = []
  queryList:Array<PluginInstallConfigData> = []
  initializeList:Array<PluginInstallConfigData> = []

  navigationList:SelectOption[] = []
  nzDisabled:boolean = false

  modalRef:NzModalRef|undefined
  refreshPage:Function|undefined
  constructor (private fb: UntypedFormBuilder, private api:ApiService, private message:EoNgMessageService,
    private navService:EoNgNavigationService) {
    this.validateForm = this.fb.group({
      server: ['', [Validators.required]]
    })
  }

  ngOnInit (): void {
  }

  checkValid () {
    let valid:boolean = true
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

      this.api.post('system/plugin/enable', data, { id: this.pluginId }).subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '启用插件成功')
          const subscription = this.navService.getMenuList().subscribe(() => {
            subscription.unsubscribe()
          })
          this.refreshPage && this.refreshPage()
        }
      })
    }
  }

  // 不插入数据
  nzCheckAddRow = () => {
    return false
  }
}
