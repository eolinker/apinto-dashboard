// 插件管理列表
export type PluginItem = {
    name:string
    extended:string
    desc:string
    updateTime:string
    operator:string
    isDelete:boolean
    isBuilt:boolean
}

// 插件管理编辑
export type PluginData = {
    name:string
    rely:string
    extended:string
    desc:string
}
