FROM golang:alpine
# Setup
RUN apk add curl bats git make gcc musl-dev findutils grep docker
RUN base="https://github.com/docker/machine/releases/latest/download" \
&& curl -L $base/docker-machine-$(uname -s)-$(uname -m) >/usr/local/bin/docker-machine \
&& chmod +x /usr/local/bin/docker-machine
RUN curl -s -L "https://github.com/hyperonecom/h1-cli/releases/latest/download/h1-alpine" -o /bin/h1 \
&& chmod +x /bin/h1
# Build
WORKDIR /src
ADD ./ /src/
RUN go build
ENV PATH=/src/:$PATH

# Run
CMD ["docker-machine", "create", "-d", "hyperone", "--help"]
