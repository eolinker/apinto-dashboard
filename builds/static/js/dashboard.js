let http = {
    post: function (url, data, success, error){
        this.ajax("POST", url, data, success, error)
    },
    put: function (url, data, success, error){
        this.ajax("PUT", url, data, success, error)
    },
    get: function (url, data, success, error){
        this.ajax("GET", url, data, success, error)
    },
    delete: function (url, data, success, error){
        this.ajax("DELETE", url, data, success, error)
    },
    patch: function (url, data, success, error){
        this.ajax("PATCH", url, data, success, error)
    },
    ajax: function (type, url, data, success, error, complete, async) {
        let options = {
            url: url,
            type: type,
            contentType: "application/json; charset=utf-8",
            dataType: 'json',
            success: function (res) {
                if (success) {
                    success(res)
                }
            },
            error: function (res) {
                if (error) {
                    error(res)
                }
            },
            complete: function (res) {
                if (complete) {
                    complete(res)
                }
            },
        }
        if (data) {
            options["data"] = JSON.stringify(data)
        }
        if(!async){
            options["async"] = false
        }
        $.ajax(options);
    },
    handleError: function (res, msg){
        if (res['msg']){
            common.message(res['msg'], "danger")
        }else {
            common.message(msg, "danger")
        }
    }
}
let dashboard = {
    get: function (url, success, error){
        http.get(url, null,  success, error)
    },
    update: function (url, data, success, error) {
        http.put(url, data,  success, error)
    },
    patch: function (url, data, success, error) {
        http.patch(url, data,  success, error)
    },
    delete:function (url, success, error){
        http.delete(url,"", success, error)
    },
    create: function (url, data, success, error){
        http.post(url, data,  success, error)
    },
    getWithAsync: function (url, success, error){
        http.ajax("GET", url, null, success, error, null, false)
    },
    getExtenders: function (success, error){
        this.get("/api/extenders/", success, error)
    },
    getExtenderInfo: function (id, success, error){
        this.get("/api/extenders/"+id, success, error)
    },
    searchSkill:function (module,skill,success,error){
        this.get(`/skill/${module}?skill=${skill}`, success, error)
    }
}
let common = {
    /**
     * ???????????????
     * @param title
     * @param msg
     * @param success
     * @param cancel
     */
    confirm :function (title, msg, success, cancel){
        let model = $("#confirmModel")
        if (model.length > 0) {
            model.remove();
        }
        let html = `<div class="modal fade" id="confirmModel" tabindex="-1" aria-labelledby="confirmModelLabel" aria-hidden="true">
                <div class="modal-dialog">
                    <div class="modal-content">
                        <div class="modal-header">
                            <h5 class="modal-title" id="confirmModelLabel">`+title+`</h5>
                            <button type="button" class="close" data-dismiss="modal" aria-label="Close">
                                <span aria-hidden="true">&times;</span>
                            </button>
                        </div>
                        <div class="modal-body">`+msg+`</div>
                        <div class="modal-footer">
                            <button type="button" id="cancel_btn" class="btn btn-secondary">??????</button>
                            <button type="button" id="confirm_btn" class="btn btn-danger" >??????</button>
                        </div>
                    </div>
                </div>
            </div>`;
        $("body").append(html);
        model = $("#confirmModel")
        model.modal("show");
        model.on("click", "button#cancel_btn",function() {
            if (cancel) {
                cancel()
            }
            $("#confirmModel").modal("hide");
        });
        model.on("click", "button#confirm_btn",function() {
            if (success) {
                success()
            }
            $("#confirmModel").modal("hide");
        });
    },
    /**
     * ???????????????
     * @param msg ????????????
     * @param type ????????????????????????bootstrap???alert??? danger,success,warning,info...
     */
    alert: function(msg, type){
        if(typeof(type) =="undefined") { // ?????????type????????????success??????????????????
            type = "success";
        }
        // ??????bootstrap???alert??????
        let divElement = $("<div></div>").addClass('alert').addClass('alert-'+type).addClass('alert-dismissible').addClass('col-md-4').addClass('col-md-offset-4');
        // let scroll = document.body.scrollTop || document.documentElement.scrollTop;
        divElement.css({ // ????????????????????????
            "position": "fixed",
            "right":"10px",
            "top":  "60px"
        });
        divElement.text(msg); // ????????????????????????
        // ?????????????????????????????????
        let closeBtn = $('<button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">??</span></button>');
        $(divElement).append(closeBtn);
        // ???????????????????????????
        $('body').append(divElement);
        divElement.css("z-index","999999")
        return divElement;
    },

    /**
     * ???????????????????????????????????????
     * @param msg ????????????
     * @param type ???????????????
     */
    message: function(msg, type) {
        let divElement = this.alert(msg, type); // ??????Alert?????????
        let isIn = false; // ???????????????????????????

        divElement.on({ // ???setTimeout????????????????????????????????????????????????
            mouseover : function(){isIn = true;},
            mouseout  : function(){isIn = false;}
        });

        // ???????????????????????????
        setTimeout(function() {
            let IntervalMS = 20; // ???????????????????????????
            let floatSpace = 60; // ???????????????(px)
            // let scroll = document.body.scrollTop || document.documentElement.scrollTop;
            let nowTop = divElement.offset().top ; // ?????????????????????top???
            let stopTop = nowTop - floatSpace;    // ??????????????????top???
            divElement.fadeOut(IntervalMS * floatSpace); // ??????????????????

            let upFloat = setInterval(function(){ // ????????????
                if (nowTop >= stopTop) { // ?????????????????????top?????????????????????????????????
                    divElement.css({"top": nowTop--}); // ????????????top??????1px
                } else {
                    clearInterval(upFloat); // ????????????
                    divElement.remove();    // ????????????
                }
            }, IntervalMS);

            if (isIn) { // ???????????????setTimeout???????????????????????????????????????????????????
                clearInterval(upFloat);
                divElement.stop();
            }

            divElement.hover(function() { // ?????????????????????????????????????????????????????????
                clearInterval(upFloat);
                divElement.stop();
            },function() {
                divElement.fadeOut(IntervalMS * (nowTop - stopTop)); // ?????????????????????????????????????????????????????????*????????????????????????
                upFloat = setInterval(function(){ // ????????????
                    if (nowTop >= stopTop) {
                        divElement.css({"top": nowTop--});
                    } else {
                        clearInterval(upFloat); // ????????????
                        divElement.remove();    // ????????????
                    }
                }, IntervalMS);
            });
        }, 1500);
    }
}

let aceEditor = {
    InitEditor: function (id){
        let editor = ace.edit(id);
        editor.setFontSize(14)
        editor.setTheme("ace/theme/crimson_editor");
        editor.session.setMode("ace/mode/json");
        editor.renderer.setScrollMargin(10, 10);
        editor.setOptions({
            autoScrollEditorIntoView: true
        });
        return editor
    }
}
let modal = {
    table: function (id){
        let target = $("#"+id)
        target.removeClass("pop_window").removeClass("pop_window_small").html("")
        target.addClass("pop_window").addClass("pop_window_small").append(`<div class="pop_window_header">
            <span class="pop_window_title" id="${id}_title"></span>
            <span class="pop_window_title" id="${id}_title"></span>
            <button class="pop_window_button btn btn_default" id="${id}_close" >??????</button>
            <br>
        </div>
        <div class="card pop_window_body">
            <div class="card-body">
                <table class="table table-bordered">
                    <thead class="thead-light">
                    <tr>
                        <th>??????</th>
                        <th>??????</th>
                    </tr>
                    </thead>
                    <tbody id="${id}_body">
                    </tbody>
                </table>
    
            </div>
        </div>`
        )
        let close = $(`#${id}_close`)
        close.unbind("click");
        close.click(function (){
            target.animate({
                width:'toggle'
            }, "fast", function () {
                $(`#${id}_title`).text("")
                $(`#${id}_body`).html("")
                if($("body").hasClass("modal-open")){
                    $('body').removeClass("modal-open")
                    $("div.modal-backdrop.fade.show").remove()
                }
            });
        })
        $("body").on('click', 'div.modal-backdrop.fade.show', function () {
            $(`#${id}_close`).click()
        })
        target.show = function (title, body) {
            $(`#${id}_title`).text(title)
            $(`#${id}_body`).html(body)
            $("body").append("<div class='modal-backdrop fade show'></div>").addClass("modal-open")
            target.animate({
                width:'toggle'
            }, "fast");
        }
        target.hide = function () {
            close.click()
        }
        target.destroy = function () {
            target.hide()
            close.unbind("click")
            target.removeClass("pop_window").removeClass("pop_window_small").html("")
        }
        return target
    }
}

let date = {
    formatDate:function (v){
        let date = new Date(v)
        let year = date.getFullYear()
        let month = date.getMonth()
        let day = date.getDate()
        let fn = function (v) {
            if (v < 10) {
                v = "0" + v
            }
            return v
        }
        return year + "-" + fn(month+1) + "-" + fn(day)
    }
}

class Ace{
    constructor(id) {
        this.id = id
        let editor = ace.edit(id);
        editor.setFontSize(14)
        editor.setTheme("ace/theme/crimson_editor");
        editor.session.setMode("ace/mode/json");
        editor.renderer.setScrollMargin(10, 10);
        editor.setOptions({
            autoScrollEditorIntoView: true
        });
        this.editor = editor
    }
    get Value(){
        return JSON.parse(this.editor.getSession().getValue())
    }
    set Value(v){
        this.editor.getSession().setValue(JSON.stringify(v, null, 2))
    }
}
class Table{
    constructor(table, options) {
        this.Table = table
        $(this.Table).bootstrapTable(options)
    }
    GetData() {
        return $(this.Table).bootstrapTable('getData')
    }
    UpdateRowDetail(index, data, callback) {
        $(this.Table).bootstrapTable('updateRow', {
            index: index,
            row: data
        })
        if (callback){
            callback()
        }
    }
    AddNewRow(data){
        // let res = []
        // res.push(data)
        this.Table.bootstrapTable('refresh')
    }
    RemoveRow(field, value, callback){
        this.Table.bootstrapTable('remove', {
            field: field,
            values: [value]
        })
        if (callback){
            callback()
        }
    }
}
