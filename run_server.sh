#!/bin/bash

podman run -d --rm --name redis -p 6379:6379 redis
podman run -it --rm --name labstats -p 50051:50051 registry.botnet/labstats:dev
