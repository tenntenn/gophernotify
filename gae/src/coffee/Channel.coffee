define('Channel',
    [
    ],
    ()->
        class Channel
            constructor:(@model, @clientId, @token)->
                channel = new goog.appengine.Channel(@token)
                @callBacks = {}
                @callBacks['onpost'] = (data) =>
                    @model.message(data.Body)
                    $.post('/channel/response',{"clientID" : @clientId});
                socket = channel.open()
                socket.onmessage = (msg) =>
                    data = $.evalJSON(msg.data)
                    call = data.call
                    args = data.args
                    @callBacks[call]?(args)
)