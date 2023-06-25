# talk-walk

Let's try to explore different ways to organize the communication between clients and servers. Here I want to expose some ideas.

## Broadcast

Let's use one server, which every 10 seconds sends to all connected to it clients a message, like "Hi there!".

No matter what time any given client connects to the server, the client eventually will receive the message.

```
Time:   -------------------------------------------------------------------------- ... -->
Server: --------- Msg! --------- Msg! --------- Msg! --------- Msg! --------- Msg! ... -->
...

Client1 - Plug -- Rec! --------- Rec! --------- Rec! -- Unplug X
Client2 ------------ Plug ------ Rec! --------- Rec! --------- Rec! ----- Unplug X
...
```

So, no matter when the client connects, within the following 10 seconds it should receive the message "Hi there!"
