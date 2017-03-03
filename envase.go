package envase

import (
	"context"

	"github.com/arielizuardi/envase/provider"
	"github.com/arielizuardi/envase/provider/docker"
	"github.com/docker/docker/client"
)

// ContainerContract defines container interface
type ContainerContract interface {
	Start() error
	Stop() error
}

type container struct {
	ContainerID string
	Image       provider.ImageProvider
}

func (dc *container) Start() error {
	hasImage, err := dc.Image.Has()
	if err != nil {
		return err
	}

	if !hasImage {
		err = dc.Image.Pull()
		if err != nil {
			return err
		}
	}

	imageCreated, imageRunning, err := dc.Image.Status()
	if err != nil {
		return err
	}

	if !imageCreated {
		dc.ContainerID, err = dc.Image.Create()
		if err != nil {
			return err
		}
	}

	if !imageRunning {
		if err := dc.Image.Start(); err != nil {
			return err
		}
	}

	return nil
}

func (dc *container) Stop() error {
	return dc.Image.Stop()
}

// NewDefaultContainer returns new instance of default container
func NewDefaultContainer(imageProvider provider.ImageProvider, containerID string) ContainerContract {
	return &container{Image: imageProvider, ContainerID: containerID}
}

// NewDockerContainer returns new instance of dockerContainer
func NewDockerContainer(ctx context.Context, dockerClient *client.Client, imageName string, host string, containerPort string, exposedPort string, containerName string, envConfig []string) ContainerContract {
	imageProvider := docker.NewDockerImageProvider(ctx, dockerClient, imageName, host, containerPort, exposedPort, containerName, envConfig)
	return &container{
		Image:       imageProvider,
		ContainerID: ``,
	}
}
