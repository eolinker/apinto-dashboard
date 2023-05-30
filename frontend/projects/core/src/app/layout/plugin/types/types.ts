export type PluginItem = {
    id:string
    name:string
    cname:string
    resume:string
    icon:string
    enable:boolean
    isInner:boolean
}

export type PluginGroupItem = {
    uuid:string
    name:string
    count:number
}

export type PluginMessage = {
    id:string
    name:string
    cname:string
    resume:string
    icon:string
    enable:boolean
    uninstall:boolean
    canDisable:boolean
}

export type PluginInstallConfigData = {
    name:string
    value:string
    desc:string
    title:string
    placeholder:string
    type?:string
}

export type PluginModuleData = {
    name:string
    server:string
    header:Array<PluginInstallConfigData>
    query:Array<PluginInstallConfigData>
    initialize:Array<PluginInstallConfigData>
}

export type PluginInstallData = {
    module:PluginModuleData,
    render:{
        internet:boolean
        invisible:boolean
        headers:Array<PluginInstallConfigData>
        querys:Array<PluginInstallConfigData>
        initialize:Array<PluginInstallConfigData>
        nameConflict:boolean
    }
}
