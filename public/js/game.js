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
            // TODO: 即座に送るのではなく、ある程度キャッシュするように修正
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
    var cacheStr = '{}';
    var cacheJSON = {};
    var updated = createEvent();
    server.get(path, function (data) {
                   cacheStr = JSON.stringify(data);
                   cacheJSON = data;
                   updated.fire(cacheJSON);
               });
    return {
        get: function() {
            return cacheJSON;
        },
        update: function (data) {
            var dataStr = JSON.stringify(data);
            if (cacheStr === dataStr) {
                return;
            }
            cache = dataStr;
            cacheJSON = data;
            server.put(path, data);
            updated.fire(cacheJSON);
        },
        register: updated.register,
    }
}

function createView(jqDom) {
    var cache = jqDom.val();
    var updated = createEvent();
    jqDom.change(function () {
                     cache = jqDom.val();
                     updated.fire(cache);
                 });
    return {
        get: function () {
            return cache;
        },
        update: function () {
            var value = (0 < arguments.length) ? arguments[0] : jqDom.val();
            if (cache === value) {
                return;
            }
            cache = value;
            jqDom.val(value);
            updated.fire(cache);
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
    (function() {
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
    (function() {
         var server = createServer();
         var model = {
             game: createModel(server, location.pathname),
         };
         var editGamePresenter = {
             nameTextBox: createView($('#gameNameTextBox')),
         };
         var game = {
             name: '',
         };
         editGamePresenter.nameTextBox.register(function (name) {
                                                    game.name = name;
                                                    model.game.update(game);
                                                });
         model.game.register(function (game) {
                                 editGamePresenter.nameTextBox.update(game.name);
                             });
         var editItemsPresenter = {

         };
     })();
    // TODO: 色々と待つ処理
    $('#loading').hide();
}
jQuery(init);
