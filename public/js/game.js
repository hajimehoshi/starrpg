// JSON Object must be copied on write!

function clone(obj) {
    return JSON.parse(JSON.stringify(obj));
}

function createServer() {
    return {
        get: function (path, callback) {
            var args = {
                url: path,
                dataType: 'json',
                type: "GET",
                success: function (data, status, jqXHR) {
                    if (jqXHR.status == 200) {
                        callback(data);
                    } else {
                        // unexpected!
                    }
                }
            };
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
            };
            $.ajax(args);
        }
    };
}

function createEvent() {
    var registeredFuncs = [];
    return {
        fire: function () {
            for (var i = 0; i < registeredFuncs.length; i++) {
                var func = registeredFuncs[i];
                if (func instanceof Function) {
                    func.apply(null, arguments);
                }
            }
        },
        register: function () {
            registeredFuncs.push(arguments[0]);
        }
    }
}

function createModel(server, path) {
    var cache = {};
    var updated = createEvent();
    var isLoaded = false;
    server.get(path, function (data) {
                   cache = data;
                   updated.fire(cache);
                   isLoaded = true;
               });
    return {
        get: function(key) {
            return cache[key];
        },
        update: function (key, value) {
            if (cache[key] === value) {
                return;
            }
            cache = clone(cache);
            cache[key] = value;
            server.put(path, cache);
            updated.fire(cache);
        },
        register: updated.register,
        isLoaded: function () {
            return isLoaded;
        },
    }
}

function createView(jqDom) {
    var updated = createEvent();
    jqDom.change(function () {
                     updated.fire(jqDom.val());
                 });
    return {
        get: function () {
            return jqDom.val();
        },
        update: function (value) {
            if (jqDom.val() === value) {
                return;
            }
            jqDom.val(value);
            updated.fire(jqDom.val());
        },
        register: updated.register,
    }
}

function init($) {
    (function () {
         var mainPanels = $('.mainPanel');
         function switchMainPanel() {
             mainPanels.hide();
             var m = this.id.match(/^(.+?)NavItem$/);
             $('#' + m[1]).show();
             return false;
         }
         $('.mainPanelNavItem').click(switchMainPanel);
         $('.mainPanelNavItem.default').click();
     })();
    var activeEditPanel = null;
    (function () {
         var editPanels = $('.editPanel');
         function switchEditPanel() {
             // calll check function?
             editPanels.hide();
             var m = this.id.match(/^(.+?)NavItem$/);
             $('#' + m[1]).show();
             activeEditPanel = m[1];
             return false;
         }
         $('.editPanelNavItem').click(switchEditPanel);
         $('.editPanelNavItem.default').click();
     })();
    (function () {
         var server = createServer();
         var game = createModel(server, location.pathname);
         var items = createModel(server, location.pathname + '/items');
         (function () {
              var nameView = createView($('#editGame *[name=name]'));
              var descriptionView = createView($('#editGame *[name=description]'));
              nameView.register(function (name) {
                                    game.update('name', name);
                                });
              descriptionView.register(function (description) {
                                           game.update('description', description);
                                       });
              game.register(function (game) {
                                nameView.update(game.name);
                                descriptionView.update(game.description);
                            });
          })();
         (function () {
              var entriesView = createView($('#editItems nav'));
          })();
         var editItemsPresenter = {
             
         };
     })();
    // TODO: 色々と待つ処理
    
    $('#loading').hide();
}
jQuery(init);
