import { Component, ViewChild, TemplateRef } from '@angular/core'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { NzTreeNodeOptions } from 'ng-zorro-antd/tree'
import { MODAL_NORMAL_SIZE } from '../../constant/app.config'
import { EoNgLogRetrievalTailComponent } from './tail-log/tail-log.component'
import { ApiService } from '../../service/api.service'
import { ClusterEnum, ClusterEnumData } from '../../constant/type'
import { LogFileData, LogOutputData } from './type'
import { LogRetrievalTableBody, LogRetrievalTableHeadName } from './conf'
import { environment } from 'projects/core/src/environments/environment'

@Component({
  selector: 'eo-ng-log-retrieval',
  templateUrl: './log-retrieval.component.html',
  styles: [
    `
    :host ::ng-deep{
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
    cluster:string,
    node:number|string
  } = {
    cluster: '',
    node: ''
  }

  clusterList:NzTreeNodeOptions[] = []
  nodeList:SelectOption[] = []

  accGroupList:LogOutputData[] = []
  accTableHeader:THEAD_TYPE[] = [...LogRetrievalTableHeadName]
  accTableBody:TBODY_TYPE[] = [...LogRetrievalTableBody]

  modalRef:NzModalRef|undefined

  constructor (private api:ApiService, private modalService:EoNgFeedbackModalService) {

  }

  disabledEdit ($event:any) {
    this.nzDisabled = $event
  }

  ngOnInit () {
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
            title: env.name,
            key: env.name,
            children: env.clusters.map((c:ClusterEnumData) => ({
              title: c.title,
              key: c.uuid
            }))
          })
        }
        this.searchData.cluster = resp.data.envs[0].clusters[0].uuid
        this.getNodeList(init)
      }
    }
    )
  }

  getNodeList (init?:boolean) {
    if (!this.searchData.cluster) {
      return
    }
    // TODO 假数据
    if (environment.production) {
      this.api.get(`cluster/${this.searchData.cluster}/nodes/simple`).subscribe((resp:{code:number, data:{nodes:Array<{id:number}>}}) => {
        if (resp.code === 0) {
          this.nodeList = resp.data.nodes.map((node:{id:number}) => {
            return ({ label: node.id, value: node.id })
          })
          this.searchData.node = resp.data.nodes[0].id
          init && this.getData()
        }
      })
    } else {
      this.nodeList = [{ label: 'node1', value: 'node1' }, { label: 'node12', value: 'node12' }, { label: 'node13', value: 'node13' }]
      this.searchData.node = 'node1'
      init && this.getData()
    }
  }

  getData () {
    // TODO 假数据
    if (environment.production) {
      this.api.get('api/log/files', { cluster: this.searchData.cluster, node: this.searchData.node }).subscribe((resp:{code:number, data:{outputs:LogOutputData[]}, msg:string}) => {
        if (resp.code === 0) {
          this.accGroupList = resp.data.outputs
          for (const g of this.accGroupList) {
            g.active = true
          }
        }
      })
    } else {
      this.accGroupList = [
        {
          name: 'output1',
          tail: '1111111111',
          files: [
            { file: 'file1', size: '1 M', mod: '2023-06-15 16:35:11', key: 'file1' },
            { file: 'file2', size: '2.2 M', mod: '2023-06-15 16:35:11', key: 'file2' },
            { file: 'file3', size: '3.33 M', mod: '2023-06-15 16:35:11', key: 'file3' },
            { file: 'file4', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file4' },
            { file: 'file41', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file41' },
            { file: 'file42', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file42' },
            { file: 'file43', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file43' },
            { file: 'file44', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file44' },
            { file: 'file45', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file45' },
            { file: 'file46', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file46' },
            { file: 'file47', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file47' },
            { file: 'file48', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file48' },
            { file: 'file49', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file49' }]
        },
        {
          name: 'output2',
          tail: '11111111112',
          files: [
            { file: 'file21', size: '1 M', mod: '2023-06-15 16:35:11', key: 'file21' },
            { file: 'file22', size: '2.2 M', mod: '2023-06-15 16:35:11', key: 'file22' },
            { file: 'file23', size: '3.33 M', mod: '2023-06-15 16:35:11', key: 'file23' },
            { file: 'file24', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file24' }]
        },
        {
          name: 'output3',
          tail: '11111112231',
          files: []
        },
        {
          name: 'output4',
          tail: '11111111144',
          files: [
            { file: 'file1', size: '1 M', mod: '2023-06-15 16:35:11', key: 'file1' },
            { file: 'file2', size: '2.2 M', mod: '2023-06-15 16:35:11', key: 'file2' },
            { file: 'file3', size: '3.33 M', mod: '2023-06-15 16:35:11', key: 'file3' },
            { file: 'file4', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file4' },
            { file: 'file41', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file41' },
            { file: 'file42', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file42' },
            { file: 'file43', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file43' },
            { file: 'file44', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file44' },
            { file: 'file45', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file45' },
            { file: 'file46', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file46' },
            { file: 'file47', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file47' },
            { file: 'file48', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file48' },
            { file: 'file49', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file49' }]
        },
        {
          name: 'output15',
          tail: '1111111111',
          files: [
            { file: 'file1', size: '1 M', mod: '2023-06-15 16:35:11', key: 'file1' },
            { file: 'file2', size: '2.2 M', mod: '2023-06-15 16:35:11', key: 'file2' },
            { file: 'file3', size: '3.33 M', mod: '2023-06-15 16:35:11', key: 'file3' },
            { file: 'file4', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file4' },
            { file: 'file41', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file41' },
            { file: 'file42', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file42' },
            { file: 'file43', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file43' },
            { file: 'file44', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file44' },
            { file: 'file45', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file45' },
            { file: 'file46', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file46' },
            { file: 'file47', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file47' },
            { file: 'file48', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file48' },
            { file: 'file49', size: '4.444 M', mod: '2023-06-15 16:35:11', key: 'file49' }]
        }
      ]
      for (const g of this.accGroupList) {
        g.active = true
      }
    }
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
  }

  downloadLog (file:LogFileData) {
    this.api.get(`log/download/${file.key}`).subscribe((resp:any) => {
      console.log(resp)
    })
  }

  show (val:any) {
    console.log(val)
  }
}
