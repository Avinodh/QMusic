$(document).ready(function($) {

});

function authSpotify() {
  $.get("/authspotify").done(function(data){
      alert(data);
      window.location = data;
      // console.log(data);
  });
}

function openSearch() {
  window.location = "search.html";
}

function openDashboard() {
  window.location = "dashboard.html";
}

function showCreatePartyModule() {
  $("#popup").load("overlay.html",function(){
    $(this).clone().appendTo("body").remove();
    $('#createPartyButton').on(function () {
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

function closeOverlay() {
  var myNode = document.getElementById("popup");
  while (myNode.firstChild) {
      myNode.removeChild(myNode.firstChild);
  }
}
