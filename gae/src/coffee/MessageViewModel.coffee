define('MessageViewModel',
    [
    ],
    ()->
        class MessageViewModel
            constructor:(msg)->
                @message = ko.observable('')
                if msg?
                    @message(msg)
)