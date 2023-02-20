package cli

import (
	"errors"
	"github.com/docker/docker/api/types/container"
	flag "github.com/spf13/pflag"
	"log"
)

type Cli struct {
	Remove     bool
	Name       string
	ConfigList *container.Config
}

func NewCli() *Cli {
	remove := flag.Bool("rm", true, "remove container after stopping")
	image := flag.String("i", "", "docker image(required)")
	name := flag.String("n", "default", "container name")
	env := flag.StringArray("e", nil, "env in key=value format")
	cmd := flag.StringArray("c", nil, "command that will be executed when container is starting")
	flag.Parse()

	if *image == "" {
		log.Fatal(errors.New("required parameter is missing: -i"))
	}

	return &Cli{
		Remove: *remove,
		Name:   *name,
		ConfigList: &container.Config{
			Image: *image,
			Env:   *env,
			Cmd:   *cmd,
		},
	}
}
