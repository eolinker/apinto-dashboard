/*
 * @Date: 2024-01-24 14:52:35
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-01-24 14:57:22
 * @FilePath: \apinto\projects\core\src\app\component\plugin-config-modal\plugin-config-modal.component.ts
 */
/*
 * @Date: 2024-01-24 14:52:35
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-01-24 14:53:08
 * @FilePath: \apinto\projects\core\src\app\component\plugin-config-modal\plugin-config-modal.component.ts
 */
import { Component, EventEmitter, Input, Output, TemplateRef, ViewChild } from '@angular/core'
import { FormGroup, UntypedFormBuilder, Validators } from '@angular/forms'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { MODAL_NORMAL_SIZE } from '../../constant/app.config'
import { ApiPluginConfigFormComponent } from '../../layout/plugin-template/config/form/form.component'
import { PluginTemplateConfigThead, PluginTemplateConfigTbody } from '../../layout/plugin-template/types/conf'

export type PluginTemplateConfigItem = {
  name:string
  config:string
  disable:boolean
  eoKey?:string
}

@Component({
  selector: 'eo-ng-plugin-config-modal',
  template: `
    <div>
      <button
        type="button"
        [disabled]="nzDisabled"
        eo-ng-button
        (click)="openDrawer()"
      >
        添加插件
      </button>
  </div>
  <div
    *ngIf="configList && configList.length > 0"
    class="mt-btnybase"
    style="width: 524px"
  >
    <eo-ng-apinto-table
      [nzTbody]="configTableBody"
      [nzThead]="configTableHeadName"
      [nzData]="configList"
      (nzDataChange)="configListChange.emit(configList)"
      [nzTrClick]="configTableClick"
      [nzMaxOperatorButton]="2"
      [nzNoScroll]="true"
    >
    </eo-ng-apinto-table>
  </div>

  <ng-template #switchTpl let-item="item">
    <eo-ng-switch
      [(ngModel)]="item.disable"
      (ngModelChange)="this.configListChange.emit(this.configList)"
      (click)="$event.stopPropagation()"
      [nzDisabled]="nzDisabled"
    ></eo-ng-switch>
  </ng-template>

  `,
  styles: [
  ]
})
export class PluginConfigModalComponent {
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
