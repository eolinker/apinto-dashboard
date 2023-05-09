import { Component, EventEmitter, Input, OnInit, Output, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { CheckBoxOptionInterface } from 'eo-ng-checkbox'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { TBODY_TYPE, THEAD_TYPE } from 'eo-ng-table'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE } from 'projects/core/src/app/constant/app.config'
import { filterTableHeadName } from '../../types/conf'
import { FilterForm, FilterShowData } from '../../types/types'
import { ServGovernanceFilterService } from '../serv-governance-filter.service'

@Component({
  selector: 'eo-ng-filter-table',
  templateUrl: './table.component.html',
  styles: [
  ]
})
export class FilterTableComponent implements OnInit {
  @ViewChild('filterContentTpl', { read: TemplateRef, static: true })
  filterContentTpl: TemplateRef<any> | undefined

  @ViewChild('filterFooterTpl', { read: TemplateRef, static: true })
  filterFooterTpl: TemplateRef<any> | undefined

  @Input()
  get filterShowList () {
    return this._filterShowList
  }

  set filterShowList (val) {
    this._filterShowList = val
    this.filterShowListChange.emit(this._filterShowList)
  }

  @Input()
  get filterNamesSet () {
    return this._filterNameSet
  }

  set filterNamesSet (val) {
    this._filterNameSet = val
    this.filterNameSetChange.emit(this._filterNameSet)
  }

  @Input() filterTableTipShowFn?:Function
  @Input() nzDisabled:boolean = false
  @Input() filterTableTip:string = ''
  @Input() drawerTitle:string = '筛选条件'
  @Output() filterShowListChange = new EventEmitter()
  @Output() filterNameSetChange = new EventEmitter()

  _filterShowList: FilterShowData [] = []
  _filterNameSet: Set<string> = new Set()
  strategyType: string = ''
  remoteSelectList: string [] = []
  remoteSelectNameList: string [] = []
  filterType: string = '' // 筛选条件类型, 当type=pattern,显示输入框, static显示一组勾选框, remote显示穿梭框
  drawerFilterRef: NzModalRef | undefined
  editFilter: FilterForm|undefined // 正在编辑的配置
  staticsList:CheckBoxOptionInterface[] = []

  filterForm: FilterForm = {
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

  filterTableHeadName: THEAD_TYPE[]= [...filterTableHeadName]
  filterTableBody:TBODY_TYPE[] = []

  constructor (public modalService: EoNgFeedbackModalService,
    private router: Router,
    private service: ServGovernanceFilterService) {
    this.strategyType = this.router.url.split('/')[2]
  }

  ngOnInit (): void {
    this.filterTableBody = this.service.createFilterTbody(this)
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  filterTableClick = (item: {data:FilterForm}) => {
    this.openDrawer('editFilter', item.data)
  }

  openDrawer = (type: string, data?: FilterForm) => {
    switch (type) {
      case 'addFilter': {
        this.editFilter = undefined
        break
      }
      case 'editFilter': {
        this.remoteSelectList = []
        this.remoteSelectNameList = []
        this.filterForm = Object.assign({}, data)
        if (data!.values?.length === 1 && data!.values[0] === 'ALL') {
          this.filterForm.allChecked = true
        }
        this.editFilter = Object.assign({}, data)
        if (data!.name === 'ip') {
          this.filterForm.values = [this.filterForm.values.join('\n')]
          this.editFilter!.values = [this.editFilter!.values.join('\n')]
        }
        break
      }
    }

    this.drawerFilterRef = this.modalService.create({
      nzTitle: type === 'addFilter' ? '配置' + this.drawerTitle : '编辑' + this.drawerTitle,
      nzWidth: MODAL_NORMAL_SIZE,
      nzContent: this.filterContentTpl,
      nzFooter: this.filterFooterTpl,
      nzComponentParams: { data: data },
      nzWrapClassName: 'filter-drawer'
    })
    this.drawerFilterRef?.afterClose.subscribe(() => {
      this.cleanFilterForm()
    })
  }

  drawerClose (val:boolean) {
    val && this.drawerFilterRef?.close()
  }

  filterDelete = (context:FilterTableComponent, item: FilterForm) => {
    this.filterNamesSet.delete(item.name)
    for (const index in context._filterShowList) {
      if (context._filterShowList[index] === item) {
        context._filterShowList.splice(Number(index), 1)
        context._filterShowList = [...context._filterShowList]
      }
    }

    context.filterShowListChange.emit(context.filterShowList)
  }

  cleanFilterForm () {
    this.filterForm = {
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
    this.editFilter = undefined
    this.remoteSelectList = []
    this.remoteSelectNameList = []
    this.staticsList = []
  }
}
