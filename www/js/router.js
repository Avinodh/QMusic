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

function viewManageParties() {
    $(".dashboard-section").css("display", "none");
    $("#manage-parties-section").css("display", "block");
}
