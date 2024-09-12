# yuka

## Overview

Yuka is designed to be an alternative to [ngrok](https://ngrok.com/) providing a means to enable network access to local applications. There are a few separate components required to enable this:

- Server: The server will be exposed publicly and will be the entrypoint for public network connections to your agent. 
- Client/Daemon/Agent: The client will be accessed via a binary providing a CLI that can be used to start/stop tunnels. 
