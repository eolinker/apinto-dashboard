import { Directive, ElementRef, HostListener, Input, Renderer2 } from '@angular/core'

@Directive({
  selector: '[eo-ng-tree-drag]'
})
export class TreeDragDirective {
  constructor (
    private renderer: Renderer2,
    private el: ElementRef
  ) {
    this.affectDom = document.getElementsByClassName('spreed-content')
  }

  @Input() minWidth?: any
  @Input() maxWidth?: any
  affectDom: any
  @HostListener('mousedown', ['$event'])
  public onMounseDown ($event: any) {
    $event.stopPropagation()
    document.onmousemove = this.moveFun
    document.onmouseup = () => {
      document.onmousemove = null
      document.onmouseup = null
    }
  }

  moveFun = ($event: any) => {
    const tmpMoveX = $event.movementX
    const offset = 0
    const tmpWidth = this.affectDom[0].clientWidth + offset + tmpMoveX
    if (tmpWidth <= this.minWidth) return
    if (tmpWidth >= this.maxWidth) return
    const elemWidth = `${tmpWidth}px`
    const tmpAffectCount = this.affectDom.length
    for (let key = 0; key < tmpAffectCount; key++) {
      this.affectDom[key].style.width = elemWidth
    }
  }
}
