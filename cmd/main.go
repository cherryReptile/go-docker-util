package main

import (
	"go-docker/cli"
	"go-docker/internal/sys"
	"go-docker/pkg"
	"os"
	"os/signal"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs)
	go sys.HandleSignal(sigs)
	c := cli.NewCli()
	docker := pkg.NewDocker(c)
	if !c.Remove {
		docker.Start()
	} else {
		docker.StartAndRemoveContainer()
	}
}
