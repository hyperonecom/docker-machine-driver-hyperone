FROM node:alpine as cli_builder
RUN apk add --no-cache libstdc++
RUN wget https://github.com/hyperonecom/h1-cli/archive/develop.tar.gz
RUN tar xvzf master.tar.gz
WORKDIR /h1-cli-master
RUN npm install .
RUN npm run build

FROM golang:alpine
# Setup
RUN apk add --no-cache libstdc++ curl bats git make gcc musl-dev findutils grep docker
RUN base="https://github.com/docker/machine/releases/latest/download" \
&& curl -L $base/docker-machine-$(uname -s)-$(uname -m) >/usr/local/bin/docker-machine \
&& chmod +x /usr/local/bin/docker-machine
COPY --from=cli_builder /h1-cli-master/dist/h1-alpine /bin/h1
# Build
WORKDIR /src
ADD ./ /src/
RUN go build
ENV PATH=/src/:$PATH

# Run
CMD ["docker-machine", "create", "-d", "hyperone", "--help"]
