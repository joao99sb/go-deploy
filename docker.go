package main

import (
	"bg_gorquestrador/utils"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type Command interface{}

type NewProviderCommand struct {
	Port string
}

type DockerHandler struct {
	portList  []string
	dockerCli *client.Client

	MessageCh chan Command
}

func NewDockerHandler(rl []resource) *DockerHandler {

	var portList []string

	for _, pl := range rl {
		portList = append(portList, pl.Port)
	}

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &DockerHandler{
		portList:  portList,
		dockerCli: dockerClient,
		MessageCh: make(chan Command),
	}

}

func (d *DockerHandler) Start() {

	go d.dockerWatcher()

}

func (d *DockerHandler) dockerWatcher() error {
	ctx := context.Background()
	defer ctx.Done()

	options := types.EventsOptions{
		Filters: filters.NewArgs(
			filters.Arg("type", "container"),
		),
	}

	eventsCh, errCh := d.dockerCli.Events(ctx, options)

	for {
		select {
		case event := <-eventsCh:
			d.handleContainerLifeTime(ctx, event)
		case err := <-errCh:
			if errors.Is(err, io.EOF) {
				slog.Error("Provider event stream closed")
			}
			return err
		case <-ctx.Done():
			fmt.Println("flw ai glr...")
			return nil
		}
	}
}

func (d *DockerHandler) handleContainerLifeTime(ctx context.Context, e events.Message) {

	inspect, _ := d.dockerCli.ContainerInspect(ctx, e.Actor.ID)

	if e.Action == "start" {
		bindList := []string{}
		for _, port := range inspect.NetworkSettings.Ports {
			for _, binding := range port {
				bindList = append(bindList, binding.HostPort)
			}
		}
		bindPort := bindList[0]

		isAllowed := utils.Contains(d.portList, bindPort)

		if isAllowed {

			d.MessageCh <- NewProviderCommand{Port: bindPort}
		}

	} else if e.Action == "die" {
		fmt.Println("container die")
	}
}
