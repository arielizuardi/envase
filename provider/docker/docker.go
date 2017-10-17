package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/arielizuardi/envase/provider"
	godocker "github.com/fsouza/go-dockerclient"
)

const (
	// DefaultDockerLibraryURL host
	DefaultDockerLibraryURL = `docker.io/library/`
)

type dockerImageProvider struct {
	Ctx           context.Context
	DockerClient  *godocker.Client
	ImageName     string
	Host          string
	ContainerPort string
	ExposedPort   string
	ContainerName string
	ContainerID   string
	EnvConfig     []string
}

func (i *dockerImageProvider) Has() (bool, error) {
	images, err := i.DockerClient.ListImages(g.ListImagesOptions{})
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
	out, err := i.DockerClient.PullImage(godocker.ListImagesOptions{}, godocker.AuthConfiguration{})

	if err != nil {
		return err
	}

	io.Copy(os.Stdout, out)
	fmt.Printf(`>>> Finished pulling image [%v]`+"\n", i.ImageName)

	return nil
}

func (i *dockerImageProvider) Status() (bool, bool, string, error) {
	imageCreated := false
	imageRunning := false

	containers, err := i.DockerClient.ListContainers(godocker.ListContainersOptions{All: true})
	if err != nil {
		return imageCreated, imageRunning, i.ContainerID, err
	}

	for _, container := range containers {
		if container.Image == i.ImageName {
			imageCreated = true

			if container.State == `running` {
				imageRunning = true
			}

			i.ContainerID = container.ID
		}
	}

	return imageCreated, imageRunning, i.ContainerID, nil
}

func (i *dockerImageProvider) Create() (string, error) {
	exposedPorts, portBindings, err := nat.ParsePortSpecs([]string{i.Host + `:` + i.ExposedPort + `:` + i.ContainerPort + `/tcp`})
	if err != nil {
		return ``, err
	}

	containerCreated, err := i.DockerClient.CreateContainer(godocker.CreateContainerOptions{
		Name: i.ContainerName,
		Config: godocker.Config{
			Image: i.ImageName,
			ExposedPorts: exposedPorts,
			Env: i.EnvConfig,
		},
		HostConfig: godocker.HostConfig{
			PortBindings: portBindings
		},
		Context: i.Ctx,
	})

	if err != nil {
		return ``, err
	}

	return containerCreated.ID, nil
}

func (i *dockerImageProvider) Start(containerID string) error {
	return i.DockerClient.StartContainer(containerID, godocker.HostConfig{})
}

func (i *dockerImageProvider) Stop(containerID string) error {
	dur := time.Duration(30) * time.Second
	if err := i.DockerClient.StopContainer(containerID, dur); err != nil {
		return err
	}

	return nil
}

// NewDockerImageProvider returns new instance of docker image provider
func NewDockerImageProvider(ctx context.Context, dockerClient *godocker.Client, imageName string, host string, containerPort string, exposedPort string, containerName string, envConfig []string) provider.ImageProvider {
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
