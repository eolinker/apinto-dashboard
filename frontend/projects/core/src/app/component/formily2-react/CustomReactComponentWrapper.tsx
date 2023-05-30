import {
  Component,
  ElementRef,
  EventEmitter,
  Input,
  Output,
  SimpleChanges,
  ViewChild
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
  styles: [
    `
      :host ::ng-deep {
        .ant-input-affix-wrapper,
        textarea,
        .ant-input-number,
        ant-picker,
        .ant-select {
          width: 346px;
          min-height: 32px;
        }

        .ant-formily-array-items {
          width: 400px;
        }

        .ant-input-affix-wrapper-disabled,
        .ant-input-affix-wrapper-disabled > .ant-input {
          background: var(--disabled-background-color);
          cursor: not-allowed;
        }

        .ant-input-affix-wrapper-disabled {
          border-color: var(--disabled-border-color);
        }

        .ant-picker-disabled:hover {
          border-color: var(--disabled-border-color);
        }

        .ant-switch {
          width: 28px;
          height: 16px;
          min-width: 28px;
          .ant-switch-handle {
            height: 12px;
            width: 12px;
          }
          &.ant-switch-checked,
          &.ant-switch-checked:hover:not(.ant-switch-disabled) {
            background: var(--primary-color);
          }
          &.ant-switch-checked .ant-switch-handle {
            inset-inline-start: calc(100% - 14px);
          }
        }

        .ant-formily-array-table {
          width: 524px;
          .ant-input-affix-wrapper,
          textarea,
          .ant-input-number,
          .ant-formily-array-items,
          ant-picker,
          .ant-select {
            width: 100%;
            min-height: 32px;
          }
        }
        .ant-formily-array-items .ant-select {
          width: unset;
        }

        .ant-btn {
          box-shadow: none !important;
          transition: none !important;
          animation: none !important;
          outline: none !important;
        }
        .ant-btn::after {
          border: 0 none;
          opacity: 0;
          animation: none 0 ease 0 1 normal;
          transition: none !important;
        }
        .ant-btn-default:not(:disabled):hover,
        .ant-btn-dashed:not(:disabled):hover {
          color: var(--button-default-hover-text-color);
          border-color: var(--button-default-hover-border-color);
          background-color: var(--button-default-hover-background-color);
        }

        .ant-btn-primary:not(:disabled):hover {
          color: var(--button-primary-text-color);
          border-color: var(--button-primary-hover-border-color);
          background-color: var(--button-primary-hover-background-color);
        }
      }
    `
  ]
})
export class CustomReactComponentWrapperComponent {
  @ViewChild(containerElementName, { static: true }) containerRef!: ElementRef

  @Output() public componentClick = new EventEmitter<void>()
  @Output() onSubmit = new EventEmitter<any>()

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
        },
        scheme: {
          title: '请求协议',
          'x-decorator': 'FormItem',
          'x-component': 'Select',
          'x-validator': [],
          'x-component-props': {},
          'x-decorator-props': {},
          required: true,
          default: 'HTTP',
          'x-reactions': ['{{useAsyncDataSource(getSkillData,"service")}}'],
          name: 'scheme',
          'x-index': 0
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
    this.render()
  }

  ngAfterViewInit() {
    this.render()
  }

  ngOnDestroy() {
    ReactDOM.unmountComponentAtNode(this.containerRef.nativeElement)
  }

  handlerSubmit = (value: any) => {
    this?.onSubmit && this?.onSubmit?.emit(value)
  }

  private render() {
    ReactDOM.render(
      <React.StrictMode>
        <div>
          <IntelligentPluginEditComponent
            ref={this.reactComponent}
            schema={this.renderSchema}
            initFormValue={this.initFormValue}
            driverSelectOptions={this.driverSelectOptions}
            editPage={this.editPage}
            demoSchema={this.demoSchema}
            demo={this.demo}
            onSubmit={this.handlerSubmit}
          />
        </div>
      </React.StrictMode>,
      this.containerRef.nativeElement
    )
  }
}
