import { Injectable } from '@angular/core'
import { DeployClusterEnvironmentConfigUpdateComponent } from './environment/config/update/update.component'
import { DeployClusterEnvironmentComponent } from './environment/environment.component'
import { DeployClusterListComponent } from './list/list.component'
import { DeployClusterPluginComponent } from './plugin/plugin.component'

@Injectable({
  providedIn: 'root'
})
export class DeployService {
  clusterName:string = ''
  clusterDesc:string = ''

  createClusterEnvUpdateThead (context:DeployClusterEnvironmentConfigUpdateComponent) {
    return [

      {
        type: 'checkbox',
        click: () => {
          context.getClusterCheckedList()
        }
      },
      { title: '集群名称', resizeable: true },
      { title: '所在环境' }
    ]
  }

  createClusterEnvUpdateTbody (context:DeployClusterEnvironmentConfigUpdateComponent) {
    return [
      {
        key: 'checked',
        type: 'checkbox',
        click: () => {
          context.getClusterCheckedList()
        }
      },
      {
        key: 'title',
        copy: true
      },
      {
        key: 'env'
      }
    ]
  }

  createClusterEnvUpdate2Thead (context:DeployClusterEnvironmentConfigUpdateComponent) {
    return [
      {
        type: 'checkbox',
        click: () => {
          context.getVarCheckedList()
        }
      },
      {
        title: 'KEY',
        resizeable: true,
        copy: true
      },
      {
        title: 'VALUE',
        resizeable: true,
        copy: true
      },
      { title: '更新时间' }
    ]
  }

  createClusterEnvUpdate2Tbody (context:DeployClusterEnvironmentConfigUpdateComponent) {
    return [
      {
        key: 'checked',
        type: 'checkbox',
        click: () => {
          context.getVarCheckedList()
        }
      },
      { key: 'key' },
      { key: 'value' },
      { key: 'updateTime' }
    ]
  }

  createClusterEnvConfigTbody (context:DeployClusterEnvironmentComponent) {
    return [
      {
        key: 'key',
        copy: true
      },
      {
        key: 'value',
        copy: true
      },
      {
        key: 'desc'
      },
      {
        key: 'publish'
      },
      {
        key: 'operator'
      },
      {
        key: 'updateTime'
      },
      {
        type: 'btn',
        right: true,
        showFn: (item:any) => {
          return item.publish !== 'DEFECT'
        },
        btns: [
          {
            title: '编辑',
            disabledFn: () => {
              return context.nzDisabled
            },
            click: (item:any) => {
              context.openDrawer('editConfig', item.data)
            }
          },
          {
            title: '删除',
            click: (item:any) => {
              context.delete(item.data)
            },
            disabledFn: () => {
              return context.nzDisabled
            }
          }
        ]
      },
      {
        type: 'btn',
        right: true,
        showFn: (item:any) => {
          return item.publish === 'DEFECT'
        },
        btns: [
          {
            title: '编辑',
            click: (item:any) => {
              context.openDrawer('editConfig', item.data)
            },
            disabledFn: () => {
              return context.nzDisabled
            }
          }
        ]
      }
    ]
  }

  createClusterTbody (context:DeployClusterListComponent) {
    return [
      {
        key: 'title',
        copy: true
      },
      {
        key: 'desc'
      },
      { key: 'env' },
      {
        key: 'status',
        title: context.clusterStatusTpl
      },
      {
        type: 'btn',
        right: true,
        btns: [
          {
            title: '查看',
            click: (item: any) => {
              context.router.navigate(['/', 'deploy', 'cluster', 'content', item.data.name])
              this.clusterName = item.data.title
              this.clusterDesc = item.data.desc
            }
          },
          {
            title: '删除',
            disabledFn: () => { return context.nzDisabled },
            click: (item:any) => {
              context.delete(item.data)
            }
          }
        ]
      }
    ]
  }

  createPluginTbody (context:DeployClusterPluginComponent) {
    return [
      {
        key: 'name',
        copy: true
      },
      {
        key: 'publish'
      },
      {
        key: 'status'
      },
      {
        key: 'config',
        json: true,
        copy: true
      },
      {
        key: 'updateTime'
      },
      {
        type: 'btn',
        right: true,
        btns: [
          {
            title: '配置',
            disabledFn: (item:any) => {
              return (item.isBuiltin) || context.nzDisabled
            },
            click: (item:any) => {
              context.openDrawer('editConfig', item.data)
            }
          }
        ]
      }
    ]
  }
}
