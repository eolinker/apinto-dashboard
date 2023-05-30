import { Component, OnInit } from '@angular/core'
import { IframeHttpService } from '../../service/iframe-http.service'
import { ApiService } from '../../service/api.service'
import { NavigationEnd, Router } from '@angular/router'
import { BaseInfoService } from '../../service/base-info.service'
import { Subscription } from 'rxjs'
import { environment } from 'projects/core/src/environments/environment'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { EoIframeComponent } from '../../component/iframe/iframe.component'

@Component({
  selector: 'eo-ng-local-plugin',
  templateUrl: '../../component/iframe/iframe.component.html',
  styles: [
    `
    :host{
      display:block;
      height:100%;
      overflow-y:hidden;
    }
    :host ::ng-deep{
      nz-spin.iframe-spin,
      nz-spin.iframe-spin >.ant-spin-container,
      #iframePanel,
      #iframePanel > iframe{
        width:100%;
        height:100%;
        border:none;
      }
    }`
  ]
})
export class LocalPluginComponent extends EoIframeComponent {}
