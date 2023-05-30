import { Pipe, PipeTransform, SecurityContext } from '@angular/core'
import { DomSanitizer, SafeHtml, SafeStyle, SafeScript, SafeUrl, SafeResourceUrl } from '@angular/platform-browser'

@Pipe({
  name: 'safe'
})
export class SafePipe implements PipeTransform {
  constructor (protected sanitizer: DomSanitizer) {}

  transform (value: string, type: string = 'url'): SafeHtml | SafeStyle | SafeScript | SafeUrl | SafeResourceUrl {
    switch (type) {
      case 'html':
        return this.sanitizer.bypassSecurityTrustHtml(value)
      case 'style':
        return this.sanitizer.bypassSecurityTrustStyle(value)
      case 'script':
        return this.sanitizer.bypassSecurityTrustScript(value)
      case 'url':
        return this.sanitizer.bypassSecurityTrustUrl(value)
      case 'resourceUrl':
        return this.sanitizer.bypassSecurityTrustResourceUrl(value)
        // this.sanitizer.bypassSecurityTrustResourceUrl(this.sanitizer.sanitize(SecurityContext.URL, value)) // TODO，待确认，防止XSS
      default:

        throw new Error(`Unable to bypass security for invalid type: ${type}`)
    }
  }
}
