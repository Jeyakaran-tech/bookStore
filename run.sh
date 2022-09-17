#!/bin/sh

# Start the proxy
./cloud_sql_proxy -instances=$INSTANCE_CONNECTION_NAME=tcp:5432  &

# wait for the proxy to spin up
sleep 10

# Start the server
WORKDIR /app
CMD ["/app/server"]
./server