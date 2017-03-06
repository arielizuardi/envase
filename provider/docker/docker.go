package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/arielizuardi/envase/provider"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	// DefaultDockerLibraryURL host
	DefaultDockerLibraryURL = `docker.io/library/`
)

type dockerImageProvider struct {
	Ctx           context.Context
	DockerClient  *client.Client
	ImageName     string
	Host          string
	ContainerPort string
	ExposedPort   string
	ContainerName string
	ContainerID   string
	EnvConfig     []string
}

func (i *dockerImageProvider) Has() (bool, error) {
	images, err := i.DockerClient.ImageList(i.Ctx, types.ImageListOptions{})
	if err != nil {
		return false, err
	}

	for _, image := range images {
		for _, t := range image.RepoTags {
			if t == i.ImageName {
				fmt.Printf(">>> Found image [%v]\n", i.ImageName)
				return true, nil
			}
		}
	}

	return false, nil
}

func (i *dockerImageProvider) Pull() error {
	fmt.Printf(`>>> Pulling image [%v] ...`+"\n", i.ImageName)
	imageURL := DefaultDockerLibraryURL + i.ImageName
	out, err := i.DockerClient.ImagePull(i.Ctx, imageURL, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, out)
	fmt.Printf(`>>> Finished pulling image [%v]`+"\n", i.ImageName)

	return nil
}

func (i *dockerImageProvider) Status() (bool, bool, error) {
	imageCreated := false
	imageRunning := false

	containers, err := i.DockerClient.ContainerList(i.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return imageCreated, imageRunning, err
	}

	for _, container := range containers {
		if container.Image == i.ImageName {
			imageCreated = true

			if container.State == `running` {
				imageRunning = true
			}
		}
	}

	return imageCreated, imageRunning, nil
}

func (i *dockerImageProvider) Create() (string, error) {
	exposedPorts, portBindings, err := nat.ParsePortSpecs([]string{i.Host + `:` + i.ExposedPort + `:` + i.ContainerPort + `/tcp`})
	if err != nil {
		return ``, err
	}

	containerCreated, err := i.DockerClient.ContainerCreate(
		i.Ctx,
		&container.Config{
			Image:        i.ImageName,
			ExposedPorts: exposedPorts,
			Env:          i.EnvConfig,
		},
		&container.HostConfig{
			PortBindings: portBindings,
		},
		nil,
		i.ContainerName,
	)

	if err != nil {
		return ``, err
	}

	return containerCreated.ID, nil
}

func (i *dockerImageProvider) Start(containerID string) error {
	return i.DockerClient.ContainerStart(i.Ctx, containerID, types.ContainerStartOptions{})
}

func (i *dockerImageProvider) Stop(containerID string) error {
	dur := time.Duration(30) * time.Second
	if err := i.DockerClient.ContainerStop(i.Ctx, containerID, &dur); err != nil {
		return err
	}

	return nil
}

// NewDockerImageProvider returns new instance of docker image provider
func NewDockerImageProvider(ctx context.Context, dockerClient *client.Client, imageName string, host string, containerPort string, exposedPort string, containerName string, envConfig []string) provider.ImageProvider {
	return &dockerImageProvider{
		Ctx:           ctx,
		DockerClient:  dockerClient,
		ImageName:     imageName,
		Host:          host,
		ContainerPort: containerPort,
		ExposedPort:   exposedPort,
		ContainerName: containerName,
		EnvConfig:     envConfig,
	}
}
