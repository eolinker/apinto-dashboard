/* eslint-disable dot-notation */
import { ComponentFixture, fakeAsync, TestBed } from '@angular/core/testing'
import { HttpClientModule } from '@angular/common/http'
import { ApiService, API_URL } from 'projects/core/src/app/service/api.service'
import { ActivatedRoute } from '@angular/router'
import { ElementRef, Renderer2, ChangeDetectorRef, LOCALE_ID } from '@angular/core'
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
import { AuditLogListComponent } from './list.component'
import { FormsModule, ReactiveFormsModule } from '@angular/forms'
import { EoNgInputModule } from 'eo-ng-input'
import { EoNgButtonModule } from 'eo-ng-button'
import { EoNgTableModule } from 'eo-ng-table'
import { EoNgDatePickerModule } from 'eo-ng-date-picker'
// import { NzDatePickerModule } from 'ng-zorro-antd/date-picker'
import { NZ_DATE_LOCALE } from 'ng-zorro-antd/i18n'

import { APP_BASE_HREF, registerLocaleData } from '@angular/common'
import zh from '@angular/common/locales/zh'
registerLocaleData(zh)

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

describe('AuditLogListComponent test', () => {
  let component: AuditLogListComponent
  let fixture: ComponentFixture<AuditLogListComponent>
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
        FormsModule,
        ReactiveFormsModule,
        BidiModule,
        NoopAnimationsModule,
        NzNoAnimationModule,
        EoNgInputModule,
        EoNgButtonModule,
        HttpClientModule,
        EoNgTableModule,
        EoNgDatePickerModule
        // NzDatePickerModule
      ],
      declarations: [AuditLogListComponent],
      providers: [
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: ElementRef, useValue: new MockElementRef() },
        { provide: Renderer2, useClass: MockRenderer },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService },
        { provide: EoNgFeedbackModalService, useClass: MockEnsureService },
        { provide: ChangeDetectorRef, useClass: ChangeDetectorRef },
        { provide: NZ_DATE_LOCALE, useValue: zh },
        {
          provide: LOCALE_ID,
          useValue: 'zh-CN'
        },
        {
          provide: ActivatedRoute,
          useValue: {
            queryParams: of({ clusterName: 'clus2' })
          }
        }
      ]
    }).compileComponents()

    fixture = TestBed.createComponent(AuditLogListComponent)
    component = fixture.componentInstance
    fixture.detectChanges()
  })

  it('##should create', () => {
    expect(component).toBeTruthy()
  })

  it('##ngOnInit should call getLogList and getTargetList ', () => {
    const spyGetLogList = jest.spyOn(component, 'getLogList')
    const spyGetTargetList = jest.spyOn(component, 'getTargetList')
    expect(spyGetLogList).not.toHaveBeenCalled()
    expect(spyGetTargetList).not.toHaveBeenCalled()

    component.ngOnInit()
    fixture.detectChanges()
    expect(spyGetLogList).toHaveBeenCalled()
    expect(spyGetTargetList).toHaveBeenCalled()
  })

  it('##getLogList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: {
          items: [{
            id: 0,
            operator: {
              user_id: 0,
              keyword: '',
              nickname: '',
              avatar: ''
            },
            operateType: '',
            kind: '',
            time: '',
            ip: ''
          }],
          total: 1
        },
        msg: ''
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getLogList()
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toBeCalledWith('audit-logs', {
      keyword: '',
      operateType: '',
      kind: '',
      end: Math.floor(new Date().getTime() / 1000),
      pageSize: 20,
      pageNum: 1,
      total: 0
    })

    expect(component.logsList).toStrictEqual([{
      id: 0,
      eoKey: 123456789,
      operator: {
        user_id: 0,
        keyword: '',
        nickname: '',
        avatar: ''
      },
      operateType: '',
      kind: '',
      time: '',
      ip: ''
    }])
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getLogList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.searchData = {
      keyword: 'test',
      operateType: 'test1',
      kind: 'test2',
      start: new Date(),
      end: null,
      pageSize: 20,
      pageNum: 1,
      total: 0
    }
    component.getLogList()
    fixture.detectChanges()

    expect(spyService).toBeCalledWith('audit-logs', {
      keyword: 'test',
      operateType: 'test1',
      kind: 'test2',
      start: Math.floor(new Date().getTime() / 1000),
      end: Math.floor(new Date().getTime() / 1000),
      pageSize: 20,
      pageNum: 1,
      total: 0
    })

    expect(component.logsList).toStrictEqual([])
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
  })

  it('##getLogList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.searchData = {
      keyword: 'test',
      operateType: 'test1',
      kind: 'test2',
      start: new Date(),
      end: new Date(1667447768000),
      pageSize: 20,
      pageNum: 1,
      total: 0
    }
    component.getLogList()
    fixture.detectChanges()

    expect(spyService).toBeCalledWith('audit-logs', {
      keyword: 'test',
      operateType: 'test1',
      kind: 'test2',
      start: Math.floor(new Date().getTime() / 1000),
      end: Math.floor(new Date(1667447768000).getTime() / 1000),
      pageSize: 20,
      pageNum: 1,
      total: 0
    })

    expect(component.logsList).toStrictEqual([])
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalledWith('获取列表数据失败！')
  })

  it('##getTargetList with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: {
          items: [{
            title: '1',
            name: '2'
          },
          {
            title: '2',
            name: '3'
          },
          {
            title: '3',
            name: '4'
          }]
        },
        msg: 'success'
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getTargetList()
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toBeCalledWith('audit-log/kinds')

    expect(component.listOfKind).toStrictEqual([{
      title: '1',
      value: '2',
      label: '1',
      name: '2'
    },
    {
      title: '2',
      value: '3',
      label: '2',
      name: '3'
    },
    {
      title: '3',
      value: '4',
      label: '3',
      name: '4'
    }])
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getTargetList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getTargetList()
    fixture.detectChanges()

    expect(spyService).toBeCalledWith('audit-log/kinds')

    expect(component.listOfKind).toStrictEqual([])
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toBeCalledWith('fail')
  })

  it('##getTargetList with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()

    component.getTargetList()
    fixture.detectChanges()

    expect(spyService).toBeCalledWith('audit-log/kinds')

    expect(component.logsList).toStrictEqual([])
    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalledWith('获取操作对象数据失败！')
  })

  it('##getLogDetail with success return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest.spyOn(httpCommonService, 'get').mockReturnValue(
      of({
        code: 0,
        data: {
          args: [{
            attr: '1',
            value: '12'
          },
          {
            attr: '2',
            value: '22'
          },
          {
            attr: '3',
            value: '32'
          }]
        },
        msg: ''
      })
    )
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    const spyOpenDrawer = jest.spyOn(component, 'openDrawer')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    expect(spyOpenDrawer).not.toHaveBeenCalled()

    component.getLogDetail({ id: 'testId' })
    fixture.detectChanges()
    expect(spyService).toHaveBeenCalled()

    expect(component.auditLogDetail).toStrictEqual([{
      attr: '1',
      value: '12',
      eoKey: 123456789
    },
    {
      attr: '2',
      value: '22',
      eoKey: 123456789
    },
    {
      attr: '3',
      value: '32',
      eoKey: 123456789
    }])
    expect(spyOpenDrawer).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalled()
    expect(spyService).toHaveBeenCalledWith('audit-log', { logId: 'testId' })
    expect(spyMessage).not.toHaveBeenCalled()
  })

  it('##getLogDetail with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1, msg: 'fail' }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.logsTableClick({ data: { id: 'testId' } })
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalledWith('fail')
  })

  it('##getLogDetail with fail return', () => {
    const httpCommonService = fixture.debugElement.injector.get(ApiService)
    const spyService = jest
      .spyOn(httpCommonService, 'get')
      .mockReturnValue(of({ code: -1 }))
    // @ts-ignore
    const spyMessage = jest.spyOn(component.message, 'error')
    expect(spyService).not.toHaveBeenCalled()
    expect(spyMessage).not.toHaveBeenCalled()
    component.logsTableBody[5].btns[0].click({ id: 'testId' })
    fixture.detectChanges()

    expect(spyService).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalled()
    expect(spyMessage).toHaveBeenCalledWith('获取日志详情失败！')
  })

  it('##openDrawer', () => {
    expect(component.drawerRef).toBeUndefined()

    component.openDrawer()
    fixture.detectChanges()
    expect(component.drawerRef).not.toBeUndefined()
  })

  it('##clearSearch test', fakeAsync(() => {
    component.searchData = {
      keyword: '123',
      operateType: '213',
      kind: '32123',
      start: null,
      end: null,
      pageSize: 20,
      pageNum: 1,
      total: 0
    }
    component.date = [1, 2]

    component.clearSearch()

    expect(component.searchData).toStrictEqual({
      keyword: '',
      operateType: '',
      kind: '',
      start: null,
      end: null,
      pageSize: 20,
      pageNum: 1,
      total: 0
    })

    expect(component.date).toStrictEqual([])
  }))

  it('##onDateRangeChange test', fakeAsync(() => {
    expect(component.searchData.start).toStrictEqual(null)
    expect(component.searchData.end).toStrictEqual(null)

    component.onDateRangeChange([new Date(111111111), new Date()])
    expect(component.searchData.start).toStrictEqual(new Date(111111111))
    expect(component.searchData.end).toStrictEqual(new Date())

    component.onDateRangeChange([new Date(22222), new Date(2222222222232)])
    expect(component.searchData.start).toStrictEqual(new Date(22222))
    expect(component.searchData.end).toStrictEqual(new Date(2222222222232))
  }))
})
