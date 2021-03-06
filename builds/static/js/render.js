const nameRule = /^[a-zA-Z\d_]+$/;

class Render {
    constructor(panel, schema, name, data,  callback) {
        const options = {
            panel: $(panel),
            schema: schema,
            name: name,
        }
        let target = $(panel)
        this.InitValue = data

        let renderHandler = new FormRender(options)

        if (data) {
            renderHandler.Value = data
        }
        console.log("render data",data)
        if (callback) {
            callback()
        }
        this.panel = renderHandler
        this.target = target

    }

    Reset(schema, name) {
        this.Destroy()
        this.panel = new FormRender({
            panel: this.target,
            schema: schema,
            name: name
        })
        this.InitValue = {}
    }

    ResetVal() {

            this.panel.Value = this.InitValue

    }

    set Value(data) {

        this.InitValue = data
        this.panel.Value = data
    }

    get Value() {
        return this.panel.Value
    }

    Check() {
        return this.panel.check()
    }

    Submit(success, error) {
        let ckr  = this.Check()
        if (ckr  === true) {
            success(this.Value)
        } else {
            error(ckr)
        }
    }

    Destroy() {
        this.panel = null
    }
}

class ProfessionRender {
    constructor(module, options) {
        this.module = module
        this.options = options
        this.ui = null
        this.generateBtn()
    }

    getDriverInfo(driver, success) {
        dashboard.get(`/profession/${this.module}/${driver}`, function (res) {
            if (res.code !== 200) {
                return http.handleError(res, "获取driver信息失败")
            }
            success(driver, res.data["render"])
        }, function (res) {
            return http.handleError(res, "获取driver信息失败")
        })
    }

    updateUi(driver, render, data) {
        render = formatProfessionRender(render)
        let btn = this.options["btns"]
        this.ui = new Render(this.options["panel"], render, driver, data,  function () {
            if ($(btn).length > 0 && !$(btn).is(":visible")) {
                $(btn).show()
            }
        })
    }

    resetUi(driver, render) {
        if (this.ui) {
            render = formatProfessionRender(render)
            this.ui.Reset(render, driver)
        }
    }

    generateBtn() {
        let o = this

        $(this.options["btns"]).on("click", "button[name=reset]", function () {
            o.resetEvent()
        })
        $(this.options["btns"]).find("button[name=submit]").bind("click", function () {
            o.submitEvent()
        })
    }

    resetEvent() {
        if (this.ui) {
            this.ui.ResetVal()
        }
    }

    submitEvent() {
    }
}

class ProfessionCreator extends ProfessionRender {
    generator(options) {
        return BaseGenerator(options)
    }

    constructor(module, options) {
        super(module, options)
        this.options = options
        this.initName()
        this.init()
    }

    initName() {

        let o = this
        $(o.options["id"]).attr("readonly", true)
        $(o.options["id"]).val(`@${o.options["profession"]}`)
        $(o.options["name"]).on("change", function (e) {
            let v = $(this).val()
            $(o.options["id"]).val(`${v}@${o.options["profession"]}`)

            if (nameRule.test(v)) {
                $(this).removeClass("is-invalid")
                $(this).addClass("is-valid")
                return true
            } else {
                $(this).removeClass("is-valid")
                $(this).addClass("is-invalid")
                return false
            }
            return true
        })

    }

    init() {
        let o = this
        dashboard.get(`/profession/${this.module}/`, function (res) {
            if (res.code !== 200) {
                return http.handleError(res, "获取driver列表失败")
            }
            let data = res.data
            let target = $(o.options["drivers"])
            target.empty()
            for (let i = 0; i < data.length; i++) {
                if (i === 0) {
                    target.append("<option value='" + data[i].name + "' selected>" + data[i].name + "</option>")
                } else {
                    target.append("<option value='" + data[i].name + "'>" + data[i].name + "</option>")
                }
            }
            o.getDriverInfo(target.val(), function (driver, render) {
                o.updateUi(driver, render, null)
            })
            target.change(function () {
                o.getDriverInfo($(this).val(), function (driver, render) {
                    o.resetUi(driver, render, null)
                })
            });
        }, function (res) {
            return http.handleError(res, "获取driver列表失败")
        })
    }

    submitEvent() {
        const o = this
        if (this.ui) {
            let url = `/api/${this.module}/`
            this.ui.Submit(function (data) {
                let name = $(o.options["name"]).val()
                if (!nameRule.test(name)) {
                    $(this).removeClass("is-valid")
                    $(this).addClass("is-invalid")
                    common.message("name invalid", "danger")
                    return
                }
                data["name"] = name
                data["driver"] = $(o.options["drivers"]).val()
                data["description"] = $(o.options["description"]).val()
                dashboard.create(url, data, function (res) {
                    if (res.code !== 200) {
                        http.handleError(res, "新增失败")
                        return
                    }
                    common.message("新增成功", "success")
                    Back()
                }, function (res) {
                    http.handleError(res, "新增失败")
                })
            }, function (err) {
                common.message("config format error:"+err, "danger")
            })
        }
    }

}

class ProfessionEditor extends ProfessionRender {
    generator(options) {
        return BaseGenerator(options)
    }

    constructor(module, options) {

        super(module, options)
        this.name = options["worker"]
        this.init()
    }

    initData(data) {
        let o = this
        $(o.options["id"]).attr("readonly", true)
        $(o.options["id"]).val(`${data["id"]}`)

        $(o.options["name"]).attr("readonly", true)
        $(o.options["name"]).val(`${data["name"]}`)

        $(o.options["drivers"]).attr("disabled", true)
        $(o.options["drivers"]).append(`<option value="${data["driver"]}">${data["driver"]}</option>`)
        $(o.options["drivers"]).val(`${data["driver"]}`)

        $(o.options["description"]).val(`${data["description"]}`)
    }


    init() {
        let o = this
        dashboard.get(`/api/${this.module}/${this.name}`, function (res) {
            if (res.code !== 200 || res.data.length < 1) {
                return http.handleError(res, "获取详情失败")
            }

            let data = res.data

            o.getDriverInfo(res.data["driver"], function (driver, render) {
                o.updateUi(driver, render, data)
            })
            o.initData(data)
        }, function (res) {
            return http.handleError(res, "获取详情失败")
        })
    }

    submitEvent() {
        if (this.ui) {
            let o = this
            let url = `/api/${this.module}/${this.name}`
            this.ui.Submit(function (data) {
                data["description"] = $(o.options["description"]).val()
                dashboard.update(url, data, function (res) {
                    if (res.code !== 200) {
                        http.handleError(res, "更新失败")
                        return
                    }
                    common.message("更新成功", "success")
                    Back()
                }, function (res) {
                    http.handleError(res, "更新失败")
                })
            }, function (err) {
                common.message("config format error:"+err, "danger")
            })
        }

    }
}

function formatProfessionRender(render) {
    const defaultFields = {"id": true, "name": true, "driver": true,"description":true}
    if (typeof render === "undefined") {
        throw "undefined"
    }
    let properties = render["properties"]
    let uiSort = render["ui:sort"]
    let newUiSort = new Array()
    for (let i in uiSort) {
        let name = uiSort[i]
        if (defaultFields[name] !== true) {
            newUiSort.push(name)
        }else{
            delete properties[name]
        }
    }
    render["properties"] = properties
    render["ui:sort"] = newUiSort
    return render

}