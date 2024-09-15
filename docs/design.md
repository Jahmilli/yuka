# Design

There's 2 major components:

- yukactl: This is responsible for enabling users to expose their host to the public internet. This will be through a binary that they will install and can execute, i.e `yukactl start`
- yuka-server: This provides the gateway between the internet and hosts that are interacting with it via `yukactl`. 


## Yukactl

### Config

Yukactl can take in config via a few different ways including as environment variables, flags or a YAML config file. 

When reading in configuration, yukactl will read in configuration in the following priority (greater number takes priority over less number).

1. Config file
2. Environment variables
3. Flags

### Commands

To see available commands, run `yukactl --help`

## Yuka server

