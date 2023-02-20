package main

import (
	"go-docker/cli"
	"go-docker/pkg"
)

func main() {
	c := cli.NewCli()
	docker := pkg.NewDocker(c)
	if !c.Remove {
		docker.Start()
	} else {
		docker.StartAndRemoveContainer()
	}
}
