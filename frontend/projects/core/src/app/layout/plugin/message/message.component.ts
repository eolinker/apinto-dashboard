import { Component, Inject, OnInit } from '@angular/core'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { NzModalRef } from 'ng-zorro-antd/modal'
import {
  MODAL_NORMAL_SIZE,
  MODAL_SMALL_SIZE
} from '../../../constant/app.config'
import { EmptyHttpResponse } from '../../../constant/type'
import { API_URL, ApiService } from '../../../service/api.service'
import { BaseInfoService } from '../../../service/base-info.service'
import { EoNgMessageService } from '../../../service/eo-ng-message.service'
import { EoNgNavigationService } from '../../../service/eo-ng-navigation.service'
import { PluginConfigComponent } from '../config/config.component'
import {
  PluginInstallConfigData,
  PluginInstallData,
  PluginMessage
} from '../types/types'
import { MarkdownService } from 'ngx-markdown'
import { Router } from '@angular/router'
import { UntypedFormBuilder, Validators } from '@angular/forms'

@Component({
  selector: 'eo-ng-plugin-message',
  template: `
    <header class="mx-[40px] my-[20px]">
      <div class="flex justify-between  mb-btnybase items-center">
        <div class="flex">
          <div class="mr-btnrbase w-[50px] h-[50px]">
            <img
              class="mr-btnrbase"
              [src]="icon"
              alt="icon"
              width="50px"
              height="50px"
            />
          </div>
          <div class="flex flex-col justify-around">
            <p class="text-[14px] font-medium">{{ title }}</p>
            <p>
              <span
                class="mr-[8px] font-medium text-[#0A89FF] bg-[#0A89FF1A] px-[4px] py-[2px] leading-[20px] rounded"
                *ngIf="enable"
                >已启用</span
              >
              <span
                class="mr-[8px] font-medium text-[#bbbbbb] bg-[#bbbbbb1A] px-[4px] py-[2px] leading-[20px] rounded"
                *ngIf="!enable"
              >
                未启用</span
              >
            </p>
          </div>
        </div>
        <div *ngIf="showBtn">
          <button
            *ngIf="!enable"
            class="ml-btnybase ant-btn-primary"
            eo-ng-button
        eoNgUserAccess="module-plugin"
            (click)="enablePlugin()"
          >
            启用
          </button>
          <button
            *ngIf="enable && canDisable"
            class="ml-btnybase"
            eo-ng-button
        eoNgUserAccess="module-plugin"
            (click)="disablePluginModal(false)"
          >
            停用
          </button>
          <button
            *ngIf="!enable && uninstall"
            class="ml-btnybase"
            eo-ng-button
        eoNgUserAccess="module-plugin"
            (click)="disablePluginModal(true)"
          >
            卸载
          </button>
        </div>
      </div>
      <p>{{ resume }}</p>
    </header>
    <section class="flex-1 p-[40px] markdown-block overflow-auto">
      <eo-ng-empty
        *ngIf="showEmpty"
        nzMainTitle="暂无数据"
        nzInputImage="simple"
      ></eo-ng-empty>
      <markdown
        class="markdown-body"
        *ngIf="!showEmpty"
        [src]="getMd()"
        [srcRelativeLink]="true"
        (load)="loadMd()"
        (error)="onError($event)"
      ></markdown>
    </section>
  `,
  styleUrls: ['./message.component.scss']
})
export class PluginMessageComponent implements OnInit {
  title: string = ''
  resume: string = ''
  icon: string = '' || './assets/default-plugin-icon.svg'
  enable: boolean = false
  uninstall: boolean = false
  pluginId: string = ''
  modalRef: NzModalRef | undefined
  mdFileName: string = ''
  showEmpty: boolean = false
  canDisable:boolean = false
  showBtn:boolean = false
  constructor (
    private message: EoNgMessageService,
    private modalService: EoNgFeedbackModalService,
    private api: ApiService,
    private baseInfo: BaseInfoService,
    private navigationService: EoNgNavigationService,
    private markdownService: MarkdownService,
    private router: Router,
    private fb: UntypedFormBuilder,
    @Inject(API_URL) public urlPrefix: string
  ) {
    this.navigationService.reqFlashBreadcrumb([
      { title: '企业插件', routerLink: ['/', 'module-plugin', 'group', 'list', ''] },
      { title: '插件详情' }
    ])
  }

  ngOnInit (): void {
    this.pluginId = this.baseInfo.allParamsInfo.pluginId
    this.mdFileName = this.baseInfo.allParamsInfo.mdFileName
    this.getPluginDetail()
    this.markdownService.renderer.link = (href, title, text) => {
      let html = ''
      if (
        href &&
        /^(?![http])[.]*/.test(href!) &&
        /^(?![#])[.]*/.test(href!) &&
        href.includes('.md')
      ) {
        html = `<a href="plugin/message/${this.pluginId}/${href}">${text}</a>`
      } else if (
        href &&
        /^(?![http])[.]*/.test(href!) &&
        /^(?![#])[.]*/.test(href!)
      ) {
        html = `<a href="plugin/message/${this.pluginId}/${href}">${text}</a>`
      } else {
        html =
          '<a  role="link"  tabindex="0" target="_blank" rel="nofollow noopener noreferrer" href="' +
          href +
          '">' +
          text +
          '</a>'
      }

      return html
    }

    this.markdownService.renderer.image = (src, title, alt) => {
      let html
      if (src && /^(?![http])[.]*/.test(src!)) {
        const newSrc = src.replace('./resource', '/resource')
        html = `<image src="${this.urlPrefix}plugin/info/${this.pluginId}${newSrc}" alt=${alt}/>`
      } else {
        html = `<image src="${src}" alt=${alt}/>`
      }
      return html
    }
  }

  getPluginDetail = () => {
    this.api
      .get('system/plugin/info', { id: this.pluginId })
      .subscribe(
        (resp: {
          code: number
          data: { plugin: PluginMessage }
          msg: string
        }) => {
          if (resp.code === 0) {
            this.title = resp.data.plugin.cname
            this.resume = resp.data.plugin.resume
            this.icon = resp.data.plugin.icon
              ? `${this.urlPrefix}plugin/icon/${this.pluginId}/${resp.data.plugin.icon}`
              : './assets/default-plugin-icon.svg'
            this.enable = resp.data.plugin.enable
            this.uninstall = resp.data.plugin.uninstall
            this.canDisable = resp.data.plugin.canDisable // true表示该插件可禁用
            this.showBtn = true
          }
        }
      )
  }

  getMd () {
    // return './assets/README.md' // 本地调试用
    return `../../plugin/md/${this.pluginId}/${this.mdFileName || ''}`
  }

  loadMd () {
    this.showEmpty = false
  }

  onError (value: any) {
    console.error('解析md文档出现问题', value)
    this.showEmpty = true
  }

  enablePlugin () {
    const params: {
      name: string
      server: string
      headerList: Array<PluginInstallConfigData>
      queryList: Array<PluginInstallConfigData>
      initializeList: Array<PluginInstallConfigData>
      showServer: boolean
      nameConflict:boolean
    } = {
      name: '',
      server: '',
      headerList: [],
      queryList: [],
      initializeList: [],
      showServer: false,
      nameConflict: false
    }
    this.api
      .get('system/plugin/enable', { id: this.pluginId })
      .subscribe(
        (resp: { code: number; data: PluginInstallData; msg: string }) => {
          if (resp.code === 30001) {
            this.message.success(resp.msg || '启用插件成功')
            const subscription = this.navigationService
              .getMenuList()
              .subscribe(() => {
                subscription.unsubscribe()
              })
            this.getPluginDetail()
          } else if (resp.code === 0) {
            params.name = resp.data.module.name
            params.server = resp.data.module.server
            params.headerList = resp.data.render.headers.map(
              (header: PluginInstallConfigData) => {
                header.placeholder = header.placeholder || '请输入'
                const currentValue = this.findConfigValue(resp.data.module.header, header.name)
                header.value = currentValue === undefined ? header.value : currentValue
                return header
              }
            )
            params.queryList = resp.data.render.querys.map(
              (query: PluginInstallConfigData) => {
                query.placeholder = query.placeholder || '请输入'
                const currentValue = this.findConfigValue(resp.data.module.query, query.name)
                query.value = currentValue === undefined ? query.value : currentValue
                return query
              }
            )
            params.initializeList = resp.data.render.initialize.map(
              (initItem: PluginInstallConfigData) => {
                initItem.placeholder = initItem.placeholder || '请输入'
                const currentValue = this.findConfigValue(resp.data.module.initialize, initItem.name)
                initItem.value = currentValue === undefined ? initItem.value : currentValue
                return initItem
              }
            )
            params.showServer = resp.data.render.internet
            params.nameConflict = resp.data.render.nameConflict
            if (
              params.showServer ||
              params.nameConflict ||
              params.headerList.length ||
              params.queryList.length ||
              params.initializeList.length
            ) {
              this.modalRef = this.modalService.create({
                nzTitle: '启用',
                nzWidth: MODAL_NORMAL_SIZE,
                nzContent: PluginConfigComponent,
                nzComponentParams: {
                  ...params,
                  refreshPage: this.getPluginDetail,
                  pluginId: this.pluginId,
                  validateForm: this.fb.group({
                    name: [params.name, [Validators.required]],
                    server: [params.server, [Validators.required]]
                  }),
                  closeModal: () => { this.modalRef?.close() }
                },
                nzOkText: '确定',
                nzCancelText: '取消',
                nzOnOk: (component: PluginConfigComponent) => {
                  component.enablePlugin()
                  return false
                }
              })
            } else {
              this.api
                .post(
                  'system/plugin/enable',
                  { name: params.name },
                  { id: this.pluginId }
                )
                .subscribe((resp: EmptyHttpResponse) => {
                  if (resp.code === 0) {
                    this.message.success(resp.msg || '启用插件成功')
                    const subscription = this.navigationService
                      .getMenuList()
                      .subscribe(() => {
                        subscription.unsubscribe()
                      })
                    this.getPluginDetail()
                    this.modalRef?.close()
                  }
                })
            }
          }
        }
      )
  }

  findConfigValue (valueList:Array<PluginInstallConfigData>, key:string) {
    for (const item of valueList) {
      if (item.name === key) {
        return item.value
      }
    }
    return undefined
  }

  disablePluginModal (deletePlugin: boolean) {
    this.modalRef = this.modalService.create({
      nzTitle: deletePlugin ? '卸载' : '停用',
      nzContent: `该插件${
        deletePlugin ? '卸载后将无法找回' : '停用后将无法再使用'
      }，请确认是否要${deletePlugin ? '卸载' : '停用'}？`,
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkText: '确定',
      nzOkDanger: true,
      nzCancelText: '取消',
      nzOnOk: () => {
        if (deletePlugin) {
          this.deletePlugin()
        } else {
          this.disablePlugin()
        }
        return false
      }
    })
  }

  // 禁用插件
  disablePlugin () {
    this.api
      .post('system/plugin/disable', {}, { id: this.pluginId })
      .subscribe((resp: EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '禁用成功')
          const subscription = this.navigationService
            .getMenuList()
            .subscribe(() => {
              subscription.unsubscribe()
            })
          this.getPluginDetail()
          this.modalRef?.close()
        }
      })
  }

  // 卸载插件
  deletePlugin () {
    this.api
      .post('system/plugin/uninstall', {}, { id: this.pluginId })
      .subscribe((resp: EmptyHttpResponse) => {
        if (resp.code === 0) {
          this.message.success(resp.msg || '卸载成功')
          const subscription = this.navigationService
            .getMenuList()
            .subscribe(() => {
              subscription.unsubscribe()
            })
          this.router.navigate(['/', 'module-plugin', 'group', 'list'])
          this.modalRef?.close()
        }
      })
  }
}
