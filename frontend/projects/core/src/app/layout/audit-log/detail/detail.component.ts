import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { ApiService } from '../../../service/api.service'
import { auditLogDetailTableBody, auditLogDetailTableHeadName } from '../types/conf'
import { AuditLogDetail } from '../types/types'

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
    <ng-template #detailTdTpl let-item="item">
    <div class="leading-[22px] flex items-center">
    <span
    class=" mr-[8px] default-tpl-td-span whitespace-pre-wrap break-all break-words"
    eoNgFeedbackTooltip
    [nzTooltipTitle]="item.value"
    [nzTooltipVisible]="false"
    [nzTooltipTrigger]="'hover'"
    [nzTooltipOverlayClassName]="item.attr === '请求内容' ? 'tooltip-json' : ''"
    [ngClass]="{'break-all':item.attr === '请求内容' ,'whitespace-pre-wrap':item.attr === '请求内容','break-words':item.attr === '请求内容'}"
    >{{ item.value}}
  </span>
    <span
      class="h-[20px] leading-[22px] text-[12px] opacity-0 cursor-pointer inline-flex items-center"
      eo-copy
      nzType="primary"
      [copyText]="item.value"
      (copyCallback)="copyCallback()"
      (click)="$event.stopPropagation()"
    >
      <svg class="iconpark-icon">
        <use href="#copy"></use>
      </svg>
    </span>
  </div>
    </ng-template>
  `,
  styles: [
  ]
})
export class AuditLogDetailComponent implements OnInit {
  @ViewChild('detailTdTpl') detailTdTpl:TemplateRef<any>|undefined
  @Input() auditLogId:string = ''
  auditLogDetail:AuditLogDetail[] = []
  auditLogDetailTableHeadName: THEAD_TYPE[] = [...auditLogDetailTableHeadName]
  auditLogDetailTableBody: TBODY_TYPE[] =[...auditLogDetailTableBody]

  constructor (private message: EoNgFeedbackMessageService,
    private api:ApiService) {}

  ngOnInit (): void {
    this.getLogDetail()
  }

  ngAfterViewInit ():void {
    this.auditLogDetailTableBody[1].title = this.detailTdTpl
  }

  // 接口返回成功才打开弹窗
  getLogDetail ():void {
    this.api.get('audit-log', { logId: this.auditLogId })
      .subscribe((resp:{code:number, data:{args:AuditLogDetail[]}, msg:string}) => {
        if (resp.code === 0) {
          this.auditLogDetail = resp.data.args
          for (const index in this.auditLogDetail) {
            if (this.auditLogDetail[index].attr === '请求内容') {
              this.auditLogDetail[index].value = JSON.stringify(JSON.parse(this.auditLogDetail[index].value), null, 4)
            }
          }
        }
      })
  }

  copyCallback () {
    this.message.success('复制成功', {
      nzDuration: 1000
    })
  }
}
