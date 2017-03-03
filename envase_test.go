package envase_test

import (
	"context"
	"testing"

	"github.com/arielizuardi/envase"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func TestStartMySQL(t *testing.T) {
	ctx := context.Background()
	dockerClient, err := client.NewEnvClient()
	assert.NoError(t, err)
	envConfig := []string{
		"MYSQL_USER=" + `user`,
		"MYSQL_ROOT_PASSWORD=" + `pass`,
		"MYSQL_DATABASE=" + `kurio_db`,
	}
	container := envase.NewDockerContainer(ctx, dockerClient, `mysql:5.7`, `127.0.0.1`, `3306`, `33060`, `papua_test`, envConfig)

	assert.NoError(t, container.Start())
}

func TestStartFluentd(t *testing.T) {
	ctx := context.Background()
	dockerClient, err := client.NewEnvClient()
	assert.NoError(t, err)
	envConfig := []string{}
	container := envase.NewDockerContainer(ctx, dockerClient, `fluent/fluentd:v0.12.32`, `127.0.0.1`, `24224`, `24224`, `charon_test`, envConfig)

	assert.NoError(t, container.Start())
}
