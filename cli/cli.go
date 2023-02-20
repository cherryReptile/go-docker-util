package cli

import (
	flag "github.com/spf13/pflag"
)

type Cli struct {
	Remove bool
	Image  string
	Name   string
	Env    []string
	Cmd    []string
}

func NewCli() *Cli {
	remove := flag.Bool("rm", true, "remove container after stopping")
	image := flag.String("i", "", "docker image")
	name := flag.String("n", "default", "container name")
	env := flag.StringArray("e", nil, "env in key=value format")
	cmd := flag.StringArray("c", nil, "command that will be executed when container is starting")
	flag.Parse()

	return &Cli{
		Remove: *remove,
		Image:  *image,
		Name:   *name,
		Env:    *env,
		Cmd:    *cmd,
	}
}
