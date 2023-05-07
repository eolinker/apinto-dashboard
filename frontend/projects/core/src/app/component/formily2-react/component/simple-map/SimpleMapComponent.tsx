import * as React from 'react'

import { Input } from '@formily/antd'
export const SimpleMapComponent = React.forwardRef(
  (props: { [k: string]: any }, ref) => {
    const {
      onChange,
      value,
      placeholderKey = '请输入Key',
      placeholderValue = '请输入Value'
    } = props

    const [kvList, setKvList] = React.useState(
      value && Object.keys(value).length > 0
        ? [
            ...Object.keys(value).map((k: string) => {
              return { key: k, value: value[k] }
            })
          ]
        : [{ key: '', value: '' }]
    )

    React.useImperativeHandle(ref, () => ({}))

    const emitNewArr = () => {
      const res: { [k: string]: any } = {}
      for (const kv of kvList) {
        res[kv.key] = kv.value
      }
      onChange(res)
    }

    const changeInputValue = (
      newValue: string,
      index: number,
      type: 'key' | 'value'
    ) => {
      const newArr = [...kvList]
      newArr[index][type] = newValue
      setKvList(newArr)
      emitNewArr()
    }

    const addLine = (index: number) => {
      const newKvList = [...kvList]
      newKvList.splice(index + 1, 0, { key: '', value: '' })
      setKvList(newKvList)
      emitNewArr()
    }

    const removeLine = (index: number) => {
      const newKvList = [...kvList]
      newKvList.splice(index, 1)
      setKvList(newKvList)
      emitNewArr()
    }

    return (
      <div>
        {kvList.map((n: any, index: any) => {
          return (
            <div
              key={n + index}
              className="flex"
              style={{ marginTop: index === 0 ? '0px' : '16px' }}
            >
              <Input
                className="mr-[8px]"
                style={{ width: '174px' }}
                value={n.key}
                onChange={(e: any) => {
                  changeInputValue(e.target.value, index, 'key')
                }}
                placeholder={placeholderKey}
              />
              <Input
                className=" mr-[8px]"
                style={{ width: '164px' }}
                value={n.value}
                onChange={(e: any) => {
                  changeInputValue(e.target.value, index, 'value')
                }}
                placeholder={placeholderValue}
              />
              <a
                className="array_item_addition ml-[10px] ant-btn-text anticon"
                onClick={() => addLine(index)}
              >
                <span>
                  <svg className="iconpark-icon">
                    <use href="#add-circle"></use>
                  </svg>
                </span>
              </a>
              {index !== 0 && (
                <a
                  className="array_item_addition ant-btn-text anticon"
                  onClick={() => removeLine(index)}
                >
                  <span>
                    <svg className="iconpark-icon">
                      <use href="#reduce-one"></use>
                    </svg>
                  </span>
                </a>
              )}
            </div>
          )
        })}
      </div>
    )
  }
)
