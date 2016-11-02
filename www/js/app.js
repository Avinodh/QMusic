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
    $('#createPlaylistForm').submit(function () {
      createPlaylist();
    });
  });

}

function createPlaylist() {
  // get all the inputs into an array.
  var $inputs = $('#createPlaylistForm :input');

  // not sure if you wanted this, but I thought I'd add it.
  // get an associative array of just the values.
  var values = {};
  $inputs.each(function() {
      values[this.name] = $(this).val();
  });
  console.log(values);
}
