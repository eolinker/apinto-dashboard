import * as React from 'react'
import AceEditor from 'react-ace'
import 'ace-builds/src-noconflict/mode-jsx' // jsx模式的包
import 'ace-builds/src-noconflict/mode-yaml' // yaml模式的包
import 'ace-builds/src-noconflict/mode-json' // yaml模式的包
import 'ace-builds/src-noconflict/theme-monokai' // monokai的主题样式
import 'ace-builds/src-noconflict/theme-xcode' // monokai的主题样式
import 'ace-builds/src-noconflict/ext-language_tools' // 代码联想

const mockData = `stages:
- exec

init_job:
resource_group: $CI_PROJECT_NAME
stage: exec
trigger:
  include: deploy.gitlab-ci.yml
  strategy: depend
`

export const CustomCodeboxComponent = React.forwardRef(
  (props: { [k: string]: any }, ref) => {
    const {
      mode = 'yaml',
      theme = 'xcode',
      fontSize,
      height,
      width = '100%',
      onChange,
      value = mockData
    } = props
    React.useImperativeHandle(ref, () => ({}))

    const handleChange = (value: string) => {
      onChange(value)
    }

    return (
      <div>
        <AceEditor
          mode={mode}
          theme={theme}
          fontSize={fontSize}
          height={height}
          width={width}
          showGutter
          onChange={(value) => {
            console.log(value) // 输出代码编辑器内值改变后的值
            handleChange(value)
          }}
          value={value}
          wrapEnabled
          enableSnippets // 启用代码段
          setOptions={{
            enableBasicAutocompletion: true, // 启用基本自动完成功能
            enableLiveAutocompletion: true, // 启用实时自动完成功能 （比如：智能代码提示）
            enableSnippets: true, // 启用代码段
            showLineNumbers: true,
            tabSize: 2
          }}
        />
      </div>
    )
  }
)
