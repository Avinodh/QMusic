$(document).ready(function(){

    $.get("/viewplaylist").done(function(data){
        var d = JSON.parse(data);
        for (var i = 0; i < d.length; i++)
            $(".track-table").append(d[i].track["name"]);
    });
});