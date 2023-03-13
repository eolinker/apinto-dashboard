/* eslint-disable dot-notation */
/* eslint-disable camelcase */
import { Component, EventEmitter, Input, OnInit, Output, SimpleChanges } from '@angular/core'
import { MonitorApiData } from '../../area/total/total.component'
export interface TableQueryData {
  page_num:number,
  page_size:number,
  total:number,
  keyword:string,
  sort?:{key:string, val:boolean}
}

@Component({
  selector: 'eo-ng-monitor-alarm-api-table',
  template: `
  <div *ngIf="showSearch"
        class="group-search-large inside-tab">
    <p  style="{font-weight:'500'; color:'#333333'}" *ngIf="title">{{ title }}</p>
        <div class="inline-block">
        <eo-ng-search-input-group  class="ml-label" [eoInputVal]="queryData.keyword" (eoClick)="queryData.keyword = '';changeDisplayData(queryData.keyword)">
        <input
            class="search"
            type="text"
            eo-ng-input
            placeholder="请输入api名称进行搜索"
            [(ngModel)]="queryData.keyword"
            (keyup.enter)="changeDisplayData(queryData.keyword)"
          />
          </eo-ng-search-input-group>
      </div >
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
      [nzTbody]="apiTableBody"
      [nzThead]="apiTableHeadName"
      [nzData]="apiTableDisplayList"
      [nzTrClick]="checkTableDetail"
      [nzMaxOperatorButton]="2"
      [nzTotal]="queryData.total"
      [nzPageIndex]="queryData.page_num"
      [nzPageSize]="queryData.page_size"
      [nzFrontPagination]="showPagination"
      [nzShowPagination]="showPagination"
      [nzScroll]="{ x: '100%' }"
      [nzTrClick]="tableClick"
      [nzMonitorDT]="nzMonitorDT"
      >
    </eo-ng-apinto-table>
  </div>
  `,
  styles: [
  ]
})
export class ApiTableComponent implements OnInit {
  @Input() title:string = ''
  @Input() type:string = 'apiTop10'
  @Input() apiTableList:Array<any> = []
  @Input() showPagination:boolean = false
  @Input() showSearch:boolean = false
  @Output() btnClick:EventEmitter<any> = new EventEmitter()
  @Output() tableConfigChange:EventEmitter<any> = new EventEmitter()
  @Input() nzMonitorDT:boolean =false

  apiTableDisplayList:Array<any> = []

   queryData:TableQueryData = {
     page_num: 1,
     page_size: 20,
     total: 0,
     keyword: ''
   }

  apiConfig:{ path :boolean, request_total :boolean, request_success :boolean,
    request_rate:boolean, proxy_total:boolean, proxy_success :boolean, proxy_rate :boolean, status_fail:boolean,
    avg_resp :boolean, max_resp :boolean, min_resp:boolean, avg_traffic :boolean, max_traffic:boolean, min_traffic:boolean, [key:string]:boolean} = {
      path: true,
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

  apiTableBody:Array<any> =
  [{
    key: 'api_name',
    left: true,
    styleFn: (item:any) => {
      if (item.is_red) {
        return 'color:red'
      }
      return ''
    }
  },
  {
    key: 'path',
    showFn: () => {
      return this.apiConfig.path === undefined ? true : this.apiConfig.path
    }
  },
  {
    key: 'request_total',
    showFn: () => {
      return this.apiConfig.request_total === undefined ? true : this.apiConfig.request_total
    }
  },
  {
    key: 'request_success',
    showFn: () => {
      return this.apiConfig.request_success === undefined ? true : this.apiConfig.request_success
    }
  },
  {
    key: 'request_rate',
    keySuffix: '%',
    showFn: () => {
      return this.apiConfig.request_rate === undefined ? true : this.apiConfig.request_rate
    }
  },
  {
    key: 'proxy_total',
    showFn: () => {
      return this.apiConfig.proxy_total === undefined ? true : this.apiConfig.proxy_total
    }
  },
  {
    key: 'proxy_success',
    showFn: () => {
      return this.apiConfig.proxy_success === undefined ? true : this.apiConfig.proxy_success
    }
  },
  {
    key: 'proxy_rate',
    keySuffix: '%',
    showFn: () => {
      return this.apiConfig.proxy_rate === undefined ? true : this.apiConfig.proxy_rate
    }
  },
  {
    key: 'status_fail',
    showFn: () => {
      return this.apiConfig.status_fail === undefined ? true : this.apiConfig.status_fail
    }
  },
  {
    key: 'avg_resp',
    showFn: () => {
      return this.apiConfig.avg_resp === undefined ? true : this.apiConfig.avg_resp
    }
  },
  {
    key: 'max_resp',
    showFn: () => {
      return this.apiConfig.max_resp === undefined ? true : this.apiConfig.max_resp
    }
  },
  {
    key: 'min_resp',
    showFn: () => {
      return this.apiConfig.min_resp === undefined ? true : this.apiConfig.min_resp
    }
  },
  {
    key: 'avg_traffic',
    showFn: () => {
      return this.apiConfig.avg_traffic === undefined ? true : this.apiConfig.avg_traffic
    }
  },
  {
    key: 'max_traffic',
    showFn: () => {
      return this.apiConfig.max_traffic === undefined ? true : this.apiConfig.max_traffic
    }
  },
  {
    key: 'min_traffic',
    showFn: () => {
      return this.apiConfig.min_traffic === undefined ? true : this.apiConfig.min_traffic
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

  apiTableHeadName:Array<any> = [
    {
      title: 'API名称',
      key: 'api_name',
      width: 86,
      left: true
    },
    {
      title: '请求路径',
      key: 'path',
      width: 80,
      showFn: () => {
        return this.apiConfig.path === undefined ? true : this.apiConfig.path
      }
    },
    {
      title: '请求总数',
      key: 'request_total',
      showSort: true,
      sortOrder: 'descend',
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.request_total - b.data.request_total,
      width: 90,
      showFn: () => {
        return this.apiConfig.request_total === undefined ? true : this.apiConfig.request_total
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
        return this.apiConfig.request_success === undefined ? true : this.apiConfig.request_success
      }
    },
    {
      title: '请求成功率',
      key: 'request_rate',
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.request_rate - b.data.request_rate,
      width: 100,
      showFn: () => {
        return this.apiConfig.request_rate === undefined ? true : this.apiConfig.request_rate
      }
    },
    {
      title: '转发总数',
      key: 'proxy_total',
      width: 90,
      showSort: true,
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.proxy_total - b.data.proxy_total,
      showFn: () => {
        return this.apiConfig.proxy_total === undefined ? true : this.apiConfig.proxy_total
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
        return this.apiConfig.proxy_success === undefined ? true : this.apiConfig.proxy_success
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
        return this.apiConfig.proxy_rate === undefined ? true : this.apiConfig.proxy_rate
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
        return this.apiConfig.status_fail === undefined ? true : this.apiConfig.status_fail
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
        return this.apiConfig.avg_resp === undefined ? true : this.apiConfig.avg_resp
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
        return this.apiConfig.max_resp === undefined ? true : this.apiConfig.max_resp
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
        return this.apiConfig.min_resp === undefined ? true : this.apiConfig.min_resp
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
        return this.apiConfig.avg_traffic === undefined ? true : this.apiConfig.avg_traffic
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
        return this.apiConfig.max_traffic === undefined ? true : this.apiConfig.max_traffic
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
        return this.apiConfig.min_traffic === undefined ? true : this.apiConfig.min_traffic
      }
    },
    {
      title: '操作',
      right: true
    }]

  tableDropdownMenu:Array<{title:string, key:string, checked:boolean, showFn?:any}> = []

  ngOnInit (): void {
    this.tableConfigChange.emit({ thead: this.apiTableHeadName, config: this.apiConfig })
    this.refreshMenu()
  }

  ngOnChanges (changes:SimpleChanges) {
    if (changes['type']) {
      if (localStorage.getItem(this.type)) {
        const storageConfig = localStorage.getItem(this.type)?.split(',')
        const keyList = Object.keys(this.apiConfig)
        for (const index in keyList) {
          this.apiConfig[keyList[index]] = storageConfig?.indexOf(keyList[index]) !== -1
        }
        this.refreshMenu()
      }
    }
    if (changes['apiTableList']) {
      this.apiTableDisplayList = [...this.apiTableList]
      this.queryData.total = this.apiTableDisplayList.length
    }
  }

  tableClick = (item: any) => {
    this.goToDetail(item)
  }

  changeMenu (target:string, event:boolean, item:{title:string, key:string, checked:boolean}) {
    this.apiConfig[item.key] = event
    this.tableConfigChange.emit({ thead: this.apiTableHeadName, config: this.apiConfig })
    this.storeConfig()
  }

  changeDisplayData (keyword:string) {
    if (!keyword) {
      this.apiTableDisplayList = [...this.apiTableList]
      this.queryData.total = this.apiTableDisplayList.length
      return
    }
    if (this.apiTableList?.length > 0) {
      this.apiTableDisplayList = this.apiTableList.filter((item:MonitorApiData) => (
        item.api_name.includes(keyword)
      ))
      this.queryData.total = this.apiTableDisplayList.length
    }
  }

  storeConfig () {
    const keyList = Object.keys(this.apiConfig)
    const res:Array<string> = []
    for (const index in keyList) {
      if (this.apiConfig[keyList[index]]) {
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
        title: '请求路径',
        key: 'path',
        checked: this.apiConfig.path
      },
      {
        title: '请求总数',
        key: 'request_total',
        checked: this.apiConfig.request_total
      },
      {
        title: '请求成功数',
        key: 'request_success',
        checked: this.apiConfig.request_success
      },
      {
        title: '请求成功率',
        key: 'request_rate',
        checked: this.apiConfig.request_rate
      },
      {
        title: '转发总数',
        key: 'proxy_total',
        checked: this.apiConfig.proxy_total
      },
      {
        title: '转发成功数',
        key: 'proxy_success',
        checked: this.apiConfig.proxy_success
      },
      {
        title: '转发成功率',
        key: 'proxy_rate',
        checked: this.apiConfig.proxy_rate

      },
      {
        title: '失败状态码数',
        key: 'status_fail',
        checked: this.apiConfig.status_fail
      },
      {
        title: '平均响应时间(ms)',
        key: 'avg_resp',
        checked: this.apiConfig.avg_resp
      },
      {
        title: '最大响应时间(ms)',
        key: 'max_resp',
        checked: this.apiConfig.max_resp
      },
      {
        title: '最小响应时间(ms)',
        key: 'min_resp',
        checked: this.apiConfig.min_resp
      },
      {
        title: '平均请求流量(KB)',
        key: 'avg_traffic',
        checked: this.apiConfig.avg_traffic
      },
      {
        title: '最大请求流量(KB)',
        key: 'max_traffic',
        checked: this.apiConfig.max_traffic
      },
      {
        title: '最小请求流量(KB)',
        key: 'min_traffic',
        checked: this.apiConfig.min_traffic
      }
    ]
  }
}
