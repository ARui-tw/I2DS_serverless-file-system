// Package main implements a server for Bulletin service.
package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	pb "github.com/ARui-tw/I2DS_serverless-file-system/SFS"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

var (
	m map[string][]int32
)

type server struct {
	pb.UnimplementedTrackingServer
}

func (s *server) Find(ctx context.Context, in *pb.String) (*pb.IDs, error) {
	return &pb.IDs{NodeID: m[in.GetMessage()]}, nil
}

func (s *server) UpdateList(ctx context.Context, in *pb.UpdateMessage) (*pb.ACK, error) {
	// update the list, but don't add if the node is already in the list
	for _, nodeID := range m[in.GetFilename()] {
		if nodeID == in.GetNodeID() {
			return &pb.ACK{Success: true}, nil
		}
	}

	m[in.GetFilename()] = append(m[in.GetFilename()], in.GetNodeID())

	return &pb.ACK{Success: true}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()

	m = make(map[string][]int32)

	pb.RegisterTrackingServer(s, &server{})
	log.Info("Single server listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: ", err)
	}
}