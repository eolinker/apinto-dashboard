/* eslint-disable no-prototype-builtins */
/* eslint-disable dot-notation */
/* eslint-disable no-useless-constructor */
import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { TransferItem, TransferChange } from 'ng-zorro-antd/transfer'
import { ApiService } from '../../../../service/api.service'

export interface FilterForm {
  name: string
  title: string
  values: Array<any>
  label: string
  text: string
  allChecked?: boolean
  showAll?: boolean
  total?: number
  groupUuid?: string
  pattern: RegExp | null
  patternIsPass: boolean
  [key: string]: any
}
@Component({
  selector: 'eo-ng-serv-governance-filter-form',
  templateUrl: './form.component.html',
  styles: [
    `
      textarea {
        min-height: 68px;
      }
      .form-input {
        width: 1000px !important;
      }
    `
  ]
})
export class FilterFormComponent implements OnInit {
  @Input() filterNamesSet: Set<string> = new Set() // 用户已选择的筛选条件放入set中,在显示筛选条件的选择器里需要过去set中存在的值
  @Input() // 双向绑定filterForm
  get filterForm () {
    return this._filterForm
  }

  set filterForm (val) {
    this._filterForm = val
    this.filterFormChange.emit(this._filterForm)
  }

  @Output() filterFormChange = new EventEmitter()

  _filterForm: FilterForm = {
    name: '',
    title: '',
    values: [],
    label: '',
    text: '',
    allChecked: false,
    showAll: false,
    total: 0,
    groupUuid: '',
    pattern: null,
    patternIsPass: true
  }

  @Input() // 双向绑定remoteSelectList
  get remoteSelectList () {
    return this._remoteSelectList
  }

  set remoteSelectList (val) {
    this._remoteSelectList = val
    this.remoteSelectListChange.emit(this._remoteSelectList)
  }

  @Output() remoteSelectListChange = new EventEmitter()

  _remoteSelectList: string[] = [] // 穿梭框内被勾选的选项uuid

  @Input() // 双向绑定remoteSelectNameList
  get remoteSelectNameList () {
    return this._remoteSelectNameList
  }

  set remoteSelectNameList (val) {
    this._remoteSelectNameList = val
    this.remoteSelectNameListChange.emit(this._remoteSelectNameList)
  }

  @Output() remoteSelectNameListChange = new EventEmitter()

  _remoteSelectNameList: string[] = [] // 穿梭框内被勾选的选项name

  @Input() // 双向绑定staticsList
  get staticsList () {
    return this._staticsList
  }

  set staticsList (val) {
    this._staticsList = val
    this.staticsListChange.emit(this._staticsList)
  }

  @Output() staticsListChange = new EventEmitter()

  _staticsList: Array<any> = [] // 穿梭框内被选中的checkbox

  @Input() // 双向绑定staticsList
  get filterType () {
    return this._filterType
  }

  set filterType (val) {
    this._filterType = val
    this.filterTypeChange.emit(this._filterType)
  }

  @Output() filterTypeChange = new EventEmitter()

  _filterType: string = '' // 筛选条件类型, 当type=pattern,显示输入框, static显示一组勾选框, remote显示穿梭框

  @Input() editFilter: any = null // 正在编辑的配置

  filterTypeMap: Map<string, any> = new Map() // 筛选条件值与类型的映射
  remoteList: TransferItem[] = []

  remoteCheckList: string[] = [] // 用户已选择的选项

  filterNamesList: Array<{ label: string; value: string }> = []

  apiGroupList: Array<any> = []
  // 穿梭框
  tableType: string = ''
  filterTitle: string = ''
  filterTPlaceholder: string = ''
  filterThead: Array<any> = [
    {
      type: 'checkbox',
      width: '40px'
    }
  ]

  filterTbody: Array<any> = [
    {
      key: 'checked',
      type: 'checkbox'
    }
  ]

  allChecked: boolean = false

  searchWord: string = ''
  searchGroup: string[] = []
  // pattern:RegExp | null = null // 当type为pattern时，输入的校验规则
  showFilterError: boolean = false
  strategyType: string = ''

  nzDisabled: boolean = false
  constructor (
    private message: EoNgFeedbackMessageService,
    private router: Router,
    private api: ApiService
  ) {
    this.strategyType = this.router.url.split('/')[2]
  }

  ngOnInit (): void {
    this.getFilterNamesList(this.editFilter !== null)
  }

  // 获取筛选条件中属性名称的可选选项
  // 如果不是配置选项页,则只显示不在filterNamesSet的选项,如果是配置选项页,则显示不在set的选项外加上配置的选项
  getFilterNamesList (edit?: boolean): void {
    this.api.get('strategy/filter-options').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.filterNamesList = []
        for (const index in resp.data.options) {
          if (
            (edit && this.filterForm.name === resp.data.options[index].name) ||
            !this.filterNamesSet.has(resp.data.options[index].name)
          ) {
            if (!edit && !this.filterForm.name) {
              this.filterForm.name = resp.data.options[index].name
              this.filterForm.title = resp.data.options[index].title
            }
            if (resp.data.options[index].name === 'api') {
              this.getApiGroupList()
            }
            resp.data.options[index].label = resp.data.options[index].title
            resp.data.options[index].value = resp.data.options[index].name
            // resp.data.options[index].options = resp.data.options[index].values ? [...resp.data.options[index].values] : []
            resp.data.options[index].total =
              resp.data.options[index].options?.length - 1 || 0
            resp.data.options[index].values =
              edit && this.editFilter.name === resp.data.options[index].name
                ? this.editFilter.values
                : []
            resp.data.options[index].allChecked =
              edit && this.editFilter.name === resp.data.options[index].name
                ? this.editFilter.allChecked
                : false
            resp.data.options[index].patternIsPass = true
            this.filterNamesList.push(resp.data.options[index])
            this.filterTypeMap.set(
              resp.data.options[index].name,
              resp.data.options[index]
            )

            if (this.filterForm.name === resp.data.options[index].name) {
              this.filterType = resp.data.options[index].type
            }
          }
        }
        this.changeFilterType(this.filterForm.name)
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  // 获取搜索远程类型的选项，参数为搜索的属性类型
  getRemoteList (name: string): void {
    this.api
      .get('strategy/filter-remote/' + name, {
        keyword: this.searchWord || '',
        group_uuid: this.searchGroup || ''
      })
      .subscribe((resp: any) => {
        if (resp.code === 0) {
          this.remoteList = []
          this.remoteSelectList = []
          this.remoteSelectNameList = []
          for (const index in resp.data[resp.data.target]) {
            resp.data[resp.data.target][index].key =
              resp.data[resp.data.target][index].uuid
            resp.data[resp.data.target][index].checked = false
            resp.data[resp.data.target][index].title =
              resp.data[resp.data.target][index].name
            resp.data[resp.data.target][index].direction =
              this.editFilter && this.filterForm.name === this.editFilter.name
                ? !!this.editFilter.values?.includes(
                    resp.data[resp.data.target][index].uuid
                  ) || this.editFilter.values[0] === 'ALL'
                    ? 'right'
                    : 'left'
                : 'left'
            if (resp.data[resp.data.target][index].direction === 'right') {
              this.remoteSelectList.push(
                resp.data[resp.data.target][index].uuid
              )
              this.remoteSelectNameList.push(
                resp.data[resp.data.target][index].name
              )
            }
            this.remoteList.push(resp.data[resp.data.target][index])
          }
          this.filterTbody = [
            {
              key: 'checked',
              type: 'checkbox'
            }
          ]
          this.filterThead = [
            {
              type: 'checkbox'
            }
          ]
          for (const index in resp.data.titles) {
            this.filterTbody.push({ key: resp.data.titles[index].field })
            this.filterThead.push({ title: resp.data.titles[index].title })
          }
          this.filterTypeMap.get(this.filterForm.name).total = resp.data.total
          this.filterForm.total = resp.data.total
        } else {
          this.message.error(resp.msg || '获取数据失败!')
        }
      })
  }

  // 获取API目录列表参数
  getApiGroupList () {
    this.api.get('router/groups').subscribe((resp: any) => {
      if (resp.code === 0) {
        this.apiGroupList = []
        this.apiGroupList = resp.data.root.groups
        this.apiGroupList = this.transferHeader(this.apiGroupList)
      } else {
        this.message.error(resp.msg || '获取数据失败!')
      }
    })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 搜索特定的远程类型数据
  getSearchRemoteList (direction:string): void {
    this.api
      .get('strategy/filter-remote/' + this.filterForm.name, {
        keyword: this.searchWord || '',
        group_uuid:
          this.searchGroup.length > 0
            ? this.searchGroup[this.searchGroup.length - 1]
            : ''
      })
      .subscribe((resp: any) => {
        if (resp.code === 0) {
          const showItemUuidArr:Array<string> = []
          if (resp.data[resp.data.target]) {
            for (const item of resp.data[resp.data.target]) {
              showItemUuidArr.push(item.uuid)
            }
          }
          for (const row of this.remoteList) {
            row.hide = row.direction === direction &&
                         row['uuid'] &&
                          (showItemUuidArr.length === 0 ||
                           showItemUuidArr.indexOf(row['uuid']) === -1)
          }
          this.remoteList = [...this.remoteList]

          this.filterForm.total = resp.data.total
        } else {
          this.message.error(resp.msg || '筛选失败!')
        }
      })
  }

  changeFilterType (value: any) {
    this.filterType = this.filterTypeMap.get(value).type || ''
    this.filterForm = Object.assign({}, this.filterTypeMap.get(value))
    switch (this.filterType) {
      case 'pattern': {
        this.filterForm.pattern = this.filterTypeMap.get(value).pattern
          ? new RegExp(this.filterTypeMap.get(value).pattern)
          : null
        this.checkPattern()
        break
      }
      case 'static': {
        this.staticsList = []
        for (const index in this.filterTypeMap.get(value).options) {
          if (this.filterTypeMap.get(value).options[index] !== 'ALL') {
            this.staticsList.push({
              label: this.filterTypeMap.get(value).options[index],
              value: this.filterTypeMap.get(value).options[index],
              checked: this.editFilter
                ? this.filterForm.values.includes(
                  this.filterTypeMap.get(value).options[index]
                )
                : false
            })
          } else {
            this.filterForm.showAll = true
          }
        }
        if (
          this.filterForm.allChecked ||
          (this.filterForm.values?.length === 1 &&
            this.filterForm.values[0] === 'ALL')
        ) {
          this.filterForm.allChecked = true
          this.updateAllChecked()
        }
        break
      }
      case 'remote': {
        this.getRemoteList(value)
        break
      }
    }
  }

  // 穿梭框内,点击穿梭按钮后,根据数据穿梭方向,增加或减少remoteSelectList,remoteSelectNameList
  change (ret: TransferChange): void {
    if (ret.list.length > 0) {
      const listKeys = ret.list.map((l) => l['key'])
      const hasOwnKey = (e: TransferItem): boolean => e.hasOwnProperty('key')
      this.remoteList = this.remoteList.map((e) => {
        if (listKeys.includes(e['key']) && hasOwnKey(e)) {
          if (ret.to === 'left') {
            delete e.hide
            const deleteIndex = this.remoteSelectList.indexOf(e['uuid'])
            this.remoteSelectList.splice(deleteIndex, 1)
            this.remoteSelectNameList.splice(deleteIndex, 1)
          } else if (ret.to === 'right') {
            e.hide = false
            this.remoteSelectList.push(e['uuid'])
            this.remoteSelectNameList.push(e['name'])
          }
        }
        return e
      })
    }
  }

  updateAllChecked () {
    this.staticsList = this.staticsList.map((item: any) => {
      item.checked = this.filterForm.allChecked
      return item
    })
  }

  updateSingleChecked (): void {
    if (this.staticsList.every((item) => !item.checked)) {
      this.filterForm.allChecked = false
    } else if (this.staticsList.every((item) => item.checked)) {
      this.filterForm.allChecked = true
    } else {
      this.filterForm.allChecked = false
    }
    this.filterForm.values = []
    for (const index in this.staticsList) {
      if (this.staticsList[index].checked) {
        this.filterForm.values.push(this.staticsList[index].value)
      }
    }
  }

  ipArray: Array<string> = []
  checkPattern () {
    if (this.filterForm.name !== 'ip') {
      if (this.filterForm.values[0] && this.filterForm.pattern) {
        this.filterForm.patternIsPass = this.filterForm.pattern.test(
          this.filterForm.values[0].trim()
        )
      }
    } else {
      this.filterForm.values[0] = this.filterForm.values[0]?.split(/[,，、]/)
        .map((x: any) => x.trim())
        .join('\n')
      this.ipArray = this.filterForm.values[0]?.split(/[\n]/) || []
      if (this.filterForm.pattern) {
        for (const index in this.ipArray) {
          if (
            this.ipArray[index] &&
            !this.filterForm.pattern.test(this.ipArray[index].trim())
          ) {
            this.filterForm.patternIsPass = false
            return
          }
          this.filterForm.patternIsPass = true
        }
      }
    }
  }

  $asTransferItems = (data: unknown): TransferItem[] => data as TransferItem[]

  transferHeader (header: any) {
    for (const index in header) {
      if (!header[index].children || header[index].children.length === 0) {
        header[index].isLeaf = true
      } else {
        header[index].children = this.transferHeader(header[index].children)
      }
    }
    return header
  }
}
