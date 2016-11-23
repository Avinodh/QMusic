function openSearch() {
  window.location = "search.html";
}

function openDashboard() {
  window.location = "dashboard.html";
}

function viewCreateParty() {
    $(".dashboard-section").css("display", "none");
    $("#create-party-section").css("display", "block");
}