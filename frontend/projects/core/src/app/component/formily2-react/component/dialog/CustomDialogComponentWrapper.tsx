import {
  Component,
  ElementRef,
  EventEmitter,
  Input,
  Output,
  SimpleChanges,
  ViewChild,
  ViewEncapsulation
} from '@angular/core'
// eslint-disable-next-line no-use-before-define
import * as React from 'react'
import * as ReactDOM from 'react-dom'
import { CustomDialogComponent } from './CustomDialogComponent'

const containerElementName = 'customReactComponentContainer'

@Component({
  selector: 'dialog-react-wrapper',
  template: `<span #${containerElementName}></span>`,
  // styleUrls: [''],
  encapsulation: ViewEncapsulation.None
})
export class CustomDialogReactComponentWrapperComponent {
  @ViewChild(containerElementName, { static: true }) containerRef!: ElementRef

  @Output() onChange: EventEmitter<any> = new EventEmitter<any>()
  // 动态渲染区域的render语句，目前后端接口传来的是对象，可以直接用，无需前端处理

  @Input() renderSchema: string = ''
  @Input() value: string = ''
  ngOnChanges(changes: SimpleChanges): void {
    this.render()
  }

  ngAfterViewInit() {
    this.render()
  }

  ngOnDestroy() {
    ReactDOM.unmountComponentAtNode(this.containerRef.nativeElement)
  }

  handleChange = (data: string) => {
    this.onChange.emit(data)
  }

  private render() {
    ReactDOM.render(
      <React.StrictMode>
        <div>
          <CustomDialogComponent
            onChange={this.handleChange}
            render={this.renderSchema}
            value={this.value}
          />
        </div>
      </React.StrictMode>,
      this.containerRef.nativeElement
    )
  }
}
