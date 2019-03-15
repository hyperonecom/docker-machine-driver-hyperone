FROM golang:alpine
ENV H1_CLI_VERSION="v1.4.0"

# Setup
RUN apk add curl bats git make gcc musl-dev findutils grep docker
RUN base=https://github.com/docker/machine/releases/download/v0.16.0 && \
    curl -L $base/docker-machine-$(uname -s)-$(uname -m) >/tmp/docker-machine && \
	  install /tmp/docker-machine /usr/local/bin/docker-machine
RUN curl -s -L "https://github.com/hyperonecom/h1-cli/releases/download/${H1_CLI_VERSION}/h1-alpine" -o /bin/h1 \
&& chmod +x /bin/h1
# Build
WORKDIR /src
ADD ./ /src/
RUN go build
ENV PATH=/src/:$PATH

# Run
CMD ["docker-machine","create","--help"]