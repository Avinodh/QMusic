$(document).ready(function($) {

});

function authSpotify() {
    $.get("/authspotify").done(function(data){
        window.location = data;
        //console.log(data);
    });
}

function showCreatePlaylistModule() {
  $("body").load("overlay.html",function(){
    $(this).clone().appendTo("body").remove();
  });
}

function createPlaylist(usrnm, pswd, secretCode, activeTime) {
  console.log($('#usrnm'));
}
