/*
 * @Author: MengjieYang yangmengjie@eolink.com
 * @Date: 2022-08-07 21:45:15
 * @LastEditors: MengjieYang yangmengjie@eolink.com
 * @LastEditTime: 2022-08-07 23:21:55
 * @FilePath: /apinto/src/app/service/api.service.spec.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// Http testing module and mocking controller
import {
  HttpClientTestingModule,
  HttpTestingController
} from '@angular/common/http/testing'

// Other imports
import { TestBed } from '@angular/core/testing'
import { HttpErrorResponse } from '@angular/common/http'
import { Data } from '@angular/router'
import { ApiService, API_URL } from './api.service'
import { Overlay } from '@angular/cdk/overlay'
import { APP_BASE_HREF } from '@angular/common'
import { EoNgFeedbackMessageService } from 'eo-ng-feedback'
import { environment } from 'projects/core/src/environments/environment'

class MockMessageService {
  error () {
    return 'error'
  }
}

describe('HttpClient testing', () => {
  let service: ApiService
  let httpTestingController: HttpTestingController

  beforeEach(() => {
    TestBed.configureTestingModule({
      imports: [HttpClientTestingModule],
      providers: [
        ApiService,
        { provide: Overlay, useClass: Overlay },
        { provide: APP_BASE_HREF, useValue: '/' },
        { provide: API_URL, useValue: environment.urlPrefix },
        { provide: EoNgFeedbackMessageService, useClass: MockMessageService }
      ]
    })

    // Inject the http service and test controller for each test
    service = TestBed.inject(ApiService)
    httpTestingController = TestBed.inject(HttpTestingController)
  })

  /// Tests begin ///

  it('can test HttpClient.get', () => {
    const testData: Data = { name: 'Test Data' }

    // Make an HTTP GET request
    service.get('clusters').subscribe((data) =>
      // When observable resolves, result should match test data
      expect(data).toEqual(testData)
    )

    // The following `expectOne()` will match the request's URL.
    // If no requests or multiple requests matched that URL
    // `expectOne()` would throw.
    const req = httpTestingController.expectOne(
      'https://mockapi.eolink.com/K25EPjsf31dac8880a551fe2672247d21218bf854cbcf60/api/clusters?namespace=default'
    )

    // Assert that the request is a GET.
    expect(req.request.method).toEqual('GET')

    // Respond with mock data, causing Observable to resolve.
    // Subscribe callback asserts that correct data was returned.
    req.flush(testData)

    // Finally, assert that there are no outstanding requests.
    httpTestingController.verify()
  })

  it('can test HttpClient.get with params', () => {
    const testData: Data = { name: 'Test Data' }

    // Make an HTTP GET request
    service
      .get('clusters', { namespace: 'test', query: { test: 1 } })
      .subscribe((data) =>
        // When observable resolves, result should match test data
        expect(data).toEqual(testData)
      )
    const req = httpTestingController.expectOne(
      'https://mockapi.eolink.com/K25EPjsf31dac8880a551fe2672247d21218bf854cbcf60/api/clusters?namespace=default&query=%7B%22test%22:1%7D'
    )

    expect(req.request.method).toEqual('GET')
    req.flush(testData)
    httpTestingController.verify()
  })

  it('can test HttpClient.post', () => {
    const testData: Data = { name: 'Test Data' }

    // Make an HTTP GET request
    service.post('clusters').subscribe((data) =>
      // When observable resolves, result should match test data
      expect(data).toEqual(testData)
    )

    const req = httpTestingController.expectOne(
      'https://mockapi.eolink.com/K25EPjsf31dac8880a551fe2672247d21218bf854cbcf60/api/clusters?namespace=default'
    )
    expect(req.request.method).toEqual('POST')
    req.flush(testData)
    httpTestingController.verify()
  })

  it('can test HttpClient.post with params', () => {
    const testData: Data = { name: 'Test Data' }

    // Make an HTTP GET request
    service.post('clusters', { namespace: 'test' }).subscribe((data) =>
      // When observable resolves, result should match test data
      expect(data).toEqual(testData)
    )

    const req = httpTestingController.expectOne(
      'https://mockapi.eolink.com/K25EPjsf31dac8880a551fe2672247d21218bf854cbcf60/api/clusters?namespace=default'
    )
    expect(req.request.method).toEqual('POST')
    req.flush(testData)
    httpTestingController.verify()
  })

  it('can test HttpClient.put', () => {
    const testData: Data = { name: 'Test Data' }

    // Make an HTTP GET request
    service.put('clusters').subscribe((data) =>
      // When observable resolves, result should match test data
      expect(data).toEqual(testData)
    )

    const req = httpTestingController.expectOne(
      'https://mockapi.eolink.com/K25EPjsf31dac8880a551fe2672247d21218bf854cbcf60/api/clusters?namespace=default'
    )

    expect(req.request.method).toEqual('PUT')

    req.flush(testData)

    httpTestingController.verify()
  })

  it('can test HttpClient.put', () => {
    const testData: Data = { name: 'Test Data' }

    // Make an HTTP GET request
    service.put('clusters', { namespace: 'test' }).subscribe((data) =>
      // When observable resolves, result should match test data
      expect(data).toEqual(testData)
    )

    const req = httpTestingController.expectOne(
      'https://mockapi.eolink.com/K25EPjsf31dac8880a551fe2672247d21218bf854cbcf60/api/clusters?namespace=default'
    )
    expect(req.request.method).toEqual('PUT')
    req.flush(testData)
    httpTestingController.verify()
  })

  it('can test HttpClient.delete', () => {
    const testData: Data = { name: 'Test Data' }

    // Make an HTTP GET request
    service.delete('clusters').subscribe((data) =>
      // When observable resolves, result should match test data
      expect(data).toEqual(testData)
    )

    const req = httpTestingController.expectOne(
      'https://mockapi.eolink.com/K25EPjsf31dac8880a551fe2672247d21218bf854cbcf60/api/clusters?namespace=default'
    )

    expect(req.request.method).toEqual('DELETE')

    req.flush(testData)

    httpTestingController.verify()
  })

  it('can test HttpClient.delete with params', () => {
    const testData: Data = { name: 'Test Data' }

    // Make an HTTP GET request
    service.delete('clusters', { namespace: 'test' }).subscribe((data) =>
      // When observable resolves, result should match test data
      expect(data).toEqual(testData)
    )
    const req = httpTestingController.expectOne(
      'https://mockapi.eolink.com/K25EPjsf31dac8880a551fe2672247d21218bf854cbcf60/api/clusters?namespace=default'
    )
    expect(req.request.method).toEqual('DELETE')
    req.flush(testData)
    httpTestingController.verify()
  })

  it('handleError', () => {
    const err1 = new HttpErrorResponse({ error: { text: 'test' }, status: 0 })
    const messagService = TestBed.inject(EoNgFeedbackMessageService)
    const spyService = jest.spyOn(messagService, 'error')
    expect(spyService).not.toBeCalled()
    service.handleError(err1)
    expect(spyService).toBeCalled()

    const err2 = new HttpErrorResponse({ error: { text: 'test' }, status: 1 })
    service.handleError(err2)
    expect(spyService).toHaveBeenCalledTimes(2)
  })
})
