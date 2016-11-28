// // When the user clicks on <div>, open the popup
function togglePopupPartyGoer() {
  $("#myPopup").load("popup.html", function(){
    $(this).clone().appendTo("body").remove();
    document.getElementById('usercode').focus();
  });
}

function closePopup() {
  var myNode = document.getElementById("myPopup");
  while (myNode.firstChild) {
      myNode.removeChild(myNode.firstChild);
  }
}
