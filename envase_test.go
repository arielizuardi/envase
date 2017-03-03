package envase_test

import (
	"errors"
	"testing"

	"github.com/arielizuardi/envase"
	"github.com/arielizuardi/envase/provider/mocks"
	"github.com/stretchr/testify/assert"
)

func TestStartWithNoImageInstalled(t *testing.T) {
	provider := &mocks.ImageProvider{}
	provider.On(`Has`).Return(false, nil)
	provider.On(`Pull`).Return(nil)
	provider.On(`Status`).Return(false, false, nil)
	provider.On(`Create`).Return(`container-id`, nil)
	provider.On(`Start`).Return(nil)

	container := envase.NewDefaultContainer(provider, ``)
	assert.NoError(t, container.Start())

	provider.AssertCalled(t, `Has`)
	provider.AssertCalled(t, `Pull`)
	provider.AssertCalled(t, `Status`)
	provider.AssertCalled(t, `Create`)
	provider.AssertCalled(t, `Start`)
}

func TestStartWithImageAlreadyInSystemAndStartTheImage(t *testing.T) {
	provider := &mocks.ImageProvider{}
	provider.On(`Has`).Return(true, nil)
	provider.On(`Status`).Return(true, false, nil)
	provider.On(`Start`).Return(nil)

	container := envase.NewDefaultContainer(provider, ``)
	assert.NoError(t, container.Start())

	provider.AssertCalled(t, `Has`)
	provider.AssertNotCalled(t, `Pull`)
	provider.AssertCalled(t, `Status`)
	provider.AssertNotCalled(t, `Create`)
	provider.AssertCalled(t, `Start`)
}

func TestStartWithImageAlreadyInSystemAndAlreadyRunning(t *testing.T) {
	provider := &mocks.ImageProvider{}
	provider.On(`Has`).Return(true, nil)
	provider.On(`Status`).Return(true, true, nil)

	container := envase.NewDefaultContainer(provider, ``)
	assert.NoError(t, container.Start())

	provider.AssertCalled(t, `Has`)
	provider.AssertNotCalled(t, `Pull`)
	provider.AssertCalled(t, `Status`)
	provider.AssertNotCalled(t, `Create`)
	provider.AssertNotCalled(t, `Start`)
}

func TestStartWithNoImageInstalledAndFailedToPull(t *testing.T) {
	provider := &mocks.ImageProvider{}
	provider.On(`Has`).Return(false, nil)
	provider.On(`Pull`).Return(errors.New(`Whoops!`))
	provider.On(`Status`).Return(false, false, nil)
	provider.On(`Create`).Return(`container-id`, nil)
	provider.On(`Start`).Return(nil)

	container := envase.NewDefaultContainer(provider, ``)
	assert.Error(t, container.Start())

	provider.AssertCalled(t, `Has`)
	provider.AssertCalled(t, `Pull`)
	provider.AssertNotCalled(t, `Status`)
	provider.AssertNotCalled(t, `Create`)
	provider.AssertNotCalled(t, `Start`)
}

func TestStartWithNoImageInstalledAndFailedToCreate(t *testing.T) {
	provider := &mocks.ImageProvider{}
	provider.On(`Has`).Return(false, nil)
	provider.On(`Pull`).Return(nil)
	provider.On(`Status`).Return(false, false, nil)
	provider.On(`Create`).Return(``, errors.New(`Whoops!`))
	provider.On(`Start`).Return(nil)

	container := envase.NewDefaultContainer(provider, ``)
	assert.Error(t, container.Start())

	provider.AssertCalled(t, `Has`)
	provider.AssertCalled(t, `Pull`)
	provider.AssertCalled(t, `Status`)
	provider.AssertCalled(t, `Create`)
	provider.AssertNotCalled(t, `Start`)
}

func TestStop(t *testing.T) {
	provider := &mocks.ImageProvider{}
	provider.On(`Stop`).Return(nil)

	container := envase.NewDefaultContainer(provider, ``)
	assert.NoError(t, container.Stop())
	provider.AssertCalled(t, `Stop`)
}

func TestStopAndGotError(t *testing.T) {
	provider := &mocks.ImageProvider{}
	provider.On(`Stop`).Return(errors.New(`Whoops!`))

	container := envase.NewDefaultContainer(provider, ``)
	assert.Error(t, container.Stop())
	provider.AssertCalled(t, `Stop`)
}
