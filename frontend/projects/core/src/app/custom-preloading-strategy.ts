import { PreloadingStrategy, Route } from '@angular/router'
import { Observable, of, delay, mergeMap } from 'rxjs'

export class CustomPreloadingStrategy implements PreloadingStrategy {
  preload (route:Route, fn:()=> Observable<boolean>): Observable<boolean> {
    return of(true).pipe(
      delay(1000),
      // eslint-disable-next-line @typescript-eslint/no-unused-vars
      mergeMap((_: boolean) => fn())
    )
  }
}
