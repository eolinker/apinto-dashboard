/* eslint-disable dot-notation */
/* eslint-disable camelcase */
import { Component, EventEmitter, Input, OnInit, Output, SimpleChanges } from '@angular/core'
import { MonitorNodeData } from '../../area/total/total.component'
import { TableQueryData } from '../api-table/api-table.component'

@Component({
  selector: 'eo-ng-monitor-alarm-node-table',
  template: `
  <div *ngIf="showSearch"
        class="group-search-large inside-tab">
        <div class="inline-block">
        <eo-ng-search-input-group  class="ml-label" [eoInputVal]="queryData.keyword" (eoClick)="queryData.keyword = '';changeDisplayData(queryData.keyword)">
    <input
        class="search"
        type="text"
        eo-ng-input
        placeholder="请输入目标节点进行搜索"
        [(ngModel)]="queryData.keyword"
        (keyup.enter)="changeDisplayData(queryData.keyword)"
      /></eo-ng-search-input-group>
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
      style="cursor: pointer"
      [title]="dropdownTitleTpl"
      [itemTmp]="dropdownMenTpl"
      [menus]="tableDropdownMenu"
      trigger="click"
      overlayClassName="eo-test-app-dropdown"
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
        class="eo-monitor-table-dropdown"
        (ngModelChange)="changeMenu('table2', $event, item)"
        eo-ng-checkbox
        [(ngModel)]="item.checked"
        >{{ item.title }}</label
      >
    </ng-template>
    </div>
    <eo-ng-apinto-table
      [nzTbody]="nodeTableBody"
      [nzThead]="nodeTableHeadName"
      [nzData]="nodeDisplayList"
      [nzTrClick]="checkTableDetail"
      [nzMaxOperatorButton]="2"
      [nzTotal]="queryData.total"
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
export class NodeTableComponent implements OnInit {
  @Input() type:string = 'nodeTotal'
  @Input() nodeTableList:Array<any> = []
  @Input() showPagination:boolean = false
  @Input() showSearch:boolean = false
  @Input() nzMonitorDT:boolean =false
  @Output() btnClick:EventEmitter<any> = new EventEmitter()
  @Output() tableConfigChange:EventEmitter<any> = new EventEmitter()
  nodeDisplayList:Array<any> = []

  queryData:TableQueryData = {
    page_num: 1,
    page_size: 20,
    total: 0,
    keyword: ''
  }

  nodeConfig:{ proxy_total:boolean, proxy_success :boolean, proxy_rate :boolean, status_fail:boolean,
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

  nodeTableBody:Array<any> =
  [{
    key: 'addr',
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
      return this.nodeConfig.proxy_total === undefined ? true : this.nodeConfig.proxy_total
    }
  },
  {
    key: 'proxy_success',
    showFn: () => {
      return this.nodeConfig.proxy_success === undefined ? true : this.nodeConfig.proxy_success
    }
  },
  {
    key: 'proxy_rate',
    keySuffix: '%',
    showFn: () => {
      return this.nodeConfig.proxy_rate === undefined ? true : this.nodeConfig.proxy_rate
    }
  },
  {
    key: 'status_fail',
    showFn: () => {
      return this.nodeConfig.status_fail === undefined ? true : this.nodeConfig.status_fail
    }
  },
  {
    key: 'avg_resp',
    showFn: () => {
      return this.nodeConfig.avg_resp === undefined ? true : this.nodeConfig.avg_resp
    }
  },
  {
    key: 'max_resp',
    showFn: () => {
      return this.nodeConfig.max_resp === undefined ? true : this.nodeConfig.max_resp
    }
  },
  {
    key: 'min_resp',
    showFn: () => {
      return this.nodeConfig.min_resp === undefined ? true : this.nodeConfig.min_resp
    }
  },
  {
    key: 'avg_traffic',
    showFn: () => {
      return this.nodeConfig.avg_traffic === undefined ? true : this.nodeConfig.avg_traffic
    }
  },
  {
    key: 'max_traffic',
    showFn: () => {
      return this.nodeConfig.max_traffic === undefined ? true : this.nodeConfig.max_traffic
    }
  },
  {
    key: 'min_traffic',
    showFn: () => {
      return this.nodeConfig.min_traffic === undefined ? true : this.nodeConfig.min_traffic
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

  nodeTableHeadName:Array<any> = [
    {
      title: '目标节点',
      key: 'addr',
      left: true,
      width: 80
    },
    {
      title: '转发总数',
      key: 'proxy_total',
      showSort: true,
      sortOrder: 'descend',
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.proxy_total - b.data.proxy_total,
      width: 90,
      showFn: () => {
        return this.nodeConfig.proxy_total === undefined ? true : this.nodeConfig.proxy_total
      }
    },
    {
      title: '转发成功数',
      width: 100,
      key: 'proxy_success',
      sortOrder: null,
      sortPriority: false,
      sortFn: (a: any, b: any) => a.data.proxy_success - b.data.proxy_success,
      showSort: true,
      showFn: () => {
        return this.nodeConfig.proxy_success === undefined ? true : this.nodeConfig.proxy_success
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
        return this.nodeConfig.proxy_rate === undefined ? true : this.nodeConfig.proxy_rate
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
        return this.nodeConfig.status_fail === undefined ? true : this.nodeConfig.status_fail
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
        return this.nodeConfig.avg_resp === undefined ? true : this.nodeConfig.avg_resp
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
        return this.nodeConfig.max_resp === undefined ? true : this.nodeConfig.max_resp
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
        return this.nodeConfig.min_resp === undefined ? true : this.nodeConfig.min_resp
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
        return this.nodeConfig.avg_traffic === undefined ? true : this.nodeConfig.avg_traffic
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
        return this.nodeConfig.max_traffic === undefined ? true : this.nodeConfig.max_traffic
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
        return this.nodeConfig.min_traffic === undefined ? true : this.nodeConfig.min_traffic
      }
    },
    {
      title: '操作',
      right: true
    }]

  tableDropdownMenu:Array<{title:string, key:string, checked:boolean}> = [
  ]

  ngOnInit (): void {
    this.tableConfigChange.emit({ thead: this.nodeTableHeadName, config: this.nodeConfig })
    this.refreshMenu()
  }

  ngOnChanges (changes:SimpleChanges) {
    if (changes['type'] && localStorage.getItem(this.type)) {
      const storageConfig = localStorage.getItem(this.type)?.split(',')
      const keyList = Object.keys(this.nodeConfig)
      for (const index in keyList) {
        this.nodeConfig[keyList[index]] = storageConfig?.indexOf(keyList[index]) !== -1
      }
      this.refreshMenu()
    }
    if (changes['nodeTableList']) {
      this.nodeDisplayList = [...this.nodeTableList]
      this.queryData.total = this.nodeDisplayList.length
    }
  }

  tableClick = (item: any) => {
    this.goToDetail(item)
  }

  changeMenu (target:string, event:boolean, item:{title:string, key:string, checked:boolean}) {
    this.nodeConfig[item.key] = event
    this.tableConfigChange.emit({ thead: this.nodeTableHeadName, config: this.nodeConfig })
    this.storeConfig()
  }

  changeDisplayData (keyword:string) {
    if (!keyword) {
      this.nodeDisplayList = [...this.nodeTableList]
      this.queryData.total = this.nodeDisplayList.length
      return
    }
    if (this.nodeTableList?.length > 0) {
      this.nodeDisplayList = this.nodeTableList.filter((item:MonitorNodeData) => (
        item['addr'].includes(keyword)
      ))
      this.queryData.total = this.nodeDisplayList.length
    }
  }

  storeConfig () {
    const keyList = Object.keys(this.nodeConfig)
    const res:Array<string> = []
    for (const index in keyList) {
      if (this.nodeConfig[keyList[index]]) {
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
        title: '转发总数',
        key: 'proxy_total',
        checked: this.nodeConfig.proxy_total
      },
      {
        title: '转发成功数',
        key: 'proxy_success',
        checked: this.nodeConfig.proxy_success
      },
      {
        title: '转发成功率',
        key: 'proxy_rate',
        checked: this.nodeConfig.proxy_rate

      },
      {
        title: '失败状态码数',
        key: 'status_fail',
        checked: this.nodeConfig.status_fail
      },
      {
        title: '平均响应时间(ms)',
        key: 'avg_resp',
        checked: this.nodeConfig.avg_resp
      },
      {
        title: '最大响应时间(ms)',
        key: 'max_resp',
        checked: this.nodeConfig.max_resp
      },
      {
        title: '最小响应时间(ms)',
        key: 'min_resp',
        checked: this.nodeConfig.min_resp
      },
      {
        title: '平均请求流量(KB)',
        key: 'avg_traffic',
        checked: this.nodeConfig.avg_traffic
      },
      {
        title: '最大请求流量(KB)',
        key: 'max_traffic',
        checked: this.nodeConfig.max_traffic
      },
      {
        title: '最小请求流量(KB)',
        key: 'min_traffic',
        checked: this.nodeConfig.min_traffic
      }]
  }
}
