import { Injectable } from '@angular/core'
import { EO_TBODY_TYPE } from 'projects/eo-ng-apinto-table/src/public-api'
import { DeployClusterEnvironmentConfigUpdateComponent } from './cluster/environment/config/update/update.component'
import { DeployClusterEnvironmentComponent } from './cluster/environment/environment.component'
import { DeployClusterListComponent } from './cluster/list/list.component'
import { DeployClusterPluginComponent } from './cluster/plugin/plugin.component'
import { DeployPluginListComponent } from './plugin/list/list.component'

@Injectable({
  providedIn: 'root'
})
export class DeployService {
  createPluginsTbody = (context:DeployPluginListComponent):EO_TBODY_TYPE[] => {
    return [
      {
        type: 'sort'
      },
      { title: context.pluginName },
      {
        key: 'extended',
        copy: true
      },
      {
        key: 'desc'
      },
      { key: 'updateTime' },
      {
        type: 'btn',
        right: true,
        btns: [
          {
            title: '编辑',
            click: (item: any) => {
              context.router.navigate(['/', 'deploy', 'plugin', 'message', item.data.name])
            }
          },
          {
            title: '删除',
            disabledFn: (item:any) => {
              return !item.isDelete || context.nzDisabled
            },
            click: (item:any) => {
              context.delete(item.data)
            }
          }
        ]
      }
    ]
  }

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
        key: 'name',
        copy: true
      },
      {
        key: 'env',
        copy: true
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
        key: 'name',
        copy: true
      },
      { key: 'env' },
      { key: 'status' },
      {
        type: 'btn',
        right: true,
        btns: [
          {
            title: '查看',
            click: (item: any) => {
              context.router.navigate(['/', 'deploy', 'cluster', 'content', item.data.name])
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
        key: 'publish',
        copy: true
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
