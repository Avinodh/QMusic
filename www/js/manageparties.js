/* Global Variables */
var partyList, partyListView;

$(document).ready(function() {
  // Party Item Model
  var PartyItem = Backbone.Model.extend({
    defaults: {
      active_time:'',
      party_name:'',
      party_location:'',
      playlist_id:'',
    }
  });

  // Parties Collection
  var PartyList = Backbone.Collection.extend({
      url : '/gethostparties'
  });

  // Instantiate a new Party Collection
  partyList = new PartyList();

  // View for one Party Item row
  var PartyItemView = Backbone.View.extend({
    model: new PartyItem(),
    tagName: 'span',
    events: {
      "click .party-object": "openPlaylist"
    },
    initialize: function() {
      this.template = _.template($('.party-object-template').html())
    },
    openPlaylist: function(e) {
      e.preventDefault();
      var playlistId = this.model.get('playlist_id');
      console.log("PLaylist Id: "+playlistId);

      window.location = "/playlist?playlist_id="+playlistId;
    },
    render: function() {
      this.$el.html(this.template(this.model.toJSON()));
      return this;
    }
  });

  // View for entire Party List
  var PartyListView = Backbone.View.extend({
    model: partyList, //collection
    el: $("#party-container"),
    initialize: function() {
      var self = this;
      /*this.model.fetch({
        success: function(response) {
          //this.render();
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
      self.$el.append("<h1> Manage Parties </h1>");
      if (this.model.toArray().length == 0) {
        self.$el.append("<h5>No parties created!</h5>");
        return this;
      }
      _.each(this.model.toArray(), function(party) {
        self.$el.append((new PartyItemView({
          model: party
        })).render().$el);
      });
      return this;
    }
  });

  partyListView = new PartyListView();

});

function viewManageParties() {
    $(".dashboard-section").css("display", "none");
    $("#manage-parties-section").css("display", "block");
    partyListView.model.fetch().done(function() {
    partyListView.render();});
}



