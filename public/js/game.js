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

function createModelFunc(server, path) {
    var cacheStr = '{}';
    var cacheJSON = {};
    var changed = createEvent();
    var func = function (data) {
        if (arguments.length === 0) {
            return cacheJSON;
        }
        var dataStr = JSON.stringify(data);
        if (cacheStr === dataStr) {
            return
        }
        cache = dataStr;
        cacheJSON = data;
        server.put(path, data);
        changed.fire(cacheJSON);
    };
    server.get(path, function (data) {
                   cacheStr = JSON.stringify(data);
                   cacheJSON = data;
                   changed.fire(cacheJSON);
               });
    return [func, changed.register];
}

function createViewFunc(jqDom) {
    var cache = jqDom.val();
    var changed = createEvent();
    var func = function () {
        var value = (0 < arguments.length) ? arguments[0] : jqDom.val();
        if (cache === value) {
            return value;
        }
        cache = value;
        jqDom.val(value);
        changed.fire(cache);
    }
    jqDom.change(function () {
                     cache = jqDom.val();
                     changed.fire(cache);
                 });
    return [func, changed.register];
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
         var funcs = createModelFunc(server, location.pathname);
         var model = {
             game: funcs[0],
             gameChangedReg: funcs[1],
         };
         var funcs = createViewFunc($('#gameNameTextBox'));
         var editGamePresenter = {
             nameTextBox: funcs[0],
             nameTextBoxChangedReg: funcs[1],
         };
         var game = {
             name: '',
         };
         editGamePresenter.nameTextBoxChangedReg(function (name) {
                                                     game.name = name;
                                                     model.game(game);
                                                 });
         model.gameChangedReg(function (game) {
                                  editGamePresenter.nameTextBox(game.name);
                              });
         var editItemsPresenter = {

         };
     })();
    // TODO: 色々と待つ処理
    $('#loading').hide();
}
jQuery(init);
