package cli

import (
	"errors"
	"fmt"
	"github.com/docker/docker/api/types/container"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Cli struct {
	Remove     bool
	Name       string
	Logs       bool
	ConfigList *container.Config
}

func NewCli() *Cli {
	remove := flag.Bool("rm", false, "remove container after stopping")
	image := flag.String("i", "", "docker image(required)")
	name := flag.String("n", "default", "container name")
	env := flag.StringArray("e", nil, "env in key=value format")
	cmd := flag.StringArray("c", nil, "command that will be executed when container is starting")
	logs := flag.Bool("l", false, "container logs")
	config := flag.String("f", "", "config file for creating image and container")
	flag.Parse()

	if *image == "" && *config == "" {
		log.Fatal(errors.New("required parameter is missing and config file doesn't set: -i"))
	}

	if *config != "" {
		setConfig(*config)
		c := readConfig()
		c.Remove = *remove
		c.Logs = *logs
		return c
	}

	return &Cli{
		Remove: *remove,
		Name:   *name,
		Logs:   *logs,
		ConfigList: &container.Config{
			Image: *image,
			Env:   *env,
			Cmd:   *cmd,
		},
	}
}

func setConfig(configFile string) {
	_, err := os.Stat(configFile)
	if err == nil {
		fmt.Println("Using User Specified Configuration file!")
		viper.SetConfigFile(configFile)
	} else {
		viper.SetConfigName(configFile)
		viper.AddConfigPath("$HOME")
		viper.AddConfigPath(".")
	}
}

func readConfig() *Cli {
	c := new(Cli)
	c.ConfigList = new(container.Config)

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Using config file: %s\n", viper.ConfigFileUsed())

	if viper.IsSet("name") {
		c.Name = viper.GetString("name")
	}
	if viper.IsSet("image") {
		c.ConfigList.Image = viper.GetString("image")
	}
	if viper.IsSet("env") {
		c.ConfigList.Env = viper.GetStringSlice("env")
	}
	if viper.IsSet("cmd") {
		c.ConfigList.Cmd = viper.GetStringSlice("cmd")
	}

	return c
}
