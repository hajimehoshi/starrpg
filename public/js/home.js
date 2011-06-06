jQuery(function ($) {
    $("#createGame").click(function () {
        var server = createServer($);
        server.post('/games', {
            title: 'New Game',
        }, function (jqXHR, data) {
            var newGameURL = jqXHR.getResponseHeader("Location");
            server.put(newGameURL + '/items/1', {}, function () {
                location.replace(newGameURL);
            });
            server.flush();
        });
        server.flush();
        return false;
    });
});

