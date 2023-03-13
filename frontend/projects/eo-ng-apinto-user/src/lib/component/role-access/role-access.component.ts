/* eslint-disable dot-notation */
import { Component, forwardRef, Inject, Input, OnInit, TemplateRef, ViewChild } from '@angular/core'
import { ControlValueAccessor, NG_VALUE_ACCESSOR } from '@angular/forms'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { API_SERVICE_ADAPTER, ApiServiceAdapter } from '../../constant/api-service-adapter'

@Component({
  selector: 'eo-ng-apinto-role-access',
  templateUrl: './role-access.component.html',
  providers: [
    {
      provide: NG_VALUE_ACCESSOR,
      useExisting: forwardRef(() => RoleAccessComponent),
      multi: true
    }
  ]
})
export class RoleAccessComponent implements OnInit, ControlValueAccessor {
  @ViewChild('operatorTpl', { read: TemplateRef, static: true }) operatorTpl: TemplateRef<any> | undefined
  @Input() accessSet:Set<string> = new Set()
  @Input() disabled :boolean = false
  accessThead:Array<any> = [
    { title: '一级功能模块' },
    { title: '二级功能模块' },
    { title: '权限配置' }
  ]

  accessTbody:Array<any> = [
    { key: 'titleParent' },
    { key: 'titleChild' },
    {
      key: ''
    }
  ]

  listOfAccess:Array<any> = []
  accessMap:Map<string, Array<string>> = new Map() // 被依赖的权限与依赖其的权限字段名
  // eslint-disable-next-line no-useless-constructor
  constructor (private message: EoNgFeedbackMessageService,
    @Inject(API_SERVICE_ADAPTER) private apiService: ApiServiceAdapter) { }

  ngOnInit (): void {
    this.getAccessList()
  }

  ngAfterViewInit () {
    this.accessTbody[2].title = this.operatorTpl
  }

  getAccessList () {
    this.apiService.get('access').subscribe((resp:any) => {
      if (resp.code === 0) {
        const list = resp.data.modules
        for (const index in list) {
          this.getMapForAccess(list[index])
          this.changeListAccessStatus()
          this.listOfAccess = [...this.listOfAccess]
        }
      } else {
        this.message.error(resp.msg || '获取权限列表失败!')
      }
    })
  }

  // 分析每个模块,并创建映射
  getMapForAccess (module:any) {
    if (module.children?.length > 0) {
      for (const indexChild in module.children) {
        const childModule = module.children[indexChild]
        this.listOfAccess.push({
          titleParent: indexChild === '0' ? module.title : '',
          titleChild: childModule.title,
          access: childModule.access
        })

        // 建立被依赖映射
        this.mapForAccess(childModule.access)
      }
    } else {
      this.listOfAccess.push({
        titleParent: module.title,
        access: module.access
      })
      // 建立被依赖映射
      this.mapForAccess(module.access)
    }
  }

  changeAccessSet (value:any) {
    // 取消勾选,同时取消依赖其的权限
    if (!value.checked) {
      this.accessSet.delete(value.key)
      const depend:any = this.accessMap.get(value.key)
      for (const index in depend) {
        this.accessSet.delete(depend[index])
      }
    } else {
      // 确认其依赖的权限是否被勾选,如未勾选则自动勾选
      if (value?.dependencies?.length > 0) {
        for (const index in value.dependencies) {
          this.accessSet.add(value.dependencies[index])
        }
      }
      this.accessSet.add(value.key)
    }
    // 根据最新的accessSet内容,更新列表以供页面显示
    this.changeListAccessStatus()
  }

  // 根据最新的accessSet内容,更新列表以供页面显示
  changeListAccessStatus () {
    for (const index in this.listOfAccess) {
      for (const indexAcc in this.listOfAccess[index].access) {
        const acc = this.listOfAccess[index].access
        acc[indexAcc].checked = this.accessSet.has(acc[indexAcc].key)
      }
    }
  }

  // 建立被依赖映射
  mapForAccess (accessList:any) {
    for (const index in accessList) {
      for (const indexA in accessList[index]?.dependencies) {
        const depend:any = accessList[index]?.dependencies[indexA]
        if (this.accessMap.get(depend)) {
          this.accessMap.get(depend)?.push(accessList[index].key)
        } else {
          this.accessMap.set(accessList[index]?.dependencies[indexA], [accessList[index].key])
        }
      }
    }
  }

  onChange: any = () => { };
  onTouch: () => void = () => null;

  // 封装组件搭配form的formControlName 使用
  writeValue (obj: any): void {
    this.accessSet = obj
    this.changeListAccessStatus()
  }

  registerOnChange (fn: any): void {
    this.onChange = fn
  }

  registerOnTouched (fn: any): void {
    this.onTouch = fn
  }

  setDisabledState (isDisabled:boolean):void {
    this.disabled = isDisabled
  }
}
