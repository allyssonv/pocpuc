let app = require('express')()
let http = require('http').Server(app)
let io = require('socket.io')(http)
var Kafka = require('no-kafka')
let log = console.log

io.on('connection', (socket) => {
    log('Connected')
    socket.on('disconnect', () => {
        log('Disconnected')
    })
})

http.listen(3000, () => {
    log('Server listen on PORT 3000')
    var consumer = new Kafka.SimpleConsumer({
        connectionString: 'localhost:9092',
        clientId: 'chart'
    })

    var dataHandler = function(messageSet, topic, partition) {
        messageSet.forEach((storedMessage) => {
            log('===========Consuming messages===========')
            if (topic == 'sales') {
                log( storedMessage.message.value.toString('utf8') )
                io.emit('sales', storedMessage.message.value.toString('utf8'))
            }
        })
    }

    return consumer.init().then(() => {
        var one = consumer.subscribe('sales', 0, dataHandler)
        return one
    })

})