/* eslint-disable dot-notation */

import { Component, OnInit, ViewChild, TemplateRef, Input, Output, EventEmitter } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE } from 'projects/core/src/app/constant/app.config'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { PluginTemplateConfigThead, PluginTemplateConfigTbody } from '../../../types/conf'
import { PluginTemplateConfigItem } from '../../../types/types'
import { ApiPluginConfigFormComponent } from '../form/form.component'

@Component({
  selector: 'eo-ng-router-plugin-config-table',
  templateUrl: './table.component.html',
  styles: [
  ]
})
export class ApiPluginConfigTableComponent implements OnInit {
  @ViewChild('configTypeTranslateTpl', { read: TemplateRef, static: true }) configTypeTranslateTpl: TemplateRef<any> | undefined
  @ViewChild('positionTranslateTpl', { read: TemplateRef, static: true }) positionTranslateTpl: TemplateRef<any> | undefined
  @ViewChild('switchTpl', { read: TemplateRef, static: true }) switchTpl: TemplateRef<any> | undefined
  @Input() nzDisabled:boolean = false
  @Input()
  get configList () {
    return this._configList
  }

  set configList (val:PluginTemplateConfigItem[]) {
    this._configList = val
    this.configListChange.emit(val)
  }

  @Output() configListChange = new EventEmitter()
  _configList:PluginTemplateConfigItem[] = []
  editData:PluginTemplateConfigItem | undefined
  validateConfigForm:FormGroup = new FormGroup({})
  modalRef:NzModalRef | undefined
  configTableHeadName:THEAD_TYPE[] = [...PluginTemplateConfigThead]
  configTableBody:EO_TBODY_TYPE[]= [...PluginTemplateConfigTbody]

  constructor (private fb: UntypedFormBuilder,
    private modalService:EoNgFeedbackModalService) {
    this.validateConfigForm = this.fb.group({
      position: ['', [Validators.required]],
      key: ['', [Validators.required, Validators.pattern('^[a-zA-Z][a-zA-Z0-9-_]*')]],
      configType: ['', [Validators.required]],
      pattern: ['']
    })
  }

  ngOnInit (): void {
    this.configTableBody[3].btns[0].click = (item:any) => {
      this.openDrawer(item.data)
    }
    this.configTableBody[3].btns[0].disabledFn = () => {
      return this.nzDisabled
    }
    this.configTableBody[3].btns[1].disabledFn = () => {
      return this.nzDisabled
    }
  }

  ngAfterViewInit () {
    this.configTableBody[1].title = this.switchTpl
  }

  configTableClick = (item:{data:PluginTemplateConfigItem}) => {
    this.openDrawer(item.data)
  }

  openDrawer (data?:any) {
    if (data) {
      this.editData = data
    }

    this.modalRef = this.modalService.create({
      nzTitle: '配置插件',
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: ApiPluginConfigFormComponent,
      nzComponentParams: {
        data: data, closeDrawer: this.closeDrawer, configList: this.configList, editData: data || undefined
      },
      nzOkDisabled: this.nzDisabled,
      nzOnOk: (component:ApiPluginConfigFormComponent) => {
        component.saveConfig()
        this.configList = component.configList
        this.configListChange.emit(this.configList)
        return false
      }
    })
  }

  closeDrawer = () => {
    this.modalRef?.close()
  }
}
