/* eslint-disable dot-notation */
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

import { EoNgSelectModule } from 'eo-ng-select'
import { ListComponent } from './list.component'
import { FuseCreateComponent } from '../fuse/create/create.component'
import { TrafficCreateComponent } from '../traffic/create/create.component'
import { GreyCreateComponent } from '../grey/create/create.component'
import { VisitCreateComponent } from '../visit/create/create.component'
import { CacheCreateComponent } from '../cache/create/create.component'
import { ReactiveFormsModule } from '@angular/forms'

class MockRenderer {
  removeAttribute (element: any, cssClass: string) {
    return cssClass + 'is removed from' + element
  }
}

class MockMessageService {
  success () {
    return 'success'
  }

  error () {
    return 'error'
  }
}

class MockEnsureService {
  create () {
    return 'modal is create'
  }
}

jest.mock('uuid', () => {
  return {
    v4: () => 123456789
  }
})

describe('ListComponent test', () => {
  let component: ListComponent
  let fixture: ComponentFixture<ListComponent>
  class MockElementRef extends ElementRef {
    constructor () {
      super(null)
    }
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        NzOverlayModule,
        EoNgSelectModule,
        ReactiveFormsModule,
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
      declarations: [ListComponent],
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

    fixture = TestBed.createComponent(ListComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('##should create', () => {
    expect(component).toBeTruthy()
  })

  it('##ngOnInit should subscribe urlQueryParams change', () => {
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
    component.clusterName = ''

    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetStrategiesList).toHaveBeenCalledTimes(1)
    expect(component.clusterName).toStrictEqual('clus2')
  })

  it('##ngOnDestroy test', () => {
    // @ts-ignore
    const spySubscription = jest.spyOn(component.subscription, 'unsubscribe')
    // @ts-ignore
    expect(spySubscription).not.toHaveBeenCalled()
    component.ngOnDestroy()
    expect(spySubscription).toHaveBeenCalled()
  })

  it('##click table btns', () => {
    component.strategyType = 'traffic'
    const spyEditStrategy = jest.spyOn(component, 'editStrategy')
    // @ts-ignore
    const spyModalService = jest.spyOn(component.modalService, 'create')
    const spyRecoverStrategy = jest.spyOn(component, 'recoverStrategy')
    expect(spyEditStrategy).toHaveBeenCalledTimes(0)
    expect(spyModalService).toHaveBeenCalledTimes(0)
    expect(spyRecoverStrategy).toHaveBeenCalledTimes(0)
    const open = jest.fn()
    Object.defineProperty(window, 'open', open)
    expect(open).not.toHaveBeenCalled()

    const item = { key: 'test' }
    component.strategiesTableBody[8].btns[0].click(item)
    expect(spyEditStrategy).toHaveBeenCalledTimes(1)
    expect(spyModalService).toHaveBeenCalledTimes(0)
    expect(spyRecoverStrategy).toHaveBeenCalledTimes(0)
    component.strategiesTableBody[8].btns[1].click(item)
    expect(spyEditStrategy).toHaveBeenCalledTimes(1)
    expect(spyModalService).toHaveBeenCalledTimes(1)
    expect(spyRecoverStrategy).toHaveBeenCalledTimes(0)

    component.strategiesTableBody[9].btns[0].click(item)
    expect(spyEditStrategy).toHaveBeenCalledTimes(2)
    expect(spyModalService).toHaveBeenCalledTimes(1)
    expect(spyRecoverStrategy).toHaveBeenCalledTimes(0)
    component.strategiesTableBody[9].btns[1].click(item)
    expect(spyEditStrategy).toHaveBeenCalledTimes(2)
    expect(spyModalService).toHaveBeenCalledTimes(1)
    expect(spyRecoverStrategy).toHaveBeenCalledTimes(1)

    component.strategiesTableClick({ data: { uuid: 'test' } })
    expect(spyEditStrategy).toHaveBeenCalledTimes(3)
  })

  it('##getPublishList with success return and all data', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(
        of({
          code: 0,
          data: {
            strategies: [1, 2, 3],
            is_publish: true,
            source: '123',
            version_name: 'test1',
            unpublish_msg: 'unpublish_test'
          }
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getPublishList()
    fixture.detectChanges()
    expect(component.publishList).toStrictEqual([1, 2, 3])
    expect(component.strategyIsPublish).toStrictEqual(true)
    expect(component.strategySource).toStrictEqual('123')
    expect(component.validateForm.controls['version_name'].value).toStrictEqual(
      'test1'
    )
    expect(component.strategyUnpulishMsg).toStrictEqual('unpublish_test')
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getPublishList with success return and no strategies', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(
        of({
          code: 0,
          data: {
            strategies: [],
            is_publish: false,
            source: '123',
            version_name: 'test1',
            unpublish_msg: ''
          }
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getPublishList()
    fixture.detectChanges()
    expect(component.publishList).toStrictEqual([])
    expect(component.strategyIsPublish).toStrictEqual(false)
    expect(component.strategySource).toStrictEqual('123')
    expect(component.validateForm.controls['version_name'].value).toStrictEqual(
      'test1'
    )
    expect(component.strategyUnpulishMsg).toStrictEqual('当前无可发布策略')
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getPublishList with success return and not publish', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(
        of({
          code: 0,
          data: {
            strategies: [1, 2, 3],
            is_publish: false,
            source: '123',
            version_name: 'test1',
            unpublish_msg: ''
          }
        })
      )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getPublishList()
    fixture.detectChanges()
    expect(component.publishList).toStrictEqual([1, 2, 3])
    expect(component.strategyIsPublish).toStrictEqual(false)
    expect(component.strategySource).toStrictEqual('123')
    expect(component.validateForm.controls['version_name'].value).toStrictEqual(
      'test1'
    )
    expect(component.strategyUnpulishMsg).toStrictEqual('当前策略不可发布')
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getPublishList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getPublishList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getStrategiesList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: {
          strategies: [
            { priority: 1, uuid: '1' },
            { priority: 2, uuid: '2' }
          ],
          is_publish: true,
          source: '123',
          version_name: 'test1',
          unpublish_msg: 'unpublish_test'
        }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.clusterName = ''
    component.getStrategiesList()
    fixture.detectChanges()
    expect(spyService).not.toHaveBeenCalled()

    component.clusterName = '1'
    component.strategiesList = []
    component.priorityMap = new Map()
    component.getStrategiesList()
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalled()

    expect(component.strategiesList).toStrictEqual([
      {
        priority: 1,
        uuid: '1'
      },
      {
        priority: 2,
        uuid: '2'
      }
    ])
    expect(component.priorityMap.get(1)).toStrictEqual([
      { priority: 1, uuid: '1' }
    ])
    expect(component.priorityMap.get(2)).toStrictEqual([
      { priority: 2, uuid: '2' }
    ])
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getStrategiesList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.clusterName = 'test'
    component.getStrategiesList()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##addStrategy & editStrategy test', () => {
    component.strategyType = 'traffic'
    // @ts-ignore
    const spyChangeToContent = jest.spyOn(component.router, 'navigate')
    expect(spyChangeToContent).not.toHaveBeenCalled()

    component.addStrategy()
    fixture.detectChanges()

    expect(spyChangeToContent).toHaveBeenCalledTimes(1)

    component.editStrategy(1)
    fixture.detectChanges()

    expect(spyChangeToContent).toHaveBeenCalledTimes(2)
  })

  it('##stopStrategy with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'patch')
      .mockReturnValue(of({ code: 0, data: {} }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.stopStrategy(1, true)
    fixture.detectChanges()

    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyGetStrategiesList).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##stopStrategy with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'patch')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.stopStrategy(1, false)
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
  })

  it('##stopStrategy with fail return without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'patch')
      .mockReturnValue(of({ code: -1, msg: '' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.stopStrategy(1, false)
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('启用策略失败!')
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    component.stopStrategy(1, true)
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('停用策略失败!')
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
  })

  it('##stopStrategy with fail return without msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'patch')
      .mockReturnValue(of({ code: -1, msg: '' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.stopStrategy(1, true)
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('停用策略失败!')
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
  })

  it('##deleteStrategy with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'delete')
      .mockReturnValue(of({ code: 0, data: {} }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    component.deleteStrategy(1)
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
  })

  it('##deleteStrategy with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'delete')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    component.deleteStrategy(1)
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
  })

  it('##deleteStrategy with fail return but no msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'delete')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    component.deleteStrategy(1)
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('删除策略失败!')
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
  })

  it('##recoverStrategy with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'patch').mockReturnValue(
      of({
        code: 0,
        data: {
          strategies: [
            { priority: 1, uuid: '1' },
            { priority: 2, uuid: '2' }
          ],
          is_publish: true,
          source: '123',
          version_name: 'test1',
          unpublish_msg: 'unpublish_test'
        }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    component.clusterName = ''
    component.recoverStrategy(1)
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyGetStrategiesList).toHaveBeenCalled()
  })

  it('##recoverStrategy with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'patch')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    component.clusterName = 'test'
    component.recoverStrategy('1')
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
  })

  it('##openDrawer and reset validateForm', () => {
    const spyGetPublishList = jest.spyOn(component, 'getPublishList')
    expect(component.drawerPublishRef).toBeUndefined()
    expect(spyGetPublishList).not.toHaveBeenCalled()
    component.validateForm.controls['desc'].setValue('test1')
    component.validateForm.controls['version_name'].setValue('test2')
    component.openDrawer('test')
    fixture.detectChanges()
    expect(component.drawerPublishRef).toBeUndefined()
    expect(spyGetPublishList).not.toHaveBeenCalled()
    expect(component.validateForm.controls['desc'].value).toStrictEqual('test1')
    expect(component.validateForm.controls['version_name'].value).toStrictEqual(
      'test2'
    )

    component.openDrawer('publish')
    fixture.detectChanges()
    expect(component.drawerPublishRef).not.toBeUndefined()
    expect(spyGetPublishList).toHaveBeenCalled()
    expect(component.validateForm.controls['desc'].value).toStrictEqual('')
    expect(component.validateForm.controls['version_name'].value).toStrictEqual(
      ''
    )
  })

  it('##cancelDrawer', () => {
    component.openDrawer('publish')
    // @ts-ignore
    component.drawerPublishRef.close = () => {
      return 'drawer is close'
    }

    // @ts-ignore
    const spyFn = jest.spyOn(component.drawerPublishRef, 'close')
    expect(spyFn).not.toHaveBeenCalled()

    component.cancelDrawer('test')
    fixture.detectChanges()
    expect(spyFn).not.toHaveBeenCalled()
    // expect(spyClearForm).toHaveBeenCalledTimes(0)

    component.cancelDrawer('publish')
    fixture.detectChanges()
    expect(spyFn).toHaveBeenCalled()
    // expect(spyClearForm).toHaveBeenCalledTimes(1)
  })

  it('##publish with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(
      of({
        code: 0,
        data: {
          strategies: [
            { priority: 1, uuid: '1' },
            { priority: 2, uuid: '2' }
          ],
          is_publish: true,
          source: '123',
          version_name: 'test1',
          unpublish_msg: 'unpublish_test'
        }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    component.openDrawer('publish')
    // @ts-ignore
    component.drawerPublishRef.close = () => {
      return 'drawer is close'
    }
    // @ts-ignore
    const spyFn = jest.spyOn(component.drawerPublishRef, 'close')
    expect(spyFn).not.toHaveBeenCalled()
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    // 表单验证失败 且未返回isPublish
    const spyOnControlVerMarkAsDirty = jest.spyOn(
      component.validateForm.controls['version_name'],
      'markAsDirty'
    )
    component.publish()
    fixture.detectChanges()
    expect(spyFn).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
    expect(spyOnControlVerMarkAsDirty).toHaveBeenCalledTimes(1)

    // 表单验证成功 未返回isPublish
    component.validateForm.controls['version_name'].setValue('test')
    component.strategyIsPublish = ''
    component.publish()
    fixture.detectChanges()
    expect(spyFn).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    // 表单验证失败 返回isPublish
    component.validateForm.controls['version_name'].setValue('')
    component.strategyIsPublish = 'test'
    component.publish()
    fixture.detectChanges()
    expect(spyFn).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
    expect(spyOnControlVerMarkAsDirty).toHaveBeenCalledTimes(2)

    // 表单验证成功 返回isPublish
    component.validateForm.controls['version_name'].setValue('test')
    component.strategyIsPublish = 'test'
    component.publish()
    fixture.detectChanges()
    expect(spyFn).toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyGetStrategiesList).toHaveBeenCalled()
  })

  it('##publish with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'post')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    component.openDrawer('publish')
    // @ts-ignore
    component.drawerPublishRef.close = () => {
      return 'drawer is close'
    }
    // @ts-ignore
    const spyFn = jest.spyOn(component.drawerPublishRef, 'close')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    component.validateForm.controls['version_name'].setValue('test')
    component.strategyIsPublish = 'test'
    component.publish()
    fixture.detectChanges()

    expect(spyFn).not.toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
  })

  it('##publish with fail return no msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'post')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')
    component.openDrawer('publish')
    // @ts-ignore
    component.drawerPublishRef.close = () => {
      return 'drawer is close'
    }
    // @ts-ignore
    const spyFn = jest.spyOn(component.drawerPublishRef, 'close')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    component.validateForm.controls['version_name'].setValue('test')
    component.strategyIsPublish = 'test'
    component.publish()
    fixture.detectChanges()

    expect(spyFn).not.toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('发布策略失败!')
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
  })

  it('##changeEditingPriority 修改优先级', fakeAsync(() => {
    component.editingPriority = ''
    const e: Event = new Event('test')
    const spyEPropagation = jest.spyOn(e, 'stopPropagation')
    expect(spyEPropagation).not.toHaveBeenCalled()
    component.changeEditingPriority(e, '2')
    expect(component.editingPriority).toStrictEqual(2)
    expect(spyEPropagation).toHaveBeenCalledTimes(1)
    component.changeEditingPriority(e, '')
    expect(component.editingPriority).toStrictEqual('NULL')
    expect(spyEPropagation).toHaveBeenCalledTimes(2)
  }))

  it('##disabledEdit test', fakeAsync(() => {
    component.nzDisabled = false
    component.disabledEdit(true)
    expect(component.nzDisabled).toStrictEqual(true)
    component.disabledEdit(false)
    expect(component.nzDisabled).toStrictEqual(false)
  }))

  it('##checkPriority 输入为空或大于999时, 放入priorityMap中key为NULL的数组中,并提示优先级为空或大于999不允许提交', fakeAsync(() => {
    const spyChangePriorityMap = jest.spyOn(component, 'changePriorityMap')
    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    const spyChangePriority = jest.spyOn(component, 'changePriority')
    expect(spyChangePriorityMap).not.toHaveBeenCalled()
    expect(spyMessageError).not.toHaveBeenCalled()
    expect(spyChangePriority).not.toHaveBeenCalled()
    component.editingPriority = ''
    component.checkPriority('', 'test')
    expect(spyChangePriorityMap).not.toHaveBeenCalled()
    expect(spyMessageError).toBeCalledWith('优先级不能为空, 请填写后提交')
    expect(spyChangePriority).not.toHaveBeenCalled()

    component.editingPriority = 5
    component.priorityMap.set(5, [{ uuid: 'test', name: 'testName' }])
    component.checkPriority('', 'test')
    expect(spyChangePriorityMap).toHaveBeenCalled()
    expect(spyMessageError).toHaveBeenCalledTimes(2)
    expect(spyMessageError).toBeCalledWith('优先级不能为空, 请填写后提交')
    expect(spyChangePriority).not.toHaveBeenCalled()

    component.checkPriority('9999', 'test')
    expect(spyChangePriorityMap).toHaveBeenCalledTimes(2)
    expect(spyMessageError).toHaveBeenCalledTimes(3)
    expect(spyMessageError).toBeCalledWith(
      '优先级范围在1-999之间，请修改后提交'
    )
    expect(spyChangePriority).not.toHaveBeenCalled()
  }))

  it('##输入不为空, 检查priorityMap中相同priority的数组,', fakeAsync(() => {
    const spyChangePriorityMap = jest.spyOn(component, 'changePriorityMap')
    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    const spyChangePriority = jest.spyOn(component, 'changePriority')
    // @ts-ignore
    const spyScrollToAnchor = jest.spyOn(component.viewportScroller, 'scrollToAnchor'
    )
    expect(spyChangePriorityMap).not.toHaveBeenCalled()
    expect(spyMessageError).not.toHaveBeenCalled()
    expect(spyChangePriority).not.toHaveBeenCalled()
    expect(spyScrollToAnchor).not.toHaveBeenCalled()

    component.priorityDangerP = []
    component.editingPriority = 5
    component.priorityMap.set(3, [{ uuid: 'test1', name: 'test1Name' }])
    component.priorityMap.set(5, [{ uuid: 'test', name: 'testName' }])
    component.checkPriority(3, 'test')
    expect(spyChangePriorityMap).toHaveBeenCalled()
    expect(spyMessageError).toHaveBeenCalledTimes(1)
    expect(spyMessageError).toBeCalledWith(
      '修改后的优先级与test1Name冲突，无法自动提交'
    )
    expect(spyChangePriority).not.toHaveBeenCalled()
    expect(component.priorityDangerP).toStrictEqual([3])

    component.priorityMap.set(6, [{ uuid: 'test', name: 'testName' }])
    component.editingPriority = 6
    component.checkPriority('3', 'test')
    expect(spyChangePriorityMap).toHaveBeenCalledTimes(2)
    expect(spyMessageError).toHaveBeenCalledTimes(2)
    expect(spyMessageError).toBeCalledWith(
      '优先级存在冲突或数值超出范围，无法自动提交'
    )
    expect(spyChangePriority).not.toHaveBeenCalled()
    expect(component.priorityDangerP).toStrictEqual([3])

    component.priorityMap.set(5, [{ uuid: 'test', name: 'testName' }])
    component.editingPriority = 5
    component.priorityDangerP = []
    component.checkPriority('3', 'test')
    expect(spyChangePriorityMap).toHaveBeenCalledTimes(3)
    expect(spyMessageError).toHaveBeenCalledTimes(3)
    expect(spyMessageError).toBeCalledWith(
      '优先级存在冲突或数值超出范围，无法自动提交'
    )
    expect(spyChangePriority).not.toHaveBeenCalled()
    expect(component.priorityDangerP).toStrictEqual([3])
  }))

  it('##checkPriority 存在冲突优先级，无法成功提交', () => {
    const spyChangePriorityMap = jest.spyOn(component, 'changePriorityMap')
    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyChangePriority = jest.spyOn(component, 'changePriority')
    // @ts-ignore
    const spyScrollToAnchor = jest.spyOn(component.viewportScroller, 'scrollToAnchor'
    )
    expect(spyChangePriorityMap).not.toHaveBeenCalled()
    expect(spyMessageError).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyChangePriority).not.toHaveBeenCalled()
    expect(spyScrollToAnchor).not.toHaveBeenCalled()

    component.priorityDangerP = []
    component.editingPriority = 5
    component.priorityMap.set(1, [{ uuid: 'test' }, { uuid: 'test2' }])
    component.priorityMap.set('NULL', [{ uuid: 'test' }, { uuid: 'test2' }])
    component.priorityMap.set(3, [])
    component.priorityMap.set(5, [{ uuid: 'test', name: 'testName' }])
    component.checkPriority(3, 'test')
    expect(spyChangePriorityMap).toHaveBeenCalled()
    expect(spyMessageError).toHaveBeenCalled()
    expect(spyMessageError).toBeCalledWith(
      '优先级存在冲突或数值超出范围，无法自动提交'
    )
    expect(spyChangePriority).not.toHaveBeenCalled()
  })

  it('##checkPriority 无冲突优先级，成功提交', () => {
    const spyChangePriorityMap = jest.spyOn(component, 'changePriorityMap')
    // @ts-ignore
    const spyMessageError = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyChangePriority = jest.spyOn(component, 'changePriority')
    // @ts-ignore
    const spyScrollToAnchor = jest.spyOn(component.viewportScroller, 'scrollToAnchor'
    )
    expect(spyChangePriorityMap).not.toHaveBeenCalled()
    expect(spyMessageError).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyChangePriority).not.toHaveBeenCalled()
    expect(spyScrollToAnchor).not.toHaveBeenCalled()

    component.priorityDangerP = []
    component.editingPriority = 5
    component.priorityMap.set(3, [])
    component.priorityMap.set(5, [{ uuid: 'test', name: 'testName' }])
    component.checkPriority(3, 'test')
    expect(spyChangePriorityMap).toHaveBeenCalled()
    expect(spyMessageError).not.toHaveBeenCalled()
    expect(spyChangePriority).toHaveBeenCalled()
  })

  it('##changePriorityMap 将策略从priorityMap的editingPriority中移除, 放入priority的数组中', fakeAsync(() => {
    component.priorityMap = new Map()
    component.priorityMap.set(1, [{ uuid: 'test1' }])
    component.priorityMap.set(2, [{ uuid: 'test2' }])
    component.priorityMap.set(3, [{ uuid: 'test3' }])

    component.changePriorityMap(1, 4, 'test1')
    expect(component.priorityMap.get(4)).toStrictEqual([{ uuid: 'test1' }])
    expect(component.priorityMap.get(1)).toStrictEqual([])

    component.changePriorityMap(2, 3, 'test2')
    expect(component.priorityMap.get(3)).toStrictEqual([
      { uuid: 'test3' },
      { uuid: 'test2' }
    ])
    expect(component.priorityMap.get(2)).toStrictEqual([])
  }))

  it('##changePriorityMap 将策略从priorityMap的editingPriority中移除, 放入priority的数组中', fakeAsync(() => {
    component.priorityDangerP = [1, 2, 3, 4, 5]
    expect(component.checkListStatus('')).toStrictEqual('error')
    expect(component.checkListStatus(2)).toStrictEqual('error')
    expect(component.checkListStatus(6)).toStrictEqual('')
  }))

  it('##checkPriorityMap test', fakeAsync(() => {
    component.priorityMap.set(1, [{ uuid: 'test' }, { uuid: 'test2' }])
    component.priorityMap.set('NULL', [{ uuid: 'test' }, { uuid: 'test2' }])
    expect(component.checkPriorityMap()).toStrictEqual(false)
    expect(component.priorityDangerP).toStrictEqual([1])
  }))

  it('##changePriority  with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'post').mockReturnValue(
      of({
        code: 0,
        data: {
          strategies: [
            { priority: 1, uuid: '1' },
            { priority: 2, uuid: '2' }
          ],
          is_publish: true,
          source: '123',
          version_name: 'test1',
          unpublish_msg: 'unpublish_test'
        }
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')

    component.priorityMap = new Map()
    component.priorityMap.set('NULL', [])
    component.priorityMap.set('2', [{ uuid: 'test1' }])
    component.priorityMap.set('21', [{ uuid: 'tets2' }, { uuid: 'test3' }])
    component.priorityMap.set(3, [{ uuid: 'test4' }])
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    component.changePriority()
    fixture.detectChanges()

    expect(component.prioritySaveMap).toStrictEqual({
      test1: 2,
      test4: 3
    })
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyMessageSuccess).toHaveBeenCalled()
    expect(spyGetStrategiesList).toHaveBeenCalled()
  })

  it('##changePriority  with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'post')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')

    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    component.changePriority()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
  })

  it('##changePriority  with fail return no msg', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'post')
      .mockReturnValue(of({ code: -1, msg: '' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    // @ts-ignore
    const spyMessageSuccess = jest.spyOn(component.message, 'success')
    const spyGetStrategiesList = jest.spyOn(component, 'getStrategiesList')

    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()

    component.changePriority()
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('修改优先级失败!')
    expect(spyMessageSuccess).not.toHaveBeenCalled()
    expect(spyGetStrategiesList).not.toHaveBeenCalled()
  })

  it('##strategyConf and getBreadcrumb test', fakeAsync(() => {
    // @ts-ignore
    const spyReqFlashBreadcrumb = jest.spyOn(component.appConfigService, 'reqFlashBreadcrumb'
    )
    expect(spyReqFlashBreadcrumb).not.toHaveBeenCalled()

    component.strategyType = 'test'
    expect(component.strategyConf()).toStrictEqual('限流规则')
    component.getBreadcrumb()
    expect(spyReqFlashBreadcrumb).lastCalledWith([{ title: '流量策略' }])

    component.strategyType = 'traffic'
    expect(component.strategyConf()).toStrictEqual('限流规则')
    component.getBreadcrumb()
    expect(spyReqFlashBreadcrumb).lastCalledWith([{ title: '流量策略' }])

    component.strategyType = 'grey'
    expect(component.strategyConf()).toStrictEqual('灰度规则')
    component.getBreadcrumb()
    expect(spyReqFlashBreadcrumb).lastCalledWith([{ title: '灰度策略' }])

    component.strategyType = 'fuse'
    expect(component.strategyConf()).toStrictEqual('熔断维度')
    component.getBreadcrumb()
    expect(spyReqFlashBreadcrumb).lastCalledWith([{ title: '熔断策略' }])

    component.strategyType = 'cache'
    expect(component.strategyConf()).toStrictEqual('缓存有效时间')
    component.getBreadcrumb()
    expect(spyReqFlashBreadcrumb).lastCalledWith([{ title: '缓存策略' }])

    component.strategyType = 'visit'
    expect(component.strategyConf()).toStrictEqual('访问规则')
    component.getBreadcrumb()
    expect(spyReqFlashBreadcrumb).lastCalledWith([{ title: '访问策略' }])
  }))
})
