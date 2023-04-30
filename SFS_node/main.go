// Package main implements a server for Bulletin service.
package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/ARui-tw/I2DS_serverless-file-system/SFS"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port        = flag.Int("port", 50052, "The server port")
	server_port = flag.Int("server_port", 50051, "The server port")
)

var (
	addr string
)

type server struct {
	pb.UnimplementedNodeServer
}

func (s *server) GetLoad(ctx context.Context, in *pb.Empty) (*pb.Load, error) {
	// TODO: get load
	return &pb.Load{Load: 0}, nil
}

func (s *server) Download(ctx context.Context, in *pb.DownloadMessage) (*pb.ACK, error) {
	// TODO: place the file to download folder
	os.MkdirAll(fmt.Sprintf("share/%d/%d", *port, in.GetNodeID()), 0755)
	_, err := copy(fmt.Sprintf("files/%d/%s", *port, in.GetFilename()), fmt.Sprintf("share/%d/%d/%s", *port, in.GetNodeID(), in.GetFilename()))
	if err != nil {
		log.Error(err)
		return &pb.ACK{Success: false}, nil
	}
	return &pb.ACK{Success: true}, nil
}

func handleDownload(fileName string) (err error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()

	c := pb.NewTrackingClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Find(ctx, &pb.String{Message: fileName})
	if err != nil {
		log.Error(err)
		return
	}

	// select a random node
	rand.Seed(time.Now().UnixNano())
	nodeID := r.GetNodeID()[rand.Intn(len(r.GetNodeID()))]

	fmt.Printf("Downloading from node %d\n", nodeID)

	// download the file
	conn, err = grpc.Dial(fmt.Sprintf("localhost:%d", nodeID), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()

	c_node := pb.NewNodeClient(conn)

	ctx_node, cancel_node := context.WithTimeout(context.Background(), time.Second)
	defer cancel_node()
	r_node, err_node := c_node.Download(ctx_node, &pb.DownloadMessage{Filename: fileName, NodeID: int32(*port)})
	if err_node != nil {
		log.Error(err_node)
		return
	}

	if !r_node.GetSuccess() {
		log.Error("Download failed!")
	}

	// copy the file to my folder
	_, error := copy(fmt.Sprintf("share/%d/%d/%s", nodeID, *port, fileName), fmt.Sprintf("files/%d/%s", *port, fileName))

	if error != nil {
		log.Error(error)
		return
	}

	fmt.Println("Download success!")

	return nil
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func updateList(filename string) (err error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()

	c := pb.NewTrackingClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UpdateList(ctx, &pb.UpdateMessage{NodeID: int32(*port), Filename: filename})
	if err != nil {
		log.Error(err)
		return err
	}

	if !r.GetSuccess() {
		log.Error("Update list failed!")
	}

	return nil
}

func PrintMenu() {
	fmt.Println("-----------------")
	fmt.Println("\nMenu:")
	fmt.Println("\t1. Download")
	fmt.Println("\tq. Exit")
	fmt.Print("> ")
}

func main() {
	flag.Parse()

	addr = fmt.Sprintf("localhost:%d", *server_port)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// start the server in a goroutine
	go func() {
		s := grpc.NewServer()
		pb.RegisterNodeServer(s, &server{})
		log.Info("Single server listening at ", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: ", err)
		}
	}()

	// create the share folder
	if err := os.MkdirAll(fmt.Sprintf("share/%d", *port), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// updateList
	files, err := ioutil.ReadDir(fmt.Sprintf("files/%d", *port))
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			if err := updateList(f.Name()); err != nil {
				log.Error(err)
			}
		}
	}

	// handle the interrupt signal
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.RemoveAll(fmt.Sprintf("share/%d", *port))
		os.Exit(1)
	}()

	// handle the user input
	buf := bufio.NewReader(os.Stdin)
	for {
		PrintMenu()

		text, err := buf.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Error(err)
			break
		}

		switch text {
		case "1\n":
			var content string

			fmt.Print("Enter the file name: ")
			_, err := fmt.Scanf("%s", &content)
			if err != nil {
				log.Error(err)
				break
			}
			handleDownload(content)
		default:
			fmt.Println("Invalid input!")
		}
	}
}
