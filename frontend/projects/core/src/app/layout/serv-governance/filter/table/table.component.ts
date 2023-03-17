import { Component, EventEmitter, Input, Output, TemplateRef, ViewChild } from '@angular/core'
import { Router } from '@angular/router'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { MODAL_NORMAL_SIZE, MODAL_SMALL_SIZE } from 'projects/core/src/app/constant/app.config'
import { FilterShowData } from '../footer/footer.component'
import { FilterForm } from '../form/form.component'

@Component({
  selector: 'eo-ng-filter-table',
  templateUrl: './table.component.html',
  styles: [
  ]
})
export class FilterTableComponent {
  @ViewChild('filterContentTpl', { read: TemplateRef, static: true })
  filterContentTpl: TemplateRef<any> | undefined

  @ViewChild('filterFooterTpl', { read: TemplateRef, static: true })
  filterFooterTpl: TemplateRef<any> | undefined

  @Input() nzDisabled:boolean = false
  @Input() filterTableTip:string = ''
  @Input() drawerTitle:string = '筛选条件'
  @Input()
  get filterShowList () {
    return this._filterShowList
  }

  set filterShowList (val) {
    this._filterShowList = val
    this.filterShowListChange.emit(this._filterShowList)
  }

  @Output() filterShowListChange = new EventEmitter()

  _filterShowList: FilterShowData [] = []

  @Input()
  get filterNamesSet () {
    return this._filterNameSet
  }

  set filterNamesSet (val) {
    this._filterNameSet = val
    this.filterNameSetChange.emit(this._filterNameSet)
  }

  @Output() filterNameSetChange = new EventEmitter()

  _filterNameSet: Set<string> = new Set()

  @Input() filterTableTipShowFn?:Function

  strategyType: string = ''
  remoteSelectList: string [] = []
  remoteSelectNameList: string [] = []
  filterType: string = '' // 筛选条件类型, 当type=pattern,显示输入框, static显示一组勾选框, remote显示穿梭框
  drawerFilterRef: NzModalRef | undefined
  editFilter: any = null // 正在编辑的配置
  staticsList: Array<any> = []

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

  filterTableHeadName: Array<any> = [
    {
      title: '属性名称'
    },
    {
      title: '属性值'
    },
    {
      title: '操作',
      right: true
    }
  ]

  filterTableBody: Array<any> = [
    { key: 'title' },
    { key: 'label' },
    {
      type: 'btn',
      right: true,
      btns: [
        {
          title: '配置',
          click: (item: any) => {
            this.openDrawer('editFilter', item.data)
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        },
        {
          title: '删除',
          click: (item: any) => {
            this.modalService.create({
              nzTitle: '删除',
              nzContent: '该数据删除后将无法找回，请确认是否删除？',
              nzClosable: true,
              nzClassName: 'delete-modal',
              nzOkDanger: true,
              nzWidth: MODAL_SMALL_SIZE,
              nzOnOk: () => {
                this.filterDelete(item.data)
              }
            })
          },
          disabledFn: () => {
            return this.nzDisabled
          }
        }
      ]
    }
  ]

  constructor (private modalService: EoNgFeedbackModalService,
    private router: Router) {
    this.strategyType = this.router.url.split('/')[2]
  }

  disabledEdit (value: any) {
    this.nzDisabled = value
  }

  filterTableClick = (item: any) => {
    this.openDrawer('editFilter', item.data)
  }

  openDrawer (type: string, data?: any) {
    switch (type) {
      case 'addFilter': {
        this.editFilter = null
        break
      }
      case 'editFilter': {
        this.remoteSelectList = []
        this.remoteSelectNameList = []
        this.filterForm = Object.assign({}, data)
        if (data.values?.length === 1 && data.values[0] === 'ALL') {
          this.filterForm.allChecked = true
        }
        this.editFilter = Object.assign({}, data)
        if (data.name === 'ip') {
          this.filterForm.values = [this.filterForm.values.join('\n')]
          this.editFilter.values = [this.editFilter.values.join('\n')]
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

  filterDelete (item: any) {
    this.filterNamesSet.delete(item.name)
    for (const index in this.filterShowList) {
      if (this.filterShowList[index] === item) {
        this.filterShowList.splice(Number(index), 1)
        this.filterShowList = [...this.filterShowList]
      }
    }
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
    this.editFilter = null
    this.remoteSelectList = []
    this.remoteSelectNameList = []
    this.staticsList = []
  }
}
