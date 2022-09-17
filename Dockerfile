FROM golang:1.12

WORKDIR /app
COPY . .

# build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /build/server ./cmd/app

# download the cloudsql proxy binary
RUN wget https://dl.google.com/cloudsql/cloud_sql_proxy.linux.amd64 -O /build/cloud_sql_proxy
RUN chmod +x /build/cloud_sql_proxy

# copy the wrapper script and credentials
COPY run.sh /build/run.sh
COPY credentials.json /build/credentials.json

#
# -- build minimal image --
#
FROM alpine:latest

WORKDIR /root

# add certificates
RUN apk --no-cache add ca-certificates

# copy everything from our build folder
COPY --from=0 /build .

CMD ["./run.sh"]