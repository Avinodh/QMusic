/* Global Variables */
var trackList, playlistView;

$(document).ready(function(){

	/*
  $.get("/viewplaylist").done(function(data){
    var d = JSON.parse(data);
    for (var i = 0; i < d.length; i++)
      $(".track-table").append(d[i].track["name"]);
  });
	*/

  // Track Model
  var TrackItem = Backbone.Model.extend({
    defaults: {
      id:'',
      name:'',
      artists:'',
    }
  });

  // Track Collection (Track List)
  var TrackList = Backbone.Collection.extend({
      url : '/viewplaylist'
  });

	// Instantiate a new Track Collection
  trackList = new TrackList();

  // View for one Track Item row
  var TrackItemView = Backbone.View.extend({
    model: new TrackItem(),
    tagName: 'tr',
    className: 'track-record',
    events: {
      "click td.remove-track": "removeTrack"
    },
    initialize: function() {
      this.template = _.template($('.playlist-template').html())
    },
    removeTrack: function(e) {
      e.preventDefault();
      $("#song-list").html('');
      var track_id = this.model.get('id');

      $.post("/removesong", {trackId:track_id}).done(function(data) {
        alert(data);
  		});
    },
    render: function() {
      this.$el.html(this.template(this.model.toJSON()));
      return this;
    }
  });

  // View for entire Track List
  var playlistView = Backbone.View.extend({
    model: trackList, //collection
    el: $(".track-table"),
    initialize: function() {
      var self = this;
    },
    render: function() {
      var self = this;
      this.$el.html('');
      this.$el.append('<tr><th class="track-column-header">Track Name</th><th class="track-column-header">Artist</th></tr>');
      var count = 1;
      if (this.model.toArray().length == 0) {
        self.$el.append("<h5>Playlist is empty!</h5>");
        return this;
      }
      _.each(this.model.toArray(), function(track) {
        var trackString = JSON.stringify(track);
        var trackJson = JSON.parse(trackString);
        var artistString ='';
        for(artist in trackJson.artists) {
          artistString += trackJson.artists[artist].name + ' ';
        }
        track.set({
            artists:artistString,
        });
        self.$el.append((new TrackItemView({
          model: track
        })).render().$el);
      });
      return this;
    }
  });

  playlistView = new playlistView();
  playlistView.model.fetch().done(function() {
    playlistView.render();});
});