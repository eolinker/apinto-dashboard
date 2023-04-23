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

  openIframe (url:string, option:{headers:{[k:string]:string}}) {
    return new Observable((observer: Subscriber<any>) => {
      let objectUrl: string|null
      console.log(`/${url}`)
      this.http
        .get('http://localhost:4200/plugin/group/list', { ...option })
        .subscribe((m:any) => {
          console.log(m)
          objectUrl = URL.createObjectURL(m.blob())
          observer.next(objectUrl)
        })

      return () => {
        if (objectUrl) {
          URL.revokeObjectURL(objectUrl)
          objectUrl = null
        }
      }
    })
  }
}
