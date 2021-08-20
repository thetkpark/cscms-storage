package main

import (
	"github.com/hashicorp/go-hclog"
	"github.com/thetkpark/cscms-temp-storage/pb"
	"github.com/thetkpark/cscms-temp-storage/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"os"
)

func main() {
	logger := hclog.Default()

	fileService := service.NewFileService(logger)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register service to gRPC server
	pb.RegisterFileServiceServer(grpcServer, fileService)
	// Enable reflection
	reflection.Register(grpcServer)

	// Create listener and start serving
	listener, err := net.Listen("tcp", ":5000")
	if err != nil {
		logger.Error("unable to create new listener", err)
		os.Exit(1)
	}
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Error("unable to start server", err)
		os.Exit(1)
	}
}
