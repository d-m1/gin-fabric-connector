# Gin Fabric Connector

> This project was created as a hobby project and only for learning and experimental purposes.

## Table of Contents

1. [Introduction](#introduction)
1. [Usage](#usage)
1. [Environment Variables](#environment-variables)
1. [Running the server](#running-the-server)
1. [TODOs](#todos)

## Introduction

Gin Fabric Connector is a development project created to facilitate interactions with the Hyperledger Fabric Test-Network. It provides a straightforward REST API built using Gin and Auth0, enabling users to easily send transactions to the chaincodes, as an alternative to the Fabric CLI tools. It also implements the fabric-gateway to communicate with peers.

## Usage

As this project is meant to be used with the fabric test-network, you'll need to deploy the network and set some environment variables before starting the server.

### Start the test-network

This is generally done using the network.sh script:

```
# From your test-network folder
./network.sh up
```

### Environment Variables

- `FABRIC_CONNECTOR_BASEDIR`: Your fabric-samples/test-network repository folder.
- `FABRIC_CONNECTOR_AUTH0_AUDIENCE`: Identifier of your app in Auth0.
- `FABRIC_CONNECTOR_AUTH0_DOMAIN`: Your Auth0 domain.

Every config value can be set independently if your setup is not standard (i.e. `PEER_TLS_CERT`)

### Running the server

To start the server, run the following command from the project root:

```bash
go run cmd/main.go
```

### TODOs

1. Async TXs
1. Private Data TXs
1. Testing
1. Linting
1. Dockerfile
1. Toggle auth
1. Error handling
1. Logging
1. Dynamic config
1. K8s template
1. Versioning