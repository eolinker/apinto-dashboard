import { ChangeDetectorRef, Component, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { RouterService } from '../../../router.service'
import { ApiData, ApiPublishItem } from '../../../types/types'
import { BaseInfoService } from 'projects/core/src/app/service/base-info.service'

@Component({
  selector: 'eo-ng-api-publish',
  templateUrl: './publish.component.html',
  styles: [
  ]
})
export class ApiPublishComponent implements OnInit {
  @ViewChild('clusterStatusTpl', { read: TemplateRef, static: true }) clusterStatusTpl: TemplateRef<any> | undefined
  apiInfo: ApiData | undefined
  apiUuid:string = ''
  publishTableBody:TBODY_TYPE[] = []
  publishTableHeadName:THEAD_TYPE[] = [...this.service.createApiPublishThead(this)]
  publishList:ApiPublishItem[] = []
  selectedNum:number = 0
  selectedClusters:Array<string> = []
  moduleName:string = ''
  closeModal:any
  nzDisabled:boolean = false
  getApisData:any
  unpublishMsg:string = ''
  returnToSdk:Function|undefined
  constructor (
    private message: EoNgFeedbackMessageService,
    private service:RouterService,
    private api:ApiService,
    private baseInfo:BaseInfoService,
    private cdref:ChangeDetectorRef) {}

  disabledEdit (value:any) {
    this.nzDisabled = value
  }

  ngOnInit (): void {
    this.getPublishList()
  }

  ngAfterViewInit () {
    this.publishTableBody = [...this.service.createApiPublishTbody(this)]
    this.cdref.detectChanges()
  }

  getPublishList () {
    this.api.get('router/online/info', { uuid: this.apiUuid }).subscribe((resp:{code:number, msg:string, data:{info:ApiData, clusters:ApiPublishItem[]}}) => {
      if (resp.code === 0) {
        this.apiInfo = resp.data.info
        if (this.apiInfo.scheme === 'websocket') {
          this.apiInfo.method = '-'
        } else if (!this.apiInfo.method.length) {
          this.apiInfo.method = 'ALL'
        }
        this.publishList = resp.data.clusters
      }
    })
  }

  tableClick = (item:any) => {
    item.checked = !item.checked
    item.data.checked = !item.data.checked
    this.checkSelectedCluster()
  }

  // 点击表头全选
  checkAll () {
    this.checkSelectedCluster()
  }

  // 点击单条数据
  clickData () {
    this.checkSelectedCluster()
  }

  checkSelectedCluster () {
    setTimeout(() => {
      this.selectedClusters = this.publishList.filter((item:any) => {
        return item.checked
      }).map((item) => {
        return item.title
      })
      this.selectedNum = this.selectedClusters.length
      this.publishList = [...this.publishList] // 表头的勾选状态需要重载数据才能刷新
    }, 0
    )
  }

  offline () {
    const cluster:Array<string> = this.publishList.filter((item) => {
      return item.checked
    }).map((item) => {
      return item.name
    })
    this.api.put('router/offline', { clusterNames: cluster }, { uuid: this.apiUuid }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.unpublishMsg = ''
        this.message.success(resp.msg)
        this.closeModal && this.closeModal()
        this.getApisData && this.getApisData()
      } else {
        this.unpublishMsg = '下线失败' + resp.data.router.map((item:any) => {
          return item.msg
        }).join('，')
      }
      this.returnToSdk && this.returnToSdk(resp)
    })
  }

  online () {
    const cluster:Array<string> = this.publishList.filter((item) => {
      return item.checked
    }).map((item) => {
      return item.name
    })
    this.api.put('router/online', { clusterNames: cluster }, { uuid: this.apiUuid }).subscribe((resp:any) => {
      if (resp.code === 0) {
        this.unpublishMsg = ''
        this.message.success(resp.msg)
        this.closeModal && this.closeModal()
        this.getApisData && this.getApisData()
      } else {
        this.unpublishMsg = '上线失败：' + resp.data.router.map((item:any) => {
          return item.msg
        }).join('，')
      }
      this.returnToSdk && this.returnToSdk(resp)
    })
  }
}
