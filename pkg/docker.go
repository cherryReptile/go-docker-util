package pkg

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"io"
	"log"
	"os"
)

type Docker struct {
	ContainerInfo ContainerInfo
	Client        *client.Client
	Context       context.Context
}

type ContainerInfo struct {
	ID    string
	Image string
}

func NewDocker(image string) *Docker {
	return &Docker{
		ContainerInfo: ContainerInfo{
			Image: image,
		},
		Client:  &client.Client{},
		Context: context.Background(),
	}
}

func (d *Docker) Start(cmd, env []string, name string) {
	var containerID string
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatal(err)
	}
	d.Client = cli
	defer d.Client.Close()

	reader, err := cli.ImagePull(d.Context, d.ContainerInfo.Image, types.ImagePullOptions{})
	if err != nil {
		log.Fatal(err)
	}
	if _, err = io.Copy(os.Stdout, reader); err != nil {
		log.Fatal(err)
	}

	c := d.getByContainerName(name)
	if c.ID == "" {
		resp, err := d.createContainer(cmd, env, name)
		if err != nil {
			log.Fatal(err)
		}
		containerID = resp.ID
	} else if c.ID != "" {
		containerID = c.ID
	}

	d.ContainerInfo.ID = containerID
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

func (d *Docker) StartAndRemoveContainer(cmd, env []string, name string) {
	d.Start(cmd, env, name)
	if err := d.Client.ContainerRemove(d.Context, d.ContainerInfo.ID, types.ContainerRemoveOptions{}); err != nil {
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

func (d *Docker) createContainer(cmd, env []string, name string) (container.CreateResponse, error) {
	resp, err := d.Client.ContainerCreate(d.Context, &container.Config{
		Image: d.ContainerInfo.Image,
		Cmd:   cmd,
		Env:   env,
	}, nil, nil, nil, name)
	if err != nil {
		return resp, err
	}

	return resp, err
}
