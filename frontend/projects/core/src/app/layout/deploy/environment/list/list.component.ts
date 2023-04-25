/* eslint-disable no-useless-constructor */
/* eslint-disable no-undef */
/* eslint-disable camelcase */
/* eslint-disable dot-notation */
/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-07-20 22:34:58
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-11-02 23:51:33
 * @FilePath: /apinto/src/app/layout/deploy/deploy-environment/deploy-environment.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { Component, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import {
  EoNgFeedbackModalService,
  EoNgFeedbackMessageService
} from 'eo-ng-feedback'
import { THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgNavigationService } from 'projects/core/src/app/service/eo-ng-navigation.service'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { DeployEnvironmentDetailComponent } from '../detail/detail.component'
import { DeployGlobalEnvTableBody, DeployGlobalEnvTableHeadName } from '../types/conf'

@Component({
  selector: 'eo-ng-deploy-environment-list',
  templateUrl: './list.component.html',
  styles: [
    `
      label {
        line-height: 32px !important;
      }

      input.ant-input:not(.w206):not(.w131):not(.w240),
      eo-ng-select.ant-select {
        width: 216px !important;
      }
    `
  ]
})
export class DeployEnvironmentListComponent {
  @ViewChild('variableStatusTpl', { read: TemplateRef, static: true })
  variableStatusTpl: TemplateRef<any> | undefined

  clusterName: string = ''

  globalEnvForms: {
    variables: Array<{
      key: string
      usage: number
      description: string
      operator: string
      createTime: string
    }>
    total: number
  } = {
    variables: [],
    total: 0
  }

  nzDisabled: boolean = false

  globalEnvTableHeadName: THEAD_TYPE[] = [...DeployGlobalEnvTableHeadName]
  globalEnvTableBody: EO_TBODY_TYPE[] = [...DeployGlobalEnvTableBody]

  editConfigDrawerRef: NzModalRef | undefined

  statusList: Array<{ label: string; value: string }> = [
    { label: '使用中', value: 'IN_USE' },
    { label: '空闲', value: 'UNUSED' }
  ]

  searchForm: { key: string; status: string } = { key: '', status: '' }

  // 环境变量分页参数
  variablePage: { pageNum: number; pageSize: number; total: number } = {
    pageNum: 1,
    pageSize: 20,
    total: 0
  }

  constructor (
    private message: EoNgFeedbackMessageService,
    private modalService: EoNgFeedbackModalService,
    private api: ApiService,
    private router: Router,
    private appConfigService: EoNgNavigationService
  ) {
    this.appConfigService.reqFlashBreadcrumb([
      { title: '环境变量', routerLink: 'deploy/variable' }
    ])
  }

  ngOnInit (): void {
    this.getVariables()
  }

  ngAfterViewInit () {
    this.globalEnvTableBody[4].title = this.variableStatusTpl
    this.globalEnvTableBody[5].btns[0].click = (item: any) => {
      this.openDrawer(item.data)
    }
    this.globalEnvTableBody[5].btns[1].click = (item: any) => {
      this.deleteModal(item.data)
    }
    this.globalEnvTableBody[5].btns[1].disabledFn = (data:any, item:any) => {
      return this.nzDisabled || item.data.status === 'IN_USE'
    }
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  deleteModal (item: any, e?: Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.delete(item)
      }
    })
  }

  getVariables () {
    this.api
      .get('variables', {
        pageNum: this.variablePage.pageNum,
        pageSize: this.variablePage.pageSize,
        key: this.searchForm?.key || '',
        status: this.searchForm?.status || ''
      })
      .subscribe((resp) => {
        if (resp.code === 0) {
          this.globalEnvForms = resp.data
          this.variablePage.total = resp.data.total
        }
      })
  }

  globalEnvTableClick = (item: any) => {
    this.openDrawer(item.data)
  }

  delete (item: any) {
    this.api.delete('variable', { key: item.key || '' }).subscribe((resp) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '删除成功', { nzDuration: 1000 })
        this.getVariables()
      }
    })
  }

  addConfig () {
    this.router.navigate(['/', 'deploy', 'variable', 'create'])
  }

  resetSearch () {
    this.searchForm.key = ''
    this.searchForm.status = ''
  }

  openDrawer (rowItem:any) {
    this.editConfigDrawerRef = this.modalService.create({
      nzTitle: '查看环境变量',
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: DeployEnvironmentDetailComponent,
      nzComponentParams: {
        envKey: rowItem.key
      },
      nzFooter: null
    })
  }
}
