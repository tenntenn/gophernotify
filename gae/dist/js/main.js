
define("main", ['MessageViewModel', 'Channel'], function(MessageViewModel, Channel) {
  var channel, model;
  model = new MessageViewModel();
  channel = new Channel(model, gophernotify.clientID, gophernotify.token);
  return ko.applyBindings(model);
});
