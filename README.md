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
  - Default setup is couple of chat rooms, default arguments in docker-compose file should be enough to run.
  - REDIS_HOST: 
    - HOST part of connection string, port is hard coded to Redis's default port since image has it baked in.  
  - BACKEND_LISTENER1_PORT, BACKEND_LISTENER2_PORT:
    - Ports where respective backend should listen for WS connections. The host part is set to `0.0.0.0`
  - ONE EXAMPLE
    - BACKEND_LISTENER1_PORT=8085 BACKEND_LISTENER2_PORT=8086 REDIS_HOST=host.docker.internal docker-compose up
- This is tested on MacOS, some settings are particular to Mac's docker setup (localhost resolution)

## Comments/Further work
- The user input should be validated/sanitzied
- Add additional layer of security by adding authentication (and perhaps ACLs) to redis and chat room integration.
- Improve UI.
- The redis service should be part of an internal/private network
  - Currently redis is using unauthenticated connection. In a real setting, this must be password protected.
- The redis client is basic. One per hub/chatroom. This might be room for using a "java equivalent" of a thread pool of 
clients. I did not investigate this in enough detail.
- Testing was done manually by running the system and running client commands. For a production product: I would go for 
following layers
  - Integration tests (in a CI pipeline): boot a chatroom and run client commands from a script which does assertions. 
This framework can be used to test most of the execution paths.