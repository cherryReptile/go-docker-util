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
		name := viper.Get("name")
		validateConfig(name, "name", c)
	}
	if viper.IsSet("image") {
		image := viper.Get("image")
		validateConfig(image, "image", c)
	}
	if viper.IsSet("env") {
		env := viper.Get("env")
		validateConfig(env, "env", c)
	}
	if viper.IsSet("cmd") {
		cmd := viper.Get("cmd")
		validateConfig(cmd, "cmd", c)
	}

	return c
}

func validateConfig(item interface{}, itemName string, c *Cli) {
	switch item.(type) {
	case []interface{}:
		iSlice := item.([]interface{})
		for _, v := range iSlice {
			switch v.(type) {
			case string:
				if itemName == "env" {
					c.ConfigList.Env = append(c.ConfigList.Env, v.(string))
				}
				if itemName == "cmd" {
					c.ConfigList.Cmd = append(c.ConfigList.Cmd, v.(string))
				}
			default:
				log.Fatal(errors.New("unknown array type in config"))
			}
		}
	case string:
		if itemName == "image" {
			c.ConfigList.Image = item.(string)
		}
		if itemName == "name" {
			c.Name = item.(string)
		}
	}
}
