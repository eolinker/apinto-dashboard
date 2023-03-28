/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable no-useless-constructor */
/* eslint-disable dot-notation */
import { Component, Input, OnInit, Output, TemplateRef, ViewChild, EventEmitter, ChangeDetectionStrategy, ChangeDetectorRef, SimpleChanges } from '@angular/core'
import { Validators } from '@angular/forms'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { Subscription } from 'rxjs'
import { MODAL_NORMAL_SIZE } from '../../constant/app.config'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { EditableEnvTableComponent } from '../editable-env-table/editable-env-table.component'
import { ApiService } from '../../service/api.service'

@Component({
  selector: 'dynamic-component',
  templateUrl: './dynamic-component.component.html',
  styles: [
    `

    .margin-dynamic-row{
      height:12px;
    }

    .ant-form-item-extra {
      min-height:16px;
    }

    .ant-form-item-extra + .margin-dynamic-row{
      height:0px;
    }

    .error-input > .ant-input{
      border-color:rgb(255, 60, 50);
      color:rgb(255, 60, 50);
    }
    `
  ],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class DynamicComponentComponent implements OnInit {
  public _itemList:any[] =[]
  @Input() baseData:any = null
  @Input() nzDisabled:boolean = false
  public _baseData:any = null
  public dataList:any = {}
  public itemList:any = []
  public _apiData:any = []
  @Input() startValidate:boolean = false
  @Input()
  get apiData () {
    return this._apiData
  }

  set apiData (val) {
    this._apiData = val
    this.apiDataChange.emit(this.apiData)
  }

  @Output() apiDataChange = new EventEmitter()

  constructor (
    private modalService:EoNgFeedbackModalService,
     private changeDetectorRef:ChangeDetectorRef,
     private api:ApiService) { }

  ngOnInit (): void {
  }

  ngOnChanges (changes:SimpleChanges): void {
    if (this.baseData) {
      this.itemList = []
      this.itemList = this.checkDataForItemList(this.dataTransfer(this.baseData))
      this.fillDataForApiData(this.itemList)
    }
    if (changes['startValidate']) {
      this.changeDetectorRef.markForCheck()
      this.changeDetectorRef.detectChanges()
    }
    // this.buildForm()
  }

  // dynamicForm:FormGroup|undefined
  rule:any = {}

  show (value:any) {
    // console.log(value)
  }

  valid (value:any, pattern:RegExp) {
    return new RegExp(pattern).test(value)
  }

  checkShow (items:any):boolean {
    if (!items['x-reactions']) {
      return true
    } else {
      if (items['x-reactions'].otherwise) {
        if (items['name'] !== 'secret' && items['name'] !== 'publicKey') {
          return !this.seachForDeps(items['x-reactions'].dependencies[0])
        } else if (items['name'] === 'secret') {
          const temres = this.seachForDeps(items['x-reactions'].dependencies[0])
          const tesresstring:string = (typeof temres === 'boolean') ? '' : temres
          return !tesresstring?.includes('HS')
        } else if (items['name'] === 'publicKey') {
          const temres = this.seachForDeps(items['x-reactions'].dependencies[0])
          const tesresstring:string = (typeof temres === 'boolean') ? '' : temres
          return !(tesresstring?.includes('RS') || tesresstring?.includes('ES'))
        }
      } else {
        if (items['name'] !== 'secret' && items['name'] !== 'publicKey') {
          return !!this.seachForDeps(items['x-reactions'].dependencies[0])
        } else if (items['name'] === 'secret') {
          const temres = this.seachForDeps(items['x-reactions'].dependencies[0])
          const tesresstring:string = (typeof temres === 'boolean') ? '' : temres
          return tesresstring?.includes('HS')
        } else if (items['name'] === 'publicKey') {
          const temres = this.seachForDeps(items['x-reactions'].dependencies[0])
          const tesresstring:string = (typeof temres === 'boolean') ? '' : temres
          return tesresstring?.includes('RS') || tesresstring?.includes('ES')
        }
      }
    }
    return true
  }

  checkDataForItemList (itemList:any) {
    if (itemList?.length > 0) {
      itemList.forEach((item:any) => {
        item.default = this.searchForData(item.name, this.apiData) || item.default
        if (item.properties?.length > 0) {
          item.properties = this.checkDataForItemList(item.properties)
        }
        if (item.type === 'array') {
          let diffNum = item.default?.length - (item.properties?.length)
          // const itemProperties = this.dataTransfer(item.items.properties)
          while (diffNum > 0) {
            const itemProperties = JSON.parse(JSON.stringify(item.items))
            itemProperties.properties = this.dataTransfer(item.items)
            itemProperties['x-index'] = item.properties.length
            // console.log(item.properties.length)
            // console.log(itemProperties)
            // item.items.properties = itemProperties
            item.properties.push(itemProperties)
            diffNum--
          }
          item.properties = this.fillDataForArray(item.properties, item)
        }
      })
    }
    return itemList
  }

  fillDataForArray (propertiesList:any, item:any) {
    const valueList = item.default
    switch (Object.prototype.toString.call(valueList)) {
      case '[object Array]': {
        for (let i = 0; i < valueList?.length; i++) {
          propertiesList[i].properties.forEach((el:any) => {
            if (el.name === item.name) {
              el.default = valueList[i]
            } else {
              el.default = valueList[i][el.name]
            }
          })
        }
        break }
      default: {
        for (let i = 0; i < valueList?.length; i++) {
          propertiesList[i].properties.forEach((el:any) => {
            el.default = valueList[i][el.name]
          })
        }
      }
    }
    return propertiesList
  }

  searchForData (key:string, apiData:any) :any {
    if (apiData) {
      const keyLists = Object.keys(apiData)
      for (let i = 0; i < keyLists.length; i++) {
        if (keyLists[i] === key) {
          return apiData[keyLists[i]]
        } else if (typeof apiData[keyLists[i]] === 'object') {
          const tmpRes = this.searchForData(key, apiData[keyLists[i]])
          if (tmpRes) {
            return tmpRes
          }
        }
      }
    }
    return null
  }

  dataTransfer (data:any): Array<any> {
    let itemList:Array<any> = []
    const baseData = JSON.parse(JSON.stringify(data))
    // console.log(baseData)
    if (baseData?.properties) {
      Object.keys(baseData?.properties).forEach(key => {
        const tempbaseData = JSON.parse(JSON.stringify(baseData))
        const items = tempbaseData.properties[key]
        items['name'] = key
        items['default'] = items['default'] || null
        // if (key === 'addr') {
        //   console.log(items)
        // }
        if (items['required'] !== null || items['minimum'] !== null || items['maximum'] !== null) {
          this.rule[key] = [items['default'], []]
          if (items['required'] !== null) {
            this.rule[key][1].push(Validators.required)
          }
          if (items['minimum'] !== null) {
            this.rule[key][1].push(Validators.minLength(items['minimum']))
          }
          if (items['maximum'] !== null) {
            this.rule[key][1].push(Validators.maxLength(items['maximum']))
          }
        }
        if (items['x-component'] === 'Input') {
          items.minimum = items.minimum || -Infinity
          items.maximum = items.maximum || Infinity
          items.minLength = items.minLength || -Infinity
          items.maxLength = items.maxLength || Infinity
        }
        let temP:any = null
        if (items.properties) {
          temP = this.dataTransfer(items)
        }
        items.properties = temP
        itemList.push(items)
      })
    }
    itemList = itemList.sort(this.compare('x-index'))
    return itemList
  }

  seachForDeps (itemName:string):string | boolean {
    // this.itemList = this.itemList
    for (let i = 0; i < this.itemList.length; i++) {
      const res = this.seachForDepsDeep(this.itemList[i], itemName)
      if (res) {
        return res
      }
    }
    return false
  }

  seachForDepsDeep (itemList:any, itemName:string):any {
    if (itemList.properties) {
      for (let i = 0; i < itemList.properties.length; i++) {
        if (itemList.properties[i].name === itemName || itemList.properties[i].name === itemName.replace(/_([a-z])/g, (p, m) => m.toUpperCase())) {
          return itemList.properties[i].default
        } else if (itemList.properties[i].properties) {
          return this.seachForDepsDeep(itemList.properties[i], itemName)
        }
      }
      return false
    } else {
      if (itemList.name === itemName || itemList.name === itemName.replace(/_([a-z])/g, (p, m) => m.toUpperCase())) {
        return itemList.default
      }
    }
  }

  compare (property:string) {
    return (o1:any, o2:any) => o1[property] - o2[property]
  }

  flushData () {
    this.fillDataForApiData(this.itemList)
    this.apiDataChange.emit(this.apiData)
  }

  fillDataForApiData (data:any) {
    data?.forEach((el:any) => {
      if (el.type !== 'object') {
        this.updateApiData(el.name, el, el.type, this.apiData)
      }
      if (el.type !== 'array' && el.properties?.length > 0) {
        this.fillDataForApiData(el.properties)
      }
      if (el.type === 'array' && el.properties?.length > 0) {
        el.properties?.properties?.forEach((eel:any) => {
          this.updateApiData(eel.name, eel.default, eel.type, this.apiData)
        })
      }
    })
  }

  // 将value里的值更新到对应data
  updateApiData (key:string, value:any, type:string, data:any) {
    const keyLists = Object.keys(data)
    // console.log(key, value, type, data, keyLists)
    // 遍历data里的key值，找到与value同名值
    for (let i = 0; i < keyLists.length; i++) {
      if (keyLists[i] === key) {
        // console.log(data[key], type)
        if (type === 'array' && data[key] !== undefined) {
          // console.log(data[key].length, [data[key][0]])
          if (data[key].length > 1) {
            const temp = data[key][0] === undefined ? '' : JSON.parse(JSON.stringify(data[key][0]))
            data[key] = [temp]
            // data[key].push(temp)
          }
          // console.log(data[key].length, [data[key][0]])
          const dataLength = data[key].length
          for (let j = 0; j < value.properties.length; j++) {
            // 当data是对象数组，且数组长度不足时，需要补齐长度
            if (j >= dataLength && data[key][0] && Object.prototype.toString.call(data[key][0]) === '[object Object]') {
              const _tmpData = JSON.parse(JSON.stringify(data[key][0]))
              data[key].push(_tmpData)
            }
            value.properties[j].properties.forEach((property:any) => {
              if (data[key][j] !== null && Object.prototype.toString.call(data[key][j]) === '[object Object]') {
                Object.keys(data[key][j]).forEach((dataInner:any) => {
                  if (property.name === dataInner) {
                    data[key][j][dataInner] = property.default
                  }
                })
              } else {
                if (property.name === key && typeof property.default !== 'object') {
                  data[key][j] = property.default
                }
              }
            })
          }
          return
        } else {
          data[keyLists[i]] = value.default
          return
        }
      } else if (keyLists[i] !== '0' && data[keyLists[i]] && Object.prototype.toString.call(data[keyLists[i]]) === '[object Object]') {
        // console.log('0')
        // console.log(keyLists[i], data[keyLists[i]])
        this.updateApiData(key, value, type, data[keyLists[i]])
      }
    }
  }

  addNewProperies (item:any, parent: any, items?:any) {
    // console.log(item, parent, items)
    if (items) {
      const _items = JSON.parse(JSON.stringify(items))
      _items.properties = this.dataTransfer(_items)
      const addListNum = item.properties.length - 1
      if (item.name !== 'params' && item.name !== 'label') {
        _items.name = item.properties[addListNum].name.replace(`${item.properties[addListNum]['x-index']}/`, '$1') + (item.properties[addListNum].properties.length)
        _items['x-index'] = item.properties[addListNum].properties.length
        item.properties[addListNum].properties.push(_items)
      } else {
        if (_items.name) {
          _items.name = item.properties[addListNum].name.replace(`${item.properties[addListNum]['x-index']}/`, '$1') + (addListNum)
        } else {
          _items.name = item.name + (addListNum) + ''
        }

        _items['x-index'] = (addListNum)
        // console.log(_items)
        item.properties.push(_items)
        // console.log(item)
      }
    } else if (item.items) {
      const newitem = JSON.parse(JSON.stringify(item.items))
      // console.log(newitem)
      item.properties.upshift(newitem)
      // console.log(item)
    }
  }

  checkDrag (item:any) {
    try {
      for (let i = 0; i < item.length; i++) {
        if (item[i]['x-component'].includes('SortHandle')) {
          return true
        }
      }
      return false
    } catch (e:any) {
      console.log(e.message)
      return false
    }
  }

  submit () {
  }

  change (value:any) {
  }

  remove (parent:any, index:number) {
    // console.log(parent, index)
    // console.log(this.apiData)
    this.flushData()
    if (parent.name !== 'params' && parent.name !== 'label') {
      const removeRow = parent.properties.length - 1
      parent.properties[removeRow].properties.splice(index, 1)
      parent.properties[removeRow].properties = [...parent.properties[removeRow].properties]
      this.removeDataFromApidata(parent.properties[removeRow].name, index, this.apiData)
    } else {
      parent.properties.splice(index, 1)
      parent.properties = [...parent.properties]
      this.removeDataFromApidata(parent.name, index, this.apiData)
    }
  }

  removeDataFromApidata (name:string, index:number, data:any) {
    // console.log(name, index, data)
    for (const key in data) {
      if (key === name) {
        if (Object.prototype.toString.call(data[key]) === '[object Array]') {
          data[key].splice(index, 1)

          this.apiDataChange.emit(this.apiData)
        }
      } else {
        if (Object.prototype.toString.call(data[key]) === '[object Object]') {
          this.removeDataFromApidata(name, index, data[key])
        }
      }
    }
  }

  drop (value:any) {
  }

  drawerAddRef:NzModalRef | undefined
  envNameForSear:string = ''
  public propertyWaitForChoose:any = null
  private subscription: Subscription = new Subscription()

  openDrawer (property:any) {
    this.propertyWaitForChoose = property

    this.drawerAddRef = this.modalService.create({
      nzTitle: '添加环境变量',
      nzContent: EditableEnvTableComponent,
      nzClosable: true,
      nzWidth: MODAL_NORMAL_SIZE,
      nzFooter: null,
      nzComponentParams: {
        ...property,
        chooseEnv: ($event:Event) => { this.chooseEnv($event) }
      }

    })
    this.subscription = this.drawerAddRef.afterClose.subscribe(() => {
      this.envNameForSear = ''
      this.propertyWaitForChoose = null
    })
  }

  ngOnDestroy () {
    this.subscription.unsubscribe()
  }

  chooseEnv = (item:any) => {
    if (!item?.data?.editing) {
      // this.createUpstreamForm.config.staticConf[0].addr = '${' + item.key + '}'
      // this.itemList[0].properties[1].properties[0].properties[0].default = '${' + item.data.key + '}'
      this.fillEnv(this.propertyWaitForChoose.properties, item.data.key)
      this.changeDetectorRef.markForCheck()
      this.changeDetectorRef.detectChanges()
      this.flushData()
      this.drawerAddRef?.close()
    }
  }

  fillEnv = (properties:any, key:string) => {
    for (let i = 0; i < properties.length; i++) {
      if (properties[i].name === 'value' || properties[i].name === 'addrsVariable') {
        properties[i].default = '${' + key + '}'
      } else if (properties.properties) {
        this.fillEnv(properties, key)
      }
    }
  }
}
