/* eslint-disable dot-notation */
import { Component, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { TransferChange, TransferItem } from 'ng-zorro-antd/transfer'
import { TransferComponent } from 'projects/core/src/app/component/transfer/transfer/transfer.component'
import { ApiGroup, ApiGroupsData, RemoteData } from 'projects/core/src/app/constant/type'
import { ApiService } from 'projects/core/src/app/service/api.service'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { CascaderOption } from 'eo-ng-cascader'

@Component({
  selector: 'eo-ng-monitor-alarm-strategy-transfer',
  templateUrl: './transfer.component.html',
  styles: [
  ]
})
export class MonitorAlarmStrategyTransferComponent implements OnInit {
  @ViewChild('renderListTpl') renderListTpl: TemplateRef<any> | undefined
  @ViewChild('transferRef') transferRef:TransferComponent |undefined
  @Input() type:'api'|'service'|''=''
  @Input() rule:'contain'|'not_contain'|'' = ''
  @Input() selectedList:Array<string> = [] // 穿梭框内被选中的选项
  tableList:TransferItem[] = [];
  backendTableList:TransferItem[] = [];
  searchWord:string = ''
  list: TransferItem[] = [];
  tableHead:THEAD_TYPE[] = []
  tableBody:TBODY_TYPE[] = []
  apiGroupList:CascaderOption[] = []
  searchGroup: string[] = []
  searchKeyword:string = ''
  // eslint-disable-next-line no-useless-constructor
  constructor (private api:ApiService, private message: EoNgFeedbackMessageService) { }

  $asTransferItems = (data: unknown): TransferItem[] => {
    return data as TransferItem[]
  }

  checkThCheckAll (dir:string) {
    const uncheckedArr:TransferItem[] = this.tableList.filter((item:TransferItem) => {
      return item.direction === dir && !item.checked
    })
    const checkedArr:TransferItem[] = this.tableList.filter((item:TransferItem) => {
      return item.direction === dir && item.checked
    })

    return uncheckedArr.length === 0 && checkedArr.length !== 0
  }

  checkThHalf (dir:string) {
    const uncheckedArr:TransferItem[] = this.tableList.filter((item:TransferItem) => {
      return item.direction === dir && !item.checked
    })
    const checkedArr:TransferItem[] = this.tableList.filter((item:TransferItem) => {
      return item.direction === dir && item.checked
    })

    return checkedArr.length !== 0 && uncheckedArr.length > 0
  }

  ngOnInit (): void {
    this.getRemoteList()
    if (this.type === 'api') {
      this.getApiGroupList()
    }
  }

  // 获取搜索远程类型的选项，参数为搜索的属性类型
  getRemoteList (): void {
    this.api
      .get('strategy/filter-remote/' + this.type)
      .subscribe((resp: {code : number, data:RemoteData}) => {
        if (resp.code === 0) {
          const newTableList:TransferItem[] = []
          for (const index in resp.data[resp.data.target]) {
            resp.data[resp.data.target][index]['key'] =
              resp.data[resp.data.target][index].uuid
            resp.data[resp.data.target][index]['checked'] = false
            resp.data[resp.data.target][index]['direction'] =
              this.selectedList && this.selectedList.indexOf(resp.data[resp.data.target][index].uuid) !== -1
                ? 'right'
                : 'left'
            newTableList.push({ ...resp.data[resp.data.target][index], title: resp.data[resp.data.target][index].name, disabled: false })
          }
          this.tableList = newTableList
          this.backendTableList = [...this.tableList]
          this.tableBody = [
            {
              key: 'checked',
              type: 'checkbox',
              click: (item:TransferItem) => {
                item.direction === 'left' ? this.transferRef?.moveToRight() : this.transferRef?.moveToLeft()
              }
            }
          ]
          this.tableHead = [
            {
              type: 'checkbox'
            }
          ]
          for (const index in resp.data.titles) {
            this.tableBody.push({ key: resp.data.titles[index].field })
            this.tableHead.push({ title: resp.data.titles[index].title })
          }
        }
      })
  }

  // 获取API目录列表参数
  getApiGroupList () {
    this.api.get('router/groups').subscribe((resp: {code:number, data:ApiGroup, msg:string}) => {
      if (resp.code === 0) {
        this.apiGroupList = this.transferHeader(resp.data.root.groups)
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  transferHeader (header: ApiGroupsData[]):CascaderOption[] {
    const res:CascaderOption[] = []
    for (const index in header) {
      res[index] = { ...header[index] }
      if (!header[index].children || header[index].children.length === 0) {
        res[index]['isLeaf'] = true
      } else {
        res[index].children = this.transferHeader(header[index].children)
      }
    }
    return header
  }

  // keyword是搜索框绑定的值,接口只提供对目录、关键词的过滤，且返回时间久，所以在不涉及目录修改的情况下完全由前端搜索，涉及目录则结合
  getSearchRemoteList (direction:string, type?:'service' |'keyword') {
    if (type) {
      this.tableList = this.backendTableList.map((list:TransferItem) => { return { ...list } })
      for (const row of this.tableList) {
        if (!row.hide && row.direction === direction) {
          row.hide = (!!this.searchKeyword && !row['title'].toLocaleLowerCase().includes(this.searchKeyword.toLocaleLowerCase()))
        }
      }
      this.tableList = [...this.tableList]
      return
    }
    this.api
      .get('strategy/filter-remote/api', {
        keyword: this.searchKeyword,
        group_uuid:
      this.searchGroup!.length > 0
        ? this.searchGroup[this.searchGroup!.length - 1]
        : ''
      })
      .subscribe((resp: {code:number, data:RemoteData, msg:string}) => {
        if (resp.code === 0) {
          const showItemUuidArr:Array<string> = []
          if (resp.data[resp.data.target]) {
            for (const item of resp.data[resp.data.target]) {
              showItemUuidArr.push(item.uuid)
            }
          }
          for (const row of this.backendTableList) {
            row.hide = row.direction === direction &&
                         row['uuid'] &&
                          (showItemUuidArr.length === 0 ||
                           showItemUuidArr.indexOf(row['uuid']) === -1)
          }
          this.backendTableList = [...this.backendTableList]
          this.tableList = this.backendTableList.map((list:TransferItem) => { return { ...list } })
        } else {
          this.message.error(resp.msg || '筛选失败!')
        }
      })
  }

  // 穿梭框内,点击穿梭按钮后,根据数据穿梭方向,增加或减少selectedList
  change (ret: TransferChange): void {
    if (ret.list.length > 0) {
      const listKeys = ret.list.map((l) => l['key'])
      // eslint-disable-next-line no-prototype-builtins
      const hasOwnKey = (e: TransferItem): boolean => e.hasOwnProperty('key')
      this.tableList = this.tableList.map((e) => {
        if (listKeys.includes(e['key']) && hasOwnKey(e)) {
          if (ret.to === 'left') {
            delete e.hide
            const deleteIndex = this.selectedList.indexOf(e['uuid'])
            this.selectedList.splice(deleteIndex, 1)
          } else if (ret.to === 'right') {
            e.hide = false
            this.selectedList.push(e['uuid'])
          }
        }
        return e
      })
    }
  }
}
