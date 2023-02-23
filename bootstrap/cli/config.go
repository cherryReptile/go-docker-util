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

type Config struct {
	Remove     bool
	Name       string
	Logs       bool
	ConfigList *container.Config
}

func NewConfigFromCli() *Config {
	c := new(Config)
	c.ConfigList = new(container.Config)
	setFlags()

	if err := viper.BindPFlags(flag.CommandLine); err != nil {
		log.Fatal(err)
	}

	if viper.GetString("img") == "" && viper.GetString("config") == "" {
		log.Fatal(errors.New("required parameter is missing and config file doesn't set"))
	}

	c.Remove = viper.GetBool("remove")
	c.Logs = viper.GetBool("logs")

	if viper.GetString("config") != "" {
		setConfigFile(viper.GetString("config"))
		readConfig(c)
		return c
	}

	if err := viper.BindPFlags(flag.CommandLine); err != nil {
		log.Fatal(err)
	}

	c.Name = viper.GetString("n")
	c.ConfigList.Image = viper.GetString("img")
	c.ConfigList.Env = viper.GetStringSlice("environment")
	c.ConfigList.Cmd = viper.GetStringSlice("command")

	return c
}

func setFlags() {
	flag.BoolP("remove", "r", false, "remove container after stopping")
	flag.StringP("img", "i", "", "docker image(required)")
	flag.String("n", "default", "container name")
	flag.StringArrayP("environment", "e", nil, "env in key=value format")
	flag.StringArrayP("command", "c", nil, "command that will be executed when container is starting")
	flag.BoolP("logs", "l", false, "container logs")
	flag.StringP("config", "f", "", "config file for creating image and container")
	flag.Parse()
}

func setConfigFile(configFile string) {
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

func readConfig(c *Config) {
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
}
