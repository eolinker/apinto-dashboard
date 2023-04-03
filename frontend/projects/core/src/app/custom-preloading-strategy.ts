/* eslint-disable @typescript-eslint/no-unused-vars */
import { PreloadingStrategy, Route } from '@angular/router'
import { Observable, of, delay, mergeMap } from 'rxjs'

export class CustomPreloadingStrategy implements PreloadingStrategy {
  preload (route:Route, fn:()=> Observable<boolean>): Observable<boolean> {
    return of(true).pipe(
      delay(1000),
      mergeMap((_: boolean) => fn())
    )
  }
}
