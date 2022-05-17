# gochat
A basic, scalable, dockerized chat room written in go. This is patched from various
web resources, mostly from [gorilla web sockets's chat example](https://github.com/gorilla/websocket/examples/chat)


## References
### Go
https://gobyexample.com
https://go.dev

### Websockets in Go
https://github.com/gorilla/websocket/examples/chat
https://yalantis.com/blog/how-to-build-websockets-in-go/

### Pub Sub
https://www.compose.com/articles/redis-go-and-how-to-build-a-chat-application/
https://dev.to/franciscomendes10866/using-redis-pub-sub-with-golang-mf9


## Run Instructions
- docker-compose up ( and hope it works ;)
  - Default setup is couple of chat rooms, default arguments in docker-compose file should be enough to run on MacOS
- This is tested on MacOS, some settings are particular to Mac's docker setup (localhost resolution)

## Comments
- The redis service should be part of an internal/private network
  - Currently redis is using unauthenticated connection. In a real setting, this must be password protected.
- The redis client is basic. One per hub/chatroom. This is room of using a "java equivalent" of a thread pool of clients.
I did not investigate this in enough detail.
- Testing was done manually by running the system. For a production product: I would go for following layers
  - Integration tests: boot a chatroom and run client commands from a script which does assertions. This framework
  can be used to test most of the execution paths.