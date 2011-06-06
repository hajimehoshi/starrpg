// JSON Object must be copied on write!

jQuery(function ($) {
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
            regUpdated: updated.register,
            isLoaded: function () {
                return isLoaded;
            },
        };
    }
    function createView(jqDom) {
        var updated = createEvent();
        jqDom.change(function () {
            updated.fire(jqDom.val());
        });
        return {
            update: function (value) {
                if (jqDom.val() === value) {
                    return;
                }
                jqDom.val(value);
                updated.fire(value);
            },
            regUpdated: updated.register,
        };
    }
    function packZeroes(num, length) {
        var zeroes = '';
        for (var i = 0; i < length; i++) {
            zeroes += '0';
        }
        return (zeroes + num).substr(-length)
    }
    function createEntriesView(jqDom) {
        var select = jqDom;
        for (var i = 1; i <= 500; i++) {
            var option = $('<option></option>').text(packZeroes(i, 3) + ': ').attr('value', i);
            select.append(option);
        }
        var options = select.children('option')
        var selected = createEvent();
        select.change(function () {
            selected.fire($(this).children('option:selected').val());
        });
        return {
            update: function (value) {
                $.each(value, function (i, val) {
                    if (!val.name) {
                        return;
                    }
                    var option = options.filter('[value="' + i + '"]');
                    if (!option) {
                        // error?
                        return;
                    }
                    var text = packZeroes(i, 3) + ': ' + val.name
                    option.text(text);
                });
            },
            regSelected: selected.register,
        };
    }
    (function () {
        var mainPanels = $('.mainPanel');
        $('.mainPanelNavItem').click(function () {
            mainPanels.hide();
            var m = this.id.match(/^(.+?)NavItem$/);
            $('#' + m[1]).show();
            return false;
        });
        $('.mainPanelNavItem.default').click();
    })();
    (function () {
        var editPanels = $('.editPanel');
        $('.editPanelNavItem').click(function () {
            // calll check function?
            editPanels.hide();
            var m = this.id.match(/^(.+?)NavItem$/);
            $('#' + m[1]).show();
            $(window).resize();
            return false;
        });
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
            titleView.regUpdated(function (title) {
                game.update('title', title);
            });
            descriptionView.regUpdated(function (description) {
                game.update('description', description);
            });
            game.regUpdated(function (game) {
                titleView.update(game.title);
                descriptionView.update(game.description);
            });
        })();
        (function () {
            var selectedIndex = -1;
            var entriesView = createEntriesView($('#editItems nav select'));
            entriesView.regSelected(function (i) {
                selectedIndex = i;
                console.log(selectedIndex);
            });
            items.regUpdated(function (items) {
                entriesView.update(items);
            });
        })();
    })();
    $(window).resize(function () {
        $('section.hasEntries nav select').each(function (i, dom) {
            var jqDom = $(dom);
            jqDom.height(jqDom.parent().innerHeight());
        });
    });
    // TODO: 色々と待つ処理
    $('#loading').hide();
});

