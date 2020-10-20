package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// CreateNewContainer creates and starts a docker container using an existing image
// defined by imageName
func CreateNewContainer(imageName string, address string, port string) (string, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		err = fmt.Errorf("Unable to create docker client: %s", err.Error())
		return "", err
	}

	hostBinding := nat.PortBinding{
		HostIP:   address,
		HostPort: port,
	}
	containerPort, err := nat.NewPort("tcp", port)
	if err != nil {
		err = fmt.Errorf("Failed to get port: %s", err.Error())
		return "", err
	}

	portBinding := nat.PortMap{containerPort: []nat.PortBinding{hostBinding}}
	cont, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: imageName,
		},
		&container.HostConfig{
			PortBindings: portBinding,
		}, nil, "")
	if err != nil {
		err = fmt.Errorf("Failed to create docker container: %s", err.Error())
		return "", err
	}

	cli.ContainerStart(context.Background(), cont.ID, types.ContainerStartOptions{})
	fmt.Printf("Container %s is started", cont.ID)
	return cont.ID, nil
}
