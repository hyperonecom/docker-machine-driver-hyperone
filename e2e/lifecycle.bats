#!/usr/bin/env bats
#
# This tests the HyperOne driver for Docker machine. The teardown function will
# delete any VM and disk with the text "machinebats" within the name.

# Required parameters
: ${HYPERONE_TOKEN:?}
: ${HYPERONE_PROJECT:?}
command -v h1 >/dev/null 2>&1 || {
    echo "'h1' must be installed" >&2
    exit 1
}

USER_VARS="--hyperone-disk-name machinebats-os-disk"
USER_VARS="${USER_VARS} --hyperone-token ${HYPERONE_TOKEN}"
USER_VARS="${USER_VARS} --hyperone-project ${HYPERONE_PROJECT}"

hyperone_has_resource() {
    h1 $1 list --project-select=${HYPERONE_PROJECT} --output=tsv | grep "machinebats" | grep "$2" -c
}

hyperone_has_machine() {
    docker-machine ls | grep "machinebats" | grep "$2" -c
}

hyperone_vm_has_ip() {
    h1 vm nic list --project-select=${HYPERONE_PROJECT} --output=tsv --vm "$1" | grep "$2" -c
}

machine_status () {
    docker-machine status "$1" | grep "$2" -c
}


teardown() {
    h1 vm list --project-select=${HYPERONE_PROJECT} --output=tsv \
        | grep machinebats \
        | awk '{print $1}' \
        | xargs -r -n 1 h1 vm delete --project-select=${HYPERONE_PROJECT} --yes --vm
    h1 disk list --project-select=${HYPERONE_PROJECT} --output=tsv \
        | grep machinebats \
        | awk '{print $1}' \
        | xargs -r -n 1 h1 disk delete --project-select=${HYPERONE_PROJECT} --yes --disk
    docker-machine ls -q | grep 'machinebats' | xargs -r docker-machine rm -y
}

@test "hyperone: create" {
    run docker-machine create --driver hyperone ${USER_VARS} machinebats-minimal
    [ "$status" -eq 0 ]
    [ "$(hyperone_has_resource "vm" "minimal")" -eq 1 ]
    [ "$(hyperone_has_machine "machinebats-minimal")" -eq 1 ]
    [ "$(machine_status machinebats-minimal Running)" -eq 1 ]

}

@test "hyperone: docker-machine env" {
    docker-machine create --driver hyperone ${USER_VARS} machinebats-env
    eval $(docker-machine env machinebats-env --shell sh)
    docker info
    [ "$(hyperone_vm_has_fqdn "machinebats-env" $(docker-machine ip machinebats-env))" -eq 1 ]
}

@test "hyperone: docker-machine ip" {
    run docker-machine create --driver hyperone ${USER_VARS} machinebats-env
    run docker-machine ip machinebats-env
    [ "$(hyperone_vm_has_ip "machinebats-env" $(docker-machine ip machinebats-env))" -eq 1 ]
}

@test "hyperone: docker-machine stop" {
    run docker-machine create --driver hyperone ${USER_VARS} machinebats-stop
    run docker-machine stop machinebats-stop
    [ "$(hyperone_has_resource "vm" "stop")" -eq 1 ]
    [ "$(machine_status machinebats-stop Stopped)" -eq 1 ]
}

@test "hyperone: docker-machine restart" {
    run docker-machine create --driver hyperone ${USER_VARS} machinebats-restart
    run docker-machine restart machinebats-restart
    [ "$(machine_status machinebats-restart Running)" -eq 1 ]
}

@test "hyperone: docker-machine rm" {
    docker-machine create --driver hyperone ${USER_VARS} machinebats-rm
    [ "$(hyperone_has_resource "disk" "os-disk")" -eq 1 ]
    run docker-machine rm -y machinebats-rm
    [ "$(hyperone_has_resource "vm" "rm")" -eq 0 ]
    [ "$(hyperone_has_resource "disk" "os-disk")" -eq 0 ]
    [ "$(docker-machine ls | grep machinebats-rm -c)" -eq 0 ]
}

@test "hyperone: docker-machine kill" {
    run docker-machine create --driver hyperone ${USER_VARS} machinebats-kill
    [ "$(hyperone_has_resource "vm" "kill")" -eq 1 ]
    run docker-machine kill machinebats-kill
    [ "$(hyperone_has_resource "vm" "kill")" -eq 1 ]
    [ "$(machine_status machinebats-kill Stopped)" -eq 1 ]
}
