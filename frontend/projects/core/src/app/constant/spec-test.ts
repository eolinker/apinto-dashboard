// 此处放置单元测试用到的参数

import { Observable } from 'rxjs'
import { ApiMessage } from '../layout/api/types/types'
import { ApiGroup } from './type'
import { ModalOptions } from 'ng-zorro-antd/modal'

export class MockRenderer {
  removeAttribute (element: any, cssClass: string) {
    return cssClass + 'is removed from' + element
  }
}

export class MockMessageService {
  success () {
    return 'success'
  }

  error () {
    return 'error'
  }
}

export class MockEnsureService {
  create (args:ModalOptions) {
    return {
      close: () => { 'modal is close' },
      afterClose: new Observable((observer) => {
        observer.next('close')
      }),
      ...args
    }
  }
}

export const MockEmptySuccessResponse = {
  code: 0,
  data: {},
  msg: 'success'
}

export const MockEmptyFailedResponse = {
  code: 0,
  data: {},
  msg: 'failed'
}

export const MockGetCommonProviderService = {
  code: 0,
  data: {
    Service: [
      {
        name: 'test1cjk_service@upstream',
        title: 'test1[http]'
      },
      {
        name: 'test2@upstream',
        title: 'test2[http]'
      }
    ]
  }
}

export const MockPluginTemplateEnum = {
  code: 0,
  data: {
    templates: [
      {
        uuid: '70623690-430f-23db-ec75-763fe7c380d9',
        name: 'test1',
        desc: '',
        operator: '',
        create_time: '',
        update_time: '',
        is_delete: false
      },
      {
        uuid: 'aa0be463-43c8-4f67-2633-ea6bdcea9709',
        name: 'test2',
        desc: '',
        operator: '',
        create_time: '',
        update_time: '',
        is_delete: true
      }
    ]
  },
  msg: 'success'
}

export const MockRouterGroups:{code:number, data:ApiGroup, msg:string} = {
  code: 0,
  data: {
    apis: [],
    root: {
      uuid: '',
      name: '',
      groups: [
        {
          uuid: '50458642-5a9f-4136-9ff1-e30d647297e8',
          name: 'test1',
          children: [
            {
              uuid: '35938ae4-1a62-4e22-ad8c-3691e111820e',
              name: 'test1-c1',
              children: [],
              isDelete: false
            },
            {
              uuid: 'b238751a-dbfb-4610-8f40-a599737ac4e5',
              name: 'test1-c2',
              children: [],
              isDelete: false
            }
          ],
          isDelete: false
        },
        {
          uuid: '00db4977-331f-4b7e-93be-b64648751a5f',
          name: 'test2',
          children: [],
          isDelete: false
        }
      ],
      isDelete: false
    }
  },
  msg: 'success'
}

export const MockAccessList:{code:number, data:{access:Array<{name:string, access:string}>}, msg:string} = {
  code: 0,
  data: {
    access: [
      {
        name: 'access-log',
        access: 'edit'
      },
      {
        name: 'visit',
        access: 'edit'
      },
      {
        name: 'ext-app',
        access: 'edit'
      },
      {
        name: 'variable',
        access: 'edit'
      },
      {
        name: 'grey',
        access: 'edit'
      },
      {
        name: 'redis',
        access: 'edit'
      },
      {
        name: 'plugin',
        access: 'edit'
      },
      {
        name: 'apispace',
        access: 'edit'
      },
      {
        name: 'upstream',
        access: 'edit'
      },
      {
        name: 'webhook',
        access: 'edit'
      },
      {
        name: 'plugin-template',
        access: 'edit'
      },
      {
        name: 'user',
        access: 'edit'
      },
      {
        name: 'module-plugin',
        access: 'edit'
      },
      {
        name: 'email',
        access: 'edit'
      },
      {
        name: 'discovery-2',
        access: 'edit'
      },
      {
        name: 'audit-log',
        access: 'edit'
      },
      {
        name: 'influxdb',
        access: 'edit'
      },
      {
        name: 'application',
        access: 'edit'
      },
      {
        name: 'discovery-3',
        access: 'edit'
      },
      {
        name: 'monitor',
        access: 'edit'
      },
      {
        name: 'traffic',
        access: 'edit'
      },
      {
        name: 'apispace-2',
        access: 'edit'
      },
      {
        name: 'cache',
        access: 'edit'
      },
      {
        name: 'cluster',
        access: 'edit'
      },
      {
        name: 'discovery',
        access: 'edit'
      },
      {
        name: 'fuse',
        access: 'edit'
      },
      {
        name: 'api',
        access: 'edit'
      }
    ]
  },
  msg: 'success'
}

export const MockModuleList:{code:number, data:{navigation:Array<any>}, msg:string} = {
  code: 0,
  data: {
    navigation: [
      {
        title: '仪表盘',
        icon: 'file-cabinet',
        modules: [
          {
            name: 'monitor',
            title: '监控告警',
            path: 'module/monitor'
          }
        ],
        default: 'monitor'
      },
      {
        title: 'API',
        icon: 'APIjiekou-7mme3dcg',
        modules: [
          {
            name: 'api',
            title: 'API管理',
            path: 'router/api'
          }
        ],
        default: 'api'
      }
    ]
  },
  msg: 'success'
}

export const MockApiWsMessage:{code:number, data:{api:ApiMessage}, msg:string} = {
  code: 0,
  data: {
    api: {
      name: 'ss',
      uuid: '569c8d47-d742-5306-c0e2-a5ae38727fa7',
      groupUuid: '50458642-5a9f-4136-9ff1-e30d647297e8',
      desc: '',
      isDisable: false,
      scheme: 'websocket',
      requestPath: '/tetetetete',
      service: 'testService',
      method: [],
      proxyPath: 'tetetetete',
      hosts: ['test1.host.addr', 'test2.host.addr'],
      timeout: 10000,
      retry: 0,
      match: [
        {
          position: 'header',
          matchType: 'PREFIX',
          key: 'ee',
          pattern: 'te'
        }
      ],
      proxyHeader: [
        {
          optType: 'DELETE',
          key: 'eee',
          value: ''
        }
      ],
      templateUuid: 'e59693df-27cc-61a7-d0f0-c17da203026a'
    }
  },
  msg: 'success'
}

export const MockApiWsMessage2:{code:number, data:{api:ApiMessage}, msg:string} = {
  code: 0,
  data: {
    api: {
      name: 'ss',
      uuid: '569c8d47-d742-5306-c0e2-a5ae38727fa7',
      groupUuid: '50458642-5a9f-4136-9ff1-e30d647297e8',
      desc: '',
      isDisable: false,
      scheme: 'websocket',
      requestPath: '{{baseUrl}}/test',
      service: 'testService',
      method: [],
      proxyPath: 'tetetetete',
      hosts: [],
      timeout: 10000,
      retry: 0,
      match: [
        {
          position: 'header',
          matchType: 'PREFIX',
          key: 'ee',
          pattern: 'te'
        }
      ],
      proxyHeader: [
        {
          optType: 'DELETE',
          key: 'eee',
          value: ''
        }
      ],
      templateUuid: 'e59693df-27cc-61a7-d0f0-c17da203026a'
    }
  },
  msg: 'success'
}
