FROM golang:alpine
ENV H1_CLI_VERSION="v1.4.0"

# Setup
RUN apk add curl bats git make gcc musl-dev findutils grep docker
RUN curl -s https://api.github.com/repos/docker/machine/releases/latest | \
		grep "browser_download_url.*Linux-$(uname -m)" | \
		cut -d : -f 2,3 | \
		tr -d '"' | \
		head -n 1 | \
		xargs curl -L >/tmp/docker-machine \
&&  install /tmp/docker-machine /usr/local/bin/docker-machine
RUN curl -s -L "https://github.com/hyperonecom/h1-cli/releases/download/${H1_CLI_VERSION}/h1-alpine" -o /bin/h1 \
&& chmod +x /bin/h1
# Build
WORKDIR /src
ADD ./ /src/
RUN go build
ENV PATH=/src/:$PATH

# Run
CMD ["docker-machine", "create", "-d", "hyperone", "--help"]