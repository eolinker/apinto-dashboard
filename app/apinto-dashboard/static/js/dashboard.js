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
