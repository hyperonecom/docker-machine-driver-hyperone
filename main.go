package main

import (
	"github.com/docker/machine/libmachine/drivers/plugin"
	hyperone "github.com/hyperonecom/docker-machine-driver-hyperone/driver"
)

func main() {
	plugin.RegisterDriver(hyperone.NewDriver("", ""))
}
