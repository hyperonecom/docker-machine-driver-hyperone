FROM golang:alpine
# Setup
RUN apk add --repository "http://dl-cdn.alpinelinux.org/alpine/edge/testing" h1-cli
RUN apk add --no-cache libstdc++ curl bats git make gcc musl-dev findutils grep docker
RUN base="https://github.com/docker/machine/releases/latest/download" \
&& curl -L $base/docker-machine-$(uname -s)-$(uname -m) >/usr/local/bin/docker-machine \
&& chmod +x /usr/local/bin/docker-machine
# Build
WORKDIR /src
ADD ./ /src/
RUN go build
ENV PATH=/src/:$PATH

# Run
CMD ["docker-machine", "create", "-d", "hyperone", "--help"]
