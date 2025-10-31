import { Injectable, NgZone, Optional, PlatformRef, SkipSelf } from '@angular/core'

@Injectable({
  providedIn: 'root'
})
export class PlatformProviderService {
  private ngZoneInstance: NgZone | undefined;

  constructor (@Optional() @SkipSelf() private platformRef: PlatformRef,
  @Optional() @SkipSelf() private parentNgZone: NgZone) {}

  getPlatformRef (): PlatformRef {
    return this.platformRef
  }

  getNgZone (): NgZone {
    if (this.ngZoneInstance) {
      return this.ngZoneInstance
    } else if (this.parentNgZone) {
      return this.parentNgZone
    } else {
      // 如果没有可用的 NgZone，可以考虑抛出一个错误，以便调试
      throw new Error('NgZone is not available.')
    }
  }

  setNgZone (ngZone: NgZone) {
    this.ngZoneInstance = ngZone
  }
}
