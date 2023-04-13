import { Component, OnInit } from '@angular/core'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'

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
  stepList:Array<
    {title:string,
    desc:Array<string>,
    status:'done'|'doing'|'undo',
    img?:string}> = [
      {
        title: '创建集群',
        desc: ['创建 Apinto 网关集群，集群用于承载网络流量'],
        status: 'done'
      },
      {
        title: '添加上游服务',
        desc: ['添加上游服务器或动态服务发现，接收网关节点转发的流量'],
        status: 'doing'
      },
      {
        title: '添加 API',
        desc: ['添加需要网关转发的 API'],
        status: 'undo'
      },
      {
        title: '发布 API',
        desc: ['发布之后就可以通过 Apinto 安全高效地访问API啦！'],
        status: 'undo'
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

    constructor (private appConfigService:EoNgNavigationService) {}
    ngOnInit (): void {
      this.appConfigService.reqFlashBreadcrumb([])
    }
}
