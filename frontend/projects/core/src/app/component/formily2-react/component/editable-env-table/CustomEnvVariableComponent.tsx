import * as React from 'react'

import { Modal } from 'antd'
import { CustomEnvVariableTableComponent } from './CustomEnvVariableTableComponent'

export const CustomEnvVariableComponent = React.forwardRef(
  (props: { [k: string]: any }, ref) => {
    const { onChange, value, title = '引用环境变量' } = props
    React.useImperativeHandle(ref, () => ({}))

    //   const [isModalOpen, setIsModalOpen] = React.useState(false)

    let modalRef: any

    const openModal = (e: any) => {
      e.preventDefault()
      const handleChooseEnv = ({ key }: any) => {
        modalRef?.destroy()
        onChange && onChange('${' + key + '}')
      }
      modalRef = Modal.confirm({
        title: '添加环境变量',
        width: '900px',
        icon: '',
        closable: true,
        footer: null,
        content: (
          <div>
            <CustomEnvVariableTableComponent chooseEnv={handleChooseEnv} />
          </div>
        )
      })
    }

    return (
      <a href="#!" onClick={openModal}>
        {title}
      </a>
    )
  }
)
