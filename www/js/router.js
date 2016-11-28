function openSearch() {
  window.location = "search.html";
}

function openDashboard() {
  window.location = "renderdashboard";
}

function viewCreateParty() {
    $(".dashboard-section").css("display", "none");
    $("#create-party-section").css("display", "block");
}