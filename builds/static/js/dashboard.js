let http = {
    post: function (url, data, success, error){
        this.ajax("POST", url, data, success, error)
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
    ajax: function (type, url, data, success, error, complete) {
        $.ajax({
            url: url,
            type: type,
            data: data,
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
        });
    },
}
let confirm = function (title, msg, success, cancel){
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
                            <button type="button" id="cancel_btn" class="btn btn-secondary">取消</button>
                            <button type="button" id="confirm_btn" class="btn btn-danger" >删除</button>
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
}