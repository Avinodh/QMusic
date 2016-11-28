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

function openPlaylist() {
     $.get("/currentplaylist", function(data){
      window.location = "/playlist?playlist_id="+data;
  });
}

function logout() {
    window.location = "/";
}