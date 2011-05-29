function init($) {
    function createGame(e) {
        var server = createServer($);
        var data = {
            name: 'New Game',
        }
        var callback = function (jqXHR, data) {
            var newGameURL = jqXHR.getResponseHeader("Location");
            var callback = function () {
                location.replace(newGameURL);                
            }
            server.put(newGameURL + '/items/1', {}, callback);
            server.flush();
        }
        server.post('/games', data, callback);
        server.flush();
        return false;
    }
    $("#createGame").click(createGame);
}
jQuery(init);
