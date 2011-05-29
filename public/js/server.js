function createServer($) {
    var server = {
        get: function (path, callback) {
            var args = {
                url: path,
                dataType: 'json',
                type: "GET",
                success: function (data, status, jqXHR) {
                    if (callback instanceof Function) {
                        callback(jqXHR, data);
                    }
                },
                error: function (jqXHR, status) {
                    // TODO: logging
                }
            };
            $.ajax(args);
        },
        post: function (path, data, callback) {
            var args = {
                url: path,
                data: JSON.stringify(data),
                contentType: 'application/json; charset=utf-8',
                dataType: 'json',
                type: 'POST',
                success: function (data, status, jqXHR) {
                    if (callback instanceof Function) {
                        callback(jqXHR, data);                        
                    }
                },
                error: function (jqXHR, status) {
                    // TODO: logging
                }
            }
            $.ajax(args);
        },
        put: function (path, data) {
            // TODO: 即座に送るのではなく、ある程度キャッシュして最適化の後、
            // 送信するように修正
            var args = {
                url: path,
                data: JSON.stringify(data),
                contentType: 'application/json; charset=utf-8',
                dataType: 'json',
                type: "PUT",
                error: function (jqXHR, status) {
                    // TODO: logging
                }
            };
            $.ajax(args);
        },
        flush: function () {
            // TODO: implement it
        }
    };
    return server;
}
