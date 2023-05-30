import { Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { EmptyHttpResponse } from 'projects/core/src/app/constant/type'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { of } from 'rxjs'
import { DeployService } from '../../deploy.service'
import { PluginsTableHeadName } from '../types/conf'
import { PluginItem } from '../types/types'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'

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
    private navigationService: EoNgNavigationService,
    private service:DeployService
  ) {
    this.navigationService.reqFlashBreadcrumb([{ title: '节点插件' }])
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
    return this.getDragCheck(oldArr)
  }

  getDragCheck (oldArr:any) {
    this.api
      .put('plugin/sort', { names: oldArr })
      .subscribe((resp:EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '删除成功', { nzDuration: 1000 })
          this.getPluginsData()
          return of(true)
        } else {
          return of(true)
        }
      })
  }

  pluginTableClick= (item:any) => {
    this.router.navigate(['/', 'deploy', 'plugin', 'message', item.data.name])
  }

  addPlugin (): void {
    this.router.navigate(['/', 'deploy', 'plugin', 'create'])
  }

  copyCallback () {
    this.message.success('复制成功', {
      nzDuration: 1000
    })
  }
}
