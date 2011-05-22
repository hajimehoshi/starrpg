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
         var server = {
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
                     contentType: 'application/json',
                     dataType: 'json',
                     type: "PUT",
                 };
                 $.ajax(args);
             }
         }
         function createModelFunc(path) {
             var cacheStr = '{}';
             var cacheJSON = {};
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
                 func.changed();
             }
             var changedFunc = null;
             func.changed = function () {
                 if (arguments.length === 0) {
                     if (changedFunc instanceof Function) {
                         changedFunc(cacheJSON);
                     }
                 } else {
                     changedFunc = arguments[0];                         
                 }
             }
             return func;
         }
         var model = {
             game: createModelFunc(location.pathname),
         };
         function createViewFunc(jqDom) {
             var cache = jqDom.val();
             var func = function () {
                 var value = (0 < arguments.length) ? arguments[0] : jqDom.val();
                 if (cache === value) {
                     return value;
                 }
                 cache = value;
                 jqDom.val(value);
                 func.changed();
             }
             var changedFunc = null;
             func.changed = function () {
                 if (arguments.length === 0) {
                     if (changedFunc instanceof Function) {
                         changedFunc(cache);
                     }
                 } else {
                     changedFunc = arguments[0];                         
                 }
             }
             return func;
         }
         var editGamePresenter = {
             nameTextBox: createViewFunc($('#gameNameTextBox')),
         };
         var game = {
             name: '',
         };
         editGamePresenter.nameTextBox.changed(function(name) {
                                                   game.name = name;
                                                   model.game(game);
                                               });
         model.game.changed(function (game) {
                                editGamePresenter.nameTextBox(game.name);
                            });
         var editItemsPresenter = {

         };
         function reportToPresenter() {
             switch (activeEditPanel) {
             case 'editGame':
                 editGamePresenter.nameTextBox();
                 break;
             }
         }
         setInterval(reportToPresenter, 1000);

         server.get(location.pathname, function (data) {
                        model.game(data);
                    })
     })();
    // TODO: 色々と待つ処理
    $('#loading').hide();
}
jQuery(init);
