let RootId = "FormRender"

let validate = {
    _validator: null,
    djv: function () {
        if (this._validator) {
            return this._validator
        }
        this._validator = new djv()
        return this._validator
    },
}
function formatterKV(v) {
    let html = ""
    for (let k in v) {
        let value = ""
        switch (typeof v[k]){
            case "object":
                value = JSON.stringify(v[k])
                break
            default:
                value = v[k]
                break
        }

        html += "<span class='px-1 border btn-sm  btn-outline-secondary'>" + k + "=" + value + "</span><br>"
    }
    html += ""
    return html
}
function configTable(uiSort,properties){
    let columns = []
    let dataFormat=new Set(["date-time","time","date","duration"])
    for (let i in uiSort) {
        let name = uiSort[i]
        let item = properties[name]
        let format = item["format"]
        if (typeof format != "undefined" && dataFormat.has(format)) {
            item["eo:type"] = "date-time"
            item["eo:format"] = format
            if (item["type"] !== "string") {
                delete (item["format"])
            }
        }
        switch (item["eo:type"]) {
            case "object": {
                columns.push({
                    title: getLabel(name, item),
                    field: name,
                    sortable: false,

                    formatter: JSON.stringify,

                })
                break
            }
            case "date-time":{
                columns.push({
                    title: getLabel(name, item),
                    field: name,
                    sortable: false,
                    formatter: function (v) {
                        if (v === 0) {
                            return "不过期"
                        }
                        return date.formatDate(v)
                    }

                })
                break
            }
            case "map": {
                columns.push({
                    title: getLabel(name, item),
                    field: name,
                    sortable: false,
                    editable: false,
                    formatter: formatterKV,
                })
                break
            }
            case "array": {
                columns.push({
                    title: getLabel(name, item),
                    field: name,
                    sortable: false,
                    editable: false,
                    formatter: formatterKV,
                })
                break
            }
            default: {

                columns.push({
                    title: getLabel(name, item),
                    field: name,
                    sortable: true,
                })

                break
            }
        }
    }
    return columns
}
function requiredMap(required) {
    if (!required) {
        return {}
    }
    let r = {}
    for (let i in required) {
        r[required[i]] = true
    }
    return r
}

function readId(path) {
    return path.replaceAll(".", "_")
}

function readGenerator(options) {
    let fn = options["generator"]
    if (!fn) {
        return BaseGenerator
    }
    return fn
}

function CheckBySchema(id, schema, value) {
    let env = validate.djv()
    if (!env.resolved.hasOwnProperty(id)) {
        env.addSchema(id, schema);
    }
    return env.validate(id, value);
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

function createEnum(schema, id, required, appendAttr) {
    let readOnly = ""
    this.schema = schema
    if (schema["readonly"] === true) {
        readOnly = "readonly"
    }
    if (typeof appendAttr === "undefined") {
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

function createText(schema, id, required, appendAttr) {
    let readOnly = ""

    if (schema["readonly"] === true) {
        readOnly = "readonly"
    }
    let idstr = ""
    if (id && id.length > 0) {
        idstr = `id="${id}"`
    }
    // let input = `<textarea class="form-control" rows="3"  type="text" id="inputDescription"></textarea>`
    let input = `<textarea ${readOnly} class="form-control form-control-sm" rows="3" ${idstr} aria-describedby="validation_${id}"`;
    if (appendAttr) {
        input += appendAttr
    }


    if (required) {
        input += ' required'
    }
    input += '></textarea>'
    return input
}

function createInput(schema, id, required, appendAttr) {
    let readOnly = ""

    if (schema["readonly"] === true) {
        readOnly = "readonly"
    }
    let idstr = ""
    if (id && id.length > 0) {
        idstr = `id="${id}"`
    }
    let input = `<input ${readOnly} class="form-control form-control-sm" ${idstr} aria-describedby="validation_${id}" `;
    if (appendAttr) {
        input += appendAttr
    }

    function readFormatForString(format) {

        switch (format) {

            case "email":{}
            case "password":{}
            case "date":{}
            case "time":{}
                case "number": {
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

function getLabel(name, schema) {
    let label = schema["label"]
    if (!label || label.trim() === "") {
        label = name
    }

    label = label.replace(label[0], label[0].toUpperCase());
    label = label.replaceAll("_", " ")
    return label
}

function createLabel(name, schema, id, required, appendAttr) {
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
    return `<label class=" col-form-label  text-nowrap" ${idFor} ${appendAttr}>${require}${getLabel(name, schema)}</label>`
}

class BaseChangeHandler {
    constructor(path) {
        this.Id = readId(path)
    }

    onChange(fn) {

        if (fn) {
            if (!this.ChangeHandler) {
                this.ChangeHandler = []
            }
            this.ChangeHandler.push(fn)
        } else {
            console.log("change:", this.Id)
            for (let i in this.ChangeHandler) {
                this.ChangeHandler[i].apply(this)
            }
        }
    }
}

class BaseValue {
    constructor(schema, target, path) {
        this.Id = readId(path)
        this.Schema = schema
        this.Target = target
        if (typeof schema["default"] !== "undefined") {
            this.Value = schema["default"]
        }
        let JsonSchema = new SchemaHandler(schema).JsonSchema
        this.InputValid(JsonSchema, target)
    }

    onChange(fn) {
        let o = this
        $(this.Target).on("change", function () {
            fn.apply(o)
        })
    }

    isOk(v) {
        return
    }

    ValidHandler(v) {
        let id = this.Id
        console.debug("ValidHandler:", id, "=", v)
        let value = v
        value = valueForType(this.Schema["type"], value)
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
        return valueForType(this.Schema["type"], $(this.Target).val())
    }

    set Value(v) {

        if (typeof v === "undefined") {
            let schema = this.Schema
            if (typeof schema["default"] !== "undefined") {
                v = schema["default"]
            }
        }
        switch ($(this.Target).get(0).tagName) {
            case "select": {
                if ($(this.Target).prop("multiple")) {
                    for (let i in v) {
                        let item = v[i]
                        $(this.Target).find(`option[value="${item}"]`).prop("selected", true);
                    }
                } else {
                    $(this.Target).find(`option[value="${v}"]`).prop("selected", true);
                }
                break
            }
            default: {
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
        super(schema, $(createEnum(schema, readId(path))), path)

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

        super(schema, $(`<input id="${readId(path)}" type="checkbox" class="form-control-sm" data-toggle="toggle" data-size="sm"/>`), path)

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
        let schema = options["schema"]
        let path = options["path"]
        let panel = options["panel"]

        super(schema, $(`<select id=${readId(path)} class="form-control form-control-sm"></select>`), path)
        this.MOptions = options
        $(panel).append(this.Target)

        let select = this.Target
        dashboard.searchSkill(ModuleName(), schema["skill"], function (res) {

            $(select).empty()
            if (options["required"]) {
                $(select).append(`<option value="">请选择</option>`)
            } else {
                if (schema["empty_label"]) {
                    $(select).append(`<option value="">${schema["empty_label"]}</option>`)
                } else {
                    $(select).append(`<option value="">不启用</option>`)
                }
            }
            for (let i in res.data) {
                let d = res.data[i]
                $(select).append(`<option value="${d.id}">${d.id}[${d.driver}]</option>`)
            }
        })
    }

    isOk(v) {
        if (!v || v.length === 0) {
            if (this.MOptions["required"]) {
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

        super(schema, $(`<select id=${readId(path)} data-live-search="true" class="selectpicker form-control form-control-sm" title="请选择" data-width="100%" data-size="5" data-selected-text-format="count>4" multiple></select>`), path)
        this.MOptions = options
        $(panel).append(this.Target)

        let $Select = this.Target
        dashboard.searchSkill(ModuleName(), schema["skill"], function (res) {

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
            for (let i in res.data) {
                let d = res.data[i]
                $($Select).append(`<option value="${d.id}">${d.id}[${d.driver}]</option>`)
            }
            $Select.selectpicker()
        })
    }

    set Value(vs) {
        this.Target.selectpicker('val', vs);
    }

    get Value() {
        return this.Target.selectpicker('val');
    }

    isOk(v) {
        if (!v || v.length === 0) {
            if (this.MOptions["required"]) {
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
        super(options["schema"], $(createInput(schema, readId(path), options["required"])), path)
        $(panel).append(this.Target)
        if (schema["description"] && schema["description"].length > 0) {
            $(panel).append(`<small id="help:${path}" class="text-muted">${schema["description"]}</small>`)
        }
    }
}

class BaseInputTextRender extends BaseValue {
    constructor(options) {
        let schema = options["schema"]
        let path = options["path"]
        let panel = options["panel"]
        super(options["schema"], $(createText(schema, readId(path), options["required"])), path)
        $(panel).append(this.Target)
        if (schema["description"] && schema["description"].length > 0) {
            $(panel).append(`<small id="help:${path}" class="text-muted">${schema["description"]}</small>`)
        }
    }
}
class FileRender extends BaseChangeHandler {
    constructor(options) {
        super(options["path"])
        let schema = options["schema"]
        let path = options["path"]
        let panel = options["panel"]
        let id = readId(path)
        this.Id = id
        this.Schema = schema
        let input = $(`<input type="file" multiple id="${id}" class="custom-file-input" />`)

        $(panel).append(input)
        const $this = this;
        this.$Value = [];


        if (schema["description"] && schema["description"].length > 0) {
            $(panel).append(`<small id="help:${path}" class="text-muted">${schema["description"]}</small>`)
        }


    }


    set Value(v) {
        if (typeof v !== "array"){
            v= []
        }
        this.$Value = v
        let initialPreview=[],initialPreviewConfig=[]
        let $this = this
        this.$Value.forEach((v,i)=>{
            let item = DecodeFile(v);
            // console.log(item.DataUrl)
            initialPreview.push(item.DataUrl)
            initialPreviewConfig.push({
                caption: item.name,
                size:item.size,
                key:`${item.size}_${item.name}`,
                filetype:item.type,
                fieldId: item.index,
            })
        })
        $(`#${this.Id}`).fileinput('destroy').fileinput({
            showUpload:false,
            uploadUrl: '#',
            fileActionSettings:{
                showRotate:false,
                showUpload:false,
            },
            generateFileId:function (file){
                if(!file){
                    return null
                }
                return file.size +"_" + _getFileName(file)
            },
            initialPreviewAsData: true,
            overwriteInitial: false,
            initialPreview: initialPreview,
            initialPreviewConfig: initialPreviewConfig,
        }).on("filebeforeload", function(event, file, index, reader) {
            console.log("try add:index="+index)
            if (!file){
                return false;
            }
            let key = file.size +"_"+file.name
            for(let v of $this.$Value){
                if (v.index === key){
                    return false
                }
            }
            return true

        }).on('fileloaded', function(event, file, previewId, fileId, index, reader) {
            console.log("add:index="+index+", fileId="+fileId)
            let f = new FileItem(_getFileName(file),file.size,file.type,reader.result);
            $this.$Value.push(f)
            $this.onChange()
        }).on('filepreremove', function(event, id, index) {
            console.log("try remove:"+index+", id="+id)
             for (let i=0;i<$this.$Value.length;i++){

                if ($this.$Value[i].index === index){
                    console.log("remove:"+index+" at "+i)
                    $this.$Value.splice(i,1)
                    break
                }
            }
            $this.onChange()
        });
    }

    get Value() {
        let vs =[]
        for (let v of this.$Value){
            vs.push(v.Value)
        }
        return vs
    }
}

//字段封装
class FieldPanel {
    constructor(name, options) {
        let panel = options["panel"]
        let schema = options["schema"]
        this.Id = readId(options["path"])
        this.$Panel = $(`<div class=""></div>`)
        this.Enable = true
        panel.append(this.$Panel)
        let valuePanel = $(`<div class=""></div>`)
        this.$Panel.append(`<div class="">${createLabel(name, schema, this.Id, options["required"])}</div>`)
        this.$Panel.append(valuePanel)
        this.$Value = readGenerator(options)({
            panel: valuePanel,
            schema: schema,
            parent: options["parent"],
            generator: readGenerator(options),
            path: options["path"],
            name: name,
            required: options["required"]
        })

        switch (schema["type"]) {
            case "array": {
                if (schema["items"]["enum"]) {
                    break
                }
            }
            case  "object": {
                valuePanel.addClass("border-top border-bottom px-3 py-1")
            }
        }
    }

    onChange(fn) {
        console.log(this.Id)
        let o = this
        this.$Value.onChange(function () {
            fn.apply(o)
        })
    }

    Show() {
        this.Enable = true
        $(this.$Panel).show()
    }

    Hide() {
        this.Enable = false
        $(this.$Panel).hide()
    }

    set Value(v) {
        this.$Value.Value = v
    }

    get Value() {
        // if ($(this.$Panel).is(":visible")) {
        //     return this.$Value.Value
        // } else {
        //     return {}
        // }
        if (this.Enable){
            return this.$Value.Value
        }
        return  undefined
    }
}

class DatetimeRender extends BaseChangeHandler {
    constructor(options) {
        super(options["path"]);
        this.Options = options
        this.Schema = options["schema"]
        this.Panel = options["panel"]
        this.Panel.append(`
            <div class="form-group">
                <div class="input-group date form_datetime col-md-5" data-date-format="dd MM yyyy - HH:ii p" id="${this.Id}_date">
                    <input class="form-control" size="16" type="text" value="" readonly id="${this.Id}_date_data">
                    <span class="input-group-addon col-md-1"><span class="glyphicon glyphicon-remove"><svg xmlns="http://www.w3.org/2000/svg" width="25" height="35" fill="currentColor" class="bi bi-calendar2-x" viewBox="0 0 16 16">
  <path d="M6.146 8.146a.5.5 0 0 1 .708 0L8 9.293l1.146-1.147a.5.5 0 1 1 .708.708L8.707 10l1.147 1.146a.5.5 0 0 1-.708.708L8 10.707l-1.146 1.147a.5.5 0 0 1-.708-.708L7.293 10 6.146 8.854a.5.5 0 0 1 0-.708z"/>
  <path d="M3.5 0a.5.5 0 0 1 .5.5V1h8V.5a.5.5 0 0 1 1 0V1h1a2 2 0 0 1 2 2v11a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2h1V.5a.5.5 0 0 1 .5-.5zM2 2a1 1 0 0 0-1 1v11a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V3a1 1 0 0 0-1-1H2z"/>
  <path d="M2.5 4a.5.5 0 0 1 .5-.5h10a.5.5 0 0 1 .5.5v1a.5.5 0 0 1-.5.5H3a.5.5 0 0 1-.5-.5V4z"/>
</svg></span></span>
                    <span class="input-group-addon col-md-1"><span class="glyphicon glyphicon-th"><svg xmlns="http://www.w3.org/2000/svg" width="25" height="35" fill="currentColor" class="bi bi-calendar-date" viewBox="0 0 16 16">
  <path d="M6.445 11.688V6.354h-.633A12.6 12.6 0 0 0 4.5 7.16v.695c.375-.257.969-.62 1.258-.777h.012v4.61h.675zm1.188-1.305c.047.64.594 1.406 1.703 1.406 1.258 0 2-1.066 2-2.871 0-1.934-.781-2.668-1.953-2.668-.926 0-1.797.672-1.797 1.809 0 1.16.824 1.77 1.676 1.77.746 0 1.23-.376 1.383-.79h.027c-.004 1.316-.461 2.164-1.305 2.164-.664 0-1.008-.45-1.05-.82h-.684zm2.953-2.317c0 .696-.559 1.18-1.184 1.18-.601 0-1.144-.383-1.144-1.2 0-.823.582-1.21 1.168-1.21.633 0 1.16.398 1.16 1.23z"/>
  <path d="M3.5 0a.5.5 0 0 1 .5.5V1h8V.5a.5.5 0 0 1 1 0V1h1a2 2 0 0 1 2 2v11a2 2 0 0 1-2 2H2a2 2 0 0 1-2-2V3a2 2 0 0 1 2-2h1V.5a.5.5 0 0 1 .5-.5zM1 4v10a1 1 0 0 0 1 1h12a1 1 0 0 0 1-1V4H1z"/>
</svg></span></span>
                </div>
            </div>
            `)
        let now = new Date()
        let formatDate = date.formatDate(now.getTime())
        $(`#${this.Id}_date`).datetimepicker({
            format: 'yyyy-mm-dd',
            minView:2,
            weekStart: 0,
            todayBtn:  1,
            todayHighlight: 1,
            startView: 2,
            forceParse: 0,
            autoclose:1,
            startDate:formatDate,
        });
    }
    set Value(v){
         if (this.Schema["type"]!=="string"){
             if (v===0){
                 v =  ""
             } else {
                 v = date.formatDate(v)
             }
        }
        $(`#${this.Id}_date_data`).val(v)
    }
    get Value() {
        let v = $(`#${this.Id}_date_data`).val()
        if (this.Schema["type"]!=="string"){
           if (v === ""){
               return 0
           }
           return new Date(v).getTime()
        }
        return v
    }
}

// 简单map
class SimpleMapRender extends BaseChangeHandler {
    constructor(options) {
        super(options["path"]);
        this.MOptions = options
        let panel = options["panel"]
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
        myPanel.on("click", `#${this.Id}_add`, function () {
            let isEmpty = false

            for (let k in o.Values) {
                let it = o.Values[k];
                if (it.key.Value.length < 1) {
                    isEmpty = true
                    $(it.key.Target)[0].focus()
                    break
                }
            }

            if (!isEmpty) {
                o.add().key.Target[0].focus()
            }
        })
    }

    add(k, v) {
        let index = this.ItemCount++
        let $item = $(`<div class="input-group input-group-sm mt-2" data-index="${index}">
                        <div class="input-group-prepend">
                            <span class="input-group-text">key</span>
                        </div>
                        </div>`)
        this.$itemPanel.append($item)
        let $keyInput = $(createInput({"type": "string", "eo:type": "string"}, `${this.Id}.item_key_${index}`, true))
        $item.append($keyInput)
        $item.append(`<div class="input-group-prepend"><span class="input-group-text">value</span></div>`)

        let keyInput = new BaseValue({"type": "string", "eo:type": "string"}, $keyInput, `${this.Id}.item_key_${index}`)
        let $valueInput = $(createInput(this.Items, `${this.Id}.item_value_${index}`, true))
        $item.append($valueInput)
        let valueInput = new BaseValue(this.Items, $valueInput, `${this.Id}.item_value_${index}`)
        $item.append(`<div class="input-group-prepend"><button type="button" class="btn btn-danger" > - </button></div>`)
        let it = {key: keyInput, value: valueInput}
        this.Values[index] = it
        let o = this
        $($item).on("click", "button", function () {
            $item.remove()
            delete o.Values[index]
        })
        if (k) {
            keyInput.Value = k
        }
        if (v) {
            valueInput.Value = v
        }
        return it
    }

    get Value() {
        let v = {}
        for (let i in this.Values) {
            let it = this.Values[i]
            let key = it.key.Value
            if (key.length > 0) {
                v[key] = it.value.Value
            }
        }
        return v
    }

    set Value(v) {
        this.$itemPanel.empty()
        this.Values = {}
        for (let k in v) {
            this.add(k, v[k])
        }
    }
}

class ObjectArrayRender extends BaseChangeHandler {
    PopPanel(options, callbackFn, v) {
        let $Panel = $(`<div class="pop_window pop_window_small p-3" id="detail_container">
    <div class="pop_window_header">
        <span class="pop_window_title">${options["title"]}</span>
<!--        <button class="pop_window_button btn btn_default close" >关闭</button>-->
        <br>
    </div>
   
</div>`)

        let $Fade = $("<div class='modal-backdrop fade show modal-open'></div>")
        let $Body = $(` <div class="pop_window_body"></div>`)
        $("body").append($Fade)
        $("body").append($Panel)
        $Panel.append($Body)
        // 判断是否为时间戳

        let $Value = readGenerator(options)({
            schema: options["schema"],
            path: options["path"],
            // render: options["render"],
            panel: $Body,
            generator: readGenerator(options)
        })
        if (v) {
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
        let close = function () {
            $Fade.remove()
            $Panel.remove()
        }
        $Panel.show()
        $Fade.click(function () {
            close()

        })
        $Panel.on("click", "button.close", close)
        $Panel.on("click", "button.form_cancel", close)
        $Panel.on("click", "button.form_submit", function () {
            callbackFn($Value.Value)
            close()
        })
    }

    constructor(options) {
        super(options["path"])
        let panel = options["panel"]
        let schema = options["schema"]
        let O = this
        let Id = this.Id;
        let items = schema["items"]
        let p = $(panel)
        let $btn = $(`
<div id="${Id}_toolbar" class="m-1">
<button id="${Id}_AddButton" type="button" class="btn btn-secondary btn-sm">Add</button>
</div>`)
        p.append($btn)
        $btn.on("click", "button", function (event) {
            O.PopPanel({
                    path: `${options["path"]}.items`,
                    schema: items,
                    name: options["name"],
                    title: `添加 ${getLabel(options["name"], options["schema"])}`
                },
                function (v) {
                    $Table.bootstrapTable("append", [v])
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


        let columns = []
        columns.push({
            title: "index",
            field: "__index",
            sortable: false,

            formatter: function (v, row, index) {
                return index + 1
            }
        })

        columns.push(...configTable(uiSort,properties))
        columns.push({
            title: "操作",
            field: "",
            sortable: false,
            editable: false,
            formatter: function (v, row, index) {
                return `<a class="edit" href="javascript:void(0)" array-row="${index}" title="edit">配置</a> <a class="remove" href="javascript:void(0)" array-row="${index}" title="remove">删除</a>`
            },
            events: {
                "click .remove": function (e, value, row, index) {
                    // let rowIndex = $(this).attr("array-row")
                    $Table.bootstrapTable('remove', {
                        field: '$index',
                        values: [index]
                    })
                    O.onChange()
                },
                "click .edit": function (e, value, row, index) {
                    O.PopPanel({
                            path: `${options["path"]}.items`,
                            schema: items,
                            name: options["name"],
                            // render: options["render"],
                            title: `修改 ${getLabel(options["name"], options["schema"])}:${index + 1}`
                        },
                        function (v) {
                            $Table.bootstrapTable("updateRow", {
                                index: index,
                                row: v
                            })
                            O.onChange()
                        }, row)
                }
            }
        })
        let tableOptions = {
            columns: columns,
            width: "100%",
            useRowAttrFunc: true,
            reorderableRows: true,
            onReorderRow: function (data) {
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
        console.log("set data:",v)


        this.Table.bootstrapTable("load", v)
    }

    get Value() {
        let data = this.Table.bootstrapTable('getData')
        console.log("get data:",data)
        return data
    }

}
// output 用到到 formatter类型
class FormatterConfigRender extends BaseChangeHandler {
    PopPanel( callbackFn, keyHas, v) {
        let title = "增加配置"
        if (v){
            title = "修改配置"
        }
        let $Panel = $(`<div class="pop_window pop_window_small p-3" id="detail_container">
    <div class="pop_window_header">
        <span class="pop_window_title">${title}</span>
        <button class="pop_window_button btn btn_default close" >关闭</button>
  
    </div>
   
</div>`)

        let $Fade = $("<div class='modal-backdrop fade show modal-open'></div>")
        let $Body = $(` <div class="pop_window_body"></div>`)
        $("body").append($Fade)
        $("body").append($Panel)
        $Panel.append($Body)

        let $KeyInput = new FieldPanel("key", {
            schema: {"type": "string", "eo:type": "string", "pattern": "[a-zA-Z0-9]+[a-zA-Z0-9_]*"},
            required: true,
            path: `Formatter.key`,
            panel: $Body
        })
        let $Value = new FieldPanel("配置",{
            schema: {"type":"array","eo:type":"array","items":{"type":"string","eo:type":"string"}},
            path: "Formatter.items",
            panel: $Body,
        })
        if (v) {
            $KeyInput.Value = v["key"]
            $Value.Value = v["value"]
            if (v["key"] === "fields"){
                $KeyInput.$Value.Target.attr("readonly", true)
            }
        }
        $Panel.append(`<div class="row justify-content-between">
                <div class="col-4">
                    <button type="button" class="btn btn-outline-secondary form_cancel">取消</button>
                </div>
                <div class="col-4" style="text-align: right">
                    <button type="button" class="btn btn-primary form_submit">提交</button>
                </div>
            </div>`)
        let close = function () {
            $Fade.remove()
            $Panel.remove()
        }
        $Panel.show()
        $Fade.click(function () {
            close()
        })
        $Panel.on("click", "button.close", close)
        $Panel.on("click", "button.form_cancel", close)
        $Panel.on("click", "button.form_submit", function () {
            let key = $KeyInput.Value
            if (key.length < 1) {
                $KeyInput.$Value.Target.addClass("is-invalid")
                return
            }
            if (keyHas && keyHas(key) === true){
                $KeyInput.$Value.Target.addClass("is-invalid")
                return
            }
            let v = $Value.Value

            callbackFn({
                "key":key,
                "value":v
            })
            close()
        })
    }

    constructor(options) {
        super(options["path"])
        let panel = options["panel"]
        let O = this
        let Id = this.Id;

        let $btn = $(`
<div id="${Id}_toolbar" class="m-1">
 
  <button id="${Id}_AddButton" type="button" class="btn btn-secondary btn-sm">Add</button>
 
  <small>
  <a href="https://help.apinto.com/docs/formatter/" target="formatter" class="btn-link" >formatter 配置说明</a>
</small>
</div>`)
        $(panel).append($btn)
        $btn.on("click", "button", function (event) {
            O.PopPanel(
                function (v) {
                    $Table.bootstrapTable("append", [v])
                    $Table.bootstrapTable('scrollTo', 'bottom')
                    O.onChange()
                }, function (k) {
                    let vs = O.Value;
                    return typeof vs[k] !== "undefined"
                })
            return false
        })

        let $Table = $(`<table  id="${Id}_items"></table>`)
        $(panel).append($Table)

        this.Table = $Table

        let columns = []
        columns.push({
            title: "name",
            field: "key",
            sortable: false,
        })
        columns.push({
            title:"config",
            field:"value",
            sortable:false,
            formatter:function  (v, row, index){
                let html = '<ul>'
                for (let i in v){
                    let it = v[i]
                    html += `<li><span></span>${it}</li>`
                }
                html +='</ul>'
                return html
            }
        })

         columns.push({
            title: "操作",
            field: "",
            sortable: false,
            editable: false,
            formatter: function (v, row, index) {
                return `<a class="edit" href="javascript:void(0)"  title="edit">配置</a> <a class="remove" href="javascript:void(0)"  title="remove">删除</a>`
            },
            events: {
                "click .remove": function (e, value, row, index) {
                    // let rowIndex = $(this).attr("array-row")
                    $Table.bootstrapTable('remove', {
                        field: '$index',
                        values: [index]
                    })
                    O.onChange()
                },
                "click .edit": function (e, value, row, index) {
                    O.PopPanel(
                        function (v) {
                            $Table.bootstrapTable("updateRow", {
                                index: index,
                                row: v
                            })

                            O.onChange()
                            return true
                        }, function (k) {
                          return false
                        }
                        , row)
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
        for (let k in vs) {
            let v = vs[k]
            list.push({
                "key":k,
                "value":v
            } )
        }
        list.sort((a, b) => {
            if (a.key === "fields"){
                return -1
            }
            return a.key - b.key
        })
        this.Table.bootstrapTable("load", list)
    }

    get Value() {
        let list = this.Table.bootstrapTable('getData')
        let vs = {}
        for (let i in list) {
            let v = Object.assign({}, list[i])
            let key = v["key"]

            vs[key] = v["value"]
        }
        return vs
    }

}

// 结构体map
class ObjectMapRender extends BaseChangeHandler {

    PopPanelObjectMap(options, callbackFn, keyHas, v) {

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

        let $KeyInput = new FieldPanel("key", {
            schema: {"type": "string", "eo:type": "string", "pattern": "[a-zA-Z0-9]+[a-zA-Z0-9_]*"},
            required: true,
            path: `${options["path"]}.key`,
            panel: $Body
        })

        let $Value = readGenerator(options)({
            schema: options["schema"],
            path: options["path"],

            panel: $Body,
            generator: readGenerator(options)
        })
        if (v) {
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
        let close = function () {
            $Fade.remove()
            $Panel.remove()
        }
        $Panel.show()
        $Fade.click(function () {
            close()
        })
        $Panel.on("click", "button.close", close)
        $Panel.on("click", "button.form_cancel", close)
        $Panel.on("click", "button.form_submit", function () {
            let key = $KeyInput.Value
            if (key.length < 1) {
                $KeyInput.$Value.Target.addClass("is-invalid")
                return
            }
            if (keyHas && keyHas(key) === true){
                $KeyInput.$Value.Target.addClass("is-invalid")
                return
            }
            let v = $Value.Value
            v["__key"] = key
            callbackFn(v)
            close()
        })
    }

    constructor(options) {
        super(options["path"])
        let panel = options["panel"]
        let schema = options["schema"]
        let O = this
        let Id = this.Id;
        let items = schema["additionalProperties"]

        let $btn = $(`
<div id="${Id}_toolbar" class="m-1">
<button id="${Id}_AddButton" type="button" class="btn btn-secondary btn-sm">Add</button>
</div>`)
        $(panel).append($btn)
        $btn.on("click", "button", function (event) {
            O.PopPanelObjectMap({
                    path: `${options["path"]}.items`,
                    schema: items,
                    name: options["name"],
                    title: `添加 ${getLabel(options["name"], options["schema"])}`
                },
                function (v) {
                    $Table.bootstrapTable("append", [v])
                    $Table.bootstrapTable('scrollTo', 'bottom')
                    O.onChange()
                }, function (k) {
                    let vs = O.Value;
                    return vs.hasOwnProperty(k)
                })
            return false
        })

        let $Table = $(`<table  id="${Id}_items"></table>`)
        $(panel).append($Table)

        this.Table = $Table
        let properties = items["properties"]
        let uiSort = items["ui:sort"]

        let columns = []
        columns.push({
            title: "key",
            field: "key",
            sortable: false,
        })
        columns.push(...configTable(uiSort,properties))

        columns.push({
            title: "操作",
            field: "",
            sortable: false,
            editable: false,
            formatter: function (v, row, index) {
                return `<a class="edit" href="javascript:void(0)"  title="edit">配置</a> <a class="remove" href="javascript:void(0)"  title="remove">删除</a>`
            },
            events: {
                "click .remove": function (e, value, row, index) {
                    // let rowIndex = $(this).attr("array-row")
                    $Table.bootstrapTable('remove', {
                        field: '$index',
                        values: [index]
                    })
                    O.onChange()
                },
                "click .edit": function (e, value, row, index) {
                    O.PopPanelObjectMap({
                            path: `${options["path"]}.items`,
                            schema: items,
                            name: options["name"],
                            title: `修改 ${getLabel(options["name"], options["schema"])}:${index + 1}`
                        },
                        function (v) {
                            $Table.bootstrapTable("updateRow", {
                                index: index,
                                row: v
                            })

                            O.onChange()
                            return true
                        }, function (k) {
                            let vs = O.Value;
                            return vs.hasOwnProperty(k)
                        }
                        , row)
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
        for (let k in vs) {
            let v = vs[k]
            v["__key"] = k
            list.push(v)
        }
        list.sort((a, b) => {
            return a.__key - b.__key
        })
        this.Table.bootstrapTable("load", list)
    }

    get Value() {
        let list = this.Table.bootstrapTable('getData')
        let vs = {}
        for (let i in list) {
            let v = Object.assign({}, list[i])
            let key = v["__key"]
            delete v["__key"]
            vs[key] = v
        }
        return vs
    }

}

// 带枚举的基础数组
class ArrayRenderEnum extends BaseChangeHandler {
    constructor(options) {
        super(options["path"])
        let panel = options["panel"]
        let schema = options["schema"]
        let Id = this.Id
        let items = schema["items"]
        this.Enum = items["enum"]

        let p = $(panel);
        let itemPanel = `<div id="${Id}_items" class="border p-sm-1 btn-toolbar form-control" role="toolbar">`

        for (let i in items["enum"]) {
            let e = items["enum"][i]
            let itemId = `${Id}_${e}`
            itemPanel += `<div class="custom-control custom-checkbox custom-control-inline"><input type="checkbox" id="${itemId}" value="${e}" name="${Id}" class="custom-control-input"> <label class="custom-control-label" for="${itemId}">${e}</label></div>`
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
class ArrayRenderSimple extends BaseChangeHandler {
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
        let $itemInput = $(createInput(this.SchemaItems, itemId, true))

        $item.prepend($itemInput)
        this.ValuesList[index] = new BaseValue(this.SchemaItems, $itemInput, this.Id)
        let o = this
        this.ValuesList[index].onChange(function () {
            o.onChange()
        })
        $item.on("click", "button", function () {
            $item.remove()
           delete o.ValuesList[index]
        })

        return $itemInput
    }

    get Value() {
        let val = []
        let indexList = Object.keys(this.ValuesList)
        indexList.sort()
        for (let i in indexList) {
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

        for (let i in vs) {
            this.add().val(vs[i])
        }
    }
}

// 结构体
class ObjectRender extends BaseChangeHandler {
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

            let field = new FieldPanel(name, {
                panel: this.$Panel,
                schema: sub,
                generator: readGenerator(options),
                path: subId,
                required: requiredData[name] === true,
                // render : options["render"],
                parent: schema
            })
            this.Fields[name] = field
            field.onChange(function () {
                o.switchPanel(name, this.Value)
                o.onChange()
            })
        }

        let Switches = {}
        for (let name in properties) {
            let sub = properties[name]
            if (sub["switch"]) {
                Switches[name] = sub["switch"]
            }
        }
        this.Switches = Switches
        for (let name in this.Fields) {
            this.switchPanel(name, this.Fields[name].Value)
        }

    }

    switchPanel(name, value) {

        switch (typeof value) {
            case "string": {
                value = `"${value}"`
                break
            }
            case "object":
            case "undefined": {
                return
            }
            default: {
                break
            }
        }

        for (let f in this.Switches) {
            let expression = this.Switches[f]
            try {
                let funcStr = `(function(${name}){return ${expression}})(${value})`
                // console.log(funcStr)
                let r = eval(funcStr)
                if (r === true) {
                    this.Fields[f].Show()
                } else {
                    this.Fields[f].Hide()
                }
            } catch (e) {

                // console.log(e)
            }
        }
    }

    get Value() {
        let vs = {}
        for (let k in this.Fields) {
            let fi = this.Fields[k]
            let v = fi.Value
            if (typeof v !== "undefined") {
                vs[k] = v
            }
        }
        return vs
    }

    set Value(v) {
        if (!v || typeof v === "undefined") {
            v = {}
        }
        for (let k in this.Fields) {
            let fi = this.Fields[k]
            fi.Value = v[k]
            this.switchPanel(k, v[k])
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
    getRenderInfo(id, success) {
        let path = id.replaceAll(":", "/")
        let n = id.replaceAll(":", "_")
        dashboard.getExtenderInfo(path, function (res) {
            if (res.code !== 200) {
                return http.handleError(res, "获取extender信息失败")
            }
            success(n, res.data["render"])
        }, function (res) {
            return http.handleError(res, "获取extender信息失败")
        })
    }

    show(exclude, callbackFn, data) {
        if (this.IsInit === true) {
            this.render(exclude, callbackFn, data)
            return
        }
        let o = this
        o.IsInit = true
        dashboard.get(`/api/plugins`, function (pluginsData) {
            o.IsInit = true
            if (pluginsData.data["plugins"]){
                o.Plugins =pluginsData.data["plugins"]
            }else {
                o.Plugins =[]
            }

            o.PluginsExtenders = new Map(
                o.Plugins.map(object => {
                    return [object.name, object.id];
                }),
            );
            o.render(exclude, callbackFn, data)
        }, function (err) {
            o.IsInit = false
            common.message(err, "danger")
        })
    }

    render(exclude, callbackFn, data) {

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

        let $PluginName = $(`<select class="form-control form-control-sm"  ></select>`)

        let $Disable = $(`<input type="checkbox" class="form-control-sm form-control" data-toggle="toggle" data-size="sm"/>`)

        $Body.append(
            $(`<div><div><label class=" col-form-label  text-nowrap" for="FormRender_test_target"><span style="color: red">*</span>选择插件</label></div></div>`)
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
        let O = this

        function renderConfig(name) {
            O.getRenderInfo(O.PluginsExtenders.get(name), function (driver, render) {
                $Config.empty()
                O.ConfigTarget = new ObjectRender({
                    path: `plugins.${name}`,
                    panel: $Config,
                    schema: render
                })
                if (data) {
                    O.ConfigTarget.Value = data.config
                }
            })
        }

        if (data) {
            if (data.disable === true || data.disable === "true") {
                $Disable.bootstrapToggle("on");
            } else {
                $Disable.bootstrapToggle("off");
            }
            $PluginName.prepend(`<option value="${data.name}" selected>${data.name}</option>`)
            renderConfig(data.name)
        } else {
            $Disable.bootstrapToggle("off");
            let set = new Set(exclude)
            let first = true
            for (let i in this.Plugins) {
                let p = this.Plugins[i]
                if (!set.has(p.name)) {
                    $PluginName.append(`<option value="${p.name}" ${first ? "selected" : ""}>${p.name}</option>`)
                    first = false
                }
            }
            if ($PluginName.children().length > 0) {
                renderConfig($PluginName.val())
            }

            $PluginName.on("change", function () {
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
        let close = function () {
            $Fade.remove()
            $Panel.remove()
        }
        $Panel.show()
        $Fade.click(function () {
            close()
        })
        $Panel.on("click", "button.close", close)
        $Panel.on("click", "button.form_cancel", close)
        $Panel.on("click", "button.form_submit", function () {
            callbackFn({name: $PluginName.val(), disable: $($Disable).prop("checked"), config: O.ConfigTarget.Value})
            close()
        })
    }
}

// 插件配置表
class PluginsRender extends BaseChangeHandler {
    constructor(options) {
        super(options["path"])
        let panel = options["panel"]
        let schema = options["schema"]
        let items = schema["additionalProperties"]
        let Id = this.Id
        this.PopPanel = new PopPanelPlugin()
        let $Table = $(`<table  id="${Id}_items"></table>`)
        $(panel).append($Table)

        this.Table = $Table
        let O = this

        let tableOptions = {
            columns: [{
                title: "插件名",
                field: "name",
                sortable: false,
            }
                , {
                    title: "状态",
                    field: "disable",
                    formatter: function (v, row, index) {
                        return v === true ? `<span class="badge badge-danger">禁用</span>` : `<span class="badge-primary"">启用</span>`
                    }
                },
                {
                    title: "配置",
                    field: "config",
                    formatter: JSON.stringify
                },
                {
                    title: "操作",
                    field: "",
                    sortable: false,
                    editable: false,
                    formatter: function (v, row, index) {
                        return `<a class="edit" href="javascript:void(0)"  title="edit">配置</a> <a class="remove" href="javascript:void(0)"  title="remove">删除</a>`
                    },
                    events: {
                        "click .remove": function (e, value, row, index) {
                            $Table.bootstrapTable('remove', {
                                field: '$index',
                                values: [index]
                            })
                            O.onChange()
                        },
                        "click .edit": function (e, value, row, index) {
                            O.PopPanel.show(Object.keys(O.Value),
                                function (v) {
                                    $Table.bootstrapTable("updateRow", {
                                        index: index,
                                        row: v
                                    })

                                    O.onChange()
                                    return true
                                }, row)
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
        $btn.on("click", "button", function (event) {
            O.PopPanel.show(Object.keys(O.Value),
                function (v) {
                    $Table.bootstrapTable("append", [v])
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
        for (let k in vs) {
            let v = vs[k]
            v["name"] = k
            list.push(v)
        }
        list.sort((a, b) => {
            return a.name - b.name
        })
        this.Table.bootstrapTable("load", list)
    }

    get Value() {
        let list = this.Table.bootstrapTable('getData')
        let vs = {}
        for (let i in list) {

            let v = list[i]
            let key = v["name"]
            vs[key] = {
                disable: v["disable"],
                config: v["config"],
            }
        }
        return vs
    }

}

// 弹出额外render窗
class PopPanelAdditionRender {
    constructor() {

    }

    show(exclude, callbackFn, data) {
        if (this.IsInit === true) {
            this.render(exclude, callbackFn, data)
            return
        }
        let o = this
        o.IsInit = true
        dashboard.additionRender(`auth`, function (res) {
            o.IsInit = true
            if (res.data){
                o.Renders = res.data
            }else {
                o.Renders =[]
            }

            o.AdditionalRenders = new Map(
                o.Renders.map(object => {
                    return [object.name, object.render];
                }),
            );
            o.render(exclude, callbackFn, data)
        }, function (err) {
            o.IsInit = false
            common.message(err, "danger")
        })
    }

    render(exclude, callbackFn, data) {

        let $Panel = $(`
<div class="pop_window pop_window_small p-3" id="detail_container">
    <div class="pop_window_header">
        <span class="pop_window_title">配置信息</span>
        <button class="pop_window_button btn btn_default close" >关闭</button>
    </div>
</div>`)

        let $Fade = $("<div class='modal-backdrop fade show modal-open'></div>")
        let $Body = $(` <div class="pop_window_body"></div>`)
        $("body").append($Fade)
        $("body").append($Panel)
        $Panel.append($Body)

        let $DriverName = $(`<select class="form-control form-control-sm"  ></select>`)

        $Body.append(
            $(`<div><div><label class=" col-form-label  text-nowrap" for="FormRender_test_target"><span style="color: red">*</span>选择驱动</label></div></div>`)
                .append($DriverName)
        )

        let $Config = $(`<div class="border-top p-1"></div>`)
        $Body.append(
            $(`<div><div><label class=" col-form-label  text-nowrap" for="FormRender_test_target"><span style="color: red">*</span>配置</label></div></div>`)
                .append($Config)
        )
        let O = this
        function renderConfig(name) {
            let render = O.AdditionalRenders.get(name)
            if (render) {
                $Config.empty()
                O.ConfigTarget = new ObjectRender({
                    path: `auth.${name}`,
                    panel: $Config,
                    schema: render,
                })
                if (data) {
                    O.ConfigTarget.Value = data
                }
            }
        }

        if (data) {
            $DriverName.prepend(`<option value="${data.type}" selected>${data.type}</option>`)
            renderConfig(data.type)
        } else {
            let set = new Set(exclude)
            let first = true
            for (let i in this.Renders) {
                let p = this.Renders[i]
                if (!set.has(p.name)) {
                    $DriverName.append(`<option value="${p.name}" ${first ? "selected" : ""}>${p.name}</option>`)
                    first = false
                }
            }
            if ($DriverName.children().length > 0) {
                renderConfig($DriverName.val())
            }

            $DriverName.on("change", function () {
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
        let close = function () {
            $Fade.remove()
            $Panel.remove()
        }
        $Panel.show()
        $Fade.click(function () {
            close()
        })
        $Panel.on("click", "button.close", close)
        $Panel.on("click", "button.form_cancel", close)
        $Panel.on("click", "button.form_submit", function () {
            let v =  O.ConfigTarget.Value
            v['type'] = $DriverName.val()
            callbackFn(v)
            close()
        })
    }
}

// 额外Render
class AdditionRender extends BaseChangeHandler {
    constructor(options) {
        super(options["path"])
        let panel = options["panel"]
        let schema = options["schema"]
        let items = schema["items"]
        let Id = this.Id
        this.PopPanel = new PopPanelAdditionRender()
        let $Table = $(`<table  id="${Id}_items"></table>`)
        $(panel).append($Table)

        this.Table = $Table
        let O = this

        let properties = items["properties"]
        let uiSort = items["ui:sort"]

        let columns = []
        columns.push({
            title: "index",
            field: "__index",
            sortable: false,

            formatter: function (v, row, index) {
                return index + 1
            }
        })

        columns.push(...configTable(uiSort,properties))
        columns.push({
            title: "操作",
            field: "",
            sortable: false,
            editable: false,
            formatter: function (v, row, index) {
                return `<a class="edit" href="javascript:void(0)" array-row="${index}" title="edit">配置</a> <a class="remove" href="javascript:void(0)" array-row="${index}" title="remove">删除</a>`
            },
            events: {
                "click .remove": function (e, value, row, index) {
                    // let rowIndex = $(this).attr("array-row")
                    $Table.bootstrapTable('remove', {
                        field: '$index',
                        values: [index]
                    })
                    this.onChange()
                },
                "click .edit": function (e, value, row, index) {
                    O.PopPanel.show(Object.keys(O.Value),
                        function (v) {
                            $Table.bootstrapTable("updateRow", {
                                index: index,
                                row: v
                            })

                            O.onChange()
                            return true
                        }, row)
                }
            }
        })
        let tableOptions = {
            columns: columns,
            width: "100%",
            useRowAttrFunc: true,
            reorderableRows: true,
            onReorderRow: function (data) {
                $Table.bootstrapTable("refresh")
            }

        }

        $Table.bootstrapTable(tableOptions);



        let $btn = $(`
<div id="${Id}_toolbar" class="m-1">
<button id="${Id}_AddButton" type="button" class="btn btn-secondary btn-sm">Add</button>
</div>`)
        $(panel).prepend($btn)
        $btn.on("click", "button", function (event) {
            O.PopPanel.show(Object.keys(O.Value),
                function (v) {
                    $Table.bootstrapTable("append", [v])
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
        for (let k in vs) {
            let v = vs[k]
            v["name"] = k
            list.push(v)
        }
        list.sort((a, b) => {
            return a.name - b.name
        })
        this.Table.bootstrapTable("load", list)
    }

    get Value() {
        let list = this.Table.bootstrapTable('getData')
        // let vs = {}
        // for (let i in list) {
        //
        //     let v = list[i]
        //     vs[v["name"]] = v
        // }
        return list
    }

}

function BaseGenerator(options) {

    let schema = options["schema"]
    let format = schema["format"]
    if ( typeof format != "undefined" ) {
        let dataFormat=new Set(["date-time","time","date","duration"])
        if (dataFormat.has(format)){
            options["schema"]["eo:type"] = "date-time"
            options["schema"]["eo:format"] = format
            if( schema["type"] !== "string"){
                delete(options["schema"]["format"])
            }
        }
    }

    switch (schema["eo:type"]) {
        case "object": {
            return new ObjectRender(options);
        }
        case "eofiles":{
            return new FileRender(options);
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
                case "array": {
                    throw `not allow type:${items["eo:type"]} in array`
                }
                case "require": {
                    return new RequireArrayRender(options)
                }
                default: {
                    throw `unknown type:${items["eo:type"]} in array`
                }
            }

        }
        case "map": {
            let item = schema["additionalProperties"];
            switch (item["type"]) {
                case "object": {
                    if (options["name"] === "plugins") {
                        return new PluginsRender(options)
                    }
                    return new ObjectMapRender(options)
                }
            }
            return new SimpleMapRender(options)
        }
        case "interface": {
            let typ = schema["type"]
            switch (typ){
                case "array":{
                    let item = schema["items"]
                    switch (item["type"]) {
                        case "object": {
                            return new AdditionRender(options)
                        }
                    }
                }
            }
            switch (item["type"]) {
                case "object": {
                    if (options["name"] === "plugins") {
                        return new PluginsRender(options)
                    }
                    return new ObjectMapRender(options)
                }
            }
        }
        case "boolean": {
            return new SwitchRender(options)
        }
        case "integer": {
        }
        case "number": {
        }
        case "text":{
        }
        case "string": {
            if (schema["enum"]) {
                return new BaseEnumRender(options)
            }
            switch (schema["format"] ) {
                case "text":{
                    return new BaseInputTextRender(options)
                }
            }
            return new BaseInputRender(options)
        }
        case "require": {
            return new RequireRender(options)
        }
        case "date-time":{
            return new DatetimeRender(options)
        }
        case "formatter" : {
            return new FormatterConfigRender(options)
        }
    }
    throw `unknown type:${schema["eo:type"]}`
}

class SchemaHandler {
    constructor(schema) {
        this.s = schema
    }

    toJsonSchema(s) {
        let schema = Object.assign({}, s)
        switch (schema["eo:type"]) {
            case "map": {
                schema["additionalProperties"] = this.toJsonSchema(schema["additionalProperties"])
                break
            }
            case "array": {
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

    get JsonSchema() {
        if (typeof this.v === "undefined" || this.v === null) {
            this.v = this.toJsonSchema(this.s)
            console.log(JSON.stringify(this.v))
        }
        return this.v
    }

    get Schema() {
        return this.s
    }
}

class FormRender {

    constructor(options) {

        this.ObjectName = `${RootId}.${options["name"]}`
        this.JsonSchema = new SchemaHandler(options["schema"])

        $(options["panel"]).empty()

        this.Object = readGenerator(options)({
            generator: options["generator"],
            panel: options["panel"],
            schema: options["schema"],
            path: this.ObjectName,
            name: options["name"],
        })

    }

    check() {

        let r = CheckBySchema(this.ObjectName, this.JsonSchema.JsonSchema, this.Value)
        console.log(`check:${this.ObjectName} = ${JSON.stringify(this.Value)} :${JSON.stringify(r)}`);
        if (typeof r === "undefined") {
            return true
        }

        return JSON.stringify(r)

    }

    get Value() {
        return this.Object.Value
    }

    set Value(v) {
        this.Object.Value = v
    }
}

