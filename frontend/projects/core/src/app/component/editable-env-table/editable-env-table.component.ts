/* eslint-disable camelcase */
/* eslint-disable no-useless-constructor */
import { Component, OnInit, Output, EventEmitter } from '@angular/core'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { ApiService } from 'projects/core/src/app/service/api.service'

@Component({
  selector: 'editable-env-table',
  templateUrl: './editable-env-table.component.html',
  styles: [
  ]
})
export class EditableEnvTableComponent implements OnInit {
  @Output() eoChooseEnv = new EventEmitter()
  envNameForSear:string = ''
  environmentTableHeadName: Array<object> = [
    { title: 'KEY' },
    { title: '描述' },
    {
      title: '操作',
      right: true
    }
  ]

  environmentTableBody: Array<any> =[
    {
      key: 'key',
      showFn: (item:any) => {
        return !item.editing
      }
    },
    {
      key: 'key',
      type: 'input',
      placeholder: '请输入KEY',
      showFn: (item:any) => {
        return item.editing
      },
      checkMode: 'change',
      check: (item: any) => {
        return item && /^[a-zA-Z][a-zA-Z0-9/_]*$/.test(item)
      },
      errorTip: '英文数字下划线任意一种，首字母必须为英文'
    },
    {
      key: 'description',
      showFn: (item:any) => {
        return !item.editing
      }
    },
    {
      key: 'description',
      type: 'input',
      placeholder: '请输入描述',
      showFn: (item:any) => {
        return item.editing
      }
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return item.editing
      },
      btns: [
        {
          title: '提交',
          click: (item:any) => {
            if (item.data.key && /^[a-zA-Z][a-zA-Z0-9/_]*$/.test(item.data.key)) {
              this.addEnv(item.data)
            }
          },
          showFn: (item:any) => {
            return !!item.key
          }
        },
        {
          title: '取消',
          click: (item:any) => {
            this.environmentList = this.environmentList.filter((el) =>
              el.id !== item.data.id
            )
          },
          showFn: (item:any) => {
            return !!item.key
          }
        }
      ]
    },
    {
      type: 'btn',
      right: true,
      showFn: (item:any) => {
        return !item.editing
      },
      btns: [
        {
          title: '添加',
          click: (item:any) => {
            this.chooseEnv({ data: item.data })
          }
        }
      ]
    }
  ]

  environmentList:Array<any> = []

  // 环境变量分页参数
  variablePage:{pageNum:number, pageSize:number, total:number}={
    pageNum: 1,
    pageSize: 15,
    total: 0
  }

  pageSizeOptions:Array<number>=[15, 20, 50, 100]

  chooseEnv = (item:any) => {
    this.eoChooseEnv.emit(item)
  }

  constructor (private message: EoNgFeedbackMessageService, private api:ApiService) { }

  ngOnInit (): void {
    this.getEnvlist()
  }

  getEnvlist (key?:string) {
    this.api.get('variables', { pageNum: this.variablePage.pageNum, pageSize: this.variablePage.pageSize, key: key || '' }).subscribe(resp => {
      if (resp.code === 0) {
        resp.data.variables.forEach((element:any) => {
          element.disabled = true
        })
        this.environmentList = resp.data.variables
        this.variablePage.total = resp.data.total
      }
    })
  }

  i:number = 0

  addEnvRow () {
    // this.environmentList.unshift({ name: '', desc: '', editing: true, disabled: false })
    this.environmentList = [
      { key: '', description: '', editing: true, id: this.i },
      ...this.environmentList
    ]
    this.i++
  }

  addEnv (item:any) {
    this.api.post('variable', { key: item.key || '', desc: item.desc || '' }).subscribe(resp => {
      if (resp.code === 0) {
        this.message.success(resp.msg || '添加成功！', { nzDuration: 1000 })
        this.getEnvlist()
      }
    })
  }
}
