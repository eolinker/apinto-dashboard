/* eslint-disable dot-notation */
import { Component, ElementRef, EventEmitter, HostListener, Input, OnInit, NgZone, Output, SimpleChanges, TemplateRef, ViewChild, ChangeDetectorRef } from '@angular/core'
import { NzTableLayout } from 'ng-zorro-antd/table'
import { EoNgTableComponent } from 'eo-ng-table'
import { Router } from '@angular/router'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
@Component({
  selector: 'eo-ng-apinto-table',
  templateUrl: './table.component.html',
  styles: [
    `
    :host ::ng-deep{

      .ant-table-tbody > tr:hover .opacity-0{
        opacity:1 !important;
      }

      eo-ng-table.cursorPointer tbody tr:not(.ant-table-placeholder){
        cursor:pointer;
      }
    }`
  ]
})
export class TableComponent extends EoNgTableComponent implements OnInit {
  @ViewChild('theadTpl', { read: TemplateRef, static: true }) theadTpl: TemplateRef<any> | undefined
  @ViewChild('tbodyTpl', { read: TemplateRef, static: true }) tbodyTpl: TemplateRef<any> | undefined
  @Input() override nzPageSizeOptions:Array<number>=[20, 50, 100]
  @Input() nzMaxOperatorButton:number = 0
  @Input() nzNoScroll:boolean = false
  @Input() override nzTableLayout :NzTableLayout = 'fixed'
  @Input() nzScrollY:number = 0
  @Input() nzDisabled:boolean = false
  @Input() nzTableTheadConfig:Array<any> = []
  @Input() nzMonitorDT:boolean = false
  @Output() nzChangeTableConfigChange:EventEmitter<{value:any, item:any}> = new EventEmitter()

  scrHeight:any;
  scrWidth:any;
  tableScrollCdk:any
  cursorPointer:boolean = false
  constructor (ngZone:NgZone, cdr:ChangeDetectorRef, private message: EoNgFeedbackMessageService, private el:ElementRef, private router:Router) {
    super(ngZone, cdr)
  }

  @HostListener('window:resize', ['$event'])
  getScreenSize (event?:any) {
    this.scrHeight = event.target.innerHeight
    this.scrWidth = event.target.innerWidth
    this.getSrollY()
  }

  override ngOnInit (): void {
    super.ngOnInit()

    this.getThead()
    this.getTbody()
    !this.nzNoScroll && !this.nzScroll.x && this.calculateScroll()
    this.scrHeight = window.innerHeight
    setTimeout(() => {
      this.getSrollY()
    }, 0)
    console.log(this)
  }

  override ngOnChanges (change:SimpleChanges) {
    super.ngOnChanges(change)
    if (change['nzData'] || change['nzScrollY']) {
      this.getSrollY()
    }
    if (change['nzTbody']) {
      this.getTbody()
    }
    if (change['nzThead']) {
      this.getThead()
    }
    if (change['nzTrClick']) {
      if (this.nzTrClick) {
        this.cursorPointer = true
      }
    }
  }

  getThead () {
    for (const index in this.nzThead) {
      if (this.nzThead[index].title === '操作') {
        this.nzThead[index].right = true
        this.nzThead[index].width = (this.nzMaxOperatorButton && !this.nzThead[index].width) ? (this.nzMaxOperatorButton > 3 ? this.nzMaxOperatorButton * 26 + 10 + 10 : 3 * 26 + 20) : (this.nzThead[index].width > 50 ? this.nzThead[index].width : 50)
        this.nzThead[index].resizeable = false
      } else if (index === (this.nzThead.length - 1).toString()) {
        this.nzThead[index].resizeable = false
      }

      if (this.nzThead[index].title !== '操作' && (this.nzThead[index].required || this.nzThead[index].tooltip)) {
        this.nzThead[index].title = this.theadTpl
      }
    }
  }

  getTbody (): void {
    for (const body of this.nzTbody) {
      if (!body.title && !body.type) {
        body.title = this.tbodyTpl
      }
      if (body?.type === 'btn') {
        for (const btn of body.btns) {
          btn.icon = this.getIconName(btn.title)
        }
      }
    }
  }

  getIconName (title:string):string {
    switch (title) {
      case '编辑': {
        return 'edit' }
      case '修改': {
        return 'edit'
      }
      case '配置': {
        return 'tool'
      }
      case '查看': {
        return 'find'
      }
      case '减少': {
        return 'reduce-one'
      }
      case '删除': {
        return 'delete'
      }
      case '移除' : {
        return 'delete'
      }
      case '添加': {
        return 'add-circle'
      }
      case '重置': {
        return 'refresh'
      }
      case '重置密码': {
        return 'refresh'
      }
      case '测试': {
        return 'lightning'
      }
      case '恢复': {
        return 'undo'
      }
      case '更新': {
        return 'refresh'
      }
      case '更新鉴权Token': {
        return 'refresh'
      }
      case '提交': {
        return 'check'
      }
      case '取消': {
        return 'close'
      }
      case '上线管理': {
        return 'shangxianguanli-new'
      }
      case '上线': {
        return 'circle-right-up-7mnlo5g9'
      }
      case '下线': {
        return 'circle-right-down-7mnlphn2'
      }
      case '禁用': {
        return 'pause'
      }
      case '启用': {
        return 'play'
      }
      case '复制Token': {
        return 'copy'
      }
    }
    return ''
  }

  transferToJson (str:string) {
    return str.replace(/(,\n)/g, ',').replace(/(,\n)/g, ',').replace(/(,)/g, ',\n').replace(/(，)/g, '，\n')
  }

  calculateScroll () {
    this.nzScroll.x = '100%'
  }

  // 动态计算虚拟滚动高度，需要判断表格是否在弹窗内，如果在弹窗内则无导航高度
  getSrollY () {
    let scrollY:number = 0
    if (this.nzMonitorDT) {
      const headerHeight = document.getElementsByClassName('list-header')[0]?.getBoundingClientRect().height || 0// list-header高度，通常为表格上方按钮
      const pagationHeight = this.nzShowPagination ? 48 : 0 // 分页高度
      scrollY = this.scrHeight - headerHeight - pagationHeight - 4 - 12 - 20 - 50 - 40
    } else {
      if (this.nzScrollY) {
        scrollY = this.nzScrollY
      } else {
        const navTop = document.getElementsByTagName('eo-ng-tabs').length > 0 ? 90 : 50 // 列表上方，通常为导航header高度加上tab高度

        const headerHeight = document.getElementsByClassName('list-header')[0]?.getBoundingClientRect().height || 0// list-header高度，通常为表格上方按钮
        const footerHeight = document.getElementsByClassName('list-footer')[0]?.getBoundingClientRect().height || 0 // list-footer高度，通常为表格下方按钮
        const pagationHeight = this.nzShowPagination ? 42 : 0 // 分页高度

        const drawerHeaderHeight = document.getElementsByClassName('drawer-list-header')[0]?.getBoundingClientRect().height || 0// list-header高度，通常为表格上方按钮
        const drawerFooterHeight = document.getElementsByClassName('drawer-list-footer')[0]?.getBoundingClientRect().height || 0 // list-footer高度，通常为表格下方按钮
        const drawerPagationHeight = this.nzShowPagination ? 42 : 0 // 分页高度
        const drawerButtonAreaHeight = document.getElementsByClassName('ant-modal-footer')[0]?.getBoundingClientRect().height || 0 // 弹窗底部

        const clusterDescHeight = document.getElementsByClassName('cluster-desc-block')[0]?.getBoundingClientRect().height || 0 // 集群环境变量里的集群描述高度
        const monitorTabsHeight = document.getElementsByClassName('monitor-total-content').length > 0 ? 42 : 0 // 监控分区内tabs高度
        const monitorTabsPieHeight = document.getElementsByClassName('eo-ng-monitor-detail-pie').length > 0 ? (document.getElementsByClassName('eo-ng-monitor-detail-pie')[0]?.getBoundingClientRect().height + 32 + 32 + 20 + 10) : 0 // 监控分区内详情页面饼图及tabs高度
        // 表格在弹窗内，需要减去弹窗标题高度53，弹窗顶部100px，弹窗内部上下padding20*2， 表头高度； 否则减去表格顶部间隙12px，底部间隙20px，表头高度
        scrollY = this.el.nativeElement.classList.contains('drawer-table') ? (this.scrHeight > 660 ? 660 - drawerHeaderHeight - drawerFooterHeight - drawerPagationHeight - 40 - 40 - 20 - 63 - drawerButtonAreaHeight : 660 - drawerHeaderHeight - drawerFooterHeight - drawerPagationHeight - 40 - 40 - 14 - 63 - 69 - 150) : (this.scrHeight - navTop - headerHeight - footerHeight - pagationHeight - clusterDescHeight - monitorTabsHeight - monitorTabsPieHeight - 14 - 6 - 40 - 1 + 10)
      }
    }

    if (this.router.url.includes('monitor-alarm/area/history')) {
      scrollY = scrollY - 42
    }

    if (this.nzData !== undefined && scrollY < this.nzData.length * 40) {
      this.nzScroll = { x: this.nzNoScroll || this.nzData.length === 0 ? undefined : this.nzScroll.x || '100%', y: scrollY > 50 ? scrollY + 'px' : '50px' }
      this.tableScrollCdk?.ngOnInit()
    } else {
      this.nzScroll = { x: this.nzNoScroll || this.nzData.length === 0 ? undefined : this.nzScroll.x || '100%', y: undefined }
    }
  }

  stopTrClick (index:number, length:number, e:any) {
    if (index === length - 1) {
      e?.stopPropagation()
    }
  }

  handlerScrollView ($event:any) {
    this.tableScrollCdk = $event
    this.getScrollViewPort.emit($event)
  }

  copyCallback () {
    this.message.success('复制成功', {
      nzDuration: 1000
    })
  }
}
