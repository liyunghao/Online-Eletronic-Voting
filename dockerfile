FROM golang:1.18.0-alpine3.15

# Metadata
LABEL name="online-electronic-voting-machine"

EXPOSE 8080

# Update container image and setup libsodium
RUN apk add --update --no-cache \
    make \
    ca-certificates \
    curl \
    libsodium \
    openssl

# Setup workdir and build code
COPY . /app
WORKDIR /app
RUN make build-server

ENTRYPOINT [ "./server" ]
