/* eslint-disable no-useless-constructor */
import { Injectable } from '@angular/core'
import { Observable } from 'rxjs'
import { MonitorAlarmChannelsData } from '../layout/monitor-alarm/types/types'
import { ApiService } from './api.service'

@Injectable({
  providedIn: 'root'
})
export class EoNgMonitorNoticeService {
  noticeChannelsList:MonitorAlarmChannelsData[] = []
  constructor (public api:ApiService) { }

  // 通过接口获取最新的通知渠道列表
  getLastedNoticeChannels ():Observable<MonitorAlarmChannelsData[]> {
    return new Observable(observer => {
      this.api.get('warn/channels').subscribe((resp:{code:number, data:{channels:MonitorAlarmChannelsData[]}, msg:string}) => {
        if (resp.code === 0) {
          this.noticeChannelsList = resp.data.channels || []
          observer.next(this.noticeChannelsList)
        }
      })
    })
  }

  // 返回最新渠道列表数据
  getNoticeChannels ():MonitorAlarmChannelsData[] {
    return this.noticeChannelsList
  }
}
