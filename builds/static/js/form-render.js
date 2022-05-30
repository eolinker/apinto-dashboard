const RootId = "FormRender"

function CheckBySchema(id, schema, value) {
    let validator = validate.Default()
    if (!validator.resolved.hasOwnProperty(id)) {
        validator.addSchema(id, schema);
    }
    let err = validator.validate(id, value)
    if (err) {
        console.log(err)
        return false
    }
    return true
}

function ValidHandler(v,schema,id) {
    console.debug("ValidHandler:",id,"=",v)
    let value = v
    value = valueForType(schema["type"], value)

    let rs = CheckBySchema(id, schema, value)
    if (rs === true) {
        $(this).removeClass("is-invalid")
        $(this).addClass("is-valid")
        return true
    } else {
        $(this).removeClass("is-valid")
        $(this).addClass("is-invalid")
        return false
    }
}

function InputValid(schema, Id) {
    $("#" + Id).on("change", function () {
        ValidHandler($(this).val(), schema,Id)
    })
}

function valueForType(t, v) {
    switch (t) {
        case "integer", "number": {
            return Number(v)
        }
        case "boolean": {
            if (typeof v === "undefined") {
                return false
            }
            switch (String(v).toLowerCase()) {
                case "on", "true": {
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

class BaseValue {
    constructor(schema, path) {
        this.Schema = schema
        this.Id = path
    }

    get Value() {
        let val = $("#" + this.Id).val()
        return valueForType(this.Schema["type"], val)
    }

    set Value(v) {

        if ( typeof v === "undefined"){
            if (typeof this.Schema["default"] !== "undefined") {
               v = this.Schema["default"]
            }
        }
        $("#" + this.Id).val(v)
    }

}

class BaseEnumRender extends BaseValue {
    constructor(panel, schema, path) {
        super(schema, path)
        $(panel).append(createEnum(schema, this.Id))
    }

}

class SwitchRender extends BaseValue {
    constructor(panel, schema, path) {
        super(schema, path)

        const Id = path
        this.Id = Id
        let $switch = $(`<input id="${Id}" type="checkbox" data-toggle="toggle" data-size="sm"/>`)

        $(panel).append($switch)
        $switch.bootstrapToggle()
        InputValid(schema, Id)
    }

    get Value() {
       return  document.getElementById(this.Id).checked
    }

    set Value(v) {
        if (v === true || v === "true") {
            $("#" + this.Id).attr('checked', 'checked');
        } else {
            $("#" + this.Id).removeAttr('checked');
        }
    }
}

class RequireRender extends BaseValue {
    constructor(panel, schema, path) {
        super(schema, path)
        this.DOM = $(`<select id=${path} class="form-controller"></select>`)

        $(panel).append(this.DOM)
    }
}

class BaseInputRender extends BaseValue {
    constructor(panel, schema, path) {
        super(schema, path)
        const Id = path


        $(panel).append(createInput(schema, Id))

        InputValid(schema, Id)
        if (typeof schema["default"] !== "undefined") {
            this.Value = schema["default"]
        }
    }
}

function getLabel(schema) {
    let label = schema["label"]
    if (!label || label.trim() === "") {
        label = schema["name"]
    }
    label = label.replace(label[0], label[0].toUpperCase());
    return label
}

function createLabel(schema, id, appendAttr) {
    if (!appendAttr) {
        appendAttr = ""
    }

    let labelHtml = '<label class="col-sm-2 col-form-label text-right"'
    if (id) {
        labelHtml += ` for="${id}"`
    }
    if (appendAttr) {
        labelHtml += ` ${appendAttr}`
    }
    labelHtml += `>`
    if (schema["required"]) {
        labelHtml += '<span style="color: red">*</span>'
    }
    labelHtml += getLabel(schema) + '</label>'
    return labelHtml
}

function createFieldPanel(panel, schema, id, appendAttr) {
    $(panel).append('<div class="form-group row mb-3"></div>')
    let fieldPanel = $(panel).children().last()
    fieldPanel.append(createLabel(schema, id, appendAttr))
    fieldPanel.append('<div class="col-sm-10"></div>')
    return fieldPanel.children().last()
}
class FieldPanel {
    constructor(panel, schema, generator,id) {
        const $FieldValuePanel = $(`<div class="col-sm-10"></div>`)
        const $FieldPanel = $(`<div class="form-group row mb-3">${createLabel(schema,id)}</div>`)
        $FieldPanel.append($FieldValuePanel)
        panel.append($FieldPanel)
        this.$Value = generator($FieldValuePanel,schema,generator,id)
    }
    set Value(v){
        this.$Value.Value = v
    }
    get Value(){
        return this.$Value.Value
    }

}
function createEnum(schema, id, appendAttr) {
    let readOnly = ""
    this.schema = schema
    if (schema["readonly"] === true) {
        readOnly = "readonly"
    }

    if (schema["enum"]) {
        let enums = schema["enum"]

        let select = `<select ${readOnly} ${appendAttr} class="form-control" id="${id}"`
        if (schema["required"]) {
            select += ` required`
        }
        select += ">"

        for (let i in enums) {
            if (typeof value != "undefined" && value === enums[i]) {
                select += '<option>' + enums[i] + '</option>'
            } else {
                select += '<option selected>' + enums[i] + '</option>'
            }
        }
        select += '</select>'

        return select
    }
}

function createInput(schema, id, appendAttr) {
    let readOnly = ""
    this.schema = schema
    if (schema["readonly"] === true) {
        readOnly = "readonly"
    }

    let input = `<input ${readOnly} class="form-control" id="${id}" aria-describedby="validation_${id}" `;
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

    switch (schema["type"]) {
        case "string": {
            input += ' type="' + readFormatForString(schema["format"]) + '"'
            if (schema["format"] === 'password') {
                input += ' autocomplete="on"'
            }
            break
        }
        case "integer": {
        }
        case "number": {
            input += ' type="number"'
            break
        }
        default: {
            throw `now allow [${schema["type"]} for input]`
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

class MapRender extends BaseValue {

    constructor(panel, schema, generator, path) {
        super(schema, path)

        this.Schema = schema
        this.Id = path;
        this.$Panel =$(`<div id='${path}_panel' class='container-fluid'></div>`)
        $(panel).append(this.$Panel)

        this.GeneratorHandler = generator

    }

    set Value(v) {

        const Items = this.Schema["items"]
        const Id = this.Id
        const keySchema = {type: "string"}

        this.$Panel.empty()
        switch (Items["type"]) {
            case "object", "map": {
                break
            }
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
                    $(this.PanelId).append('<div class="input-group input-group-sm m-2"></div>')
                    let itemPanel = $(this.PanelId).children().last()
                    itemPanel.append(`
<div class="input-group-prepend">
    <button class="btn btn-danger" id="btnGroupAddon_${Id}_key_${k}" type="button" data-itemId="${Id}_key"> - </button>
</div>`)
                    let childItemKey = new BaseInputRender(itemPanel, keySchema, `${Id}_key_${k}`)
                    childItemKey.Value = k
                    itemPanel.append('<div class="input-group-prepend"><div class="input-group-text  btn" >=</div></div>')
                    let childItemValue = new BaseInputRender(itemPanel, Items, Id + "_value_" + k)
                    childItemValue.Value = v[k]
                }
                break
            }
        }

    }

    get Value() {
        return {}
    }

}

class InnerObjectRender{
    constructor(panel, schema, generator, path) {

        const Id = path;
        this.Id = Id;
        this.Schema = schema
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
        let lastDetailRow = undefined
        let lastField = undefined

        function DetailFormatterHandler(fieldIndex) {

            this.detailFormatter = function (index, row, $element) {

                if (typeof lastDetailRow !== "undefined" && lastDetailRow !== index) {
                    $Table.bootstrapTable('collapseRow', lastDetailRow)
                }
                lastDetailRow = index
                lastField = fieldIndex
                let item = properties[fieldIndex]
                let child = generator($element, item, generator, path + "_" + item["name"])
                child.Value = row[item.name]
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

        for (let i in properties) {
            const item = properties[i]
            switch (item["type"]) {
                case "object": {
                    columns.push({
                        title: getLabel(item),
                        field: item["name"],
                        sortable: false,
                        editable: false,
                        formatter: formatterKV,
                        detailFormatter: new DetailFormatterHandler(i).detailFormatter
                    })
                    break
                }
                case "map": {
                    columns.push({
                        title: getLabel(item),
                        field: item["name"],
                        sortable: false,
                        editable: false,
                        formatter: formatterKV,
                        detailFormatter: new DetailFormatterHandler(i).detailFormatter
                    })
                    break
                }
                case "array": {
                    columns.push({
                        title: getLabel(item),
                        field: item["name"],
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
                            title: getLabel(item),
                            field: item["name"],
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
                        if (item["type"] === "number" || item["type"] === "integer") {
                            typeInput = "number"
                        }
                        columns.push({

                            title: getLabel(item),
                            field: item["name"],
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

class InnerMapRender extends BaseValue{
    constructor(panel, schema, generator, path) {
        super(schema,path)
    }

    set Value(v) {

    }

    get Value() {
        return {}
    }
}

class ArrayRenderEnum {
    constructor(panel, schema, generator, path) {
       this.Schema = schema
        const Id = path;
        this.Id = Id;
        const items = schema["items"]
        this.Enum = items["enum"]

        let p = $(panel);
        let itemPanel = `<div id="${ Id }_items" class="border p-sm-1 btn-toolbar " role="toolbar">`

        for (let i in items["enum"]) {
            let e = items["enum"][i]
            let itemId = Id + '_' + e
            itemPanel +=`<div class="custom-control custom-checkbox custom-control-inline"><input type="checkbox" id="${itemId}" value="${e}" name="${Id}" class="custom-control-input"> <label class="custom-control-label" for="${itemId}">${e}</label></div>`
        }
        itemPanel += '</div>'
        p.append(itemPanel)
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

class ArrayRenderSimple {
    constructor(panel, schema, generator, path) {

        this.Schema = schema
        const items = schema["items"]
        const Id = path;
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
            if (items["type"] === "integer" || items["type"] === "number") {
                v = Number(value)
            }
            if (v !== "") {
                if (CheckBySchema(Id, items, v)) {
                    add(v)
                    $(this).val("")
                }
            }
            return false
        })

        let lastIndex = 0

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
            return false
        }

        this.add = add
        $itemPanel.delegate('button', "click", function () {
            let itemId = $(this).attr('data-itemId')
            $itemPanel.children(`#array-item_${itemId}`).remove()
        })

    }


    get Value() {
        let val = []
        $("[array-for='" + this.Id + "']").each(function () {
            val.push($(this).val())
        })
        return val
    }

    set Value(vs) {
        if (!Array.isArray(vs)) {
            return
        }
        const arrayId = "[array-for='" + this.Id + "']"
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

class ObjectRender {
    constructor(panel, schema, generator, path) {
        this.Fields = {}
        this.Id = path
        if (schema["type"] !== "object") {
            return
        }
        let properties = schema["properties"]
        for (let i in properties) {
            let sub = properties[i]
            let name = sub["name"]
            const subId = path + "." + name
            // const fieldPanel = createFieldPanel(panel, sub, subId)
            this.Fields[name] = new FieldPanel(panel, sub, generator, subId)
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


function BaseGenerator(panel, schema, generator, path) {
    switch (schema["type"]) {
        case "object": {
            return new ObjectRender(panel, schema, generator, path)
        }
        case "array": {
            const items = schema["items"]
            switch (items["type"]) {
                case "object": {
                    return new InnerObjectRender(panel, schema, generator, path)
                }
                case "map": {
                    return new InnerMapRender(panel, schema, generator, path)
                }
                case "integer": {
                }
                case "number": {
                }
                case "string": {
                    if (items["enum"]) {
                        return new ArrayRenderEnum(panel, schema, generator, path)
                    }
                    return new ArrayRenderSimple(panel, schema, generator, path)
                }
            }
            throw `not allow type:${schema["type"]} in array`
        }
        case "map": {
            return new MapRender(panel, schema, generator, path)
        }
        case "boolean": {
            return new SwitchRender(panel, schema, path)
        }
        case "integer": {
        }
        case "number": {
        }
        case "string": {
            if (schema["enum"]) {
                return new BaseEnumRender(panel, schema, path)
            }
            return new BaseInputRender(panel, schema, path)
        }
        case "require": {
            return new RequireRender(panel, schema, path)
        }
    }
    throw `unknown type:${schema["type"]}`
}

class FormRender {
    toJsonSchema(schema){
        switch (schema["type"]){
            case "map":{
                schema["type"] = "object"
                schema["additionalProperties"]=this.toJsonSchema(schema["items"])
                schema.delete("items")
                break
            }
            case "object":{
                let properties = schema["properties"]
                schema["properties"] = {
                }
                 for (let i in properties){
                    let n = this.toJsonSchema(properties[i])
                     schema.properties[n.name] = n
                }
            }
        }

        return schema
    }
    constructor(panel, schema, generator,name) {
        this.Schema = schema
        if (!generator || typeof generator !== "function") {
            generator = BaseGenerator
        }
        $(panel).empty()
        this.ObjectName = `${RootId}.${name}`
        this.Object = generator($(panel), schema, generator, this.ObjectName)
    }

    check() {
        if (typeof this.JsonSchema === undefined || this.JsonSchema === null){
            this.JsonSchema = this.toJsonSchema(this.Schema.cloneNode(true))
        }
        // if (this.djv === null || typeof this.djv === "undefined"){
        //     this.djv = validate.Create()
        // }
        return  CheckBySchema(this.ObjectName,this.JsonSchema,this.Object.Value)

    }

    get Value() {
        return this.Object.Value
    }

    set Value(v) {
        this.Object.Value = v
    }
}

