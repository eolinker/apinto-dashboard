import { AfterViewInit, Directive, ElementRef, HostListener, Input, Renderer2 } from '@angular/core'

@Directive({
  selector: '[eoNgScrollDom]'
})
export class EoNgScrollDomDirective implements AfterViewInit {
  @Input() eoBeforeDom : HTMLDivElement|undefined // 位于滚动元素上方、不参与滚动、与滚动元素为兄弟节点的dom，如果无该元素则滚动元素的高度=父元素高度
  @Input() eoAfterDom : HTMLDivElement|undefined // 位于滚动元素下方、不参与滚动、与滚动元素为兄弟节点的dom，如果无该元素则滚动元素的高度=父元素高度
  @Input() eoParentDom : HTMLDivElement|undefined // 父元素
  domMaxHeight:number = 0 // 滚动元素的最大高度 = 父元素高度 - lastDom高度
  resizeObserver = new ResizeObserver((entries:any) => {
    for (const entry of entries) {
      if (entry.contentBoxSize) {
        this.getElementHeight()
      }
    }
  })

  constructor (
    private el:ElementRef, private renderer2:Renderer2) {
  }

    @HostListener('scroll', ['$event'])
  public onScroll () {
    if (this.eoBeforeDom && this.el.nativeElement.scrollTop) {
      this.renderer2.addClass(this.eoBeforeDom, 'scroll-top-box-shadow')
    } else {
      this.eoBeforeDom && this.renderer2.removeClass(this.eoBeforeDom, 'scroll-top-box-shadow')
    }
    if (this.eoAfterDom && (this.el.nativeElement.scrollHeight - this.el.nativeElement.scrollTop) !== this.el.nativeElement.clientHeight) {
      this.renderer2.addClass(this.eoAfterDom, 'scroll-bottom-box-shadow')
    } else {
      this.eoAfterDom && this.renderer2.removeClass(this.eoAfterDom, 'scroll-bottom-box-shadow')
    }
  }

  @HostListener('window:resize', ['$event'])
    getElementHeight () {
      this.domMaxHeight = (this.eoParentDom?.clientHeight || this.el.nativeElement.parentElement.clientHeight) - (this.eoBeforeDom?.clientHeight || 0) - (this.eoAfterDom?.clientHeight || 0)
      this.renderer2.setStyle(this.el.nativeElement, 'max-height', this.domMaxHeight + 'px')
      setTimeout(() => {
        this.onScroll()
      }, 0)
    }

  ngAfterViewInit (): void {
    this.renderer2.setStyle(this.el.nativeElement, 'overflow-y', 'auto')
    setTimeout(() => {
      this.getElementHeight()
    }, 100)
    this.resizeObserver.observe(this.el.nativeElement)
  }
}
