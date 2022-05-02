FROM golang:1.18.0-buster

# Metadata
LABEL name="online-electronic-voting-machine"

EXPOSE 8080
EXPOSE 9000

# Update container image and setup sqlite, build-essential, and libsodium, python3
RUN apt-get update && apt-get install -y sqlite3 build-essential libsodium-dev python3

# Setup workdir and build code
COPY . /app
WORKDIR /app
RUN make build-server

ENTRYPOINT [ "./server" ]
