package main

import (
	"go-docker/pkg"
)

func main() {
	docker := pkg.NewDocker()
	docker.Start()
}
