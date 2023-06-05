import { Component, Input, OnInit, ViewChild } from '@angular/core'
import { CustomReactComponentWrapperComponent } from '../../formily2-react/CustomReactComponentWrapper'
import { SelectOption } from 'eo-ng-select'
import { EoNgFeedbackMessageService, EoNgFeedbackModalService } from 'eo-ng-feedback'
import { ApiService } from '../../../service/api.service'
import { EoIntelligentPluginService } from '../intelligent-plugin.service'

@Component({
  selector: 'eo-ng-intelligent-plugin-create',
  templateUrl: './create.component.html',
  styles: [
  ]
})
export class EoIntelligentPluginCreateComponent implements OnInit {
  @ViewChild('formily', { static: true }) formily!: CustomReactComponentWrapperComponent
  @Input() renderSchema: { [k: string]: any } = {}
  @Input() editPage: boolean = false
  @Input() initFormValue: { [k: string]: any } = {}
  @Input() driverSelectOptions: SelectOption[] = []

  form:any|undefined
  moduleName:string = ''
  uuid:string = ''

  constructor (
    private message: EoNgFeedbackMessageService,
    private service:EoIntelligentPluginService,
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
    this.api.get(`dynamic/${this.moduleName}/info/${this.uuid}`)
      .subscribe((resp:{code:number, msg:string, data:{[k:string]:any}}) => {
        if (resp.code === 0) {
          const res:{[k:string]:any} = {}
          for (const key of Object.keys(resp.data)) {
            if (resp.data[key] !== '') {
              res[key] = resp.data[key]
            }
          }
          this.form.setValues(res)
        }
      })
  }
}
