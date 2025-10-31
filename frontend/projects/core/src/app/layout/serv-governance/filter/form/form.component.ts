/* eslint-disable dot-notation */
import { Component, Input, OnInit, Output, EventEmitter } from '@angular/core'
import { Router } from '@angular/router'
import { CascaderOption } from 'eo-ng-cascader'
import { CheckBoxOptionInterface } from 'eo-ng-checkbox'
import { SelectOption } from 'eo-ng-select'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { ApiGroup } from 'projects/core/src/app/constant/type'
import { ApiService } from '../../../../service/api.service'
import { FilterForm, FilterOption, FilterRemoteOption, RemoteApiItem, RemoteAppItem, RemoteServiceItem } from '../../types/types'

@Component({
  selector: 'eo-ng-serv-governance-filter-form',
  templateUrl: './form.component.html',
  styles: [
    `
      textarea {
        min-height: 68px;
      }
      .transfer-section{
        border-radius:var(--border-radius);
      }
      .tips {
        color: red;
        font-size: 12px;
        margin-top: 4px;
      }

      :host ::ng-deep{
        .label {
              width: 78px !important;
              display: inline-block;
            }
            .ant-drawer-body {
              padding-right: 44px;
              padding-bottom: 0px;
            }
            [eo-ng-transfer-search].ant-input-affix-wrapper.ant-transfer-list-search {
              input.ant-input {
                height: 20px !important;
              }
            }
          }
    `
  ]
})
export class FilterFormComponent implements OnInit {
  @Input() // 双向绑定filterForm
  get filterForm () { return this._filterForm }

  set filterForm (val) {
    this._filterForm = val
    this.filterFormChange.emit(this._filterForm)
  }

  @Input() // 双向绑定remoteSelectList
  get remoteSelectList () { return this._remoteSelectList }

  set remoteSelectList (val) {
    this._remoteSelectList = val
    this.remoteSelectListChange.emit(this._remoteSelectList)
  }

  @Input() // 双向绑定remoteSelectNameList
  get remoteSelectNameList () { return this._remoteSelectNameList }

  set remoteSelectNameList (val) {
    this._remoteSelectNameList = val
    this.remoteSelectNameListChange.emit(this._remoteSelectNameList)
  }

  @Input() // 双向绑定staticsList
  get staticsList () { return this._staticsList }

  set staticsList (val) {
    this._staticsList = val
    this.staticsListChange.emit(this._staticsList)
  }

  @Input() // 双向绑定staticsList
  get filterType () { return this._filterType }

  set filterType (val) {
    this._filterType = val
    this.filterTypeChange.emit(this._filterType)
  }

  @Input() editFilter?: FilterForm // 正在编辑的配置
  @Input() filterNamesSet: Set<string> = new Set() // 用户已选择的筛选条件放入set中,在显示筛选条件的选择器里需要过去set中存在的值

  @Output() filterFormChange:EventEmitter<FilterForm> = new EventEmitter()
  @Output() remoteSelectListChange:EventEmitter<string[]> = new EventEmitter()
  @Output() remoteSelectNameListChange:EventEmitter<string[]> = new EventEmitter()
  @Output() staticsListChange:EventEmitter<CheckBoxOptionInterface[]> = new EventEmitter()
  @Output() filterTypeChange:EventEmitter<string> = new EventEmitter()
  @Input() type: string = ''
  @Input() nzDisabled: boolean = false
  filterTypeMap: Map<string, any> = new Map() // 筛选条件值与类型的映射
  remoteList: RemoteAppItem[] | RemoteApiItem[] | RemoteServiceItem[] = []
  _remoteSelectList: string[] = [] // 穿梭框内被勾选的选项uuid
  filterNamesList: SelectOption[] = []
  _filterType: string = '' // 筛选条件类型, 当type=pattern,显示输入框, static显示一组勾选框, remote显示穿梭框
  _remoteSelectNameList: string[] = [] // 穿梭框内被勾选的选项name
  _staticsList: CheckBoxOptionInterface[]= [] // 穿梭框内被选中的checkbox
  apiGroupList: CascaderOption[] = []
  allChecked: boolean = false
  searchWord: string = ''
  searchGroup: string[] = []
  showFilterError: boolean = false
  strategyType: string = ''
  // 时间范围数组，每个元素包含开始和结束时间
  timeRanges: Array<{
    start: Date | null
    end: Date | null
    error?: string // 用于存储错误信息
    startDefaultOpenValue?: any
    endDefaultOpenValue?: any
  }> = []

  valueName: string = 'uuid' // 数据keyName，由接口传入，如果接口未传入则默认为uuid
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

  // 穿梭框
  filterThead: THEAD_TYPE[] = [
    {
      type: 'checkbox',
      width: 40
    }
  ]

  filterTbody: TBODY_TYPE[] = [
    {
      key: 'checked',
      type: 'checkbox'
    }
  ]

  ipArray: Array<string> = []
  originDataLength: number = 0 // 未经筛选的数据列表长度
  originRemoteList: any[] = [] // 未经筛选的数据列表

  constructor (
    private router: Router,
    private api: ApiService
  ) {
    this.strategyType = this.router.url.split('/')[this.router.url.split('/').indexOf('serv-governance') + 1]
  }

  // 添加新的时间范围
  addTimeRange (): void {
    this.timeRanges.push({
      start: null,
      end: null,
      startDefaultOpenValue: new Date(0, 0, 0, 0, 0, 0),
      endDefaultOpenValue: new Date(0, 0, 0, 0, 0, 0),
      error: ''
    })
    this.updateTimeValues()
  }

  // 删除时间范围
  removeTimeRange (index: number): void {
    this.timeRanges.splice(index, 1)
    this.updateTimeValues()
  }

  // 验证时间范围
  private validateTimeRange (range: { start: Date | null; end: Date | null }): string {
    if (!range.start || !range.end) {
      return '请选择完整的时间范围'
    }

    // 检查时间格式
    const startTimeStr = this.formatTime(range.start)
    const endTimeStr = this.formatTime(range.end)
    const timeRangeStr = `${startTimeStr} - ${endTimeStr}`
    // 确保 pattern 是字符串类型
    const patternStr = typeof this.filterForm.pattern === 'string'
      ? this.filterForm.pattern
      : '^((?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d|24:00:00) - ((?:[01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d|24:00:00)$'
    if (!new RegExp(patternStr).test(timeRangeStr)) {
      return '时间格式不正确，应为 HH:mm:ss - HH:mm:ss'
    }

    // 检查结束时间是否大于开始时间
    if (range.start >= range.end) {
      return '结束时间必须大于开始时间'
    }

    return ''
  }

  // 格式化时间为 HH:mm:ss
  private formatTime (date: Date): string {
    if (!date) return ''
    const hours = date.getHours().toString().padStart(2, '0')
    const minutes = date.getMinutes().toString().padStart(2, '0')
    const seconds = date.getSeconds().toString().padStart(2, '0')
    return `${hours}:${minutes}:${seconds}`
  }

  // 更新时间范围值到filterForm.values
  updateTimeValues (): void {
    this.checkPattern()
  }

  // 初始化时间范围
  private initTimeRanges (): void {
    if (this.filterForm.values?.length && this.filterForm.name === 'time') {
      this.timeRanges = this.filterForm.values.map((timeRange) => {
        const [startStr, endStr] = timeRange.split('-')
        const [startHours, startMinutes, startSeconds] = startStr.split(':').map(Number)
        const [endHours, endMinutes, endSeconds] = endStr.split(':').map(Number)

        const start = new Date()
        start.setHours(startHours, startMinutes, startSeconds || 0)

        const end = new Date()
        end.setHours(endHours, endMinutes, endSeconds || 0)

        const range = { start, end, error: '', startDefaultOpenValue: new Date(0, 0, 0, 0, 0, 0), endDefaultOpenValue: new Date(0, 0, 0, 0, 0, 0) }
        range.error = this.validateTimeRange(range)
        return range
      })
    } else {
      // 默认添加一个空的时间范围
      this.timeRanges = [{ start: null, end: null, startDefaultOpenValue: new Date(0, 0, 0, 0, 0, 0), endDefaultOpenValue: new Date(0, 0, 0, 0, 0, 0) }]
    }
  }

  ngOnInit (): void {
    this.getFilterNamesList(!!this.editFilter)
    this.initTimeRanges()
  }

  // 获取筛选条件中属性名称的可选选项
  // 如果不是配置选项页,则只显示不在filterNamesSet的选项,如果是配置选项页,则显示不在set的选项外加上配置的选项
  getFilterNamesList (edit?: boolean): void {
    this.api.get('strategy/filter-options')
      .subscribe((resp: {code:number, data:{options:FilterOption[]}, msg:string}) => {
        if (resp.code === 0) {
          this.filterNamesList = []
          for (const index in resp.data.options) {
            if (resp.data.options[index].name === 'time' && this.type !== 'apiVisit') {
              continue
            }
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
              resp.data.options[index]['total'] =
              resp.data.options[index].options?.length - 1 || 0
              resp.data.options[index]['values'] =
              edit && this.editFilter!.name === resp.data.options[index].name
                ? this.editFilter!.values
                : []
              resp.data.options[index]['allChecked'] =
              edit && this.editFilter!.name === resp.data.options[index].name
                ? this.editFilter!.allChecked
                : false
              resp.data.options[index]['patternIsPass'] = true
              this.filterNamesList.push(resp.data.options[index] as SelectOption)
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
        }
      })
  }

  // 获取搜索远程类型的选项，参数为搜索的属性类型
  getRemoteList (name: string): void {
    this.api
      .get('strategy/filter-remote/' + name, {
        keyword: this.searchWord || '',
        groupUuid: this.searchGroup?.length > 0 ? this.searchGroup : ''
      })
      .subscribe((resp: {code:number, data:FilterRemoteOption, msg:string}) => {
        if (resp.code === 0) {
          this.remoteList = []
          this.remoteSelectList = []
          this.remoteSelectNameList = []
          const list:Array<RemoteApiItem|RemoteServiceItem|RemoteAppItem> = resp.data.list
          this.valueName = resp.data.key
          for (const index in list) {
            list[index].checked = this.editFilter && this.filterForm.name === this.editFilter.name
              ? !!(!!this.editFilter.values?.includes(
                  list[index][this.valueName]
                ) || this.editFilter.values[0] === 'ALL')
              : false
            if (list[index].checked) {
              this.remoteSelectList.push(
                list[index][this.valueName]
              )
              this.remoteSelectNameList.push(
                list[index].title
              )
            }
            this.remoteList.push(list[index] as any)
          }
          this.remoteList.sort((a, b) => {
            if (a.checked === b.checked) {
              return 0
            }
            return a.checked ? -1 : 1
          })
          this.originRemoteList = [...this.remoteList]
          this.filterTbody = [
            {
              key: 'checked',
              type: 'checkbox',
              click: () => {
                this.getNewRemotesStatus()
              },
              disabledFn: () => {
                return this.nzDisabled
              }

            }
          ]
          this.filterThead = [
            {
              type: 'checkbox',
              click: () => {
                this.getNewRemotesStatus()
              },
              disabled: this.nzDisabled

            }
          ]
          for (const index in resp.data.titles) {
            this.filterTbody.push({ key: (resp.data.titles[index].field).replace(/_([a-z])/g, (p, m) => m.toUpperCase()) })
            this.filterThead.push({ title: resp.data.titles[index].title })
          }
          this.filterTypeMap.get(this.filterForm.name).total = resp.data.total
          this.filterForm.total = resp.data.total
          this.originDataLength = resp.data.total
        }
      })
  }

  clickItem = (item:any) => {
    item.checked = !item.checked
    item.data.checked = !item.data.checked
    this.getNewRemotesStatus()
  }

  getNewRemotesStatus () {
    setTimeout(() => {
      for (const item of this.remoteList) {
        if (item.checked) {
          if (this._remoteSelectList.indexOf(item[this.valueName]) === -1) {
            this._remoteSelectList.push(item[this.valueName])
            this._remoteSelectNameList.push(item.title)
          }
        } else {
          if (this._remoteSelectList.indexOf(item[this.valueName]) !== -1) {
            this._remoteSelectList.splice(this._remoteSelectList.indexOf(item[this.valueName]), 1)
            this._remoteSelectNameList.splice(this._remoteSelectNameList.indexOf(item.title), 1)
          }
        }
      }
      this.remoteSelectListChange.emit(this._remoteSelectList)
      this.remoteSelectNameListChange.emit(this._remoteSelectNameList)
    })
  }

  // 获取API目录列表参数
  getApiGroupList () {
    this.api.get('router/groups').subscribe((resp:{code:number, data:ApiGroup, msg:string}) => {
      if (resp.code === 0) {
        this.apiGroupList = []
        this.apiGroupList = resp.data.root.groups
        this.apiGroupList = this.transferHeader(this.apiGroupList)
      }
    })
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  // 搜索特定的远程类型数据
  getSearchRemoteList (): void {
    this.api
      .get('strategy/filter-remote/' + this.filterForm.name, {
        keyword: this.searchWord || '',
        groupUuid:
          this.searchGroup.length > 0
            ? this.searchGroup[this.searchGroup.length - 1]
            : ''
      })
      .subscribe((resp: {code:number, data:FilterRemoteOption, msg:string}) => {
        if (resp.code === 0) {
          this.valueName = resp.data.key
          this.remoteList = resp.data.list as any
          this.remoteList = this.remoteList.map((item:any) => {
            item.checked = this._remoteSelectList.indexOf(item.uuid) !== -1
            return item
          })
        }
      })
  }

  // 搜索远程数据（不调接口
  searchRemoteList () {
    this.remoteList = this.originRemoteList.filter((item:any) => {
      return item.title.includes(this.searchWord)
    })
  }

  changeFilterType (value: string) {
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

  updateAllChecked () {
    this.staticsList = this.staticsList.map((item: CheckBoxOptionInterface) => {
      item.checked = this.filterForm.allChecked
      return item
    })
  }

  updateSingleChecked (): void {
    if (this.staticsList.every((item:CheckBoxOptionInterface) => !item.checked)) {
      this.filterForm.allChecked = false
    } else if (this.staticsList.every((item:CheckBoxOptionInterface) => item.checked)) {
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

  checkPattern () {
    // 处理时间范围验证
    if (this.filterType === 'pattern' && this.filterForm.name === 'time') {
      // 验证所有时间范围
      let hasError = false
      this.timeRanges.forEach(range => {
        range.error = this.validateTimeRange(range)
        if (range.error) {
          hasError = true
        }
      })

      // 只保存验证通过的时间范围
      this.filterForm.values = this.timeRanges
        .filter(range => !range.error)
        .map(range => `${this.formatTime(range.start!)}-${this.formatTime(range.end!)}`)
      // 设置验证状态
      this.filterForm.patternIsPass = !hasError && this.filterForm.values.length > 0
      return
    }

    // 原有的IP和模式验证逻辑
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

  transferHeader (header: CascaderOption[]) {
    for (const index in header) {
      header[index].label = header[index]['name']
      header[index].value = header[index]['uuid']
      if (!header[index].children || header[index].children!.length === 0) {
        header[index].isLeaf = true
      } else {
        header[index].children = this.transferHeader(header[index].children as CascaderOption[])
      }
    }
    return header
  }
}
