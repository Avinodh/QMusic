/* Global Variables */
var trackList, trackListView;

$(document).ready(function() {
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
      url : '/searchsong'
  });

  // Instantiate a new Track Collection
  trackList = new TrackList();

  // View for one Track Item row
  var TrackItemView = Backbone.View.extend({
    model: new TrackItem(),
    tagName: 'tr',
    className: 'track-record',
    events: {
      "click td.add-track": "addTrack"
    },
    initialize: function() {
      this.template = _.template($('.track-list-template').html())
    },
    addTrack: function(e) {
      e.preventDefault();
      var track_id = this.model.get('id');
      $.post("/addsong", {trackId:track_id}).done(function(data){
      //$(this.$el).css("background-color","green");
        alert(data);
      // console.log(data);
      });
      this.$el.hide();
    },
    render: function() {
      this.$el.html(this.template(this.model.toJSON()));
      return this;
    }
  });

  // View for entire Track List
  var TrackListView = Backbone.View.extend({
    model: trackList, //collection
    el: $(".track-table"),
    initialize: function() {
      var self = this;
      /*this.model.fetch({
        success: function(response) {
        },
        error: function() {
          //console.log('Failed to get apps.');
        }
      }).done(function() {
        self.render();
      });*/
    },
    render: function() {
      var self = this;
      this.$el.html('');
      this.$el.append('<tr><th class="track-column-header">Track Name</th><th class="track-column-header">Artist</th></tr>');
      var count = 1;
      if (this.model.toArray().length == 0) {
        self.$el.append("<h5>0 tracks found!</h5>");
        return this;
      }
      _.each(this.model.toArray(), function(track) {
        var trackString = JSON.stringify(track);
        var trackJson = JSON.parse(trackString);
        var artistString ='';
        for(artist in trackJson.artists){
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

  trackListView = new TrackListView();
});


function searchSong() {
  var queryString = $("#search-song").val();
  trackListView.model.fetch({data: {searchsong: queryString}, processData: true}).done(function() {
    trackListView.render();});

  /*$.get("/searchsong", {searchsong:queryString}).done(function(data){
    drawTable(data);
  });*/
}
