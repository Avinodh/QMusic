$(document).ready(function($) {
    
});

function authSpotify() {
    $.get("/authspotify").done(function(data){
        window.location = data;
        //console.log(data);
    });
}