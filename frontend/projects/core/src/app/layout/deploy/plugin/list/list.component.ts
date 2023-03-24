import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'
import { DeployService } from '../../deploy.service'
import { PluginsTableHeadName } from '../types/conf'
import { PluginItem } from '../types/types'

@Component({
  selector: 'eo-ng-deploy-plugin-list',
  templateUrl: './list.component.html',
  styles: [
  ]
})
export class DeployPluginListComponent implements OnInit {
  @ViewChild('pluginName', { read: TemplateRef, static: true }) pluginName: TemplateRef<any> | undefined

  nzDisabled:boolean = false

  pluginsList: PluginItem[]= []
  pluginsTableHeadName:THEAD_TYPE[] = [...PluginsTableHeadName]
  pluginsTableBody:TBODY_TYPE[] = []
  constructor (
    private message: EoNgFeedbackMessageService,
    private modalService: EoNgFeedbackModalService,
    private api: ApiService,
    public router: Router,
    private appConfigService: AppConfigService,
    private service:DeployService
  ) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '插件管理' }])
  }

  ngOnInit (): void {
    this.getPluginsData()
  }

  ngAfterViewInit () {
    this.pluginsTableBody = this.service.createPluginsTbody(this)
  }

  getPluginsData () {
    this.api.get('plugins').subscribe((resp:{code:number, data:{plugins:PluginItem[]}, msg:string}) => {
      if (resp.code === 0) {
        this.pluginsList = resp.data.plugins
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
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
        this.deletePlugin(item)
      }
    })
  }

  deletePlugin (item: PluginItem) {
    this.api
      .delete('plugin', { name: item.name })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '删除成功', { nzDuration: 1000 })
          this.getPluginsData()
        } else {
          this.message.error(resp.msg || '删除失败!')
        }
      })
  }

  // 此处拖拽先按照每次item拖拽到currentitem下方来计算，后根据组件计算
  dragPlugin = (preItem:any, currentItem:any) => {
    let preIndex:number = 0
    let currentIndex:number = 0
    const oldArr:Array<string> = this.pluginsList.map((item:PluginItem) => {
      if (item.name === preItem.data.name) {
        preIndex = this.pluginsList.indexOf(item)
      } else if (item.name === currentItem.data.name) {
        currentIndex = this.pluginsList.indexOf(item)
      }
      return item.name
    })

    oldArr.splice(preIndex, 1)
    oldArr.splice(currentIndex, 0, preItem.data.name)
    return this.api
      .put('plugin', { uuids: oldArr })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '删除成功', { nzDuration: 1000 })
          this.getPluginsData()
        } else {
          this.message.error(resp.msg || '删除失败!')
        }
      })
  }

  pluginTableClick= (item:any) => {
    this.router.navigate(['/', 'deploy', 'plugin', 'message', item.name])
  }

  addPlugin (): void {
    this.router.navigate(['/', 'deploy', 'plugin', 'create'])
  }
}
