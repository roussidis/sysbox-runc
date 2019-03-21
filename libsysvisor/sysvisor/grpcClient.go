package sysvisor

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/opencontainers/runc/libsysvisor/sysvisor-protobuf"
	"google.golang.org/grpc"
)

const sysvisor_address = "localhost:50052"

//
// Establishes grpc connection to sysvisor-fs' remote-end.
//
func sysvisorfs_connect() *grpc.ClientConn {

	// Set up a connection to the server.
	// TODO: Secure me through TLS.
	conn, err := grpc.Dial(sysvisor_address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect to Sysvisorfs: %v", err)
		return nil
	}

	return conn
}

//
// Registers container creation in sysvisor-fs. Notice that this
// is a blocking call that can potentially have a minor impact
// on container's boot-up speed.
//
func SendContainerRegistration(data *pb.ContainerData) error {

	// Set up sysvisorfs pipeline.
	conn := sysvisorfs_connect()
	if conn == nil {
		return fmt.Errorf("failed to connect with sysvisor-fs")
	}
	defer conn.Close()

	cntrChanIntf := pb.NewContainerStateChannelClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := cntrChanIntf.ContainerRegistration(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to register container with sysvisor-fs: %v", err)
	}

	return nil
}

//
// Unregisters container from Sysvisorfs.
//
func SendContainerUnregistration(data *pb.ContainerData) error {

	// Set up sysvisorfs pipeline.
	conn := sysvisorfs_connect()
	if conn == nil {
		return fmt.Errorf("failed to connect with sysvisor-fs")
	}
	defer conn.Close()

	cntrChanIntf := pb.NewContainerStateChannelClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Generate a container-unregistration message to Sysvisorfs
	_, err := cntrChanIntf.ContainerUnregistration(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to unregister container with sysvisor-fs: %v", err)
	}

	return nil
}

//
// Sends a container-update message to sysvisor-fs end. At this point, we are
// only utilizing this message for a particular case, update the container
// creation-time attribute, but this function can serve more general purposes
// in the future.
//
func SendContainerUpdate(data *pb.ContainerData) error {

	// Set up sysvisorfs pipeline.
	conn := sysvisorfs_connect()
	if conn == nil {
		return fmt.Errorf("failed to connect with sysvisor-fs")
	}
	defer conn.Close()

	cntrChanIntf := pb.NewContainerStateChannelClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	_, err := cntrChanIntf.ContainerUpdate(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to send container-update message to ",
			"sysvisor-fs: %v", err)
	}

	return nil
}
