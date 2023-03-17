import { Component, OnInit } from '@angular/core'

@Component({
  selector: 'eo-ng-audit-log-detail',
  template: `
    <div class="drawer-list-content">
      <eo-ng-apinto-table
        class="drawer-table auth-log-table"
        [nzTbody]="auditLogDetailTableBody"
        [nzThead]="auditLogDetailTableHeadName"
        [nzData]="auditLogDetail"
        [nzNoScroll]="true"
      >
      </eo-ng-apinto-table>
    </div>
  `,
  styles: [
  ]
})
export class AuditLogDetailComponent implements OnInit {
  auditLogDetail:Array<{attr:string, value:string}> = []
  auditLogDetailTableHeadName: Array<object> = [
    {
      title: '属性',
      resizeable: true
    },
    { title: '配置' }
  ]

  auditLogDetailTableBody: Array<any> =[
    { key: 'attr' },
    {
      key: 'value',
      styleFn: (item:any) => {
        if (item.attr === '请求内容') {
          return 'white-space: pre-line;word-wrap:break-word; word-break:break-all'
        } else {
          return 'white-space: unset;word-wrap:break-word; word-break:break-all'
        }
      },
      ellipsis: false
    }
  ]

  ngOnInit (): void {
  }
}
