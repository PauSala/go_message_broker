# go_message_broker
A toy message_broker implemented in GO just for fun
## Broker Package
This package provides functionality for a broker server that handles messaging using a custom protocol over TCP.
It uses a worker pool to handle conections and implements topics as message queues. 

### Overview
The Broker struct manages incoming connections, message dispatching, and message parsing. It maintains queues for different message topics and handles various types of messages including publishing (Pub), subscribing (Sub), and queue setup (Set).

### Protocol
A custom naive, highly insecure, no error tolerant  protocol
#### Actions
- Set: sets a topic
- Pub: publish a message to a topic
- Sub: for now it does nothing

## Client Package
A TCP/HTTP client to send packages to the broker
- Reads  POST
- Parses the body to a valid message
- Sends it to the broker 

  
