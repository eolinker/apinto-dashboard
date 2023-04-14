import { Component, OnInit } from '@angular/core'
import { EoNgNavigationService } from '../../service/eo-ng-navigation.service'
import { ApiService } from '../../service/api.service'
import { Router } from '@angular/router'
import { GuideStepList, StepItem, TutorialItem, TutorialsList } from './conf'

@Component({
  selector: 'eo-ng-guide',
  templateUrl: './guide.component.html',
  styles: [
    `
    :host ::ng-deep{
      height: 100%;
      width: 100%;
      display: block;
      background-color: #f5f7fa;
      overflow: hidden;
      overflow-y:auto;
    }`
  ]
})
export class GuideComponent implements OnInit {
  stepList:Array<StepItem> = [...GuideStepList]
  tutorialsList:Array<TutorialItem> = [...TutorialsList]
  btnLoading:boolean = true

  constructor (private appConfigService:EoNgNavigationService, private api:ApiService, private router:Router) {}
  ngOnInit (): void {
    this.appConfigService.reqFlashBreadcrumb([])
    this.getStepStatus()
  }

  getStepStatus () {
    this.btnLoading = true
    this.api.get('system/quick_step').subscribe((resp:{code:number, msg:string, data:{cluster:boolean, upstream:boolean, api:boolean, publishApi:boolean}}) => {
      this.btnLoading = false
      if (resp.code === 0) {
        for (let i = 0; i < this.stepList.length; i++) {
          this.stepList[i].status = resp.data[this.stepList[i].name] ? 'done' : (i > 0 && this.stepList[i - 1].status === 'done' ? 'doing' : 'undo')
        }
      }
    })
  }

  goToStep (step:StepItem) {
    console.log(step)
    if (step.status === 'done') {
      this.router.navigate([step.doneUrl])
    } else {
      this.router.navigate([step.toDoUrl])
    }
  }

  goToGithub (url?:string) {
    window.open(`https://github.com/eolinker/apinto${url || ''}`)
  }
}
