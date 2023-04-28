import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { THEAD_TYPE, TBODY_TYPE } from 'eo-ng-table'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { RouterService } from '../../router.service'
import { PluginTemplateTableHeadName } from '../../types/conf'
import { PluginTemplateItem } from '../../types/types'

@Component({
  selector: 'eo-ng-api-plugin-template-list',
  templateUrl: './list.component.html',
  styles: [
  ]
})
export class ApiPluginTemplateListComponent implements OnInit {
  nzDisabled:boolean = false

  pluginsList: PluginTemplateItem[]= []
  pluginsTableHeadName:THEAD_TYPE[] = [...PluginTemplateTableHeadName]
  pluginsTableBody:TBODY_TYPE[] = [...this.service.createPluginTemplateTbody(this)]
  constructor (
    private message: EoNgFeedbackMessageService,
    private modalService: EoNgFeedbackModalService,
    private api: ApiService,
    public router: Router,
    private navigationService: EoNgNavigationService,
    private service:RouterService
  ) {
    this.navigationService.reqFlashBreadcrumb([{ title: '插件模板' }])
  }

  ngOnInit (): void {
    this.getPluginsData()
  }

  ngAfterViewInit () {
  }

  getPluginsData () {
    this.api.get('plugin/templates').subscribe((resp:{code:number, data:{templates:PluginTemplateItem[]}, msg:string}) => {
      if (resp.code === 0) {
        this.pluginsList = resp.data.templates
      }
    })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  delete (item:any) {
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deletePluginTemplate(item.data)
      }
    })
  }

  deletePluginTemplate (item: PluginTemplateItem) {
    this.api
      .delete('plugin/template', { uuid: item.uuid })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '删除成功', { nzDuration: 1000 })
          this.getPluginsData()
        }
      })
  }

  pluginTableClick= (item:any) => {
    this.router.navigate(['/', 'router', 'plugin-template', 'content', item.data.uuid])
  }

  addPluginTemplate (): void {
    this.router.navigate(['/', 'router', 'plugin-template', 'create'])
  }
}
