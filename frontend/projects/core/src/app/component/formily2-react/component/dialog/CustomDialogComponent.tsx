import * as React from 'react'
import { Stringify, createSchemaField } from '@formily/react'
import { Button } from 'antd'
import {
  FormItem,
  Space,
  ArrayItems,
  DatePicker,
  Editable,
  FormButtonGroup,
  Input,
  Radio,
  Select,
  Submit,
  Cascader,
  Form,
  FormGrid,
  FormLayout,
  Upload,
  ArrayCollapse,
  ArrayTable,
  ArrayTabs,
  Checkbox,
  FormCollapse,
  FormDialog,
  FormDrawer,
  FormStep,
  FormTab,
  NumberPicker,
  Password,
  PreviewText,
  Reset,
  SelectTable,
  Switch,
  TimePicker,
  Transfer,
  TreeSelect,
  ArrayCards
} from '@formily/antd'
import { CustomCodeboxComponent } from '../codebox/CustomCodeboxComponent'
import { CustomEnvVariableComponent } from '../editable-env-table/CustomEnvVariableComponent'
import { SimpleMapComponent } from '../simple-map/SimpleMapComponent'

const SchemaField = createSchemaField({
  components: {
    ArrayCards,
    ArrayCollapse,
    ArrayItems,
    ArrayTable,
    ArrayTabs,
    Cascader,
    Checkbox,
    DatePicker,
    Editable,
    Form,
    FormButtonGroup,
    FormCollapse,
    FormDialog,
    FormDrawer,
    FormGrid,
    FormItem,
    FormLayout,
    FormStep,
    FormTab,
    Input,
    NumberPicker,
    Password,
    PreviewText,
    Radio,
    Reset,
    Select,
    SelectTable,
    Space,
    Submit,
    Switch,
    TimePicker,
    Transfer,
    TreeSelect,
    Upload,
    CustomCodeboxComponent,
    CustomEnvVariableComponent,
    SimpleMapComponent
  }
})

export const CustomDialogComponent = React.forwardRef(
  (props: { [k: string]: any }, ref) => {
    const { onChange, title, value, render } = props
    React.useImperativeHandle(ref, () => ({}))
    console.log(value, render)
    let editPage: boolean = false
    try {
      editPage = Object.keys(JSON.parse(JSON.stringify(value))).length > 0
    } catch {}

    return (
      <FormDialog.Portal>
        <Button
          type="text"
          onClick={() => {
            const dialog = FormDialog(
              editPage ? `编辑${title || ''}` : `新建${title || ''}`,
              () => {
                return (
                  <FormLayout labelCol={6} wrapperCol={10} form={value}>
                    <SchemaField schema={JSON.parse(render)} />
                  </FormLayout>
                )
              }
            )
            dialog
              .forOpen((payload, next) => {
                next({
                  initialValues: value
                })
              })
              .forConfirm((payload, next) => {
                next(payload)
              })
              .forCancel((payload, next) => {
                next(payload)
              })
              .open()
              .then(onChange)
          }}
        >
          <svg style={{ width: '16px', height: '16px' }}>
            <use href="#tool"></use>
          </svg>
        </Button>
      </FormDialog.Portal>
    )
  }
)
