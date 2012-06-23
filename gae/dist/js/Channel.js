
define('Channel', [], function() {
  var Channel;
  return Channel = (function() {

    function Channel(model, clientId, token) {
      var channel, socket,
        _this = this;
      this.model = model;
      this.clientId = clientId;
      this.token = token;
      channel = new goog.appengine.Channel(this.token);
      this.callBacks = {};
      this.callBacks['onpost'] = function(data) {
        _this.model.message(data.Body);
        return $.post("/" + _this.clientId + "/response", {
          "clientID": _this.clientId
        });
      };
      socket = channel.open();
      socket.onmessage = function(msg) {
        var args, call, data, _base;
        data = $.evalJSON(msg.data);
        call = data.call;
        args = data.args;
        return typeof (_base = _this.callBacks)[call] === "function" ? _base[call](args) : void 0;
      };
    }

    return Channel;

  })();
});
