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