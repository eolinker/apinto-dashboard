/*
 * @Date: 2023-11-20 10:07:17
 * @LastEditors: maggieyyy
 * @LastEditTime: 2023-11-21 11:10:31
 * @FilePath: \apinto\projects\core\src\app\service\plugin-slot-hub.service.ts
 */
import { Injectable } from '@angular/core'

@Injectable({
  providedIn: 'root'
})
export class PluginSlotHubService {
  private slotMap:Map<string, any> = new Map()
  addSlot = (name:string, content:any) => {
    this.slotMap.set(name, content)
  }

  addSlotArr = (name:string, content:any[]) => {
    this.slotMap.get(name) ? this.slotMap.set(name, this.slotMap.get(name).push(content)) : this.slotMap.set(name, content)
  }

  removeSlot = (name:string) => {
    this.slotMap.delete(name)
  }

  getSlot = (name:string) => {
    return this.slotMap.get(name)
  }
}
