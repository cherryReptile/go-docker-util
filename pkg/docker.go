package pkg

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"go-docker/bootstrap/cli"
	"io"
	"log"
	"os"
)

type Docker struct {
	CliInfo     *cli.Config
	ContainerID string
	Client      *client.Client
	Context     context.Context
}

func NewDocker(cli *cli.Config) *Docker {
	return &Docker{
		CliInfo: cli,
		Client:  &client.Client{},
		Context: context.Background(),
	}
}

func (d *Docker) Start() {
	var containerID string
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	d.Client = dockerClient
	defer d.Client.Close()

	reader, err := dockerClient.ImagePull(d.Context, d.CliInfo.ConfigList.Image, types.ImagePullOptions{})
	if err != nil {
		log.Fatal(err)
	}
	if _, err = io.Copy(os.Stdout, reader); err != nil {
		log.Fatal(err)
	}

	c := d.getByContainerName(d.CliInfo.Name)
	if c.ID == "" {
		resp, err := d.createContainer()
		if err != nil {
			log.Fatal(err)
		}
		containerID = resp.ID
	} else if c.ID != "" {
		containerID = c.ID
	}

	d.ContainerID = containerID
	if err = d.Client.ContainerStart(d.Context, containerID, types.ContainerStartOptions{}); err != nil {
		log.Fatal(err)
	}

	out, err := d.Client.ContainerLogs(d.Context, containerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Details: true})
	if err != nil {
		log.Fatal(err)
	}

	if _, err = stdcopy.StdCopy(os.Stdout, os.Stderr, out); err != nil {
		log.Fatal(err)
	}

	if d.CliInfo.Logs {
		statusCh, errCh := d.Client.ContainerWait(d.Context, containerID, container.WaitConditionNotRunning)
		select {
		case err = <-errCh:
			if err != nil {
				log.Fatal(err)
			}
		case status := <-statusCh:
			fmt.Println(fmt.Sprintf("exit status: %v", status.StatusCode))
		}
	}
}

func (d *Docker) StartAndRemoveContainer() {
	d.Start()
	if err := d.Client.ContainerRemove(d.Context, d.ContainerID, types.ContainerRemoveOptions{}); err != nil {
		log.Fatal(err)
	}
}

func (d *Docker) getByContainerName(name string) types.Container {
	list, err := d.Client.ContainerList(d.Context, types.ContainerListOptions{All: true})
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range list {
		for _, n := range t.Names {
			if n == "/"+name {
				return t
			}
		}
	}
	return types.Container{}
}

func (d *Docker) createContainer() (container.CreateResponse, error) {
	resp, err := d.Client.ContainerCreate(d.Context, d.CliInfo.ConfigList, nil, nil, nil, d.CliInfo.Name)
	if err != nil {
		return resp, err
	}

	return resp, err
}
