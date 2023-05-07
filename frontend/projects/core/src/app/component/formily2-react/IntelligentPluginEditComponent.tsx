// eslint-disable-next-line no-use-before-define
import * as React from 'react'
import { action } from '@formily/reactive'
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
import { createForm } from '@formily/core'
import {
  FormProvider,
  RecursionField,
  createSchemaField,
  observer,
  useField,
  useForm
} from '@formily/react'
import { CustomCodeboxComponent } from './component/codebox/CustomCodeboxComponent'
import { CustomEnvVariableComponent } from './component/editable-env-table/CustomEnvVariableComponent'
import axios from 'axios'
import { SimpleMapComponent } from './component/simple-map/SimpleMapComponent'
import { CustomDialogComponent } from './component/dialog/CustomDialogComponent'

const DynamicRender = observer(() => {
  const field = useField()
  const form = useForm()
  const [schema, setSchema] = React.useState({})

  React.useEffect(() => {
    form.clearFormGraph(`${field.address}.*`)
    setSchema(DYNAMIC_INJECT_SCHEMA[form.values.driver])
  }, [form.values.driver])

  return (
    <RecursionField
      basePath={field.address}
      schema={schema}
      onlyRenderProperties
    />
  )
})

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
    DynamicRender,
    CustomCodeboxComponent,
    CustomEnvVariableComponent,
    SimpleMapComponent,
    CustomDialogComponent
  }
})

let DYNAMIC_INJECT_SCHEMA: any

export const IntelligentPluginEditComponent = React.forwardRef(
  (props: { [k: string]: any }, ref) => {
    const {
      schema,
      initFormValue,
      driverSelectOptions,
      demo,
      demoSchema,
      editPage = false,
      onSubmit
    } = props
    React.useImperativeHandle(ref, () => ({ form, submitRef }))

    const submitRef = React.createRef()
    DYNAMIC_INJECT_SCHEMA = schema

    const form = createForm({ validateFirst: editPage })
    form.setInitialValues(initFormValue)
    const pluginEditSchema = {
      type: 'object',
      properties: {
        id: {
          type: 'string',
          title: 'ID',
          required: true,
          'x-decorator': 'FormItem',
          'x-decorator-props': {
            labelCol: 6,
            wrapperCol: 10
          },
          'x-component': 'Input',
          'x-component-props': {
            placeholder: 'ID'
          }
        },
        title: {
          type: 'string',
          title: '名称',
          required: true,
          'x-decorator': 'FormItem',
          'x-decorator-props': {
            labelCol: 6,
            wrapperCol: 10
          },
          'x-component': 'Input',
          'x-component-props': {
            placeholder: '名称'
          }
        },
        driver: {
          type: 'string',
          title: 'Driver',
          required: true,
          'x-decorator': 'FormItem',
          'x-decorator-props': {
            labelCol: 6,
            wrapperCol: 10
          },
          'x-component': 'Select',
          'x-component-props': {
            disabled: editPage
          },
          'x-display': driverSelectOptions.length > 1 ? 'visible' : 'hidden',
          enum: [...driverSelectOptions]
        },
        description: {
          type: 'string',
          required: true,
          title: '描述',
          'x-decorator': 'FormItem',
          'x-decorator-props': {
            labelCol: 6,
            wrapperCol: 10
          },
          'x-component': 'Input.TextArea',
          'x-component-props': {
            placeholder: 'ID'
          }
        },
        container: {
          type: 'void',
          'x-component': 'DynamicRender',
          'x-component-props': {
            schema: JSON.stringify(schema)
          }
        }
      }
    }

    const submit = (value: any) => {
      onSubmit && onSubmit(value)
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const getSkillData = async (skill: string) => {
      return new Promise((resolve) => {
        axios.get(`api/common/provider/${skill}`).then((resp) => {
          if (resp.data.code === 0) {
            const dataList: Array<{ label: string; value: string }> =
              resp.data.data[skill].map((item: any) => {
                return {
                  label: item.title,
                  value: item.name
                }
              })
            resolve(dataList)
          }
        })
      })
    }

    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    const useAsyncDataSource =
      (service: any, skill: string) => (field: any) => {
        field.loading = true
        service(skill).then(
          action.bound &&
            action.bound((data: any) => {
              field.dataSource = data
              field.loading = false
            })
        )
      }
    return (
      <FormProvider form={form} layout="vertical">
        <SchemaField
          schema={demo ? demoSchema : pluginEditSchema}
          scope={{ useAsyncDataSource, getSkillData }}
        />
        {demo && demoSchema && (
          <FormButtonGroup>
            <Submit ref={submitRef} onSubmit={submit}>
              提交
            </Submit>
          </FormButtonGroup>
        )}
      </FormProvider>
    )
  }
)
