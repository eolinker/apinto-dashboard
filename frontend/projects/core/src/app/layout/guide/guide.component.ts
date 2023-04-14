import { Component, OnInit } from '@angular/core'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { ApiService } from '../../service/api.service'
import { Router } from '@angular/router'

type StepItem = {name:'cluster'|'service'|'api'|'publishApi',
title:string,
desc:Array<string>,
status:'undo'|'doing'|'done'
img?:string
toDoUrl:string
doneUrl:string}

@Component({
  selector: 'eo-ng-guide',
  templateUrl: './guide.component.html',
  styles: [
    `
    :host ::ng-deep{
      height: 100%;
      width: 100%;
      display: block;
      background-color: #f5f7fa;
      overflow: hidden;
      overflow-y:auto;
    }`
  ]
})
export class GuideComponent implements OnInit {
  stepList:Array<StepItem> = [
    {
      name: 'cluster',
      title: '创建集群',
      desc: ['创建 Apinto 网关集群，集群用于承载网络流量'],
      status: 'undo',
      toDoUrl: 'deploy/cluster/create',
      doneUrl: 'deploy/cluster'
    },
    {
      name: 'service',
      title: '添加上游服务',
      desc: ['添加上游服务器或动态服务发现，接收网关节点转发的流量'],
      status: 'undo',
      toDoUrl: 'upstream/upstream/create',
      doneUrl: 'upstream/upstream'
    },
    {
      name: 'api',
      title: '添加 API',
      desc: ['添加需要网关转发的 API'],
      status: 'undo',
      toDoUrl: 'router/api/create',
      doneUrl: 'router/api'
    },
    {
      name: 'publishApi',
      title: '发布 API',
      desc: ['发布之后就可以通过 Apinto 安全高效地访问API啦！'],
      status: 'undo',
      toDoUrl: 'router/api/group/list',
      doneUrl: 'router/api/group/list'
    }
  ]

    tutorialsList:Array<{
      title:string,
      content:Array<
        {text:string, url:string}
      >}> = [
        {
          title: '安全防护',
          content: [
            { text: '为API设置鉴权/身份认证', url: '' },
            { text: 'xxx', url: '' },
            { text: 'xxx', url: '' }
          ]
        },
        {
          title: '安全防护',
          content: [
            { text: '为API设置鉴权/身份认证', url: '' },
            { text: 'xxx', url: '' },
            { text: 'xxx', url: '' }
          ]
        },
        {
          title: '安全防护',
          content: [
            { text: '为API设置鉴权/身份认证', url: '' },
            { text: 'xxx', url: '' },
            { text: 'xxx', url: '' }
          ]
        },
        {
          title: '安全防护',
          content: [
            { text: '为API设置鉴权/身份认证', url: '' },
            { text: 'xxx', url: '' },
            { text: 'xxx', url: '' }
          ]
        }
      ]

      btnLoading:boolean = false
      constructor (private appConfigService:EoNgNavigationService, private api:ApiService, private router:Router) {}
      ngOnInit (): void {
        this.appConfigService.reqFlashBreadcrumb([])
        this.getStepStatus()
      }

      getStepStatus () {
        this.btnLoading = false
        this.api.get('system/quick_step').subscribe((resp:{code:number, msg:string, data:{cluster:boolean, service:boolean, api:boolean, publishApi:boolean}}) => {
          this.btnLoading = false
          if (resp.code === 0) {
            for (let i = 0; i < this.stepList.length; i++) {
              this.stepList[i].status = resp.data[this.stepList[i].name] ? 'done' : (i > 0 && this.stepList[i - 1].status === 'done' ? 'doing' : 'undo')
            }
          }
        })
      }

      goToStep (step:StepItem) {
        console.log(step)
        if (step.status === 'done') {
          this.router.navigate([step.doneUrl])
        } else {
          this.router.navigate([step.toDoUrl])
        }
      }
}
