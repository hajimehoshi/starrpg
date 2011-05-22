function init($) {
    function createGame(e) {
        var args = {
            url: "/games",
            type: "POST",
            success: function(data, status, jqXHR) {
                if (jqXHR.status === 201) {
                    location.replace(jqXHR.getResponseHeader("Location"));
                } else {
                    // unexpected                             
                }
            }
        };
        $.ajax(args);
        return false;
    }
    $("#createGame").click(createGame);
}
jQuery(init);
