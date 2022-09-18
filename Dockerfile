FROM golang:1.17-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download


# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o server ./cmd/app

#
# -- build minimal image --
#
FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*


RUN ./cloud_sql_proxy -instances=$INSTANCE_CONNECTION_NAME=tcp:5432  &

COPY --from=builder /app/server /app/server

# Run the web service on container startup.
WORKDIR /app
CMD ["/app/server"]