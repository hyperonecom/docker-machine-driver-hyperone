FROM golang:alpine
ENV H1_CLI_VERSION="v1.4.0"
ENV MACHINE_VERSION="v0.16.1"

# Setup
RUN apk add curl bats git make gcc musl-dev findutils grep docker
RUN base="https://github.com/docker/machine/releases/download/${MACHINE_VERSION}" \
&& curl -L $base/docker-machine-$(uname -s)-$(uname -m) >/usr/local/bin/docker-machine \
&& chmod +x /usr/local/bin/docker-machine
RUN curl -s -L "https://github.com/hyperonecom/h1-cli/releases/download/${H1_CLI_VERSION}/h1-alpine" -o /bin/h1 \
&& chmod +x /bin/h1
# Build
WORKDIR /src
ADD ./ /src/
RUN go build
ENV PATH=/src/:$PATH

# Run
CMD ["docker-machine", "create", "-d", "hyperone", "--help"]