/* eslint-disable multiline-ternary */
import * as React from 'react'
import {
  Button,
  Form,
  Input,
  InputNumber,
  Table,
  Typography,
  message
} from 'antd'
import axios from 'axios'
import { environment } from 'projects/core/src/environments/environment'

interface DataType {
  key: string
  description: string
}

interface EditableCellProps extends React.HTMLAttributes<HTMLElement> {
  editing: boolean
  dataIndex: string
  title: any
  inputType: 'number' | 'text'
  record: DataType
  index: number
  children: React.ReactNode
}

const EditableCell: React.FC<EditableCellProps> = (props: any) => {
  const {
    editing,
    dataIndex,
    title,
    inputType,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    record,
    // eslint-disable-next-line @typescript-eslint/no-unused-vars
    index,
    children,
    ...restProps
  } = props
  const inputNode = inputType === 'number' ? <InputNumber /> : <Input />
  return (
    <td {...restProps}>
      {editing ? (
        <Form.Item
          name={dataIndex}
          style={{ margin: 0 }}
          rules={[
            {
              required: dataIndex === 'key',
              message:
                dataIndex === 'key'
                  ? '英文数字下划线任意一种，首字母必须为英文'
                  : '',
              pattern:
                dataIndex === 'key' ? /^[a-zA-Z][a-zA-Z0-9/_]*$/ : undefined
            }
          ]}
        >
          {inputNode}
        </Form.Item>
      ) : (
        children
      )}
    </td>
  )
}

type Props = {
  chooseEnv: any
}

export const CustomEnvVariableTableComponent: React.FunctionComponent<Props> = (
  props: Props
) => {
  const { chooseEnv } = props

  React.useEffect(() => {
    getEnvlist()
  }, [])

  // const handleChange = (value: string) => {}

  const [dataSource, setDataSource] = React.useState<DataType[]>([
    {
      key: '0',
      description: 'Edward King 0'
    },
    {
      key: '1',
      description: 'Edward King 1'
    }
  ])

  const [form] = Form.useForm()
  const [editing, setEditing] = React.useState(false)
  const [keyword, setKeyword] = React.useState('')
  const [pageSize, setPageSize] = React.useState(15)
  const [pageNum, setPageNum] = React.useState(1)
  const [tableLoading, setTableLoading] = React.useState(true)
  const [total, setTotal] = React.useState(0)
  const editingKey = ''
  const columns = [
    {
      title: 'KEY',
      dataIndex: 'key',
      key: 'key',
      editable: true
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      editable: true
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: DataType) => {
        const editable = isEditing(record)
        return editable ? (
          <span>
            <Typography.Link onClick={() => save()} style={{ marginRight: 8 }}>
              Save
            </Typography.Link>
            <Typography.Link onClick={() => cancel()}>Cancle</Typography.Link>
          </span>
        ) : (
          <Typography.Link onClick={() => chooseEnv(record)}>
            Choose
          </Typography.Link>
        )
      }
    }
  ]

  const handleAdd = () => {
    if (editing) {
      return
    }
    const newData: DataType = {
      key: '',
      description: ''
    }
    setDataSource([newData, ...dataSource])
    setEditing(true)
  }

  const isEditing = (record: DataType) => record.key === editingKey

  const cancel = () => {
    setEditing(false)
    const newData = [...dataSource]
    newData.shift()
    setDataSource([...newData])
  }

  const getEnvlist = () => {
    setEditing(false)
    axios
      .get(`${environment.urlPrefix}api/variables`, {
        params: {
          pageNum: pageNum || 0,
          pageSize: pageSize || 15,
          key: keyword || ''
        }
      })
      .then((resp) => {
        setTableLoading(false)
        if (resp.data.code === 0) {
          setDataSource(resp.data.data.variables)
          setTotal(resp.data.total)
        }
      })
  }

  const tableDataChange = (...pagination: any) => {
    setPageNum(pagination[0])
    setPageSize(pagination[1])
    // getEnvlist()
  }

  const keywordChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setKeyword(e.target.value)
  }

  const save = async () => {
    try {
      const row = (await form.validateFields()) as DataType
      axios
        .post(`${environment.urlPrefix}api/variable`, {
          key: row.key || '',
          desc: row.description || ''
        })
        .then(({ data }) => {
          if (data.code === 0) {
            message.success(data.msg)
            getEnvlist()
          }
        })
    } catch (errInfo) {
      console.log('表单校验失败:', errInfo)
    }
  }

  const mergedColumns = columns.map((col) => {
    if (!col.editable) {
      return col
    }
    return {
      ...col,
      onCell: (record: DataType) => ({
        record,
        inputType: col.dataIndex === 'age' ? 'number' : 'text',
        dataIndex: col.dataIndex,
        title: col.title,
        editing: isEditing(record)
      })
    }
  })

  return (
    <div>
      <div className="pl-btnbase pr-btnrbase pb-btnybase flex flex-nowrap items-center justify-between">
        <Button onClick={handleAdd} type="primary" style={{ marginBottom: 16 }}>
          添加变量
        </Button>
        <Input
          prefix={
            <svg className="iconpark-icon">
              <use href="#search"></use>
            </svg>
          }
          allowClear
          className="w-SEARCH rounded-SEARCH_RADIUS"
          value={keyword}
          onChange={keywordChange}
          // onPressEnter={getEnvlist}
        />
      </div>
      <Form form={form} component={false}>
        <Table
          components={{
            body: {
              cell: EditableCell
            }
          }}
          loading={tableLoading}
          bordered
          dataSource={dataSource}
          columns={mergedColumns}
          rowClassName="editable-row"
          pagination={{
            onChange: tableDataChange,
            current: pageNum,
            pageSize: pageSize,
            pageSizeOptions: [15, 20, 50, 100],
            showTotal: (total) => `共${total}条`,
            showSizeChanger: true,
            showQuickJumper: true
          }}
        />
      </Form>
    </div>
  )
}
