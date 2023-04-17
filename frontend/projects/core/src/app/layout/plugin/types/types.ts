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
}

export type PluginInstallData = {
    module:{
        name:string
        server:string
        header:Array<PluginInstallConfigData>
        query:Array<PluginInstallConfigData>
        initialize:Array<PluginInstallConfigData>
    }
    render:{
        internet:boolean
        invisible:boolean
        headers:Array<PluginInstallConfigData>
        querys:Array<PluginInstallConfigData>
        initialize:Array<PluginInstallConfigData>
        nameConflict:boolean
    }
}
