import * as React from 'react'
import { FormDialog, FormItem, FormLayout, Input } from '@formily/antd'
import { createSchemaField } from '@formily/react'
import { Button } from 'antd'

const SchemaField = createSchemaField({
  components: {
    FormItem,
    Input
  }
})

export const CustomDialogComponent = React.forwardRef(
  (props: { [k: string]: any }, ref) => {
    const { onChange, title, value, render } = props
    React.useImperativeHandle(ref, () => ({}))
    console.log(value)
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
                    <SchemaField schema={render} />
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
