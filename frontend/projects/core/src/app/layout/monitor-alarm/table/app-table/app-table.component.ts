/* eslint-disable no-useless-constructor */
/* eslint-disable dot-notation */
/* eslint-disable camelcase */
import { Component, EventEmitter, Input, OnInit, Output, SimpleChanges } from '@angular/core'
import { MonitorAppData } from '../../area/total/total.component'
import { TableQueryData } from '../api-table/api-table.component'

@Component({
  selector: 'eo-ng-monitor-alarm-app-table',
  template: `
  <div *ngIf="showSearch"
        class="group-search-large inside-tab flex justify-between">
        <p  style="{font-weight:'500'; color:'#333333'}" *ngIf="title">{{ title }}</p>
    <div class="inline-block">
        <div class="inline-block">
        <eo-ng-search-input-group  class="ml-label" [eoInputVal]="queryData.keyword" (eoClick)="queryData.keyword = '';changeDisplayData(queryData.keyword)">
    <input
        class="search"
        type="text"
        eo-ng-input
        placeholder="请输入应用名称进行搜索"
        [(ngModel)]="queryData.keyword"
        (keyup.enter)="changeDisplayData(queryData.keyword)"
      />
      </eo-ng-search-input-group>
  </div>
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
      [nzTbody]="appTableBody"
      [nzThead]="appTableHeadName"
      [nzData]="appTableDisplayList"
      [nzTrClick]="checkTableDetail"
      [nzMaxOperatorButton]="2"
      [nzTotal]="appTableDisplayList.length"
      [nzPageIndex]="queryData.page_num"
      [nzPageSize]="queryData.page_size"
      [nzFrontPagination]="showPagination"
      [nzShowPagination]="showPagination"
      [nzTrClick]="tableClick"
      [nzMonitorDT]="nzMonitorDT"
    >
    </eo-ng-apinto-table>
  </div>
  `,
  styles: [
  ]
})
export class AppTableComponent implements OnInit {
  @Input() nzMonitorDT:boolean =false
  @Input() type:string = 'appTop10'
  @Input() title:string = ''
  @Input() appTableList:Array<any> = []
  @Input() showPagination:boolean = false
  @Input() showSearch:boolean = false
  @Output() btnClick:EventEmitter<any> = new EventEmitter()
  @Output() tableConfigChange:EventEmitter<any> = new EventEmitter()
  appTableDisplayList:Array<any> = []

  queryData:TableQueryData = {
    page_num: 1,
    page_size: 20,
    total: 0,
    keyword: ''
  }

  appConfig:{ app_id:boolean, request_total :boolean, request_success :boolean,
    request_rate:boolean, proxy_total:boolean, proxy_success :boolean, proxy_rate :boolean, status_fail:boolean,
    avg_resp :boolean, max_resp :boolean, min_resp:boolean, avg_traffic :boolean, max_traffic:boolean, min_traffic:boolean, [key:string]:boolean} = {
      app_id: true,
      request_total: true,
      request_success: true,
      request_rate: true,
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

  appTableBody:Array<any> =
  [{
    key: 'app_name',
    left: true,
    styleFn: (item:any) => {
      if (item.is_red) {
        return 'color:red'
      }
      return ''
    }
  },
  {
    key: 'app_id',
    showFn: () => {
      return this.appConfig.app_id === undefined ? true : this.appConfig.app_id
    }
  },
  {
    key: 'request_total',
    showFn: () => {
      return this.appConfig.request_total === undefined ? true : this.appConfig.request_total
    }
  },
  {
    key: 'request_success',
    showFn: () => {
      return this.appConfig.request_success === undefined ? true : this.appConfig.request_success
    }
  },
  {
    key: 'request_rate',
    keySuffix: '%',
    showFn: () => {
      return this.appConfig.request_rate === undefined ? true : this.appConfig.request_rate
    }
  },
  {
    key: 'proxy_total',
    showFn: () => {
      return this.appConfig.proxy_total === undefined ? true : this.appConfig.proxy_total
    }
  },
  {
    key: 'proxy_success',
    showFn: () => {
      return this.appConfig.proxy_success === undefined ? true : this.appConfig.proxy_success
    }
  },
  {
    key: 'proxy_rate',
    keySuffix: '%',
    showFn: () => {
      return this.appConfig.proxy_rate === undefined ? true : this.appConfig.proxy_rate
    }
  },
  {
    key: 'status_fail',
    showFn: () => {
      return this.appConfig.status_fail === undefined ? true : this.appConfig.status_fail
    }
  },
  {
    key: 'avg_resp',
    showFn: () => {
      return this.appConfig.avg_resp === undefined ? true : this.appConfig.avg_resp
    }
  },
  {
    key: 'max_resp',
    showFn: () => {
      return this.appConfig.max_resp === undefined ? true : this.appConfig.max_resp
    }
  },
  {
    key: 'min_resp',
    showFn: () => {
      return this.appConfig.min_resp === undefined ? true : this.appConfig.min_resp
    }
  },
  {
    key: 'avg_traffic',
    showFn: () => {
      return this.appConfig.avg_traffic === undefined ? true : this.appConfig.avg_traffic
    }
  },
  {
    key: 'max_traffic',
    showFn: () => {
      return this.appConfig.max_traffic === undefined ? true : this.appConfig.max_traffic
    }
  },
  {
    key: 'min_traffic',
    showFn: () => {
      return this.appConfig.min_traffic === undefined ? true : this.appConfig.min_traffic
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

  appTableHeadName:Array<any> = [
    {
      title: '应用名称',
      width: 80,
      key: 'app_name',
      left: true
    },
    {
      title: '应用ID',
      key: 'app_id',
      width: 80,
      showFn: () => {
        return this.appConfig.app_id === undefined ? true : this.appConfig.app_id
      }
    },
    {
      title: '请求总数',
      key: 'request_total',
      width: 90,
      showSort: true,
      sortOrder: 'descend',
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.request_total - b.data.request_total,
      showFn: () => {
        return this.appConfig.request_total === undefined ? true : this.appConfig.request_total
      }
    },
    {
      title: '请求成功数',
      key: 'request_success',
      width: 100,
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.request_success - b.data.request_success,
      showFn: () => {
        return this.appConfig.request_success === undefined ? true : this.appConfig.request_success
      }
    },
    {
      title: '请求成功率',
      key: 'request_rate',
      width: 100,
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.request_rate - b.data.request_rate,
      showFn: () => {
        return this.appConfig.request_rate === undefined ? true : this.appConfig.request_rate
      }
    },
    {
      title: '转发总数',
      width: 100,
      key: 'proxy_total',
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.proxy_total - b.data.proxy_total,
      showFn: () => {
        return this.appConfig.proxy_total === undefined ? true : this.appConfig.proxy_total
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
        return this.appConfig.proxy_success === undefined ? true : this.appConfig.proxy_success
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
        return this.appConfig.proxy_rate === undefined ? true : this.appConfig.proxy_rate
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
        return this.appConfig.status_fail === undefined ? true : this.appConfig.status_fail
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
        return this.appConfig.avg_resp === undefined ? true : this.appConfig.avg_resp
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
        return this.appConfig.max_resp === undefined ? true : this.appConfig.max_resp
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
        return this.appConfig.min_resp === undefined ? true : this.appConfig.min_resp
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
        return this.appConfig.avg_traffic === undefined ? true : this.appConfig.avg_traffic
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
        return this.appConfig.max_traffic === undefined ? true : this.appConfig.max_traffic
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
        return this.appConfig.min_traffic === undefined ? true : this.appConfig.min_traffic
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
    this.tableConfigChange.emit({ thead: this.appTableHeadName, config: this.appConfig })
    this.refreshMenu()
  }

  ngOnChanges (changes:SimpleChanges) {
    if (changes['type'] && localStorage.getItem(this.type)) {
      const storageConfig = localStorage.getItem(this.type)?.split(',')
      const keyList = Object.keys(this.appConfig)
      for (const index in keyList) {
        this.appConfig[keyList[index]] = storageConfig?.indexOf(keyList[index]) !== -1
      }
      this.refreshMenu()
    }
    if (changes['appTableList'] && this.appTableList !== undefined) {
      this.appTableDisplayList = [...this.appTableList]
      this.queryData.total = this.appTableDisplayList.length
    }
  }

  tableClick = (item: any) => {
    this.goToDetail(item)
  }

  changeMenu (target:string, event:boolean, item:{title:string, key:string, checked:boolean}) {
    this.appConfig[item.key] = event
    this.tableConfigChange.emit({ thead: this.appTableHeadName, config: this.appConfig })
    this.storeConfig()
  }

  changeDisplayData (keyword:string) {
    if (!keyword) {
      this.appTableDisplayList = [...this.appTableList]
      this.queryData.total = this.appTableDisplayList.length
      return
    }
    if (this.appTableList?.length > 0) {
      this.appTableDisplayList = this.appTableList.filter((item:MonitorAppData) => (
        item.app_name.includes(keyword)
      ))
      this.queryData.total = this.appTableDisplayList.length
    }
  }

  storeConfig () {
    const keyList = Object.keys(this.appConfig)
    const res:Array<string> = []
    for (const index in keyList) {
      if (this.appConfig[keyList[index]]) {
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
    this.tableDropdownMenu = [
      {
        title: '应用ID',
        key: 'app_id',
        checked: this.appConfig.app_id
      },
      {
        title: '请求总数',
        key: 'request_total',
        checked: this.appConfig.request_total
      },
      {
        title: '请求成功数',
        key: 'request_success',
        checked: this.appConfig.request_success
      },
      {
        title: '请求成功率',
        key: 'request_rate',
        checked: this.appConfig.request_rate
      },
      {
        title: '转发总数',
        key: 'proxy_total',
        checked: this.appConfig.proxy_total
      },
      {
        title: '转发成功数',
        key: 'proxy_success',
        checked: this.appConfig.proxy_success
      },
      {
        title: '转发成功率',
        key: 'proxy_rate',
        checked: this.appConfig.proxy_rate
      },
      {
        title: '失败状态码数',
        key: 'status_fail',
        checked: this.appConfig.status_fail
      },
      {
        title: '平均响应时间(ms)',
        key: 'avg_resp',
        checked: this.appConfig.avg_resp
      },
      {
        title: '最大响应时间(ms)',
        key: 'max_resp',
        checked: this.appConfig.max_resp
      },
      {
        title: '最小响应时间(ms)',
        key: 'min_resp',
        checked: this.appConfig.min_resp
      },
      {
        title: '平均请求流量(KB)',
        key: 'avg_traffic',
        checked: this.appConfig.avg_traffic
      },
      {
        title: '最大请求流量(KB)',
        key: 'max_traffic',
        checked: this.appConfig.max_traffic
      },
      {
        title: '最小请求流量(KB)',
        key: 'min_traffic',
        checked: this.appConfig.min_traffic
      }]
  }
}
