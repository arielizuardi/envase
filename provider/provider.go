package provider

// ImageProvider defines Image Provider Contract
type ImageProvider interface {
	// Has returns bool flag indicates image already exists or not
	Has() (bool, error)
	// Status returns bool flag indicates image creation status, image running status
	Status() (bool, bool, string, error)
	// Pull pulls image from the hub
	Pull() error
	// Create create image and returns identifier for image
	Create() (string, error)
	// Start start image
	Start(containerID string) error
	// Stop the image from running
	Stop(containerID string) error
}
