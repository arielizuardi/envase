package envase

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	// DefaultDockerLibraryURL host
	DefaultDockerLibraryURL = `docker.io/library/`
)

// Container defines container interface
type Container interface {
	Start() error
	Stop() error
}

type dockerContainer struct {
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

func (d *dockerContainer) Start() error {
	fmt.Printf(`>>> Checking image [%v] ...`+"\n", d.ImageName)
	hasImage, err := d.hasImage()
	if err != nil {
		return err
	}

	if !hasImage {
		err = d.pullImage()
		if err != nil {
			return err
		}
	}

	imageCreated, imageRunning, err := d.getImageStatus()
	if err != nil {
		return err
	}

	if !imageCreated {
		d.ContainerID, err = d.createImage()
		if err != nil {
			return err
		}
	}

	if !imageRunning {
		if err := d.DockerClient.ContainerStart(d.Ctx, d.ContainerID, types.ContainerStartOptions{}); err != nil {
			return err
		}
	}

	fmt.Printf(`>>> Image [%v] is now running`+"\n", d.ImageName)

	return nil
}

func (d *dockerContainer) Stop() error {
	return nil
}

// hasImage
func (d *dockerContainer) hasImage() (bool, error) {
	images, err := d.DockerClient.ImageList(d.Ctx, types.ImageListOptions{})
	if err != nil {
		return false, err
	}

	for _, image := range images {
		for _, t := range image.RepoTags {
			if t == d.ImageName {
				fmt.Printf(">>> Found image [%v]\n", d.ImageName)
				return true, nil
			}
		}
	}

	return false, nil
}

func (d *dockerContainer) pullImage() error {
	fmt.Printf(`>>> Pulling image [%v] ...`+"\n", d.ImageName)
	imageURL := DefaultDockerLibraryURL + d.ImageName
	out, err := d.DockerClient.ImagePull(d.Ctx, imageURL, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	io.Copy(os.Stdout, out)
	fmt.Printf(`>>> Finished pulling image [%v]`+"\n", d.ImageName)

	return nil
}

func (d *dockerContainer) getImageStatus() (bool, bool, error) {
	imageCreated := false
	imageRunning := false

	containers, err := d.DockerClient.ContainerList(d.Ctx, types.ContainerListOptions{All: true})
	if err != nil {
		return imageCreated, imageRunning, err
	}

	for _, container := range containers {
		if container.Image == d.ImageName {
			imageCreated = true

			if container.State == `running` {
				imageRunning = true
			}
		}
	}

	return imageCreated, imageRunning, nil
}

func (d *dockerContainer) createImage() (string, error) {
	exposedPorts, portBindings, err := nat.ParsePortSpecs([]string{d.Host + `:` + d.ExposedPort + `:` + d.ContainerPort + `/tcp`})
	if err != nil {
		return ``, err
	}

	containerCreated, err := d.DockerClient.ContainerCreate(
		d.Ctx,
		&container.Config{
			Image:        d.ImageName,
			ExposedPorts: exposedPorts,
			Env:          d.EnvConfig,
		},
		&container.HostConfig{
			PortBindings: portBindings,
		},
		nil,
		d.ContainerName,
	)

	if err != nil {
		return ``, err
	}

	return containerCreated.ID, nil
}

// NewDockerContainer returns new instance of dockerContainer
func NewDockerContainer(ctx context.Context, dockerClient *client.Client, imageName string, host string, containerPort string, exposedPort string, containerName string, envConfig []string) Container {

	return &dockerContainer{
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
