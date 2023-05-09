import { Directive, ElementRef } from '@angular/core'

@Directive({
  selector: '[eoNgAutoFocus]'
})
export class AutoFocusDirective {
  constructor (
    private el:ElementRef) { }

  ngAfterViewInit ():void {
    setTimeout(
      () => {
        this.el.nativeElement.focus()
      }, 0
    )
  }
}
