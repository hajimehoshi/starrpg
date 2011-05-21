$(function ($) {
      function createGame(e) {
          var args = {
              url: "/games",
              type: "POST",
              success: function(data, status, jqxhr) {
                  if (jqxhr.status === 201) {
                      location.replace(jqxhr.getResponseHeader("Location"));
                  } else {
                      // unexpected                             
                  }
              }
          };
          $.ajax(args);
          return false;
      }
      $("#createGame").click(createGame);
  });
