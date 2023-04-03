/* eslint-disable no-useless-constructor */
import { Directive, ViewContainerRef } from '@angular/core'

@Directive({
  selector: '[eoNgTabHost]'
})
export class TabHostDirective {
  constructor (public viewContainerRef: ViewContainerRef) { }
}
