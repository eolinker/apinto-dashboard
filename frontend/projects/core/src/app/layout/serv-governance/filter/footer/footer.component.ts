import { Component, EventEmitter, Input, Output } from '@angular/core'
import { Router } from '@angular/router'
import { CheckBoxOptionInterface } from 'eo-ng-checkbox'
import { FilterForm, FilterShowData } from '../../types/types'

@Component({
  selector: 'eo-ng-filter-footer',
  templateUrl: './footer.component.html',
  styles: [
  ]
})
export class FilterFooterComponent {
  @Input() // 双向绑定filterForm
  get filterForm () { return this._filterForm }

  set filterForm (val) {
    this._filterForm = val
    this.filterFormChange.emit(this._filterForm)
  }

  @Input()
  get filterShowList () { return this._filterShowList }

  set filterShowList (val) {
    this._filterShowList = val
    this.filterShowListChange.emit(this._filterShowList)
  }

  @Input()
  get filterNamesSet () { return this._filterNamesSet }

  set filterNamesSet (val) {
    this._filterNamesSet = val
    this.filterNamesSetChange.emit(this._filterNamesSet)
  }

  @Input() remoteSelectNameList: string[] = [] // 穿梭框内被勾选的选项name
  @Input() remoteSelectList: string[] = [] // 穿梭框内被勾选的选项uuid
  @Input() staticsList: CheckBoxOptionInterface[]= []
  @Input() editFilter?: FilterForm // 正在编辑的配置
  @Input() filterType: string = '' // 筛选条件类型, 当type=pattern,显示输入框, static显示一组勾选框, remote显示穿梭框
  @Output() filterFormChange:EventEmitter<FilterForm> = new EventEmitter()
  @Output() filterShowListChange:EventEmitter<FilterShowData [] > = new EventEmitter()
  @Output() filterNamesSetChange:EventEmitter<Set<string> > = new EventEmitter()
  @Output() drawerClose:EventEmitter<boolean> = new EventEmitter()
  submitButtonLoading:boolean = false
  strategyType:string = '' // 策略类型
  _filterNamesSet: Set<string> = new Set() // 穿梭框内被勾选的选项uuid
  _filterShowList: FilterShowData [] = []
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

  constructor (private router:Router) {
    this.strategyType = this.router.url.split('/')[this.router.url.split('/').indexOf('serv-governance') + 1]
  }

  saveFilter () {
    switch (this.filterType) {
      case 'remote': {
        this.filterForm.values = []
        this.filterForm.text = ''
        if (this.remoteSelectList.length === this.filterForm.total) {
          this.filterForm.values = ['ALL']
          this.filterForm.text = `所有${this.filterForm.title}`
        } else {
          this.filterForm.values = this.remoteSelectList
          this.filterForm.text = this.remoteSelectNameList.join(',')
        }
        break
      }
      case 'pattern': {
        this.filterForm.text = this.filterForm.values[0]
        break
      }
      default:
        this.filterForm.values = []
        this.filterForm.text = ''
        if (!this.filterForm.allChecked) {
          for (const index in this.staticsList) {
            if (this.staticsList[index].checked) {
              this.filterForm.values.push(this.staticsList[index].value)
            }
          }
          this.filterForm.text = this.filterForm.values.join(',')
        } else {
          this.filterForm.values = ['ALL']
          this.filterForm.text = `所有${this.filterForm.title}`
        }
        break
    }
    if (this.editFilter) {
      this.filterShowList = this.filterShowList.filter((item: FilterShowData) => {
        return item.name !== this.editFilter!.name
      })
      this.filterShowList = [
        ...this.filterShowList,
        {
          title: this.filterForm.title,
          name: this.filterForm.name,
          label: this.filterForm.text,
          values: this.filterForm.values
        }
      ]
      if (this.editFilter.name !== this.filterForm.name) {
        this.filterNamesSet.delete(this.editFilter.name)
        this.filterNamesSet.add(this.filterForm.name)
      }
    } else {
      this.filterShowList = [
        ...this.filterShowList,
        {
          title: this.filterForm.title,
          name: this.filterForm.name,
          label: this.filterForm.text,
          values: this.filterForm.values
        }
      ]
      this.filterNamesSet.add(this.filterForm.name)
    }
    this.drawerClose.emit(true)
  }

  // 是否禁用提交按钮
  checkFilterToSave (): boolean {
    switch (this.filterType) {
      case 'static': {
        return (
          !this.filterForm.allChecked && this.filterForm.values.length === 0
        )
      }
      case 'pattern': {
        return !this.filterForm.values[0]
          ? true
          : !this.filterForm.pattern
              ? false
              : !(this.filterForm.pattern && this.filterForm.patternIsPass)
      }
      default: {
        return !this.filterForm.allChecked && this.remoteSelectList.length === 0
      }
    }
  }
}
