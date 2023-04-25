import { Component, Input, OnInit, ViewChild } from '@angular/core'
import { CustomReactComponentWrapperComponent } from '../../formily2-react/CustomReactComponentWrapper'
import { SelectOption } from 'eo-ng-select'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { ApiService } from '../../../service/api.service'
import { IntelligentPluginService } from '../intelligent-plugin.service'

@Component({
  selector: 'eo-ng-intelligent-plugin-create',
  template: `
       <formily2-react-wrapper
        #formily
        [renderSchema]="renderSchema"
        [editPage]="editPage"
        [initFormValue]="initFormValue"
        [driverSelectOptions]="driverSelectOptions"
       ></formily2-react-wrapper>
  `,
  styles: [
  ]
})
export class IntelligentPluginCreateComponent implements OnInit {
  @ViewChild('formily', { static: true }) formily!: CustomReactComponentWrapperComponent
  @Input() renderSchema: { [k: string]: any } = {}
  @Input() editPage: boolean = false
  @Input() initFormValue: { [k: string]: any } = {}
  @Input() driverSelectOptions: SelectOption[] = []

  form:{[k:string]:any} = {}
  moduleName:string = ''
  uuid:string = ''

  constructor (
    private message: EoNgFeedbackMessageService,
    private service:IntelligentPluginService,
    private modalService:EoNgFeedbackModalService,
    private api:ApiService) {

  }

  ngOnInit (): void {
    this.editPage && this.getMessage()
  }

  ngAfterViewInit () {
    this.form = this.formily.reactComponent.current.form
  }

  getMessage () {
    this.api.get(`/dynamic/${this.moduleName}/info/${this.uuid}`)
      .subscribe((resp:{code:number, msg:string, data:{[k:string]:any}}) => {
        if (resp.code === 0) {
          this.initFormValue = resp.data
        }
      })
  }
}
