/* eslint-disable no-useless-constructor */
/* eslint-disable camelcase */
/* eslint-disable dot-notation */
import { Component, EventEmitter, Input, OnInit, Output, SimpleChanges } from '@angular/core'
import { MonitorUpstreamData } from '../../area/total/total.component'
import { TableQueryData } from '../api-table/api-table.component'

@Component({
  selector: 'eo-ng-monitor-alarm-upstream-table',
  template: `

  <div *ngIf="showSearch"
        class="group-search-large inside-tab">
        <div class="inline-block">
        <eo-ng-search-input-group  class="ml-label" [eoInputVal]="queryData.keyword" (eoClick)="queryData.keyword = '';changeDisplayData(queryData.keyword)">
    <input
        class="search"
        type="text"
        eo-ng-input
        placeholder="请输入上游服务名称进行搜索"
        [(ngModel)]="queryData.keyword"
        (keyup.enter)="changeDisplayData(queryData.keyword)"
      />
    </eo-ng-search-input-group>
  </div>
  </div>
  <div class="monitor-table-block">
  <div
  class="eo-test-table-opr-btns"
  style="
    display: flex;
    margin-top: 6px;
    position: absolute;
    z-index: 4;
    right: 20px;
  "
  dir="rtl"
  >
    <eo-ng-dropdown
      style=" cursor: pointer"
      [title]="dropdownTitleTpl"
      [itemTmp]="dropdownMenTpl"
      [menus]="tableDropdownMenu"
      trigger="click"
      overlayClassName="eo-monitor-table-dropdown"
    >
    </eo-ng-dropdown>
    <ng-template #dropdownTitleTpl>
      <svg class="iconpark-icon">
        <use href="#setting"></use>
      </svg>
    </ng-template>
    <ng-template #dropdownMenTpl let-item="item">
      <label
      style="width:100%"
        (click)="$event.stopPropagation()"
        class="eo-test-app-checkbox"
        (ngModelChange)="changeMenu('table2', $event, item)"
        eo-ng-checkbox
        [(ngModel)]="item.checked"
        >{{ item.title }}</label
      >
    </ng-template>
    </div>
    <eo-ng-apinto-table
      [nzTbody]="upstreamTableBody"
      [nzThead]="upstreamTableHeadName"
      [nzData]="upstreamTableList"
      [nzTrClick]="checkTableDetail"
      [nzMaxOperatorButton]="2"
      [nzTotal]="queryData.total"
      [nzPageIndex]="queryData.page_num"
      [nzPageSize]="queryData.page_size"
      [nzFrontPagination]="showPagination"
      [nzShowPagination]="showPagination"
      [nzTrClick]="tableClick"
    >
    </eo-ng-apinto-table>
  </div>
  `,
  styles: [
  ]
})
export class UpstreamTableComponent implements OnInit {
  @Input() nzScrollY:number = 0
  @Input() type:string = 'upstreamTop10'
  @Input() upstreamTableList:Array<any> = []
  @Input() totalItemsSize:number = 0
  @Input() showPagination:boolean = false
  @Input() showSearch:boolean = false
  @Output() btnClick:EventEmitter<any> = new EventEmitter()
  @Output() tableConfigChange:EventEmitter<any> = new EventEmitter()
  upstreamDisplayTableList:Array<any> = []
  queryData:TableQueryData = {
    page_num: 1,
    page_size: 20,
    total: 0,
    keyword: ''
  }

  upstreamConfig:{proxy_total:boolean, proxy_success :boolean, proxy_rate :boolean, status_fail:boolean,
    avg_resp :boolean, max_resp :boolean, min_resp:boolean, avg_traffic :boolean, max_traffic:boolean, min_traffic:boolean, [key:string]:boolean} = {
      proxy_total: true,
      proxy_success: true,
      proxy_rate: true,
      status_fail: true,
      avg_resp: true,
      max_resp: true,
      min_resp: true,
      avg_traffic: true,
      max_traffic: true,
      min_traffic: true
    }

  upstreamTableBody:Array<any> =
  [{
    key: 'service_name',
    left: true,
    styleFn: (item:any) => {
      if (item.is_red) {
        return 'color:red'
      }
      return ''
    }
  },
  {
    key: 'proxy_total',
    showFn: () => {
      return this.upstreamConfig.proxy_total === undefined ? true : this.upstreamConfig.proxy_total
    }
  },
  {
    key: 'proxy_success',
    showFn: () => {
      return this.upstreamConfig.proxy_success === undefined ? true : this.upstreamConfig.proxy_success
    }
  },
  {
    key: 'proxy_rate',
    keySuffix: '%',
    showFn: () => {
      return this.upstreamConfig.proxy_rate === undefined ? true : this.upstreamConfig.proxy_rate
    }
  },
  {
    key: 'status_fail',
    showFn: () => {
      return this.upstreamConfig.status_fail === undefined ? true : this.upstreamConfig.status_fail
    }
  },
  {
    key: 'avg_resp',
    showFn: () => {
      return this.upstreamConfig.avg_resp === undefined ? true : this.upstreamConfig.avg_resp
    }
  },
  {
    key: 'max_resp',
    showFn: () => {
      return this.upstreamConfig.max_resp === undefined ? true : this.upstreamConfig.max_resp
    }
  },
  {
    key: 'min_resp',
    showFn: () => {
      return this.upstreamConfig.min_resp === undefined ? true : this.upstreamConfig.min_resp
    }
  },
  {
    key: 'avg_traffic',
    showFn: () => {
      return this.upstreamConfig.avg_traffic === undefined ? true : this.upstreamConfig.avg_traffic
    }
  },
  {
    key: 'max_traffic',
    showFn: () => {
      return this.upstreamConfig.max_traffic === undefined ? true : this.upstreamConfig.max_traffic
    }
  },
  {
    key: 'min_traffic',
    showFn: () => {
      return this.upstreamConfig.min_traffic === undefined ? true : this.upstreamConfig.min_traffic
    }
  },
  {
    type: 'btn',
    right: true,
    btns: [
      {
        title: '查看',
        click: (item:any) => {
          this.goToDetail(item)
        }
      }
    ]
  }]

  upstreamTableHeadName:Array<any> = [
    {
      title: '上游服务名称',
      key: 'service_name',
      left: true,
      width: 120
    },
    {
      title: '转发总数',
      width: 90,
      key: 'proxy_total',
      showSort: true,
      sortOrder: 'descend',
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.proxy_total - b.data.proxy_total,
      showFn: () => {
        return this.upstreamConfig.proxy_total === undefined ? true : this.upstreamConfig.proxy_total
      }
    },
    {
      title: '转发成功数',
      width: 100,
      key: 'proxy_success',
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.proxy_success - b.data.proxy_success,
      showFn: () => {
        return this.upstreamConfig.proxy_success === undefined ? true : this.upstreamConfig.proxy_success
      }
    },
    {
      title: '转发成功率',
      width: 100,
      key: 'proxy_rate',
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.proxy_rate - b.data.proxy_rate,
      showFn: () => {
        return this.upstreamConfig.proxy_rate === undefined ? true : this.upstreamConfig.proxy_rate
      }
    },
    {
      title: '失败状态码数',
      width: 120,
      key: 'status_fail',
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.status_fail - b.data.status_fail,
      showFn: () => {
        return this.upstreamConfig.status_fail === undefined ? true : this.upstreamConfig.status_fail
      }
    },
    {
      title: '平均响应时间(ms)',
      width: 148,
      key: 'avg_resp',
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.avg_resp - b.data.avg_resp,
      showFn: () => {
        return this.upstreamConfig.avg_resp === undefined ? true : this.upstreamConfig.avg_resp
      }
    },
    {
      title: '最大响应时间(ms)',
      width: 148,
      key: 'max_resp',
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.max_resp - b.data.max_resp,
      showFn: () => {
        return this.upstreamConfig.max_resp === undefined ? true : this.upstreamConfig.max_resp
      }
    },
    {
      title: '最小响应时间(ms)',
      width: 148,
      key: 'min_resp',
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.min_resp - b.data.min_resp,
      showFn: () => {
        return this.upstreamConfig.min_resp === undefined ? true : this.upstreamConfig.min_resp
      }
    },
    {
      title: '平均请求流量(KB)',
      width: 148,
      key: 'avg_traffic',
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.avg_traffic - b.data.avg_traffic,
      showFn: () => {
        return this.upstreamConfig.avg_traffic === undefined ? true : this.upstreamConfig.avg_traffic
      }
    },
    {
      title: '最大请求流量(KB)',
      width: 148,
      key: 'max_traffic',
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.max_traffic - b.data.max_traffic,
      showFn: () => {
        return this.upstreamConfig.max_traffic === undefined ? true : this.upstreamConfig.max_traffic
      }
    },
    {
      title: '最小请求流量(KB)',
      width: 148,
      key: 'min_traffic',
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.min_traffic - b.data.min_traffic,
      showFn: () => {
        return this.upstreamConfig.min_traffic === undefined ? true : this.upstreamConfig.min_traffic
      }
    },
    {
      title: '操作',
      right: true
    }]

  tableDropdownMenu:Array<{title:string, key:string, checked:boolean}> = [

  ]

  constructor () { }

  ngOnInit (): void {
    this.tableConfigChange.emit({ thead: this.upstreamTableHeadName, config: this.upstreamConfig })
    this.refreshMenu()
  }

  ngOnChanges (changes:SimpleChanges) {
    if (changes['type'] && localStorage.getItem(this.type)) {
      const storageConfig = localStorage.getItem(this.type)?.split(',')
      const keyList = Object.keys(this.upstreamConfig)
      for (const index in keyList) {
        this.upstreamConfig[keyList[index]] = storageConfig?.indexOf(keyList[index]) !== -1
      }
      this.refreshMenu()
    }
    if (changes['upstreamTableList']) {
      this.upstreamDisplayTableList = [...this.upstreamTableList]
      this.queryData.total = this.upstreamDisplayTableList.length
    }
  }

  tableClick = (item: any) => {
    this.goToDetail(item)
  }

  changeMenu (target:string, event:boolean, item:{title:string, key:string, checked:boolean}) {
    this.upstreamConfig[item.key] = event
    this.tableConfigChange.emit({ thead: this.upstreamTableHeadName, config: this.upstreamConfig })
    this.storeConfig()
  }

  changeDisplayData (keyword:string) {
    if (!keyword) {
      this.upstreamDisplayTableList = [...this.upstreamTableList]
      this.queryData.total = this.upstreamDisplayTableList.length
      return
    }
    if (this.upstreamTableList?.length > 0) {
      this.upstreamDisplayTableList = this.upstreamTableList.filter((item:MonitorUpstreamData) => (
        item.service_name.includes(keyword)
      ))
      this.queryData.total = this.upstreamDisplayTableList.length
    }
  }

  storeConfig () {
    const keyList = Object.keys(this.upstreamConfig)
    const res:Array<string> = []
    for (const index in keyList) {
      if (this.upstreamConfig[keyList[index]]) {
        res.push(keyList[index])
      }
    }
    localStorage.setItem(this.type, res.join(','))
  }

  checkTableDetail (value:any) {
    this.btnClick.emit(value)
  }

  goToDetail (value:any) {
    this.btnClick.emit(value)
  }

  refreshMenu () {
    this.tableDropdownMenu = [{
      title: '转发总数',
      key: 'proxy_total',
      checked: this.upstreamConfig.proxy_total
    },
    {
      title: '转发成功数',
      key: 'proxy_success',
      checked: this.upstreamConfig.proxy_success
    },
    {
      title: '转发成功率',
      key: 'proxy_rate',
      checked: this.upstreamConfig.proxy_rate

    },
    {
      title: '失败状态码数',
      key: 'status_fail',
      checked: this.upstreamConfig.status_fail
    },
    {
      title: '平均响应时间(ms)',
      key: 'avg_resp',
      checked: this.upstreamConfig.avg_resp
    },
    {
      title: '最大响应时间(ms)',
      key: 'max_resp',
      checked: this.upstreamConfig.max_resp
    },
    {
      title: '最小响应时间(ms)',
      key: 'min_resp',
      checked: this.upstreamConfig.min_resp
    },
    {
      title: '平均请求流量(KB)',
      key: 'avg_traffic',
      checked: this.upstreamConfig.avg_traffic
    },
    {
      title: '最大请求流量(KB)',
      key: 'max_traffic',
      checked: this.upstreamConfig.max_traffic
    },
    {
      title: '最小请求流量(KB)',
      key: 'min_traffic',
      checked: this.upstreamConfig.min_traffic
    }]
  }
}
