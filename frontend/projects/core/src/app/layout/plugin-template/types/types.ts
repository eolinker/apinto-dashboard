/*
 * @Date: 2023-12-12 18:57:19
 * @LastEditors: maggieyyy
 * @LastEditTime: 2024-01-24 14:59:28
 * @FilePath: \apinto\projects\core\src\app\layout\plugin-template\types\types.ts
 */

import { PluginTemplateConfigItem } from '../../../component/plugin-config-modal/plugin-config-modal.component'

export type PluginTemplateItem = {
    uuid:string
    name:string
    desc:string
    createTime:string
    updateTime:string
    operator:string
    isDelete:boolean
}

export type PluginTemplateData = {
    name:string
    desc:string
    plugins:PluginTemplateConfigItem[]
}
