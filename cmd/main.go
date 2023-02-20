package main

import (
	"go-docker/cli"
	"go-docker/pkg"
)

func main() {
	c := cli.NewCli()
	docker := pkg.NewDocker(c.Image)
	if !c.Remove {
		docker.Start(c.Cmd, c.Env, c.Name)
	} else {
		docker.StartAndRemoveContainer(c.Cmd, c.Env, c.Name)
	}
}
