import { HttpClient } from '@angular/common/http'
import { Inject, Injectable } from '@angular/core'
import { Observable, Subscriber } from 'rxjs'
import { API_URL } from './api.service'

@Injectable({
  providedIn: 'root'
})
export class IframeHttpService {
  constructor (private http:HttpClient,
    @Inject(API_URL) public urlPrefix:string) { }

  openIframe (url:string, option?:{headers?:Array<{name:string, value:string}>, initialize?:Array<{name:string, value:string, type:string}>}) {
    return new Observable((observer: Subscriber<any>) => {
      let objectUrl: string|null
      const header:{[k:string]:any} = {}
      if (option?.headers?.length) {
        for (const item of option.headers) {
          header[item.name] = item.value
        }
      }
      const init:{[k:string]:any} = {}

      if (option?.initialize?.length) {
        for (const item of option.initialize) {
          init[item.name] = item.value
        }
        this.http
          .post(url, { ...init }, { ...header })
          .subscribe((m:any) => {
            objectUrl = URL.createObjectURL(new Blob([m.blob()], { type: 'application/json' }))
            observer.next(objectUrl)
          })
      } else {
        this.http
          .get(url, { ...header })
          .subscribe((m:any) => {
            objectUrl = URL.createObjectURL(new Blob([m.blob()], { type: 'application/json' }))
            observer.next(objectUrl)
          })
      }
      return () => {
        if (objectUrl) {
          URL.revokeObjectURL(objectUrl)
          objectUrl = null
        }
      }
    })
  }
}
