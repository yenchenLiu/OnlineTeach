
$('.lesson_checkbox').click(function() {
            $.ajax({
                type: "POST",
                url: '',
                data: {lesson_checkbox:$(this).attr('value')}, //--> send id of checked checkbox on other page
                dataType: "json",
                success: function(data) {
                    $("#message").empty();
                    $("#message").append(data.data);
                },
                 error: function() {
                },
                complete: function() {
                }
            });
      });