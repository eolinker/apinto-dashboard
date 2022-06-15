let RootId = "FormRender"

let validate = {
    _validator: null,
    djv: function () {
        if (this._validator){
            return this._validator
        }
        this._validator = new djv()
        return this._validator
    },
}
function requiredMap(required){
    if (!required){
        return {}
    }
    let r = {}
    for (let i in required){
        r[required[i]] = true
    }
    return r
}
function readId(path){
    return path.replaceAll(".","_")
}
function readGenerator(options){
    let fn = options["generator"]
    if (!fn){
        return BaseGenerator
    }
    return  fn
}
function CheckBySchema(id, schema, value) {
    let env = validate.djv()
    if (!env.resolved.hasOwnProperty(id)) {
        env.addSchema(id, schema);
    }
    return  env.validate(id, value);
    // let validator = validate.Default()
    // if (!validator.resolved.hasOwnProperty(id)) {
    //     validator.addSchema(id, schema);
    // }
    // let err = validator.validate(id, value)
    // if (err) {
    //     console.log(err)
    //     return false
    // }
    // return true
}

function valueForType(t, v) {
    switch (t) {
        case "integer":
        case "number": {
            return Number(v)
        }
        case "boolean": {
            if (typeof v === "undefined") {
                return false
            }
            switch (String(v).toLowerCase()) {
                case "on":
                case "true": {
                    return true
                }
                default: {
                    return false
                }
            }
        }
    }

    return v
}
function createEnum(schema, id,required, appendAttr) {
    let readOnly = ""
    this.schema = schema
    if (schema["readonly"] === true) {
        readOnly = "readonly"
    }
    if (typeof appendAttr === "undefined"){
        appendAttr = ""
    }
    if (schema["enum"]) {
        let enums = schema["enum"]
        let require = ""
        if (required) {
            require = `required`
        }
        let $select = $(`<select ${readOnly} ${appendAttr} class="form-control form-control-sm" id="${id}" ${require}></select>`)

        for (let i in enums) {
            $select.append(`<option>${enums[i]}</option>`)
            // if (typeof value !== "undefined" && value === enums[i]) {
            //     $select.append(`<option>${enums[i]}</option>`)
            // } else {
            //     $select.append(`<option selected>${enums[i]}</option>`)
            // }
        }

        return $select
    }
}

function createInput(schema, id,required, appendAttr) {
    let readOnly = ""

    if (schema["readonly"] === true) {
        readOnly = "readonly"
    }
    let idstr = ""
    if (id && id.length>0){
        idstr = `id="${id}"`
    }
    let input = `<input ${readOnly} class="form-control form-control-sm" ${idstr} aria-describedby="validation_${id}" `;
    if (appendAttr) {
        input += appendAttr
    }

    function readFormatForString(format) {

        switch (format) {

            case "email", "password", "date", "time", "number": {
                return format
            }
            case "idn-email": {
                return "email"
            }
            case "date-time": {
                return 'datetime-local'
            }
            case "boolean": {
                return "checkbox"
            }
            default: {
                return "text"
            }
        }
    }

    switch (schema["eo:type"]) {
        case "string": {
            input += ' type="' + readFormatForString(schema["format"]) + '"'
            if (schema["format"] === 'password') {
                input += ' autocomplete="on"'
            }
            break
        }
        case "integer":
        case "number": {
            input += ' type="number"'
            input += ' value="0"'
            break
        }
        default: {
            throw `now allow [${schema["eo:type"]} for input]`
        }
    }
    if (schema["maxLength"]) {
        input += ' maxlength="' + schema["maxLength"] + '"'
    }

    if (schema["minLength"]) {
        input += ' minLength="' + schema["minLength"] + '"'
    }

    if (required) {
        input += ' required'
    }
    input += '/>'
    return input
}

function getLabel(name,schema) {
    let label = schema["label"]
    if (!label || label.trim() === "") {
        label = name
    }

    label = label.replace(label[0], label[0].toUpperCase());
    label = label.replaceAll("_"," ")
    return label
}

function createLabel(name,schema, id,required, appendAttr) {
    if (!appendAttr) {
        appendAttr = ""
    }
    if (typeof appendAttr === "undefined") {
        appendAttr = ""
    }
    let idFor = ""
    if (id) {
        idFor = ` for="${id}"`
    }
    let require = ""
    if (required === true) {
        require = '<span style="color: red">*</span>'
    }
    return `<label class=" col-form-label  text-nowrap" ${idFor} ${appendAttr}>${require}${getLabel(name,schema)}</label>`
}
class BaseChangeHandler {
    constructor(path) {
        this.Id = readId(path)
    }
    onChange(fn){

        if (fn){
            if(!this.ChangeHandler){
                this.ChangeHandler = new Array()
            }
            this.ChangeHandler.push(fn)
        }else{
            console.log("change:",this.Id)
            for(let i in this.ChangeHandler){
                this.ChangeHandler[i].apply(this)
            }
        }
    }
}
class BaseValue {
    constructor(schema, target,path) {
        this.Id = readId(path)
        this.Schema = schema
        this.Target = target
        if (typeof schema["default"] !== "undefined") {
            this.Value = schema["default"]
        }
        let JsonSchema = new SchemaHandler(schema).JsonSchema
        this.InputValid(JsonSchema,target)
    }
    onChange(fn){
        let o = this
        $(this.Target).on("change",function (){
            fn.apply(o)
        })
    }
    isOk(v){
        return
    }
    ValidHandler(v) {
        let id =  this.Id
        console.debug("ValidHandler:",id,"=",v)
        let value = v
        value = valueForType(this.Schema["eo:type"], value)
        let rs = this.isOk(v)
        if (typeof rs === "undefined") {
            rs = CheckBySchema(id, this.Schema, value)
        }
        if (typeof rs === "undefined") {
            $(this.Target).removeClass("is-invalid")
            $(this.Target).addClass("is-valid")
        } else {
            $(this.Target).removeClass("is-valid")
            $(this.Target).addClass("is-invalid")
        }
    }
    InputValid(schema, target) {
        let o = this
        $(target).on("change", function () {
            o.ValidHandler($(this).val(), schema)
        })
    }

    get Value() {
        return valueForType(this.Schema["eo:type"], $(this.Target).val())
    }

    set Value(v) {

        if ( typeof v === "undefined"){
            let schema = this.Schema
            if (typeof schema["default"] !== "undefined") {
               v = schema["default"]
            }
        }
        switch($(this.Target).get(0).tagName ){
            case "select":{
                if ($(this.Target).prop("multiple")){
                    for (let i in v){
                        let item = v[i]
                        $(this.Target).find(`option[value="${item}"]`).prop("selected", true);
                    }
                }else{
                    $(this.Target).find(`option[value="${v}"]`).prop("selected", true);
                }
                break
            }
            default:{
                $(this.Target).val(v)
            }
        }
    }

}

class BaseEnumRender extends BaseValue {
    constructor(options) {
        let schema = options["schema"]
        let path = options["path"]
        let panel = options["panel"]
        super(schema, $(createEnum(schema, readId(path))),path)

        $(panel).append(this.Target)

        // this.Target.on("click",function (){
        //     $(this).trigger("change")
        // })
    }

}

class SwitchRender extends BaseValue {
    constructor(options) {
        let schema = options["schema"]
        let path = options["path"]
        let panel = options["panel"]

        super(schema, $(`<input id="${readId(path)}" type="checkbox" class="form-control-sm" data-toggle="toggle" data-size="sm"/>`),path)

        $(panel).append(this.Target)
        this.Target.bootstrapToggle({
            on: '开启',
            off: '关闭'
        })
        this.Value = false
    }

    get Value() {
        return $(this.Target).prop("checked")
    }

    set Value(v) {
        if (v === true || v === "true") {
            this.Target.bootstrapToggle("on");
        } else {
            this.Target.bootstrapToggle("off");
        }
    }
}
class PopPanel  {
    constructor(options,callbackFn,v) {

        let $Panel = $(`<div class="pop_window pop_window_small p-3" id="detail_container">
    <div class="pop_window_header">
        <span class="pop_window_title">${options["title"]}</span>
        <button class="pop_window_button btn btn_default close" >关闭</button>
        <br>
    </div>
   
</div>`)

        let $Fade = $("<div class='modal-backdrop fade show modal-open'></div>")
        let $Body = $(` <div class="pop_window_body"></div>`)
        $("body").append($Fade)
        $("body").append($Panel)
        $Panel.append($Body)

        let $Value = readGenerator(options)({
            schema:options["schema"],
            path:options["path"],
            panel:$Body,
            generator:readGenerator(options)
        })
        if (v){
            $Value.Value = v
        }
         $Panel.append(`<div class="row justify-content-between">
                <div class="col-4">
                    <button type="button" class="btn btn-outline-secondary form_cancel">取消</button>
                </div>
                <div class="col-4" style="text-align: right">
                    <button type="button" class="btn btn-primary form_submit">提交</button>
                </div>
            </div>`)
        let close = function (){
            $Fade.remove()
            $Panel.remove()
        }
        $Panel.show()
        $Fade.click(function (){
            close()

        })
        $Panel.on("click","button.close",close)
        $Panel.on("click","button.form_cancel",close)
        $Panel.on("click","button.form_submit",function (){
            callbackFn($Value.Value)
            close()
        })
    }
}
function PopPanelMap (options,callbackFn,keyHas,v) {

        let $Panel = $(`<div class="pop_window pop_window_small p-3" id="detail_container">
    <div class="pop_window_header">
        <span class="pop_window_title">${options["title"]}</span>
        <button class="pop_window_button btn btn_default close" >关闭</button>
  
    </div>
   
</div>`)

        let $Fade = $("<div class='modal-backdrop fade show modal-open'></div>")
        let $Body = $(` <div class="pop_window_body"></div>`)
        $("body").append($Fade)
        $("body").append($Panel)
        $Panel.append($Body)

        let $KeyInput = new FieldPanel("key",{
            schema:{"type":"string","eo:type":"string","pattern":"[a-zA-Z0-9]+[a-zA-Z0-9_]*"},
            required:true,
            path:`${options["path"]}.key`,
            panel:$Body
        })
        let $Value = readGenerator(options)({
            schema:options["schema"],
            path:options["path"],
            panel:$Body,
            generator:readGenerator(options)
        })
        if (v){
            $KeyInput.Value = v["__key"]
            $Value.Value = v
            $KeyInput.$Value.Target.attr("readonly", true)

        }
        $Panel.append(`<div class="row justify-content-between">
                <div class="col-4">
                    <button type="button" class="btn btn-outline-secondary form_cancel">取消</button>
                </div>
                <div class="col-4" style="text-align: right">
                    <button type="button" class="btn btn-primary form_submit">提交</button>
                </div>
            </div>`)
        let close = function (){
            $Fade.remove()
            $Panel.remove()
        }
        $Panel.show()
        $Fade.click(function (){
            close()
        })
        $Panel.on("click","button.close",close)
        $Panel.on("click","button.form_cancel",close)
        $Panel.on("click","button.form_submit",function (){
            let key = $KeyInput.Value
            if (key.length<1){
                $KeyInput.$Value.Target.addClass("is-invalid")
                return
            }
            let v = $Value.Value
            v["__key"]= key
            callbackFn(v)
            close()
        })
    }

class RequireRender extends BaseValue {
    constructor(options) {
        let schema = options["schema"]
        let path = options["path"]
        let panel = options["panel"]

        super(schema, $(`<select id=${readId(path)} class="form-control form-control-sm"></select>`),path)
        this.MOptions = options
        $(panel).append(this.Target)

        let select = this.Target
        dashboard.searchSkill(ModuleName(),schema["skill"],function (res){

            $(select).empty()
            if(options["required"]){
                $(select).append(`<option value="">请选择</option>`)
            }else{
                if (schema["empty_label"]){
                    $(select).append(`<option value="">${schema["empty_label"]}</option>`)
                }else {
                    $(select).append(`<option value="">不启用</option>`)
                }
            }
            for (let i in res.data){
                let d = res.data[i]
                $(select).append(`<option value="${d.id}">${d.id}[${d.driver}]</option>`)
            }
        })
    }
    isOk(v) {
        if (!v || v.length === 0){
            if (this.MOptions["required"]){
                return "请选择"
            }
        }
    }
}
class RequireArrayRender extends BaseValue {
    constructor(options) {
        let schema = options["schema"]
        let path = options["path"]
        let panel = options["panel"]

        super(schema, $(`<select id=${readId(path)} data-live-search="true" class="selectpicker form-control form-control-sm" title="请选择" data-width="100%" data-size="5" data-selected-text-format="count>4" multiple></select>`),path)
        this.MOptions = options
        $(panel).append(this.Target)

        let $Select = this.Target
        dashboard.searchSkill(ModuleName(),schema["skill"],function (res){

            $($Select).empty()
            // if(options["required"]){
            //     $($Select).append(`<option value="">请选择</option>`)
            // }else{
            //     if (schema["empty_label"]){
            //         $($Select).append(`<option value="">${schema["empty_label"]}</option>`)
            //     }else {
            //         $($Select).append(`<option value="">不启用</option>`)
            //     }
            // }
            for (let i in res.data){
                let d = res.data[i]
                $($Select).append(`<option value="${d.id}">${d.id}[${d.driver}]</option>`)
            }
            $Select.selectpicker()
        })
    }
    set Value(vs){
        this.Target.selectpicker('val', vs);
    }
    get Value(){
        return this.Target.selectpicker('val');
    }
    isOk(v) {
        if (!v || v.length === 0){
            if (this.MOptions["required"]){
                return "请选择"
            }
        }
    }
}

class BaseInputRender extends BaseValue {
    constructor(options) {
        let schema = options["schema"]
        let path = options["path"]
        let panel = options["panel"]
        super(options["schema"], $(createInput(schema, readId(path),options["required"])),path)
        $(panel).append(this.Target)
        if (schema["description"] && schema["description"].length >0){
            $(panel).append(`<small id="help:${path}" class="text-muted">${schema["description"]}</small>`)
        }
    }
}
//字段封装
class FieldPanel {
    constructor(name,options) {
        let panel=options["panel"]
        let schema=options["schema"]
        this.Id = readId(options["path"])
        this.$Panel = $(`<div class=""></div>`)
        panel.append(this.$Panel)
        let valuePanel = $(`<div class=""></div>`)
        this.$Panel.append(`<div class="">${createLabel(name,schema,this.Id,options["required"])}</div>`)
        this.$Panel.append(valuePanel)
        this.$Value = readGenerator(options)({panel:valuePanel,schema:schema,parent:options["parent"],generator:readGenerator(options),path:options["path"],name:name,required:options["required"]})

        switch (schema["type"]){
            case "array":{
                if (schema["items"]["enum"]){
                    break
                }
            }
            case  "object": {
                valuePanel.addClass("border-top border-bottom px-3 py-1")
            }
        }
    }
    onChange(fn){
        console.log(this.Id)
        let o = this
        this.$Value.onChange(function (){
            fn.apply(o)
        })
    }
    Show(){

        $(this.$Panel).show()
    }
    Hide(){
        $(this.$Panel).hide()
    }
    set Value(v){
        this.$Value.Value = v
    }
    get Value(){
        if ($(this.$Panel).is(":visible")){
            return this.$Value.Value
        }else{
            return {}
        }

    }
}

// 简单map
class SimpleMapRender extends BaseChangeHandler{
    constructor(options) {
        super(options["path"]);
        this.MOptions = options
        let panel=options["panel"]
        this.Schema = options["schema"]
        this.Items = this.Schema["additionalProperties"]
        let myPanel = $(`<div class="pt-1"></div>`)
        this.$itemPanel = $(`<div id="${this.Id}_maps" class="pt-1" role="toolbar"></div>`)
        panel.append(myPanel)
        myPanel.append(this.$itemPanel)

        this.ItemCount = 0
        this.Values = {}
        myPanel.append(`<div class="pt-1"> <button id="${this.Id}_add" type="button" class="btn btn-success btn-sm">+</button></div>`)
        let o = this
        myPanel.on("click",`#${this.Id}_add`,function (){
            let isEmpty = false

            for (let k in o.Values){
                let it = o.Values[k];
                if(it.key.Value.length<1){
                    isEmpty = true
                    $(it.key.Target)[0].focus()
                    break
                }
            }

            if (!isEmpty){
                o.add().key.Target[0].focus()
            }
        })
    }
    add(k,v) {
        let index = this.ItemCount++
        let $item = $(`<div class="input-group input-group-sm mt-2" data-index="${index}">
                        <div class="input-group-prepend">
                            <span class="input-group-text">key</span>
                        </div>
                        </div>`)
        this.$itemPanel.append($item)
        let $keyInput = $(createInput({"type": "string", "eo:type": "string"},`${this.Id}.item_key_${index}`,true))
        $item.append($keyInput)
        $item.append(`<div class="input-group-prepend"><span class="input-group-text">value</span></div>`)

        let keyInput = new BaseValue({"type": "string", "eo:type": "string"},$keyInput,`${this.Id}.item_key_${index}`)
        let $valueInput = $(createInput(this.Items,`${this.Id}.item_value_${index}`,true))
        $item.append($valueInput)
        let valueInput = new BaseValue(this.Items,$valueInput,`${this.Id}.item_value_${index}`)
        $item.append(`<div class="input-group-prepend"><button type="button" class="btn btn-danger" > - </button></div>`)
        let it = {key: keyInput, value: valueInput}
        this.Values[index] = it
        let o = this
        $($item).on("click", "button", function () {
            $item.remove()
            delete o.Values[index]
        })
        if (k){
            keyInput.Value = k
        }
        if (v){
            valueInput.Value = v
        }
        return it
    }
    get Value(){
        let v={}
        for (let i in this.Values){
            let it = this.Values[i]
            let key = it.key.Value
            if (key.length>0){
                v[key] = it.value.Value
            }
        }
        return v
    }
    set Value(v){
        this.$itemPanel.empty()
        this.Values={}
        for (let k in v){
            this.add(k,v[k])
        }
    }
}
class ObjectArrayRender extends BaseChangeHandler{
    constructor(options) {
        super(options["path"])
        let panel=options["panel"]
        let schema=options["schema"]

        let Id = this.Id;
        let items = schema["items"]
        let p = $(panel)
        let $btn = $(`
<div id="${Id}_toolbar" class="m-1">
<button id="${Id}_AddButton" type="button" class="btn btn-secondary btn-sm">Add</button>
</div>`)
        p.append($btn)
        $btn.on("click","button", function (event) {
            new PopPanel({
                    path:`${options["path"]}.items`,
                    schema:items,
                    name:options["name"],
                    title:`添加 ${getLabel(options["name"],options["schema"])}`
                },
                function (v){
                    $Table.bootstrapTable("append",[v])
                    $Table.bootstrapTable('scrollTo', 'bottom')
                    O.onChange()
                })
            return false
        })

        let $Table = $(`<table  id="${Id}_items"></table>`)
        p.append($Table)

        this.Table = $Table
        let properties = items["properties"]
        let uiSort = items["ui:sort"]

        function formatterKV(v) {
            let html = ""
            for (let k in v) {
                html += "<span class='px-1 border btn-sm  btn-outline-secondary'>" + k + "=" + v[k] + "</span>\n"
            }
            html += ""
            return html
        }

        let columns = []
        columns.push({
            title: "index",
            field: "__index",
            sortable: false,

            formatter: function (v, row, index) {
                return index+1
            }
        })
        let O = this
        for (let i in uiSort) {
            let name = uiSort[i]
            let item = properties[name]
            switch (item["eo:type"]) {
                case "object": {
                    columns.push({
                        title: getLabel(name,item),
                        field: item["name"],
                        sortable: false,

                        formatter: JSON.stringify,

                    })
                    break
                }
                case "map": {
                    columns.push({
                        title: getLabel(name,item),
                        field: name,
                        sortable: false,
                        editable: false,
                        formatter: formatterKV,
                    })
                    break
                }
                case "array": {
                    columns.push({
                        title: getLabel(name,item),
                        field: name,
                        sortable: false,
                        editable: false,
                        formatter: formatterKV,
                    })
                    break
                }
                default: {

                    columns.push({
                        title: getLabel(name,item),
                        field: name,
                        sortable: true,
                    })

                    break
                }
            }
        }
        columns.push({
            title: "操作",
            field: "",
            sortable: false,
            editable: false,
            formatter: function (v, row, index) {
                return `<a class="edit" href="javascript:void(0)" array-row="${index}" title="edit">配置</a> <a class="remove" href="javascript:void(0)" array-row="${index}" title="remove">删除</a>`
            },
            events:{
                "click .remove":function (e,value,row,index) {
                    // let rowIndex = $(this).attr("array-row")
                    $Table.bootstrapTable('remove', {
                        field: '$index',
                        values: [index]
                    })
                    O.onChange()
                },
                "click .edit":function (e,value,row,index){
                    new PopPanel({
                            path:`${options["path"]}.items`,
                            schema:items,
                            name:options["name"],
                            title:`修改 ${getLabel(options["name"],options["schema"])}:${index+1}`
                        },
                        function (v){
                            $Table.bootstrapTable("updateRow",{
                                index:index,
                                row:v
                            })
                            O.onChange()
                        },row)
                }
            }
        })

        let tableOptions = {
            columns: columns,
            width: "100%",
            useRowAttrFunc: true,
            reorderableRows:true,
            onReorderRow:function (data){
                $Table.bootstrapTable("refresh")
            }

        }
        $Table.bootstrapTable(tableOptions);
        this.Value = []
    }

    set Value(v) {
        if (!v || typeof v == "undefined") {
            v = []
        }
        this.Table.bootstrapTable("load", v)
    }

    get Value() {
        return this.Table.bootstrapTable('getData')
    }

}
// 结构体map
class ObjectMapRender extends BaseChangeHandler{
    constructor(options) {
        super(options["path"])
        let panel=options["panel"]
        let schema=options["schema"]

        let Id = this.Id;
        let items = schema["additionalProperties"]

        let $btn = $(`
<div id="${Id}_toolbar" class="m-1">
<button id="${Id}_AddButton" type="button" class="btn btn-secondary btn-sm">Add</button>
</div>`)
        $(panel).append($btn)
        $btn.on("click","button", function (event) {
             PopPanelMap({
                   path:`${options["path"]}.items`,
                   schema:items,
                   name:options["name"],
                   title:`添加 ${getLabel(options["name"],options["schema"])}`
               },
               function (v){
                   $Table.bootstrapTable("append",[v])
                   $Table.bootstrapTable('scrollTo', 'bottom')
                   O.onChange()
               },function (k){
                    let vs = O.Value;
                    return  vs.hasOwnProperty(k)
               })
            return false
        })

        let $Table = $(`<table  id="${Id}_items"></table>`)
        $(panel).append($Table)

        this.Table = $Table
        let properties = items["properties"]
        let uiSort = items["ui:sort"]
        
        function formatterKV(v) {
            let html = ""
            for (let k in v) {
                html += "<span class='px-1 border btn-sm  btn-outline-secondary'>" + k + "=" + v[k] + "</span>\n"
            }
            html += ""
            return html
        }
        let columns = []
        columns.push({
            title: "key",
            field: "__key",
            sortable: false,
        })
        let O = this
        for (let i in uiSort) {
            let name = uiSort[i]
            let item = properties[name]
            switch (item["eo:type"]) {
                case "object": {
                    columns.push({
                        title: getLabel(name,item),
                        field: item["name"],
                        sortable: false,

                        formatter: JSON.stringify,

                    })
                    break
                }
                case "map": {
                    columns.push({
                        title: getLabel(name,item),
                        field: name,
                        sortable: false,
                        editable: false,
                        formatter: formatterKV,
                    })
                    break
                }
                case "array": {
                    columns.push({
                        title: getLabel(name,item),
                        field: name,
                        sortable: false,
                        editable: false,
                        formatter: JSON.stringify,
                    })
                    break
                }
                default: {

                        columns.push({
                            title: getLabel(name,item),
                            field: name,
                            sortable: true,
                        })

                    break
                }
            }
        }
        columns.push({
            title: "操作",
            field: "",
            sortable: false,
            editable: false,
            formatter: function (v, row, index) {
                return `<a class="edit" href="javascript:void(0)"  title="edit">配置</a> <a class="remove" href="javascript:void(0)"  title="remove">删除</a>`
            },
            events:{
                "click .remove":function (e,value,row,index) {
                    // let rowIndex = $(this).attr("array-row")
                    $Table.bootstrapTable('remove', {
                        field: '$index',
                        values: [index]
                    })
                    O.onChange()
                },
                "click .edit":function (e,value,row,index){
                    PopPanelMap({
                        path:`${options["path"]}.items`,
                        schema:items,
                        name:options["name"],
                            title:`修改 ${getLabel(options["name"],options["schema"])}:${index+1}`
                        },
                    function (v){
                        $Table.bootstrapTable("updateRow",{
                            index:index,
                            row:v
                        })

                        O.onChange()
                        return true
                    },function (k){
                            let vs = O.Value;
                            return  vs.hasOwnProperty(k)
                        }
                        ,row)
                }
            }
        })

        let tableOptions = {
            columns: columns,
            width: "100%",
        }
        $Table.bootstrapTable(tableOptions);

    }

    set Value(vs) {
        if (!vs || typeof vs == "undefined") {
            vs = {}
        }
        let list = new Array()
        for (let k in vs){
            let v = vs[k]
            v["__key"] = k
            list.push(v)
        }
        list.sort((a,b)=>{
            return a.__key - b.__key
        })
        this.Table.bootstrapTable("load", list)
    }

    get Value() {
        let list =  this.Table.bootstrapTable('getData')
        let vs = {}
        for (let i in list){
            let v = Object.assign({},list[i])
            let key = v["__key"]
            delete v["__key"]
            vs[key] = v
        }
        return vs
    }

}

class ArrayRenderEnum extends BaseChangeHandler{
    constructor(options) {
        super(options["path"])
        let panel=options["panel"]
        let schema=options["schema"]
        let Id = this.Id
        let items = schema["items"]
        this.Enum = items["enum"]

        let p = $(panel);
        let itemPanel = `<div id="${Id}_items" class="border p-sm-1 btn-toolbar form-control" role="toolbar">`

        for (let i in items["enum"]) {
            let e = items["enum"][i]
            let itemId = `${Id}_${e}`
            itemPanel +=`<div class="custom-control custom-checkbox custom-control-inline"><input type="checkbox" id="${itemId}" value="${e}" name="${Id}" class="custom-control-input"> <label class="custom-control-label" for="${itemId}">${e}</label></div>`
        }
        itemPanel += '</div>'
        p.append(itemPanel)
    }


    get Value() {
        let list = []
        $(`input[name="${this.Id}"]`).each(function () {
            if ($(this).is(':checked')) {
                list.push($(this).val())
            }
        })
        return list
    }

    set Value(vs) {
        if (!vs) {
            vs = []
        }
        let list = vs

        for (let i in this.Enum) {
            let v = this.Enum[i]
            let itemId = this.Id + '_' + v
            $("#" + itemId).attr("checked", list.includes(v))
        }
    }

}

// 基础类型的数组
class ArrayRenderSimple extends BaseChangeHandler{
    constructor(options) {
        super(options["path"])
        let panel = options["panel"]
        let schema = options["schema"]
        this.Options = options
        let Id = this.Id
        this.ValuesList = {}
        this.SchemaItems = schema["items"]
        let myPanel = $(`<div class=""></div>`)
        let $itemPanel = $(`<div id="${Id}_items" class="" role="toolbar"></div>`)
        panel.append(myPanel)
        myPanel.append($itemPanel)
        myPanel.append(`<div class=""> <button id="${Id}_add" type="button" class="btn btn-success btn-sm">+</button></div>`)
        myPanel.on("click", `#${Id}_add`, function () {
            let isEmpty = false
            $itemPanel.find(`input`).each(function () {
                if ($(this).val().length < 1) {
                    isEmpty = true
                    $(this)[0].focus()
                    return false
                }
            })
            if (!isEmpty) {
               o.add()[0].focus()
            }
        })
        this.ItemPanel = $itemPanel
        this.lastIndex = 0
        let o = this
    }
    add() {
        let index = this.lastIndex++
        let itemId = `${this.Id}_item_${index}`
        // let appendAtt = ` array-for="${this.Id}" aria-describedby="btnGroupAddon_${itemId}" `
        let $item = $(`
<div class="input-group input-group-sm m-2" >
    <div class="input-group-prepend">
        <button class="btn btn-danger"  type="button" aria-describedby="btnGroupAddon_${itemId}" data-itemId="${itemId}"> - </button>
    </div>
</div>`)
        this.ItemPanel.append($item)
        let $itemInput =$(createInput(this.SchemaItems, itemId, true))

        $item.prepend($itemInput)
        this.ValuesList[index] = new BaseValue(this.SchemaItems,$itemInput , this.Id)
        let o = this
        this.ValuesList[index].onChange(function () {
            o.onChange()
        })
        $item.on("click", "button", function () {
            $item.remove()
            this.ValuesList[index] = null
            delete this.ValuesList[index]
        })

        return $itemInput
    }
    get Value() {
        let val = []
        let indexList= Object.keys(this.ValuesList)
        indexList.sort()
        for (let i in indexList){
            let ind = indexList[i]
            val.push(this.ValuesList[ind].Value)
        }
        return val
    }

    set Value(vs) {
        if (!Array.isArray(vs)) {
            return
        }
        this.ValuesList = {}
        this.ItemPanel.empty()

        for (let i in vs){
            this.add().val(vs[i])
        }
    }
}
// 结构体
class ObjectRender extends BaseChangeHandler{
    constructor(options) {
        super(options["path"])
        this.Options = options
        let panel = options["panel"]
        let schema = options["schema"]
        let Id = this.Id

        this.$Panel = panel
        this.Fields = {}
        // $(panel).append(this.$Panel)
        if (schema["eo:type"] !== "object") {
            return
        }
        let o = this
        let properties = schema["properties"]
        let sorts = schema["ui:sort"]
        let requiredData = requiredMap(schema["required"])
        for (let i in sorts) {
            let name = sorts[i]
            let sub = properties[name]
            let subId = `${Id}.${name}`

           let field = new FieldPanel(name,{
                panel:this.$Panel,
                schema:sub,
                generator:readGenerator(options),
                path:subId,
                required:requiredData[name]===true,
                parent:schema
            })
            this.Fields[name] =field
            field.onChange(function (){
                o.switchPanel(name,this.Value)
                o.onChange()
            })
        }

        let Switches = {}
        for (let name in properties){
            let sub = properties[name]
            if (sub["switch"]){
                Switches[name] = sub["switch"]
            }
        }
        this.Switches = Switches
        for (let name in this.Fields){
            this.switchPanel(name,this.Fields[name].Value)
        }

    }
    switchPanel(name,value){

        switch (typeof value  ){
            case "string":{
                value = `"${value}"`
                break
            }
            case "object":
            case "undefined":{
                return
            }
            default:{
                break
            }
        }

        for (let f in this.Switches){
            let expression = this.Switches[f]
            try {
                let funcStr = `(function(${name}){return ${expression}})(${value})`
                // console.log(funcStr)
                let r = eval( funcStr)
                if (r === true){
                    this.Fields[f].Show()
                }else{
                    this.Fields[f].Hide()
                }
            }catch (e) {

                // console.log(e)
            }
        }
    }
    get Value() {
        let vs = {}
        for (let k in this.Fields) {
            let fi = this.Fields[k]
            let v = fi.Value
            if (typeof v !== "undefined"){
                vs[k] = v
            }
        }
        return vs
    }

    set Value(v) {
        if (!v || typeof v === "undefined"){
            v = {}
        }
        for (let k in this.Fields) {
            let fi = this.Fields[k]
            fi.Value = v[k]
            this.switchPanel(k,v[k])
        }
    }

}
// 插件弹出窗
class PopPanelPlugin {
    constructor() {

    }

    /**
     * 更新ui面板
     * @param id
     * @param success
     */
    getRenderInfo(id, success){
        let path = id.replaceAll(":","/")
        let n = id.replaceAll(":","_")
        dashboard.getExtenderInfo(path, function (res) {
            if(res.code !== 200){
                return http.handleError(res, "获取extender信息失败")
            }
            success(n,res.data["render"])
        }, function (res){
            return http.handleError(res, "获取extender信息失败")
        })
    }
    show(exclude,callbackFn,data){
        if (this.IsInit === true){
            this.render(exclude,callbackFn,data)
            return
        }
        let o = this
        o.IsInit = true
        dashboard.get(`/api/plugins`,function (pluginsData){
            o.IsInit = true
            o.Plugins = pluginsData.data["plugins"]
            o.PluginsExtenders = new Map(
                o.Plugins.map(object => {
                    return [object.name, object.id];
                }),
            );
            o.render(exclude,callbackFn,data)
        },function (err){
            o.IsInit = false
            common.message(err,"danger")
        })
    }

    render(exclude,callbackFn,data){

        let $Panel = $(`
<div class="pop_window pop_window_small p-3" id="detail_container">
    <div class="pop_window_header">
        <span class="pop_window_title">插件信息</span>
        <button class="pop_window_button btn btn_default close" >关闭</button>
    </div>
</div>`)

        let $Fade = $("<div class='modal-backdrop fade show modal-open'></div>")
        let $Body = $(` <div class="pop_window_body"></div>`)
        $("body").append($Fade)
        $("body").append($Panel)
        $Panel.append($Body)

        let $PluginName  = $(`<select class="form-control form-control-sm"  ></select>`)

        let $Disable = $(`<input type="checkbox" class="form-control-sm form-control" data-toggle="toggle" data-size="sm"/>`)

        $Body.append(
            $(`<div><div><label class=" col-form-label  text-nowrap" for="FormRender_test_target"><span style="color: red">*</span>目标服务</label></div></div>`)
                .append($PluginName)
        )
        $Body.append(
            $(`<div><div><label class=" col-form-label  text-nowrap" for="FormRender_test_target"><span style="color: red">*</span>是否禁用</label></div></div>`)
                .append($Disable)
        )
        $Disable.bootstrapToggle({
            on: '禁用',
            off: '启用'
        })
        let $Config = $(`<div class="border-top p-1"></div>`)
        $Body.append(
            $(`<div><div><label class=" col-form-label  text-nowrap" for="FormRender_test_target"><span style="color: red">*</span>插件配置</label></div></div>`)
                .append($Config)
        )
        let O= this
        function renderConfig(name){
            O.getRenderInfo(O.PluginsExtenders.get(name),function (driver,render){
                $Config.empty()
                O.ConfigTarget = new ObjectRender({
                    path:'plugins.${name}',
                    panel:$Config,
                    schema:render
                })
                if (data){
                    O.ConfigTarget.Value = data.config
                }
            })
        }
        if (data){
            if (data.disable === true || data.disable === "true") {
                $Disable.bootstrapToggle("on");
            } else {
                $Disable.bootstrapToggle("off");
            }
            $PluginName.prepend(`<option value="${data.name}" selected>${data.name}</option>`)
            renderConfig(data.name)
        }else{
            $Disable.bootstrapToggle("off");
            let set = new Set(exclude)
            let first = true
            for (let i in this.Plugins){
                let p = this.Plugins[i]
                if (!set.has(p.name)){
                    $PluginName.append(`<option value="${p.name}" ${first?"selected":""}>${p.name}</option>`)
                    first = false
                }
            }
            if ($PluginName.children().length>0){
                renderConfig($PluginName.val())
            }

            $PluginName.on("change",function (){
                let name = $(this).val()
                renderConfig(name)
            })

        }
        $Panel.append(`<div class="row justify-content-between">
                <div class="col-4">
                    <button type="button" class="btn btn-outline-secondary form_cancel">取消</button>
                </div>
                <div class="col-4" style="text-align: right">
                    <button type="button" class="btn btn-primary form_submit">提交</button>
                </div>
            </div>`)
        let close = function (){
            $Fade.remove()
            $Panel.remove()
        }
        $Panel.show()
        $Fade.click(function (){
            close()
        })
        $Panel.on("click","button.close",close)
        $Panel.on("click","button.form_cancel",close)
        $Panel.on("click","button.form_submit",function (){
            callbackFn({name:$PluginName.val(),disable: $($Disable).prop("checked"),config:O.ConfigTarget.Value})
            close()
        })
    }
}

// 插件配置表
class PluginsRender extends BaseChangeHandler{
    constructor(options) {
        super(options["path"])
        let panel=options["panel"]
        let schema=options["schema"]
        let items =schema["additionalProperties"]
        let Id = this.Id
        this.PopPanel = new PopPanelPlugin()
        let $Table = $(`<table  id="${Id}_items"></table>`)
        $(panel).append($Table)

        this.Table = $Table
        let O = this

        let tableOptions = {
            columns:  [{
                title: "插件名",
                field: "name",
                sortable: false,
            }
                ,{
                    title:"状态",
                    field:"disable",
                    formatter:function (v, row, index){
                        return v === true ?`<span class="btn-danger btn btn-sm">禁用</span>`:`<span class="btn btn-sm btn-outline-primary">启用</span>`
                    }
                },
                {
                    title:"配置",
                    field:"config",
                    formatter:JSON.stringify
                },
                {
                    title: "操作",
                    field: "",
                    sortable: false,
                    editable: false,
                    formatter: function (v, row, index) {
                        return `<a class="edit" href="javascript:void(0)"  title="edit">配置</a> <a class="remove" href="javascript:void(0)"  title="remove">删除</a>`
                    },
                    events:{
                        "click .remove":function (e,value,row,index) {
                            $Table.bootstrapTable('remove', {
                                field: '$index',
                                values: [index]
                            })
                            O.onChange()
                        },
                        "click .edit":function (e,value,row,index){
                            O.PopPanel.show(Object.keys(O.Value),
                                function (v){
                                    $Table.bootstrapTable("updateRow",{
                                        index:index,
                                        row:v
                                    })

                                    O.onChange()
                                    return true
                                },row)
                        }
                    }
                }],
            width: "100%",
        }

        $Table.bootstrapTable(tableOptions);


            let $btn = $(`
<div id="${Id}_toolbar" class="m-1">
<button id="${Id}_AddButton" type="button" class="btn btn-secondary btn-sm">Add</button>
</div>`)
            $(panel).prepend($btn)
            $btn.on("click","button", function (event) {
                O.PopPanel.show(Object.keys(O.Value),
                    function (v){
                        $Table.bootstrapTable("append",[v])
                        $Table.bootstrapTable('scrollTo', 'bottom')
                        O.onChange()
                    })
                return false
            })
    }

    set Value(vs) {
        if (!vs || typeof vs == "undefined") {
            vs = {}
        }
        let list = new Array()
        for (let k in vs){
            let v = vs[k]
            v["name"] = k
            list.push(v)
        }
        list.sort((a,b)=>{
            return a.name - b.name
        })
        this.Table.bootstrapTable("load", list)
    }

    get Value() {
        let list =  this.Table.bootstrapTable('getData')
        let vs = {}
        for (let i in list){

            let v = list[i]
            let key = v["name"]
            vs[key] = {
                disable:v["disable"],
                config:v["config"],
            }
        }
        return vs
    }

}

function BaseGenerator(options) {

    let schema = options["schema"]
    switch (schema["eo:type"]) {
        case "object": {
            return new ObjectRender(options)
        }
        case "array": {
            let items = schema["items"]
            switch (items["eo:type"]) {
                case "object": {
                    return new ObjectArrayRender(options)
                }
                case "map": {
                    throw "now support map in array"
                    // return new InnerMapRender(options)
                }
                case "integer": {
                }
                case "number": {
                }
                case "string": {
                    if (items["enum"]) {
                        return new ArrayRenderEnum(options)
                    }
                    return new ArrayRenderSimple(options)
                }
                case "array":{
                    throw `not allow type:${items["eo:type"]} in array`
                }
                case "require":{
                    return new RequireArrayRender(options)
                }
                default:{
                    throw `unknown type:${items["eo:type"]} in array`

                }
            }

        }
        case "map": {
            let item = schema["additionalProperties"];
            switch (item["type"]){
                case "object":{
                    if(options["name"] === "plugins"){
                        return new PluginsRender(options)
                    }
                    return new ObjectMapRender(options)
                }
            }
            return new SimpleMapRender(options)
        }
        case "boolean": {
            return new SwitchRender(options)
        }
        case "integer": {
        }
        case "number": {
        }
        case "string": {
            if (schema["enum"]) {
                return new BaseEnumRender(options)
            }
            return new BaseInputRender(options)
        }
        case "require": {
            return new RequireRender(options)
        }
        case "formatter" :{
            return null
        }
    }
    throw `unknown type:${schema["eo:type"]}`
}
class SchemaHandler {
    constructor(schema) {
        this.s = schema
    }
    toJsonSchema(s){
        let schema = Object.assign({},s)
        switch (schema["eo:type"]){
            case "map":{
                schema["additionalProperties"] = this.toJsonSchema(schema["additionalProperties"])
                break
            }
            case "array":{
                schema["items"] = this.toJsonSchema(schema["items"])
                break
            }
        }
        delete schema["ui:sort"]
        delete schema["eo:type"]
        delete schema["switch"]
        delete schema["skill"]
        delete schema["empty_label"]

        return schema
    }

    get JsonSchema(){
        if (typeof this.v === "undefined" || this.v === null){
            this.v = this.toJsonSchema(this.s)
            console.log(JSON.stringify(this.v))
        }
        return this.v
    }
    get Schema(){
        return this.s
    }
}
class FormRender {

    constructor(options) {

        this.ObjectName = `${RootId}.${options["name"]}`
        this.JsonSchema = new SchemaHandler(options["schema"])

        $(options["panel"]).empty()

        this.Object = readGenerator(options)({
            generator:options["generator"],
            panel:options["panel"],
            schema:options["schema"],
            path: this.ObjectName,
            name:options["name"],
        })

    }

    check() {

        let r =  CheckBySchema(this.ObjectName,this.JsonSchema.JsonSchema,this.Value)
        console.log(`check:${this.ObjectName} = ${JSON.stringify(this.Value)} :${JSON.stringify(r)}`);
        if(typeof r === "undefined"){
            return true
        }
  
        return  JSON.stringify(r)

    }

    get Value() {
        return this.Object.Value
    }

    set Value(v) {
        this.Object.Value = v
    }
}

