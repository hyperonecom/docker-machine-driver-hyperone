# docker-machine-driver-hyperone

HyperOne Driver Plugin for docker-machine.

## Setup

### Go
```shell
# install latest (git) version of docker-machine-driver-hyperone in your $GOPATH/bin (depends on Golang and docker-machine)
$ go get -u github.com/hyperonecom/docker-machine-driver-hyperone
```

### Binary

You can find sources and pre-compiled binaries [here](https://github.com/hyperonecom/docker-machine-driver-hyperone/releases/latest)

```shell
# Download the binary (this example downloads the binary for darwin amd64)
$ wget https://github.com/hyperonecom/docker-machine-driver-hyperone/releases/download/v0.0.1/docker-machine-driver-hyperone_0.0.1_darwin_amd64.zip
$ unzip unzip docker-machine-driver-hyperone_0.0.1_darwin_amd64.zip

# Make it executable and copy the binary in a directory accessible with your $PATH
$ chmod +x docker-machine-driver-hyperone
$ sudo cp docker-machine-driver-hyperone /usr/local/bin/
```

# Usage

### 1. driver helper
```shell
$ docker-machine create -d hyperone -h
Usage: docker-machine create [OPTIONS] [arg...]

Create a machine

Description:
   Run 'docker-machine create --driver name --help' to include the create flags for that driver in the help text.

Options:
   
   --driver, -d "virtualbox"										Driver to create machine with. [$MACHINE_DRIVER]
   --engine-env [--engine-env option --engine-env option]						Specify environment variables to set in the engine
   --engine-insecure-registry [--engine-insecure-registry option --engine-insecure-registry option]	Specify insecure registries to allow with the created engine
   --engine-install-url "https://get.docker.com"							Custom URL to use for engine installation [$MACHINE_DOCKER_INSTALL_URL]
   --engine-label [--engine-label option --engine-label option]						Specify labels for the created engine
   --engine-opt [--engine-opt option --engine-opt option]						Specify arbitrary flags to include with the created engine in the form flag=value
   --engine-registry-mirror [--engine-registry-mirror option --engine-registry-mirror option]		Specify registry mirrors to use [$ENGINE_REGISTRY_MIRROR]
   --engine-storage-driver 										Specify a storage driver to use with the engine
   --hyperone-disk-name "os-disk"									HyperOne VM OS Disk Name [$HYPERONE_DIKE_NAME]
   --hyperone-disk-size "20"										HyperOne VM OS Disk Size [$HYPERONE_DIKE_SIZE]
   --hyperone-disk-type "ssd"										HyperOne VM OS Disk Type [$HYPERONE_DIKE_TYPE]
   --hyperone-image "debian"										HyperOne Image [$HYPERONE_IMAGE]
   --hyperone-project 											HyperOne Project [$HYPERONE_PROJECT]
   --hyperone-ssh-user "guru"										SSH Username [$HYPERONE_SSH_USER]
   --hyperone-token 											HyperOne Token [$HYPERONE_TOKEN]
   --hyperone-type "a1.micro"										HyperOne VM Type [$HYPERONE_TYPE]
   --swarm												Configure Machine to join a Swarm cluster
   --swarm-addr 											addr to advertise for Swarm (default: detect and use the machine IP)
   --swarm-discovery 											Discovery service to use with Swarm
   --swarm-experimental											Enable Swarm experimental features
   --swarm-host "tcp://0.0.0.0:3376"									ip/socket to listen on for Swarm master
   --swarm-image "swarm:latest"										Specify Docker image to use for Swarm [$MACHINE_SWARM_IMAGE]
   --swarm-join-opt [--swarm-join-opt option --swarm-join-opt option]					Define arbitrary flags for Swarm join
   --swarm-master											Configure Machine to be a Swarm master
   --swarm-opt [--swarm-opt option --swarm-opt option]							Define arbitrary flags for Swarm master
   --swarm-strategy "spread"										Define a default scheduling strategy for Swarm
   --tls-san [--tls-san option --tls-san option]							Support extra SANs for TLS certs
```

### 2. Create your machine

```shell
$ docker-machine create --driver hyperone --hyperone-token TOKEN --hyperone-project PROJECT vm
Running pre-create checks...
Creating machine...
(vm) Creating HyperOne VM...
Waiting for machine to be running, this may take a few minutes...
Detecting operating system of created instance...
Waiting for SSH to be available...
Detecting the provisioner...
Provisioning with debian...
Copying certs to the local machine directory...
Copying certs to the remote machine...
Setting Docker configuration on the remote daemon...
Checking connection to Docker...
Docker is up and running!
To see how to connect your Docker Client to the Docker Engine running on this virtual machine, run: docker-machine env vm
```