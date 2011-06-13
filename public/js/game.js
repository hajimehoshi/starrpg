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
                    func.apply(null, arguments);
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
        server.get(path, function (jqXHR, data) {
            cache = data;
            updated.fire(cache);
        });
        return {
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
        };
    }
    function createCollectionModel(server, path) {
        var cache = {};
        var updated = createEvent();
        var entryUpdated = createEvent();
        server.get(path, function (jqXHR, data) {
            cache = data;
            updated.fire(cache);
        });
        return {
            getEntry: function (id) {
                return cache[id];
            },
            regUpdated: updated.register,
            updateEntry: function (id, key, value) {
                if (cache[id] && cache[id][key] === value) {
                    return;
                }
                //cache = clone(cache);
                if (!cache[id]) {
                    cache[id] = {};
                } 
                cache[id][key] = value;
                server.put(path + '/' + id, cache[id]);
                entryUpdated.fire(id, cache[id]);
            },
            regEntryUpdated: entryUpdated.register,
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
            enable: function () {
                jqDom.removeAttr('disabled');
                return this;
            },
            disable: function () {
                jqDom.attr('disabled', 'disabled');
                return this;
            },
        };
    }
    function createEntriesView(jqDom) {
        function getEntryText(id, value) {
            function packZeroes(num, length) {
                var zeroes = '';
                for (var i = 0; i < length; i++) {
                    zeroes += '0';
                }
                return (zeroes + num).substr(-length)
            }
            var text = packZeroes(id, 3) + ': ';
            if (value && value.name) {
                text += value.name;
            }
            return text;
        }
        var select = jqDom;
        for (var i = 1; i <= 500; i++) {
            var option = $('<option></option>').text(getEntryText(i, null)).attr('value', i);
            select.append(option);
        }
        i = undefined;
        var options = select.children('option')
        var selected = createEvent();
        select.change(function () {
            selected.fire($(this).children('option:selected').val());
        });
        return {
            update: function (values) {
                $.each(values, function (id, value) {
                    if (!value.name) {
                        return;
                    }
                    var option = options.filter('[value="' + id + '"]');
                    option.text(getEntryText(id, value));
                });
            },
            updateEntry: function (id, value) {
                var option = options.filter('[value="' + id + '"]');
                option.text(getEntryText(id, value));
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
        var items = createCollectionModel(server, path + '/items');
        (function () {
            var titleView = createView($('#editGame *[name="title"]'));
            var descriptionView = createView($('#editGame *[name="description"]'));
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
            var entriesView = createEntriesView($('#editMaps nav select'));
        })();
        (function () {
            var selectedIndex = -1;
            var entriesView = createEntriesView($('#editItems nav select'));
            var nameView = createView($('#editItems *[name="name"]')).disable();
            entriesView.regSelected(function (i) {
                selectedIndex = i;
                var item = items.getEntry(i);
                nameView.enable().update((item && item.name) ? item.name : '');
            });
            nameView.regUpdated(function (name) {
                items.updateEntry(selectedIndex, 'name', name);
            });
            items.regUpdated(function (items) {
                entriesView.update(items);
            });
            items.regEntryUpdated(function (id, item) {
                entriesView.updateEntry(id, item);
            });
        })();
        (function () {
            var entriesView = createEntriesView($('#editActors nav select'));
        })();
        (function () {
            var entriesView = createEntriesView($('#editSkills nav select'));
        })();
        (function () {
            var entriesView = createEntriesView($('#editEnemies nav select'));
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

