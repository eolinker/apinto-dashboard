import { Component, EventEmitter, OnInit } from '@angular/core'
import { moveItemInArray, transferArrayItem, CdkDragDrop } from '@angular/cdk/drag-drop'
import { ApiService } from '../../service/api.service'
import { EoNgMessageService } from '../../service/eo-ng-message.service'
import { MiddlewaresItem } from './types/types'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { MODAL_SMALL_SIZE } from '../../constant/app.config'
import { EmptyHttpResponse } from '../../constant/type'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'

@Component({
  selector: 'eo-ng-interceptor',
  templateUrl: './interceptor.component.html',
  styles: [
    `
    div svg{
      color:var(--TIP_TEXT_COLOR);
    }`
  ]
})
export class InterceptorComponent implements OnInit {
  draggableTitle:string = ''
  draggableData:Array<any> = []
  groupData:Map<string, Array<any>> = new Map()
  groups:Array<{prefix:string, editing?:boolean}> = []
  groutDataChange:EventEmitter<Array<any>> = new EventEmitter()
  currentGroupData:Array<MiddlewaresItem> = []
  currentGroup:string = ''
  editingGroup:boolean = false
  editingCurrentMiddlewares:boolean = false
  totalMiddlewares:Array<MiddlewaresItem> = []
  constructor (private api:ApiService, private message:EoNgMessageService, private modalService:EoNgFeedbackModalService, private appConfigService: EoNgNavigationService) {
    this.appConfigService.reqFlashBreadcrumb([
      { title: '拦截器管理' }
    ])
  }

  ngOnInit (): void {
    this.getAllMiddlewares()
  }

  // 获取全部的拦截器信息，接口会返回每个分组对应的拦截器列表和所有拦截器信息列表
  getAllMiddlewares (reset?:boolean) {
    this.groups = []
    this.groupData = new Map()
    this.totalMiddlewares = []
    this.api.get('system/middleware')
      .subscribe((resp:{code:number, msg:string, data:{groups:Array<{prefix:string, middlewares:Array<string>}>, middlewares:Array<MiddlewaresItem>}}) => {
        if (resp.code === 0) {
          for (const group of resp.data.groups) {
            this.groups.push({ prefix: group.prefix })
            this.groupData.set(group.prefix, group.middlewares)
            this.totalMiddlewares = resp.data.middlewares
          }
          this.sortGroups()
          if (!this.currentGroup || reset) {
            this.currentGroup = this.groups[0].prefix
            this.getCurrentMiddlewares()
          }
        }
      })
  }

  // 对groups进行排序
  sortGroups () {
    this.groups.sort((a, b) => {
      return b.prefix.localeCompare(a.prefix)
    })
  }

  // 获取当前分组对应的拦截器完整信息与可用拦截器列表
  getCurrentMiddlewares () {
    this.currentGroupData = []
    this.draggableData = []
    if (!this.groupData.get(this.currentGroup)) {
      this.draggableData = this.totalMiddlewares
      return
    }
    for (const middleware of this.totalMiddlewares) {
      if (this.groupData.get(this.currentGroup)?.indexOf(middleware.name) !== -1) {
        this.currentGroupData.push(middleware)
      } else {
        this.draggableData.push(middleware)
      }
    }
  }

  drop (event: CdkDragDrop<MiddlewaresItem[]>) {
    if (event.previousContainer === event.container) {
      moveItemInArray(event.container.data, event.previousIndex, event.currentIndex)
    } else {
      transferArrayItem(
        event.previousContainer.data,
        event.container.data,
        event.previousIndex,
        event.currentIndex
      )
    }
  }

  changeCurrentGroup (group:{prefix:string, uuid?:string, editing?:boolean}) {
    if (this.currentGroup !== group.prefix) {
      this.groupData.set(this.currentGroup, this.currentGroupData.map((data:MiddlewaresItem) => {
        return data.name
      }))
      this.currentGroup = group.prefix
      this.getCurrentMiddlewares()
    }
  }

  deleteGroup (prefix:string, e?:Event) {
    e?.stopPropagation()
    this.modalService.create({
      nzTitle: '删除',
      nzContent: `确定要删除${prefix}吗？`,
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkText: '确定',
      nzOkDanger: true,
      nzCancelText: '取消',
      nzOnOk: () => {
        let index:number = -1
        for (let i = 0; i < this.groups.length; i++) {
          if (this.groups[i].prefix === prefix) {
            index = i
            break
          }
        }
        if (index !== -1) {
          this.groups.splice(index, 1)
        }
        if (this.currentGroup === prefix) {
          this.currentGroup = index < this.groups.length ? this.groups[index].prefix : this.groups[index - 1].prefix
        }
        this.groupData.delete(prefix)
        return true
      }
    })
  }

  addNewGroupDiv () {
    if (this.editingGroup) {
      return
    }
    this.editingGroup = true
    this.groups.push({ prefix: '', editing: true })
  }

  addNewGroup (group:any) {
    if (!group.prefix) {
      this.editingGroup = false
      group.editing = false
      this.groups.pop()
    }
    this.groupData.set(group.prefix, [])
    this.editingGroup = false
    group.editing = false
    this.sortGroups()
  }

  resetGroup () {
    this.getAllMiddlewares(true)
  }

  submit () {
    this.groupData.set(this.currentGroup, this.currentGroupData.map((data:MiddlewaresItem) => {
      return data.name
    }))
    const newGroup:Array<{prefix:string, middlewares:Array<string>}> = []
    for (const k of this.groupData.keys()) {
      newGroup.push({ prefix: k, middlewares: this.groupData.get(k) as string[] })
    }
    this.api.post('system/middleware', { groups: newGroup }).subscribe((resp:EmptyHttpResponse) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '更新拦截器成功')
        this.getAllMiddlewares()
      }
    })
  }
}
