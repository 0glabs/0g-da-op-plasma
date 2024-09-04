# Optimism Alt-DA x 0G DA:

## Overview:

This repository implements a 0gDA `da-server` for Alt-DA mode using generic
commitments.

The `da-server` connects to a 0G DA client, which runs as a sidecar process.

0G DA da-server accepts the following flags for 0g DA storage using
[0G DA OPEN API](https://docs.0g.ai/0g-doc/docs/0g-da/rpc-api/api-1)

````
    --zg.server    (default: "localhost:51001") 
        0G DA client server endpoint
    
    --addr
        server listening address
    
    --port
        server listening port
````


## Deployment

### Build DA Server

```bash
    make da-server
```

### Run DA Server
```bash
    ./bin/da-server --addr 127.0.0.1 --port 3100 --zg.server 127.0.0.1:51001
```

For guidance on setting up a 0G DA client, refer to the [documentation](https://docs.0g.ai/0g-doc/run-a-node/da-client).


## Run with Docker
Build the Docker image

```bash
docker build -t 0g-da-op-plasma .
```

Run the Docker container

Adjust commands and parameters as required for your setup:

```bash
docker run -p 3100:3100 0g-da-op-plasma:latest da-server --addr 0.0.0.0 --port 3100 --zg.server 127.0.0.1:51001
```

## Guidance to run OP Stack with 0G DA
[How to Use the OP Stack with 0G DA](./OP%20Stack%20integration.md)
