/* eslint-disable no-useless-constructor */
// api上线、应用上线、服务发现上线、上游服务上线报错时，需要打开弹窗并支持打开新窗口跳转到后端传来的新链接
import { Injectable } from '@angular/core'
import { EoNgFeedbackModalService } from 'eo-ng-feedback'
import { MODAL_SMALL_SIZE } from '../constant/app.config'

@Injectable({
  providedIn: 'root'
})
export class PublishFailService {
  constructor (
    private modalService:EoNgFeedbackModalService) { }

  openModal (msg:string, type:string, routerName:string, routerParam:{[k:string]:any}, footer?:any) {
    this.modalService.create({
      nzTitle: '提示',
      nzIconType: 'exclamation-circle',
      nzContent: `${msg}，请点击跳转至相关链接。`,
      nzClosable: true,
      nzCancelText: '取消',
      nzWidth: MODAL_SMALL_SIZE,
      nzOkText: '跳转',
      nzOnOk: () => {
        this.viewSolution(routerName, routerParam)
      }
    })
  }

  openFooterModal (msg:string, type:string, footer?:any) {
    return this.modalService.create({
      nzTitle: '提示',
      nzIconType: 'exclamation-circle',
      nzContent: `${msg}，请点击跳转至相关链接。`,
      nzClosable: true,
      nzWidth: MODAL_SMALL_SIZE,
      nzFooter: footer
    })
  }

  viewSolution (solutionRouter:string, solutionParam:{[k:string]:any}) {
    const routerS:string = '/' + solutionRouter
    const routerArr:Array<string> = routerS.split('/')
    const contentIndex:number = routerArr.indexOf('content') + 1
    if (Object.keys(solutionParam).length > 0) {
      for (const index in Object.keys(solutionParam)) {
        routerArr.splice(contentIndex, 0, solutionParam[Object.keys(solutionParam)[index]])
      }
    }
    window.open(routerArr.join('/'))
  }
}
