/* Global Variables */
var trackList, playlistView;
var GLOBAL_LIST= [];

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
      track:'',
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
      $.post("/removetrack", {trackId:track_id}).done(function(data) {
  		});
      this.$el.hide();
    },
    render: function() {
      this.$el.html(this.template(this.model.toJSON()));
      return this;
    }
  });

  // View for entire Track List
  var PlaylistView = Backbone.View.extend({
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
      _.each(this.model.toArray(), function(o) {
        var respString = JSON.stringify(o);
        var respJson = JSON.parse(respString);
        var track = respJson['track'];
        var t_id = track['id'];
        var t_name = track['name'];
        var t_artists ='';
        for(aIndex in track.artists){
          t_artists += track.artists[aIndex].name + ' ';
        }
        o.set({
            artists:t_artists,
            id:t_id,
            name:t_name,
        });
        self.$el.append((new TrackItemView({
          model: o
        })).render().$el);
      });
      return this;
    }
  });

  playlistView = new PlaylistView();
  playlistView.model.fetch().done(function() {
    playlistView.render();});
});

