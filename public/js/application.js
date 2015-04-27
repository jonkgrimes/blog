$(function() {
  window.admin = {
    publishPost: function(e) {
      e.preventDefault();
      var $target = $(e.target),
          id = $target.attr("data-id"),
          url = "/admin/posts/" + id + "/publish";

      var xhr = $.ajax({
        url: url,
        method: "POST"
      });

      xhr.done(function() {
        var alertBox =
          "<div data-alert class='alert-box success radius'>" +
            "Post published!" +
            "<a href='#' class='close'>&times;</a>" +
          "</div>";
        $(".admin-list").prepend(alertBox);
        $(document).foundation('alert', 'reflow');
      });
    },
    init: function() {
      $("a.publish-post").on("click", admin.publishPost);
    }
  };
});
