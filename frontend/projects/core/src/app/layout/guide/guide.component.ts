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
    }

    .right-border-box::after {
      content:' ';
      height:100%;
      color:black;
    }

`
  ]
})
export class GuideComponent implements OnInit {
  stepList:Array<StepItem> = [...GuideStepList]
  tutorialsList:Array<TutorialItem> = [...TutorialsList]
  btnLoading:boolean = true

  constructor (private navigationService:EoNgNavigationService, private api:ApiService, private router:Router) {}
  ngOnInit (): void {
    this.navigationService.reqFlashBreadcrumb([])
    this.getStepStatus()
  }

  getTitleStyle (step:StepItem) {
    return step.status === 'doing' ? `text-guide_${step.name} font-bold text-[16px] ` : (step.status === 'undo' ? 'text-DESC_TEXT  font-bold text-[16px]' : 'font-bold text-[16px]')
  }

  getDoingBtnClass (step:StepItem) {
    return `text-guide_${step.name} border-guide_${step.name}`
  }

  getStepStatus () {
    this.btnLoading = true

    this.api.get('system/quick_step').subscribe((resp:{code:number, msg:string, data:{cluster:boolean, upstream:boolean, api:boolean, publishApi:boolean}}) => {
      this.btnLoading = false
      if (resp.code === 0) {
        let doingFlag = false // 只有一个前往
        for (let i = 0; i < this.stepList.length - 1; i++) {
          this.stepList[i].status = resp.data[this.stepList[i].name as 'cluster'|'upstream'|'api'|'publishApi'] ? 'done' : ((i === 0 || (i > 0 && this.stepList[i - 1].status === 'done')) && !doingFlag ? 'doing' : 'undo')
          if (this.stepList[i].status === 'doing') {
            doingFlag = true
          }
        }
      }
    })
  }

  goToStep (step:StepItem) {
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
