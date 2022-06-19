# This Dockerfile is for development use.
# One can invoke this dev build in the following way:
#
# - Vscode Remote Contianer
# - Self Build
#   - `docker build -t <image_name> -f build/package/dev.dockerfile .`
#   - `docker run -it -p 8080:8080 -p 9000:9000 -v $(pwd):/app -w /app <image_name>`

FROM golang:1.18.0-buster

# Metadata
LABEL name="online-electronic-voting-machine"

EXPOSE 8080
EXPOSE 9000

# Update container image and setup sqlite, build-essential, and libsodium, python3
RUN apt-get update && apt-get install -y \
    sqlite3 build-essential libsodium-dev python3 \
    zsh curl

RUN sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)" && \
    chsh -s $(which zsh)
