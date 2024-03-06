# go message broker
A toy message_broker implemented in GO just for fun
## Broker Package
This package provides functionality for a broker server that handles messaging using a custom protocol over TCP.  
It uses a worker pool to handle conections and implements topics as message queues. 

### Overview
The Broker struct manages incoming connections, message dispatching, and message parsing. It maintains queues for different message topics and handles various types of messages including publishing (Pub), subscribing (Sub), and queue setup (Set).  
- Each topic is managed by a queue.  
- Each request is managed by a worker, which parses the message and sends it to broker's query channel.
- The query channel listener gets the queue associated with the topic and sends a message to its corresponding channel (pubC, pullC, subC).
- Finally, the queue listener handles the message (pushing to the queue, adding a subscriber or pulling and sending messages to subscribers)
- The broker has a timer which sends messages to queues in order to pull data periodically

### Protocol
A custom naive, highly insecure, no error tolerant  protocol
#### Actions
- Set: sets a topic
- Pub: publish a message to a topic
- Sub: Subscribes to a topic queue

## Client Package
A TCP/HTTP client to send packages to the broker
- Reads an HTTP POST
- Parses the body to a valid message
- Sends it to the broker via TCP

  
