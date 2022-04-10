FROM golang:1.18.0-alpine3.15

# Metadata
LABEL name="online-electronic-voting-machine"

EXPOSE 8080

# Update container image and setup sqlite, make, python3
RUN apk add --update \
    sqlite make python3

# Setup workdir and build code
COPY . /app
WORKDIR /app
RUN make build-server && \
    python3 scripts/setup_sqlite_schema.py

ENTRYPOINT [ "./server" ]
