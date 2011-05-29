function init($) {
    function createGame(e) {
        var server = createServer($);
        var data = {
            name: 'New Game',
        }
        var callback = function (jqXHR, data) {
            var newGameURL = jqXHR.getResponseHeader("Location");
            server.put(newGameURL + '/items', {});
            server.flush();
            location.replace(newGameURL);
        }
        server.post('/games', data, callback);
        server.flush();
        return false;
    }
    $("#createGame").click(createGame);
}
jQuery(init);
