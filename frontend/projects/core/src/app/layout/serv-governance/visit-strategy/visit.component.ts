import { Component } from '@angular/core'
import { Router } from '@angular/router'

@Component({
  selector: 'eo-ng-visit',
  template: `
  <router-outlet></router-outlet>
  `,
  styles: [
    `:host ::ng-deep{

      .form-row {
        flex-flow: row nowrap;
      }

        .limit-bg {
          background-color: var(--bar-background-color);
          padding: 20px;
          padding-bottom: 4px;
          border-radius: var(--border-radius);
          .label {
            width: 22%;
          }

          .ant-form-item {
            margin-bottom: 16px !important;
          }

          .ant-form-item-extra,
          .ant-form-item-explain-error {
            margin-left: var(--LAYOUT_PADDING);
          }
          .arrayItem .ant-table-tbody tr {
            height: 44px !important;
            padding-bottom: 12px !important;
            td {
              vertical-align: top;
              border: none;

              padding-bottom: 12px !important;
            }
          }

          .arrayItem .ant-table-tbody tr:last-child {
            height: 32px !important;
            padding-bottom: 0px !important;
            td {
              padding-bottom: 0px !important;
            }
          }
        }

        .ant-input-affix-wrapper,
        .ant-input:not(textarea) {
          // width: 346px;
          height: 32px;
        }
        .ant-input-affix-wrapper input.ant-input:not(textarea) {
          height: 22px !important;
          border: none;
        }

        .ant-input-affix-wrapper:has(input[disabled]) {
          color: var(--disabled-text-color);
          background-color: var(--disabled-background-color);
          cursor: not-allowed;
          opacity: 1;
        }

        .ant-input-affix-wrapper:has(input[disabled]):hover {
          border: 1px solid var(--border-color);
        }
    }`
  ]
})
export class VisitComponent {
  constructor (private router:Router) {
    if (this.router.url === '/serv-governance/visit') {
      this.router.navigate(['/', 'serv-governance', 'visit', 'group'])
    }
  }

  ngDoCheck () {
    if (this.router.url === '/serv-governance/visit') {
      this.router.navigate(['/', 'serv-governance', 'visit', 'group'])
    }
  }
}
