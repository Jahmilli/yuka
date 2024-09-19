# yuka

## Overview

Yuka is designed to be an alternative to [ngrok](https://ngrok.com/) providing a means to enable network access to local applications. There are a few separate components required to enable this:

- Server: The server will be exposed publicly and will be the entrypoint for public network connections to your agent. 
- Client/Daemon/Agent: The client will be accessed via a binary providing a CLI that can be used to start/stop tunnels. 


## Running locally

Currently using the following version of go

```bash
go version go1.23.1 darwin/arm64
```

Run local infra: `make setup`


Build and run yukactl: `make yukactl-build && ./bin/yukactl-mac-arm64 help`
Build and run yuka-apiserver: `make build && ./bin/yuka-api-server-mac-arm64`


## TODO


Phase 1:
- Initial goal will be building a server with that you can hit and it forwards requests to client which is proxied to x port and then gets the response back

Phase 2:
- Support subdomains, we can probably do this locally and it would be good to try. However, as part of this would be good to deploy it and see if it works. Deploying will require:
  - Domain name
  - Instance
  - TLS setup
- Support multiple subdomains per host

Phase 3: Support a distributed design, i.e can have multiple instances of server behind loadbalancers working 

- Build database design
- Create tunnel
    - [x] Build start command in ctl. This doesn't do anything for now and just logs foo
    - [x] Build route on apiserver that ctl interacts with
    - Setup proxy that is initialised on start and forwards requests from server to the proxied endpoint
    - Validate messages can go back/forth between client and server
    - Build "detached" mode for start
- Build stop command (only required if in "detached" mode)
- Build status command (only required if in "detached" mode)
- Setup authentication
  -  Setup database design for authentication
  -  Build basic authentication between client and server (i.e token used in config)


Read on TCP 1, 2, 3
Read on QUIC
Understand what to use between TCP, QUIC, Websockets, gRPC for connection between client/server.


## Improvements

- `streaming_handler.go` is coupled to websockets and returns a websocket connection which isn't ideal. Would be ideal to wrap the websocket connections in an interface that returns read/write and keeps it as simple as that. This way we can later move to QUIC or alternative streaming implementations without much change.


## Conventions

### Documentation

- All language should be in American English
- ReadMes and any other documentation regardless of title/paragraph etc should always be sentence cased.
  - Comments are an exception to this where we prefix a comment name with the name of the function

### Logging

Current prefer to add `slogger` to each handler. This is a Zap sugared logger and just simplifies logging. Whilst this reduces performance slightly, it simplifies development with logging for now and is considered okay. 

Logs should all start with **lower case**

Use `%v` for errors and always add the error at the end of the log after a `:`, i.e `("error occured in handler: %v, err")`



