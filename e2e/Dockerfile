FROM node:alpine as cli_builder
RUN apk add --no-cache libstdc++
RUN wget https://github.com/hyperonecom/h1-cli/archive/develop.tar.gz
RUN tar xvzf develop.tar.gz
WORKDIR /h1-cli-develop
RUN npm install .
RUN npx pkg -c package.json -t "node12-alpine" -o "./dist/h1" "./bin/h1"

FROM golang:alpine
# Setup
COPY --from=cli_builder /h1-cli-develop/dist/h1 /bin/h1
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
