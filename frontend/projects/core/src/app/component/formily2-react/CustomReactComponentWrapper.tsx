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
import { IntelligentPluginEditComponent } from './IntelligentPluginEditComponent'
import { SelectOption } from 'eo-ng-select'

const containerElementName = 'customReactComponentContainer'

@Component({
  selector: 'formily2-react-wrapper',
  template: `<span #${containerElementName}></span>`,
  // styleUrls: [''],
  encapsulation: ViewEncapsulation.None
})
export class CustomReactComponentWrapperComponent {
  @ViewChild(containerElementName, { static: true }) containerRef!: ElementRef

  @Output() public componentClick = new EventEmitter<void>()

  // 动态渲染区域的render语句，目前后端接口传来的是对象，可以直接用，无需前端处理

  mockRenderSchema: { [k: string]: any } = {
    type_1: {
      type: 'void',
      properties: {
        aa: {
          type: 'string',
          title: 'AA',
          'x-decorator': 'FormItem',
          'x-decorator-props': {
            labelCol: 6,
            wrapperCol: 10
          },
          'x-component': 'Input',
          'x-component-props': {
            placeholder: 'Input'
          }
        },
        fomatter: {
          type: 'string',
          title: '配置',
          'x-decorator': 'FormItem',
          'x-decorator-props': {
            labelCol: 6,
            wrapperCol: 10
          },
          'x-component': 'CustomCodeboxComponent',
          'x-component-props': {
            mode: 'yaml'
          }
        },
        env_addr: {
          type: 'string',
          title: '引用环境变量',
          'x-decorator': 'FormItem',
          'x-decorator-props': {
            labelCol: 6,
            wrapperCol: 10
          },
          'x-component': 'CustomEnvVariableComponent',
          'x-component-props': {
            title: '引用环境变量'
          }
        }
      }
    },
    type_2: {
      type: 'void',
      properties: {
        aa: {
          type: 'string',
          title: 'AA',
          'x-decorator': 'FormItem',
          'x-decorator-props': {
            labelCol: 6,
            wrapperCol: 10
          },
          enum: [
            {
              label: '111',
              value: '111'
            },
            { label: '222', value: '222' }
          ],
          'x-component': 'Select',
          'x-component-props': {
            placeholder: 'Select'
          }
        },
        bb: {
          type: 'string',
          title: 'BB',
          'x-decorator': 'FormItem',
          'x-decorator-props': {
            labelCol: 6,
            wrapperCol: 10
          },
          'x-component': 'Input'
        }
      }
    }
  }

  @Input() renderSchema: { [k: string]: any } = { ...this.mockRenderSchema }

  // 编辑表单的标志
  @Input() editPage: boolean = false

  // 初始表单值
  @Input() initFormValue: { [k: string]: any } = {
    name: 'Aston',
    desc: 'Martin',
    driver: 'type_1',
    aa: 'waw'
  }

  // render用的选项
  @Input() driverSelectOptions: SelectOption[] = [
    { label: '类型1', value: 'type_1' },
    { label: '类型2', value: 'type_2' }
  ]

  @Input() demoSchema: any

  @Input() demo: boolean = false

  mockDriverSelectOptions: SelectOption[] = [
    { label: '类型1', value: 'type_1' },
    { label: '类型2', value: 'type_2' }
  ]

  reactComponent: React.RefObject<any> = React.createRef()
  constructor() {
    this.handleDivClicked = this.handleDivClicked.bind(this)
  }

  public handleDivClicked() {
    if (this.componentClick) {
      this.componentClick.emit()
      this.render()
    }
  }

  ngOnChanges(changes: SimpleChanges): void {
    // eslint-disable-next-line dot-notation
    if (changes['demoSchema']) {
      console.log(this.demoSchema)
    }
    console.log(this.demoSchema)
    this.render()
  }

  ngAfterViewInit() {
    this.render()
  }

  ngOnDestroy() {
    ReactDOM.unmountComponentAtNode(this.containerRef.nativeElement)
  }

  private render() {
    ReactDOM.render(
      <React.StrictMode>
        <div>
          <IntelligentPluginEditComponent
            ref={this.reactComponent}
            schema={this.mockRenderSchema}
            initFormValue={this.initFormValue}
            driverSelectOptions={this.mockDriverSelectOptions}
            editPage={this.editPage}
            demoSchema={this.demoSchema}
            demo={this.demo}
          />
        </div>
      </React.StrictMode>,
      this.containerRef.nativeElement
    )
  }
}
