FROM golang:1.17-buster as builder

WORKDIR /app
COPY go.* ./
RUN go mod download


# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -v -o server ./cmd/app

# download the cloudsql proxy binary
RUN wget https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64 -O cloud_sql_proxy
RUN chmod +x cloud_sql_proxy
COPY build/credentials.json credentials.json
#
# -- build minimal image --
#
FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*


RUN ./cloud_sql_proxy -instances=$INSTANCE_CONNECTION_NAME=tcp:5432 -credential_file=credentials.json &

COPY --from=builder /app/server /app/server

# Run the web service on container startup.
WORKDIR /app
CMD ["/app/server"]