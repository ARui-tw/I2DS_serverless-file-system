// Package main implements a server for Bulletin service.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"time"

	pb "github.com/ARui-tw/I2DS_serverless-file-system/SFS"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port        = flag.Int("port", 50051, "The server port")
	config_file = flag.String("config_file", "config.json", "The server config file")
)

var (
	m      map[string][]int32
	md5    map[string]string
	config Config
)

type Config struct {
	Nodes   []NodeInfo `json:"nodes"`
	Primary int        `json:"primary"`
}

type NodeInfo struct {
	ID      int           `json:"id"`
	Latency []LatencyInfo `json:"latency"`
}

type LatencyInfo struct {
	Dest int `json:"dest"`
	Lat  int `json:"latency"`
}

type server struct {
	pb.UnimplementedTrackingServer
}

func (s *server) Find(ctx context.Context, in *pb.String) (*pb.IDs, error) {
	return &pb.IDs{NodeID: m[in.GetMessage()], Md5: md5[in.GetMessage()]}, nil
}

func (s *server) UpdateList(ctx context.Context, in *pb.UpdateMessage) (*pb.ACK, error) {
	// check if the file is already in the list
	if _, ok := md5[in.GetFilename()]; !ok {
		md5[in.GetFilename()] = in.GetMd5()
	} else {
		if md5[in.GetFilename()] != in.GetMd5() {
			log.Error("MD5 mismatch")
			return &pb.ACK{Success: false}, nil
		}
	}

	// update the list, but don't add if the node is already in the list
	for _, nodeID := range m[in.GetFilename()] {
		if nodeID == in.GetNodeID() {
			return &pb.ACK{Success: true}, nil
		}
	}

	m[in.GetFilename()] = append(m[in.GetFilename()], in.GetNodeID())

	return &pb.ACK{Success: true}, nil
}

func getListUpdate() {
	for _, node := range config.Nodes {
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", node.ID), grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Error(err)
			return
		}
		defer conn.Close()

		c := pb.NewNodeClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.GetList(ctx, &pb.Empty{})
		if err != nil {
			log.Info("Node ", node.ID, " is not on")
			continue
		}

		if !r.GetSuccess() {
			log.Fatal("Get list failed")
			return
		}

		log.Info(fmt.Sprint("Node ", node.ID, " list updated"))
	}
}

func startServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTrackingServer(s, &server{})
	log.Info("Node listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: ", err)
	}
}

func main() {
	flag.Parse()
	var wg sync.WaitGroup
	jsonFile, err := os.Open(*config_file)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened json config")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &config)

	m = make(map[string][]int32)
	md5 = make(map[string]string)

	wg.Add(1)
	go func(addr string) {
		defer wg.Done()
		startServer(fmt.Sprintf(":%d", *port))
	}(fmt.Sprintf(":%d", *port))

	getListUpdate()

	wg.Wait()
}
