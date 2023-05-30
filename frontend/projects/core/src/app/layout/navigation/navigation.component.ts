import { Component, OnInit } from '@angular/core'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { MODAL_SMALL_SIZE } from '../../constant/app.config'
import { NavigationCreateComponent } from './create/create.component'
import { moveItemInArray, CdkDragDrop } from '@angular/cdk/drag-drop'
import { ApiService } from '../../service/api.service'
import { NavigationItem } from './types/types'
import { EmptyHttpResponse } from '../../constant/type'
import { NzModalRef } from 'ng-zorro-antd/modal'
import { EoNgMessageService } from '../../service/eo-ng-message.service'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'

@Component({
  selector: 'eo-ng-navigation',
  templateUrl: './navigation.component.html',
  styles: [
    `:host ::ng-deep{
      .ant-list-bordered{
        border:1px solid var(--border-color);
      }

      .ant-list-split .ant-list-item{
        border-bottom:1px solid var(--border-color);
      }

      .draggable-list {
        display: block;
        border-radius: 4px;
        overflow: hidden;
      }

      .draggable-item {
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: space-between;
        cursor: move;
      }

      .cdk-drag-placeholder {
        opacity: 0;
      }

      .cdk-drag-animating {
        transition: transform 250ms cubic-bezier(0, 0, 0.2, 1);
      }

      .draggable-item:last-child {
        border: none;
      }
    }

    `
  ]
})
export class NavigationComponent implements OnInit {
  navigationList:Array<NavigationItem > = []
  modalRef:NzModalRef|undefined

  constructor (private modalService:EoNgFeedbackModalService, private api:ApiService, private message:EoNgMessageService,
    private navigationService: EoNgNavigationService) {
    this.navigationService.reqFlashBreadcrumb([
      { title: '导航管理' }
    ])
  }

  ngOnInit () {
    this.getNavigationList()
  }

  getNavigationList () {
    this.api.get('system/navigation').subscribe((resp:{code:number, msg:string, data:{navigations:Array<NavigationItem>}}) => {
      if (resp.code === 0) {
        this.navigationList = resp.data.navigations
      }
    })
  }

  disableNavigationModal (uuid:string) {
    this.modalRef = this.modalService.create({
      nzTitle: '删除',
      nzContent: '该数据删除后将无法找回，请确认是否删除？',
      nzClosable: true,
      nzClassName: 'delete-modal',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkText: '确定',
      nzOkDanger: true,
      nzCancelText: '取消',
      nzOnOk: () => {
        this.disableNavigation(uuid)
        return false
      }
    })
  }

  // 删除分组
  disableNavigation (uuid:string) {
    this.api.delete(`system/navigation/${uuid}`).subscribe((resp:EmptyHttpResponse) => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '删除成功')
        this.modalRef?.close()
      }
    })
  }

  openNavigationModal (add:boolean, uuid?:string) {
    this.modalRef = this.modalService.create({
      nzTitle: add ? '新建导航' : '编辑导航',
      nzWidth: MODAL_SMALL_SIZE,
      nzContent: NavigationCreateComponent,
      nzComponentParams: { ...(!add ? { editPage: !add, navigationUuid: uuid } : {}), modalRef: this.modalRef },
      nzOnOk: (component:NavigationCreateComponent) => {
        component.submit()
        return false
      }
    })
  }

  drop (event: CdkDragDrop<string[]>) {
    const tmpNavigationList:Array<NavigationItem> = [...this.navigationList]
    moveItemInArray(tmpNavigationList, event.previousIndex, event.currentIndex)
    this.api.put('system/navigation', {
      navigations: tmpNavigationList.map((nav:NavigationItem) => {
        return nav.uuid
      })
    }).subscribe((resp:EmptyHttpResponse) => {
      if (resp.code === 0) {
        moveItemInArray(this.navigationList, event.previousIndex, event.currentIndex)
        this.message.success(resp.msg || '操作成功')
      }
    })
  }
}
