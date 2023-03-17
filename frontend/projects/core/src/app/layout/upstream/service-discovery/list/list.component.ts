/*
 * @Author:
 * @Date: 2022-08-14 22:48:39
 * @LastEditors:
 * @LastEditTime: 2022-08-22 00:16:26
 * @FilePath: /apinto/src/app/layout/upstream/service-discovery-list/service-discovery-list.component.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
/* eslint-disable no-useless-constructor */
import { Component, OnInit } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { AppConfigService } from 'projects/core/src/app/service/app-config.service'

@Component({
  selector: 'eo-ng-service-discovery-list',
  templateUrl: './list.component.html',
  styles: [
  ]
})
export class ServiceDiscoveryListComponent implements OnInit {
  constructor (private message: EoNgFeedbackMessageService,
               private modalService:EoNgFeedbackModalService,
               private api:ApiService,
               private router:Router,
               private appConfigService:AppConfigService) {
    this.appConfigService.reqFlashBreadcrumb([{ title: '服务发现', routerLink: 'upstream/serv-discovery' }])
  }

  nzDisabled:boolean = false
  serviceName:string = ''
  serviceNameForSear:string = ''
  servicesList : Array<object> = []
  servicesTableHeadName: Array<object> = [
    { title: '服务名称' },
    { title: '服务类型' },
    { title: '描述' },
    { title: '更新时间' },
    {
      title: '操作',
      right: true
    }
  ]

  servicesTableBody: Array<any> =[
    { key: 'name' },
    { key: 'driver' },
    { key: 'desc' },
    { key: 'update_time' },
    {
      type: 'btn',
      right: true,
      btns: [
        {
          title: '查看',
          click: (item:any) => {
            this.router.navigate(['/', 'upstream', 'serv-discovery', 'content', item.data.name])
          }
        },
        {
          title: '删除',
          disabledFn: (data:any, item:any) => {
            return this.nzDisabled || !item.data.is_delete
          },
          click: (item:any) => {
            this.delete(item.data)
          }
        }
      ]
    }

  ]

  ngOnInit (): void {
    this.getServicesList()
  }

  getServicesList () {
    this.api.get('discoveries', { name: this.serviceNameForSear }).subscribe(resp => {
      if (resp.code === 0) {
        this.servicesList = resp.data.discoveries
        this.serviceName = this.serviceNameForSear
      } else {
        this.message.error(resp.msg || '获取列表数据失败!')
      }
    })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  serviceTableClick = (item:any) => {
    this.router.navigate(['/', 'upstream', 'serv-discovery', 'content', item.data.name])
  }

  addService () {
    this.router.navigate(['/', 'upstream', 'serv-discovery', 'create'])
  }

  delete (item:any, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkDanger: true,
      nzOnOk: () => {
        this.deleteDiscovery(item.name)
      }
    })
  }

  deleteDiscovery (name:string) {
    this.api.delete('discovery', { name: name }).subscribe(resp => {
      if (resp.code === 0) {
        this.getServicesList()
        this.message.success(resp.msg || '删除成功!', { nzDuration: 1000 })
      } else {
        this.message.error(resp.msg || '删除失败!')
      }
    })
  }
}
