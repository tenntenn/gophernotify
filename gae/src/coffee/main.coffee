define("main",
    [
        'MessageViewModel'
        'Channel'
    ],
    (MessageViewModel, Channel)->
        model = new MessageViewModel()
        channel = new Channel(model, gophernotify.clientID, gophernotify.token)
        ko.applyBindings(model);
)