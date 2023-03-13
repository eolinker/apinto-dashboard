import { Injectable } from '@angular/core'
import { Subject } from 'rxjs'

@Injectable({
  providedIn: 'root'
})
export class DataService {
  private flashFlag: Subject<boolean> = new Subject<boolean>()
  private flashMenu: Subject<boolean> = new Subject<boolean>()

  reqFlashGroup () {
    this.flashFlag.next(true)
  }

  repFlashGroup () {
    return this.flashFlag.asObservable()
  }

  reqFlashMenu () {
    this.flashMenu.next(true)
  }

  repFlashMenu () {
    return this.flashMenu.asObservable()
  }
}
