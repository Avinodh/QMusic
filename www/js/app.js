$(document).ready(function($) {

});

function authSpotify() {
  $.get("/authspotify").done(function(data){
      window.location = data;
      //console.log(data);
  });
}

function showCreatePartyModule() {
  $("body").load("overlay.html",function(){
    $(this).clone().appendTo("body").remove();
    $('#createPartyForm').submit(function () {
      createParty();
    });
  });

}

function createParty() {
  // get all the inputs into an array.
  var $inputs = $('#createPartyForm :input');

  // not sure if you wanted this, but I thought I'd add it.
  // get an associative array of just the values.
  var values = {};
  $inputs.each(function() {
      values[this.name] = $(this).val();
  });
  console.log(values);
}
