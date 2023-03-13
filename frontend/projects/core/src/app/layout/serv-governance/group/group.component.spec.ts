/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-14 22:56:33
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-10-27 11:35:45
 * @FilePath: /apinto/src/app/layout/upstream/service-discovery-content/service-discovery-content.component.spec.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'
import { NzOutletModule } from 'ng-zorro-antd/core/outlet'
import { NzDrawerModule } from 'ng-zorro-antd/drawer'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { ActivatedRoute, RouterModule } from '@angular/router'
import { ElementRef, Renderer2, ChangeDetectorRef } from '@angular/core'
import { APP_BASE_HREF } from '@angular/common'
import { NoopAnimationsModule } from '@angular/platform-browser/animations'
import { NzNoAnimationModule } from 'ng-zorro-antd/core/no-animation'
import { NzOverlayModule } from 'ng-zorro-antd/core/overlay'
import { environment } from 'projects/core/src/environments/environment'
import { BidiModule } from '@angular/cdk/bidi'
import { Overlay } from '@angular/cdk/overlay'
import { of } from 'rxjs'
import {
  EoNgFeedbackMessageService,
  EoNgFeedbackModalService
} from 'eo-ng-feedback'
import { FormsModule } from '@angular/forms'

import { EoNgSelectModule } from 'eo-ng-select'
import { CacheCreateComponent } from '../cache/create/create.component'
import { FuseCreateComponent } from '../fuse/create/create.component'
import { GreyCreateComponent } from '../grey/create/create.component'
import { TrafficCreateComponent } from '../traffic/create/create.component'
import { VisitCreateComponent } from '../visit/create/create.component'
import { GroupComponent } from './group.component'
import { NzTreeNode } from 'ng-zorro-antd/tree'

class MockRenderer {
  removeAttribute(element: any, cssClass: string) {
    return cssClass + 'is removed from' + element
  }
}

class MockMessageService {
  success() {
    return 'success'
  }

  error() {
    return 'error'
  }
}

class MockEnsureService {
  create() {
    return 'modal is create'
  }
}

jest.mock('uuid', () => {
  return {
    v4: () => 123456789
  }
})

describe('GroupComponent test as editPage is false', () => {
  let component: GroupComponent
  let fixture: ComponentFixture<GroupComponent>
  class MockElementRef extends ElementRef {
    constructor() {
      super(null)
    }
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
        FormsModule,
        EoNgSelectModule,
        BidiModule,
        NoopAnimationsModule,
        NzNoAnimationModule,
        NzDrawerModule,
        NzOutletModule,
        HttpClientModule,
        RouterModule.forRoot([
          {
            path: 'serv-governance/traffic/create',
            component: TrafficCreateComponent
          },
          {
            path: 'serv-governance/fuse/create',
            component: FuseCreateComponent
          },
          {
            path: 'serv-governance/grey/create',
            component: GreyCreateComponent
          },
          {
            path: 'serv-governance/visit/create',
            component: VisitCreateComponent
          },
          {
            path: 'serv-governance/cache/create',
            component: CacheCreateComponent
          }
        ])
      ],
      declarations: [GroupComponent],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: Renderer2, useClass: MockRenderer },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: EoNgFeedbackModalService, useClass: MockEnsureService },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef },
        {
          provide: ActivatedRoute,
          useValue: {
            queryParams: of({ cluster_name: 'clus2' })
          }
        }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(GroupComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })
  it('##should create', () => {
    expect(component).toBeTruthy()
  })

  it('##ngOnInit should call getGroupList()', () => {
    const spyGetGroupList = jest.spyOn(component, 'getGroupList')
    const spyGetGroupItemSelected = jest.spyOn(
      component,
      'getGroupItemSelected'
    )
    expect(spyGetGroupList).not.toHaveBeenCalled()
    expect(spyGetGroupItemSelected).not.toHaveBeenCalled()
    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetGroupList).toHaveBeenCalledTimes(1)
    expect(spyGetGroupItemSelected).toHaveBeenCalledTimes(1)
  })

  it('##ngOnDestroy test', () => {
    // @ts-ignore
    const spyUnsubscribe = jest.spyOn(component.subscription, 'unsubscribe')
    expect(spyUnsubscribe).not.toHaveBeenCalled()
    component.ngOnDestroy()
    fixture.detectChanges()
    expect(spyUnsubscribe).toHaveBeenCalledTimes(1)
  })

  it('##getGroupList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: { api: { group_uuid: 123456 }, root: { groups: [] } }
      })
    )
    const spyNodesTransfer = jest.spyOn(component, 'nodesTransfer')
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyNodesTransfer).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getGroupList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyNodesTransfer).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getGroupList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getGroupList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getGroupList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getGroupList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('获取数据失败!')
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##nodesTransfer & clustersTransfer test', () => {
    const spyClusterTransfer = jest.spyOn(component, 'clustersTransfer')
    expect(spyClusterTransfer).not.toBeCalled()
    const val1: any = [
      {
        clusters: [{ name: 'clu1' }, { name: 'clu2' }, { name: 'clu3' }],
        name: 'env1'
      },
      {
        clusters: [{ name: 'clu4' }, { name: 'clu5' }, { name: 'clu6' }],
        name: 'env2'
      },
      { clusters: [], name: 'env3' }
    ]

    let res = component.nodesTransfer([])
    fixture.detectChanges()
    expect(res).toStrictEqual([])

    component.activatedNode = null
    component.clusterName = ''
    component.clusterKey = ''
    component.strategyType = 'traffic'
    res = component.nodesTransfer(val1)
    fixture.detectChanges()
    expect(res).toStrictEqual([
      {
        key: 'env1',
        title: 'env1',
        name: 'env1',
        clusters: [
          {
            key: 'clu1_env1',
            name: 'clu1',
            title: 'clu1',
            isLeaf: true,
            selected: true
          },
          { key: 'clu2_env1', title: 'clu2', isLeaf: true, name: 'clu2' },
          { key: 'clu3_env1', title: 'clu3', isLeaf: true, name: 'clu3' }
        ],
        children: [
          {
            key: 'clu1_env1',
            name: 'clu1',
            title: 'clu1',
            isLeaf: true,
            selected: true
          },
          { key: 'clu2_env1', title: 'clu2', isLeaf: true, name: 'clu2' },
          { key: 'clu3_env1', title: 'clu3', isLeaf: true, name: 'clu3' }
        ],
        expanded: true
      },
      {
        key: 'env2',
        title: 'env2',
        name: 'env2',
        clusters: [
          { key: 'clu4_env2', name: 'clu4', title: 'clu4', isLeaf: true },
          { key: 'clu5_env2', name: 'clu5', title: 'clu5', isLeaf: true },
          { key: 'clu6_env2', name: 'clu6', title: 'clu6', isLeaf: true }
        ],
        children: [
          { key: 'clu4_env2', name: 'clu4', title: 'clu4', isLeaf: true },
          { key: 'clu5_env2', name: 'clu5', title: 'clu5', isLeaf: true },
          { key: 'clu6_env2', name: 'clu6', title: 'clu6', isLeaf: true }
        ]
      },
      {
        key: 'env3',
        title: 'env3',
        name: 'env3',
        clusters: [],
        isLeaf: true
      }
    ])
  })

  it('##getGroupItemSelected test', fakeAsync(() => {
    component.nodesList = [
      {
        title: 'title1',
        key: 'title1',
        children: [
          { title: 'title11', name: 'title11', key: 'title11', selected: true }
        ],
        expanded: true
      },
      {
        title: 'title2',
        key: 'title2',
        children: [
          { title: 'title21', name: 'title21', key: 'title21' },
          { title: 'title22', name: 'title22', key: 'title22' },
          { title: 'title23', name: 'title23', key: 'title23' }
        ],
        expanded: false
      },
      {
        title: 'title3',
        key: 'title3',
        children: [{ title: 'title31', name: 'title31', key: 'title31' }],
        expanded: false
      }
    ]
    component.clusterName = 'title22'
    component.getGroupItemSelected()

    expect(component.nodesList).toStrictEqual([
      {
        title: 'title1',
        key: 'title1',
        children: [
          { title: 'title11', name: 'title11', key: 'title11', selected: false }
        ],
        expanded: true
      },
      {
        title: 'title2',
        key: 'title2',
        children: [
          {
            title: 'title21',
            name: 'title21',
            key: 'title21',
            selected: false
          },
          { title: 'title22', name: 'title22', key: 'title22', selected: true },
          { title: 'title23', name: 'title23', key: 'title23', selected: false }
        ],
        expanded: true
      },
      {
        title: 'title3',
        key: 'title3',
        children: [
          { title: 'title31', name: 'title31', key: 'title31', selected: false }
        ],
        expanded: false
      }
    ])
  }))

  it('##activeNode test', () => {
    component.clusterKey = ''
    component.clusterName = ''
    component.activatedNode = null
    component.showList = false
    component.showApiPage = false
    component.strategyType = 'traffic'

    component.eoNgTreeDefault.getTreeNodeByKey = () => {
      return null
    }
    // @ts-ignore
    const data: any = {
      keys: ['env1'],
      node: { isExpanded: false, origin: { name: 'test', children: [] } }
    }

    component.activeNode(data)
    fixture.detectChanges()

    expect(component.clusterKey).toStrictEqual('env1')
    expect(component.clusterName).toStrictEqual('test')
    expect(component.activatedNode).toStrictEqual({
      isExpanded: false,
      origin: { name: 'test', children: [] }
    })
    expect(component.showList).toStrictEqual(true)
    expect(component.showApiPage).toStrictEqual(false)

    component.clusterKey = 'env1'
    component.clusterName = ''
    component.activatedNode = null
    component.showList = false
    component.showApiPage = false

    const data2: any = {
      keys: ['env1'],
      node: { isExpanded: false, origin: { name: 'test', children: ['1'] } }
    }

    component.activeNode(data2)
    fixture.detectChanges()

    expect(component.clusterKey).toStrictEqual('env1')
    expect(component.clusterName).toStrictEqual('')
    expect(component.activatedNode).toStrictEqual(null)
    expect(component.showList).toStrictEqual(false)
    expect(component.showApiPage).toStrictEqual(false)
    expect(data2.node.isExpanded).toStrictEqual(true)

    component.clusterKey = 'env2'
    component.clusterName = ''
    component.activatedNode = null
    component.showList = false
    component.showApiPage = false
    component.eoNgTreeDefault.getTreeNodeByKey = () => {
      return new NzTreeNode(
        { title: 'test', key: 'env2', isSelected: true },
        null,
        null
      )
    }

    component.activeNode(data2)
    fixture.detectChanges()
    expect(data2.node.isExpanded).toStrictEqual(false)
  })
})
