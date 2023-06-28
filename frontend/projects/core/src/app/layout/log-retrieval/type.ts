
export type LogFileData = {
    file:string
    size:string
    mod:string
    key:string
}

export type LogOutputData = {
        name:string
        tail:string
        files:Array<LogFileData>
        active?:boolean // 折叠面板
}
