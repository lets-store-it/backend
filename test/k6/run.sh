#!/bin/bash

docker pull grafana/k6:1.0.0

docker run --rm -i \
    --network=host \
    -v ${PWD}/scripts:/scripts \
    grafana/k6:1.0.0 run /scripts/highload.ts
