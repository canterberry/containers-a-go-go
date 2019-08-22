package main

import (
	"context"
	"fmt"
	"log"
	"syscall"
	"time"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
	"github.com/containerd/containerd/oci"
)

func main() {
	err := redisExample()

	if err != nil {
		log.Fatal(err)
	}
}

func redisExample() error {
	containerdSocketPath := "/run/containerd/containerd.sock"
	namespace := "example"
	containerName := "redis-server"
	imageName := "docker.io/library/redis:alpine"
	snapshotName := "redis-server-snapshot"

	client, err := containerd.New(containerdSocketPath)

	if err != nil {
		return err
	}

	defer client.Close()

	realm := namespaces.WithNamespace(context.Background(), namespace)

	image, err := client.Pull(realm, imageName, containerd.WithPullUnpack)

	if err != nil {
		return err
	}

	snapshot := containerd.WithNewSnapshot(snapshotName, image)
	// snapshotContext = containerd.WithSnapshot(snapshotName)

	imageConfig := oci.WithImageConfig(image)

	spec := containerd.WithNewSpec(imageConfig)

	wrappedImage := containerd.WithImage(image)

	container, err := client.NewContainer(realm, containerName, wrappedImage, snapshot, spec)

	if err != nil {
		return err
	}

	defer container.Delete(realm, containerd.WithSnapshotCleanup)

	task, err := container.NewTask(realm, cio.NewCreator(cio.WithStdio))

	if err != nil {
		return err
	}

	defer task.Delete(realm)

	exitStatusCode, err := task.Wait(realm)

	if err != nil {
		fmt.Println(err)
	}

	err = task.Start(realm)

	if err != nil {
		return err
	}

	time.Sleep(3 * time.Second)

	err = task.Kill(realm, syscall.SIGTERM)

	if err != nil {
		return err
	}

	status := <-exitStatusCode

	code, _, err := status.Result()

	if err != nil {
		return err
	}

	fmt.Printf("%s exited with status: %d\n", containerName, code)

	return nil
}
