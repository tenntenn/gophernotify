
define('MessageViewModel', [], function() {
  var MessageViewModel;
  return MessageViewModel = (function() {

    function MessageViewModel(msg) {
      this.message = ko.observable('');
      if (msg != null) this.message(msg);
    }

    return MessageViewModel;

  })();
});
