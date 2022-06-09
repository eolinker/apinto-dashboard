const RootId = "FormRender"

const validate = {
    _validator: null,
    djv: function () {
        if (this._validator){
            return this._validator
        }
        this._validator = new djv()
        return this._validator
    },
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
function createEnum(schema, id, appendAttr) {
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
        if (schema["required"]) {
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

function createInput(schema, id, appendAttr) {
    let readOnly = ""
    this.schema = schema
    if (schema["readonly"] === true) {
        readOnly = "readonly"
    }

    let input = `<input ${readOnly} class="form-control form-control-sm" id="${id}" aria-describedby="validation_${id}" `;
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

    if (schema["required"]) {
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

function createLabel(name,schema, id, appendAttr) {
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
    if (schema["required"]) {
        require = '<span style="color: red">*</span>'
    }
    return `<label class="col-sm-3 col-form-label text-right text-nowrap" ${idFor} ${appendAttr}>${require}${getLabel(name,schema)}:</label>`
}
class BaseChangeHandler {
    constructor(id) {
        this.Id = id
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
class BaseValue  {
    constructor(schema, target) {

        this.Schema = schema
        this.Target = target
        if (typeof schema["default"] !== "undefined") {
            this.Value = schema["default"]
        }
        const JsonSchema = new SchemaHandler(schema).JsonSchema
        this.InputValid(JsonSchema,target)
    }
    onChange(fn){
        const o = this
        $(this.Target).on("change",function (){
            fn.apply(o)
        })
    }
    ValidHandler(v) {
        const id =  this.Target.attr("id")
        console.debug("ValidHandler:",id,"=",v)
        let value = v
        value = valueForType(this.Schema["eo:type"], value)

        let rs = CheckBySchema(id, this.Schema, value)
        if (typeof rs === "undefined") {
            $(this).removeClass("is-invalid")
            $(this).addClass("is-valid")

        } else {
            $(this).removeClass("is-valid")
            $(this).addClass("is-invalid")

        }
    }
    InputValid(schema, target) {
        const o = this
        $(target).on("change", function () {
            o.ValidHandler($(this).val(), schema)
        })
    }

    get Value() {
        let val =  valueForType(this.Schema["eo:type"], $(this.Target).val())
        console.log(`get value:${this.Target.attr("id")}[${typeof val}]=${val}`)
        return val
    }

    set Value(v) {

        if ( typeof v === "undefined"){
            let schema = this.Schema
            if (typeof schema["default"] !== "undefined") {
               v = schema["default"]
            }
        }
        $(this.Target).val(v)
    }

}

class BaseEnumRender extends BaseValue {
    constructor(options) {
        const schema = options["schema"]
        const path = options["path"]
        const panel = options["panel"]
        super(schema, $(createEnum(schema, path)))

        $(panel).append(this.Target)

        // this.Target.on("click",function (){
        //     $(this).trigger("change")
        // })
    }

}

class SwitchRender extends BaseValue {
    constructor(options) {
        const schema = options["schema"]
        const path = options["path"]
        const panel = options["panel"]

        super(schema, $(`<input id="${path}" type="checkbox" class="form-control-sm" data-toggle="toggle" data-size="sm"/>`))

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

class RequireRender extends BaseValue {
    constructor(options) {
        const schema = options["schema"]
        const path = options["path"]
        const panel = options["panel"]

        super(schema, $(`<select id=${path} class="form-control form-control-sm">
 
</select>`))
        $(panel).append(this.Target)

        const select = this.Target
        dashboard.searchSkill(ModuleName(),schema["skill"],function (res){
            let lastValue =
            $(select).empty()
            if(schema["required"]){
                $(select).append(`<option value="">请选择</option>`)
            }else{
                $(select).append(`<option value="">不启用</option>`)
            }
            for (let i in res.data){
                let d = res.data[i]
                $(select).append(`<option value="${d.id}">${d.id}[${d.driver}]</option>`)
            }
        })
    }
}

class BaseInputRender extends BaseValue {
    constructor(options) {
        const schema = options["schema"]
        const path = options["path"]
        const panel = options["panel"]
        super(options["schema"], $(createInput(schema, path)))
        $(panel).append(this.Target)
        if (schema["description"] && schema["description"].length >0){
            $(panel).append(`<small id="help:${path}" class="text-muted">${schema["description"]}</small>`)
        }
    }
}

class FieldPanel {
    constructor(name,options,parent) {
        const panel=options["panel"]
        const schema=options["schema"]
        const generator = options["generator"]
        const id = options["path"]
        this.Parent = parent
        this.Id = id
        if (schema["eo:type"] === "object"){
            this.$Panel = $(`<div class="form-group row mb-1 overflow-hidden">${createLabel(name,schema,id)}</div>`)
            panel.append(this.$Panel)
//             const $FieldValuePanel = $(`
// <div class="form-group row mb-1 overflow-hidden">
//     <div class="col-sm-11 offset-sm-1 border p-sm-1">
//     </div>
// </div>`)
            const $FieldValuePanel = $(`
    <div class="col-sm-11 offset-sm-1 border p-sm-1">
    </div>`)
            this.$Panel.append($FieldValuePanel)
             this.$Value = generator({panel:$FieldValuePanel,schema:schema,generator: generator,path:id})

        }else{
            const $FieldValuePanel = $(`<div class="col-sm-9"></div>`)
            this.$Panel  = $(`<div class="form-group row mb-1 overflow-hidden">${createLabel(name,schema,id)}</div>`)
            this.$Panel.append($FieldValuePanel)
            panel.append(this.$Panel)
            this.$Value = generator({panel:$FieldValuePanel,schema:schema,generator: generator,path:id})
        }
    }
    onChange(fn){
        console.log(this.Id)
        const o = this
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

class MapRender extends BaseChangeHandler {

    constructor(options) {
        super(options["path"])
        const panel=options["panel"]
        const schema=options["schema"]
        const generator = options["generator"]
        const path = options["path"]

        this.Id = path;
        this.Schema = schema

        this.$Panel =$(`<div id='${path}_panel' class='container-fluid'></div>`)
        $(panel).append(this.$Panel)

        this.GeneratorHandler = generator

    }

    set Value(v) {

        const Items = this.Schema["additionalProperties"]
        const Id = this.Id
        const keySchema = {type: "string","eo:type":"string"}

        this.$Panel.empty()
        switch (Items["eo:type"]) {
            case "object":
            case "map":
                break;
            default: {
                this.$Panel.append(`
<div class="input-group input-group-sm m-2">
    <div class="input-group-prepend">
        <div class="input-group-text  btn" id="btnGroupAddon_${Id}_new">+</div>
    </div>
    ${createInput(keySchema, `${Id}_key`, `aria-describedby="btnGroupAddon_${Id}_new" placeholder="Input new key" `)}
    <div class="input-group-prepend ">
        <div class="input-group-text  btn" id="btnGroupAddon_${Id}_eq">=</div>
    </div>
    ${createInput(Items, `${Id}_value`, `aria-describedby="btnGroupAddon_${Id}_eq" placeholder="Input new value" `)}
</div>`)

                for (let k in v) {
                    $(this.$Panel).append('<div class="input-group input-group-sm m-2"></div>')
                    let itemPanel = $(this.$Panel).children().last()
                    itemPanel.append(`
<div class="input-group-prepend">
    <button class="btn btn-danger" id="btnGroupAddon_${Id}_key_${k}" type="button" data-itemId="${Id}_key"> - </button>
</div>`)

                    let childItemKey = new BaseInputRender({panel:itemPanel, schema:keySchema, path:`${Id}_key_${k}`})
                    childItemKey.Value = k
                    itemPanel.append('<div class="input-group-prepend"><div class="input-group-text  btn" >=</div></div>')
                    let childItemValue = new BaseInputRender({panel:itemPanel, schema:Items, path:`${Id}_value_${k}`})
                    childItemValue.Value = v[k];
                }
                break
            }
        }

    }

    get Value() {
        return {}
    }

}

class InnerObjectRender extends BaseChangeHandler{
    constructor(options) {
        super(options["path"])
        const panel=options["panel"]
        const schema=options["schema"]
        const generator = options["generator"]
        const Id = options["path"];

        this.Id = Id;

        const items = schema["items"]
        const p = $(panel)
        const $btn = $(`
<div id="${Id}_toolbar">
<button id="${Id}_AddButton" type="button" class="btn btn-secondary">Add</button>
</div>`)
        p.append($btn)
        $btn.on("click","button", function (event) {
            $Table.bootstrapTable('append', [{}])
            $Table.bootstrapTable('scrollTo', 'bottom')
            return false
        })

        const $Table = $(`<table  id="${Id}_items"></table>`)
        p.append($Table)
        $Table.delegate("a.remove", "click", function (event) {
            let rowIndex = $(this).attr("array-row")
            $Table.bootstrapTable('remove', {
                field: '$index',
                values: [Number(rowIndex)]
            })
        })
        this.Table = $Table
        const properties = items["properties"]
        const uiSort = items["ui:sort"]
        let lastDetailRow = undefined
        let lastField = undefined
        const o = this
        function DetailFormatterHandler(fieldIndex) {

            this.detailFormatter = function (index, row, $element) {

                if (typeof lastDetailRow !== "undefined" && lastDetailRow !== index) {
                    $Table.bootstrapTable('collapseRow', lastDetailRow)
                }
                lastDetailRow = index
                lastField = fieldIndex
                let name = uiSort[fieldIndex]
                let item = properties[name]

                let child = generator({panel:$element, schema:item, generator:generator, path:`${Id}_${name}`})
                child.Value = row[name]
                return ""
            }
            return this
        }

        function NotDetailFormatterMap(index, row, $element) {
            if (typeof lastDetailRow !== "undefined") {
               $Table.bootstrapTable('collapseRow')
                lastDetailRow = undefined
                lastField = undefined
            }
            $Table.bootstrapTable('collapseAllRows')

            return ''
        }

        function formatterKV(v) {
            let html = ""
            for (let k in v) {
                html += "<span class='btn btn-outline-secondary btn-sm'>" + k + "=" + v[k] + "</span>\n"
            }
            html += ""
            return html
        }


        const columns = []
        columns.push({
            title: "",
            field: "__index",
            sortable: false,
            editable: false,
            detailFormatter: NotDetailFormatterMap,
            formatter: function (v, row, index) {
                return index
            }
        })

        for (let i in uiSort) {
            let name = uiSort[i]
            const item = properties[name]
            switch (item["eo:type"]) {
                case "object": {
                    columns.push({
                        title: getLabel(name,item),
                        field: item["name"],
                        sortable: false,
                        editable: false,
                        formatter: formatterKV,
                        detailFormatter: new DetailFormatterHandler(i).detailFormatter,

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
                        detailFormatter: new DetailFormatterHandler(i).detailFormatter
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
                        detailFormatter: new DetailFormatterHandler(i).detailFormatter,
                    })
                    break
                }
                default: {
                    if (item["enum"]) {
                        columns.push({
                            title: getLabel(name,item),
                            field: name,
                            sortable: true,
                            detailFormatter: NotDetailFormatterMap,
                            editable: {
                                type: "select",
                                options: {
                                    items: item["enum"]
                                }
                            },
                        })
                    } else {
                        let typeInput = "text"
                        if (item["eo:type"] === "number" || item["eo:type"] === "integer") {
                            typeInput = "number"
                        }
                        columns.push({

                            title: getLabel(name,item),
                            field: name,
                            sortable: true,
                            width: 200,
                            detailFormatter: NotDetailFormatterMap,
                            editable: {
                                type: typeInput
                            }
                        })
                    }
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
                return `<a class="remove" href="javascript:void(0)" array-row="${index}" title="remove">删除</a>`
            }
        })

        const tableOptions = {
            columns: columns,
            editable: true,


            detailView: true,
            detailViewByClick: true,
            detailViewIcon: false,
            width: "100%",
            onEditorShown: function (field, row, $el, editable) {
                $Table.bootstrapTable('collapseAllRows')
                return true;
            },
            onEditorSave: function (field, row, oldValue, $el) {
                // let data = $Table.bootstrapTable('getData')
                // if (field !== "__index" ) {
                //     const rowIndex = $el.parent().data("index")
                //     data.Data[rowIndex][field] = row
                // }
                // return true;

                o.onChange()
            },
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

class InnerMapRender extends BaseChangeHandler{
    constructor(options) {
        super(options["path"])
        // super(options["schema"],path)
    }

    set Value(v) {

    }

    get Value() {
        return {}
    }
}

class ArrayRenderEnum {
    constructor(options) {
        const panel=options["panel"]
        const schema=options["schema"]
        const Id = options["path"];

        this.Id = Id;
        const items = schema["items"]
        this.Enum = items["enum"]

        let p = $(panel);
        let itemPanel = `<div id="${Id}_items" class="border p-sm-1 btn-toolbar " role="toolbar">`

        for (let i in items["enum"]) {
            let e = items["enum"][i]
            let itemId = `${Id}_${e}`
            itemPanel +=`<div class="custom-control custom-checkbox custom-control-inline"><input type="checkbox" id="${itemId}" value="${e}" name="${Id}" class="custom-control-input"> <label class="custom-control-label" for="${itemId}">${e}</label></div>`
        }
        itemPanel += '</div>'
        p.append(itemPanel)
    }

    onChange(fn){

    }
    get Value() {
        const list = []
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
        const list = vs

        for (let i in this.Enum) {
            let v = this.Enum[i]
            const itemId = this.Id + '_' + v
            $("#" + itemId).attr("checked", list.includes(v))
        }
    }

}

// 基础类型的数组
class ArrayRenderSimple extends BaseChangeHandler{
    constructor(options) {
        super(options["path"])
        const panel=options["panel"]
        const schema=options["schema"]
        const generator = options["generator"]
        const Id = options["path"];

        const JsonSchema = new SchemaHandler(schema["items"])

        const items = schema["items"]

        this.Id = Id;
        let p = $(panel);
        const $itemPanel = $(`<div id="${Id}_items" class="border p-sm-1 btn-toolbar " role="toolbar"></div>`)
        p.append($itemPanel)
        const $newInput = $(createInput(items, `${Id}_new`, `aria-describedby="btnGroupAddon_${Id}_new" placeholder="Input new" `) )
        const $newItem = $(`
<div class="input-group input-group-sm m-2">
    <div class="input-group-prepend ">
        <div class="input-group-text  btn" id="btnGroupAddon_${Id}_new">+</div>
    </div>
</div>`)
        $newItem.append($newInput)
        $itemPanel.append($newItem)
        $newInput.on("change", function () {
            let v = $(this).val()
            if (items["eo:type"] === "integer" || items["eo:type"] === "number") {
                v = Number(value)
            }
            if (v !== "") {
               let ckr = CheckBySchema(Id, JsonSchema.JsonSchema, v)
                if (typeof ckr === "undefined") {
                    add(v)
                    $(this).val("")
                }
            }
            return false
        })

        let lastIndex = 0
        const o = this
        function add(value) {
            const itemId = `${Id}_item_${lastIndex++}`

            const appendAtt = ` array-for="${Id}" aria-describedby="btnGroupAddon_${itemId}" `
            const $itemInput = $( createInput(items, itemId, appendAtt))
            const $item = $(`
<div class="input-group input-group-sm m-2" id="array-item_${itemId}">
    <div class="input-group-prepend">
        <button class="btn btn-danger" id="btnGroupAddon_${itemId}" type="button" aria-describedby="btnGroupAddon_${itemId}" data-itemId="${itemId}"> - </button>
    </div>
</div>`)
            $item.append($itemInput)
            $itemPanel.append($item)
            $itemInput.val(value)
            o.onChange()
            return false
        }


        $itemPanel.delegate('button', "click", function () {
            let itemId = $(this).attr('data-itemId')
            $itemPanel.children(`#array-item_${itemId}`).remove()
            o.onChange()
        })

    }

    get Value() {
        let val = []
        $(`[array-for='${this.Id}']`).each(function () {
            val.push($(this).val())
        })
        return val
    }

    set Value(vs) {
        if (!Array.isArray(vs)) {
            return
        }
        const arrayId = `[array-for='${this.Id}']`
        let list = $(arrayId)
        if (list.length < vs.length) {
            let num = vs.length - list.length
            for (let i = 0; i < num; i++) {
                this.add()
            }
        } else if (list.length > vs.length) {
            let num = vs.length - list.length
            for (let i = 0; i < num; i++) {
                list.last().remove()
            }
        }
        let index = 0;
        $(arrayId).each(function () {
            $(this).val(vs[index])
            index++;
        })
    }
}
// 结构体
class ObjectRender extends BaseChangeHandler{
    constructor(options) {
        super(options["path"])
        this.Options = options
        const panel = options["panel"]
        const schema = options["schema"]
        const Id = options["path"]
        const generator = options["generator"]
        this.Fields = {}

        if (schema["eo:type"] !== "object") {
            return
        }
        const o = this
        let properties = options["schema"]["properties"]
        let sorts = options["schema"]["ui:sort"]
        for (let i in sorts) {
            let name = sorts[i]
            let sub = properties[name]
            const subId = `${Id}.${name}`

           let field = new FieldPanel(name,{
                panel:panel,
                schema:sub,
                generator:generator,
                path:subId
            },this)
            this.Fields[name] =field
            field.onChange(function (){
                o.switch(name,this.Value)
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
            this.switch(name,this.Fields[name].Value)
        }

    }
    switch(name,value){

        switch (typeof value  ){
            case "string":{
                value = `"{value}"`
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
                console.log(funcStr)
                let r = eval( funcStr)
                if (r === true){
                    this.Fields[f].Show()
                }else{
                    this.Fields[f].Hide()
                }
            }catch (e) {
                console.log(e)
            }
        }
    }
    get Value() {
        let v = {}
        for (let k in this.Fields) {
            let fi = this.Fields[k]
            v[k] = fi.Value
        }
        return v
    }

    set Value(v) {
        if (!v || typeof v === "undefined"){
            v = {}
        }
        for (let k in this.Fields) {
            let fi = this.Fields[k]
            fi.Value = v[k]
        }
    }

}


function BaseGenerator(options) {
    const schema = options["schema"]
    switch (schema["eo:type"]) {
        case "object": {
            return new ObjectRender(options)
        }
        case "array": {
            const items = schema["items"]
            switch (items["eo:type"]) {
                case "object": {
                    return new InnerObjectRender(options)
                }
                case "map": {
                    return new InnerMapRender(options)
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
            }
            throw `not allow type:${schema["eo:type"]} in array`
        }
        case "map": {
            return new MapRender(options)
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
        delete schema["name"]
        delete schema["switch"]
        delete schema["skill"]

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
        const panel = options["panel"]
        const schema = options["schema"]
        // const generator = options["generator"]
        const name = options["name"]
        const newOption = Object.assign({

        },options)
        if (!newOption.generator){
            newOption.generator = BaseGenerator
        }
        this.JsonSchema = new SchemaHandler(schema)
        // if (!generator || typeof generator !== "function") {
        //     options["generator"] = BaseGenerator
        // }
        $(panel).empty()
        this.ObjectName = `${RootId}.${name}`
        newOption["path"] = this.ObjectName

        this.Object = newOption["generator"](newOption)
    }

    check() {


        let r =  CheckBySchema(this.ObjectName,this.JsonSchema.JsonSchema,this.Value)
        console.log(`check:${this.ObjectName} = ${JSON.stringify(this.Value)} :${JSON.stringify(r)}`);
        if(typeof r === "undefined"){
            return true
        }
        common.message(
            r
        )
        return  r

    }

    get Value() {
        return this.Object.Value
    }

    set Value(v) {
        this.Object.Value = v
    }
}

