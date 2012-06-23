define('MessageViewModel',
    [
    ],
    ()->
        class MessageViewModel
            constructor:()->
                @message = ko.observable('')
)