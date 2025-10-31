import { Component, ViewChild, TemplateRef } from '@angular/core'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE } from '../../constant/app.config'
import { EoNgLogRetrievalTailComponent } from './tail-log/tail-log.component'
import { ApiService } from '../../service/api.service'
import { ClusterEnum, ClusterEnumData } from '../../constant/type'
import { LogFileData, LogOutputData } from './type'
import { LogRetrievalTableBody, LogRetrievalTableHeadName } from './conf'
import { CascaderOption } from 'eo-ng-cascader'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import * as moment from 'moment'

@Component({
  selector: 'eo-ng-log-retrieval',
  templateUrl: './log-retrieval.component.html',
  styles: [
    `
    :host ::ng-deep{

      nz-spin > div{
        height:100%;
      }

      .ant-table-thead{
        height:0;
      }

      eo-ng-apinto-table:not(.arrayItem) {
        .ant-table-thead > tr > th{
          height:0px;
        }

        nz-table-inner-scroll table:first-child, eo-ng-apinto-table:not(.arrayItem) nz-table-inner-default table:first-child {
          border-top: none;
      }

    }

    .ant-collapse-borderless>.ant-collapse-item{
      background:var(--collapse-header-background-color);
      >.ant-collapse-content>.ant-collapse-content-box{
      padding-top:0px;
    }
  }

  .ant-collapse>.ant-collapse-item>.ant-collapse-header{
    background-color:transparent;

  }

      .ant-collapse>.ant-collapse-item{
        border-bottom:none;
      }

      .ant-collapse-item:last-child, .ant-collapse-item:last-child .ant-collapse-content{
        border-radius:var(--border-radius) !important;
      }
      .ant-collapse-content-box{
        padding:0;
      }

      nz-table-inner-scroll .ant-table-tbody > tr:last-child > td{
        border-bottom:none !important;
      }
    }`
  ]
})
export class LogRetrievalComponent {
  @ViewChild('tailLogModalFooterTpl') tailLogModalFooterTpl:TemplateRef<any> | undefined
 nzDisabled:boolean = false
  searchData:{
    cluster:Array<string>,
    node:number|string
  } = {
    cluster: [],
    node: ''
  }

  clusterList:CascaderOption[] = []
  nodeList:SelectOption[] = []

  accGroupList:LogOutputData[] = []
  accTableHeader:THEAD_TYPE[] = [...LogRetrievalTableHeadName]
  accTableBody:TBODY_TYPE[] = [...LogRetrievalTableBody]
  accGroupLoading:boolean = false
  modalRef:NzModalRef|undefined
  start:boolean = false

  constructor (private api:ApiService, private modalService:EoNgFeedbackModalService, private navigationService:EoNgNavigationService) {

  }

  disabledEdit ($event:any) {
    this.nzDisabled = $event
  }

  ngOnInit () {
    this.navigationService.reqFlashBreadcrumb([{ title: '日志检索' }])
    this.getClusterList(true)
    for (const btn of this.accTableBody[this.accTableBody.length - 1].btns) {
      if (btn.title === '下载') {
        btn.click = (item:{data:LogFileData}) => {
          this.downloadLog(item.data)
        }
        return
      }
    }
  }

  getClusterList (init?:boolean) {
    this.api.get('cluster/enum').subscribe((resp: {code:number, data:{ envs:ClusterEnum[]}, msg:string}) => {
      if (resp.code === 0) {
        this.clusterList = []
        for (const env of resp.data.envs) {
          this.clusterList.push({
            label: env.name,
            value: env.name,
            children: env.clusters.map((c:ClusterEnumData) => ({
              label: c.title,
              value: c.name,
              isLeaf: true
            }))
          })
        }
        this.searchData.cluster = [resp.data.envs[0].name, resp.data.envs[0].clusters[0].name]
        this.getNodeList(init)
      }
    }
    )
  }

  getNodeList (init?:boolean) {
    if (!this.searchData.cluster) {
      return
    }
    this.api.get(`cluster/${this.searchData.cluster[1]}/nodes/simple`).subscribe((resp:{code:number, data:{nodes:Array<{name:number}>}}) => {
      if (resp.code === 0) {
        this.nodeList = resp.data.nodes.map((node:{name:number}) => {
          return ({ label: node.name, value: node.name })
        })
        this.searchData.node = resp.data.nodes[0].name
        init && this.getData()
      }
    })
  }

  getData () {
    this.accGroupLoading = true
    this.api.get('log/files', { cluster: this.searchData.cluster[1], node: this.searchData.node }).subscribe((resp:{code:number, data:{output:LogOutputData[]}, msg:string}) => {
      if (resp.code === 0) {
        this.accGroupList = resp.data.output.filter((x:LogOutputData) => (x)).map((x:LogOutputData) => {
          return { ...x, files: this.filesDataTransfer(x.files), active: true }
        })
      }
      this.start = true
      this.accGroupLoading = false
    })
  }

  // 文件名按z-a排序
  filesDataTransfer (files:Array<LogFileData>) {
    const newFiles =
        files.sort((a:LogFileData, b:LogFileData) => ((b.file + '').localeCompare(a.file + '')))
          .map((x:LogFileData) => {
            x.mod = moment(x.mod).format('yyyy-MM-DD HH:mm:ss')
            return x
          })
    return newFiles
  }

  getTail (e:Event, panel:any) {
    e?.stopPropagation()

    this.modalRef = this.modalService.create({
      nzTitle: `日志详情：${panel.name}`,
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: EoNgLogRetrievalTailComponent,
      nzComponentParams: {
        outputName: panel.name,
        tailKey: panel.tail
      },
      nzFooter: this.tailLogModalFooterTpl
    })
    this.modalRef.afterClose.subscribe(() => {
      this.modalRef = undefined
    })
  }

  downloadLog (file:LogFileData) {
    window.location.href = `api/log/download/${file.key}`
  }
}
