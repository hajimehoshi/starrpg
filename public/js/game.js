// JSON Object must be copied on write!

function clone(obj) {
    return JSON.parse(JSON.stringify(obj));
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
    server.get(path, function (jqXHR, data) {
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
    if (jqDom && (jqDom.change instanceof Function)) {
        jqDom.change(function () {
            updated.fire(jqDom.val());
        });
    }
    return {
        update: function (value) {
            if (jqDom && (jqDom.val instanceof Function)) {
                if (jqDom.val() === value) {
                    return;
                }
                jqDom.val(value);
            }
            updated.fire(value);
        },
        register: updated.register,
    }
}

function packZeroes(num, length) {
    var zeroes = '';
    for (var i = 0; i < length; i++) {
        zeroes += '0';
    }
    return (zeroes + num).substr(-length)
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
    (function () {
        var editPanels = $('.editPanel');
        function switchEditPanel() {
            // calll check function?
            editPanels.hide();
            var m = this.id.match(/^(.+?)NavItem$/);
            $('#' + m[1]).show();
            $(window).resize();
            return false;
        }
        $('.editPanelNavItem').click(switchEditPanel);
        $('.editPanelNavItem.default').click();
    })();
    (function () {
        var path = location.pathname;
        var server = createServer($);
        var game = createModel(server, path);
        var items = createModel(server, path + '/items');
        (function () {
            var titleView = createView($('#editGame *[name=title]'));
            var descriptionView = createView($('#editGame *[name=description]'));
            titleView.register(function (title) {
                game.update('title', title);
            });
            descriptionView.register(function (description) {
                game.update('description', description);
            });
            game.register(function (game) {
                titleView.update(game.title);
                descriptionView.update(game.description);
            });
        })();
        $(window).resize(function () {
            $('section.hasEntries nav select').each(function (i, dom) {
                var jqDom = $(dom);
                jqDom.height(jqDom.parent().innerHeight());
            });
        });
        (function () {
            var entriesSelect = $('#editItems nav select');
            $(window).resize(function () {
                entriesSelect.height(entriesSelect.parent().innerHeight());
            });
            for (var i = 1; i <= 500; i++) {
                var option = $('<option></option>').text(packZeroes(i, 3) + ': ').attr('value', i);
                entriesSelect.append(option);
            }
            var entriesView = createView();
            entriesView.register(function (value) {
            });
            items.register(function (items) {
                entriesView.update(items);
            });
        })();
    })();
    // TODO: 色々と待つ処理
    $('#loading').hide();
}
jQuery(init);
