package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	pb "github.com/ARui-tw/I2DS_serverless-file-system/SFS"
	"github.com/codingsince1985/checksum"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port              = flag.Int("port", 50052, "The server port")
	server_port       = flag.Int("server_port", 50051, "The server port")
	config_file       = flag.String("config_file", "config.json", "The server config file")
	modify_file_count = flag.Int("modify_file_count", 0, "The count of files to be modified")
)

var (
	addr     string
	config   Config
	m        map[int32]int
	load     int32
	waitTime int
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
	pb.UnimplementedNodeServer
}

func (s *server) GetLoad(ctx context.Context, in *pb.Empty) (*pb.Load, error) {
	return &pb.Load{Load: load}, nil
}

func (s *server) GetList(ctx context.Context, in *pb.Empty) (*pb.ACK, error) {
	files, err := ioutil.ReadDir(fmt.Sprintf("files/%d", *port))

	if err != nil {
		log.Error(err)
		return &pb.ACK{Success: false}, nil
	}

	for _, f := range files {
		if !f.IsDir() {
			if err := updateList(f.Name()); err != nil {
				log.Error(err)
			}
		}
	}

	return &pb.ACK{Success: true}, nil
}

func (s *server) Download(ctx context.Context, in *pb.DownloadMessage) (*pb.Load, error) {
	nodeID := in.GetNodeID()

	os.MkdirAll(fmt.Sprintf("share/%d/%d", *port, nodeID), 0755)
	_, err := copy(fmt.Sprintf("files/%d/%s", *port, in.GetFilename()), fmt.Sprintf("share/%d/%d/%s", *port, nodeID, in.GetFilename()))
	if err != nil {
		log.Error(err)
		return &pb.Load{Load: -1}, nil
	}

	// atomic.AddInt32(&load, 1)
	load++
	waitTime += m[nodeID]
	log.Info(fmt.Sprintf("Downloading from node %d, delay %d ms\n", nodeID, m[nodeID]))
	go func() {
		time.Sleep(time.Duration(waitTime) * time.Millisecond)
		// atomic.AddInt32(&load, -1)
		load--
		waitTime -= m[nodeID]
	}()

	return &pb.Load{Load: int32(waitTime)}, nil
}

func findNode(nodes []int32) int32 {
	if len(nodes) == 0 {
		return -1
	}

	rand.Seed(time.Now().UnixNano())
	nodeID := nodes[rand.Intn(len(nodes))]
	return nodeID
}

func handleGetLoad(nodeID int32) {
	conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", nodeID), grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Error(err)
		return
	}
	defer conn.Close()

	c := pb.NewNodeClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.GetLoad(ctx, &pb.Empty{})
	if err != nil {
		log.Error(err)
		return
	}

	fmt.Println(r.GetLoad())
}

func handleDownload(fileName string) (err error) {
	var r *pb.IDs

	for success := false; !success; {
		conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Error(err)
			return err
		}
		defer conn.Close()

		c := pb.NewTrackingClient(conn)

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err = c.Find(ctx, &pb.String{Message: fileName})
		if err != nil {
			// log.Error(err)
			log.Info("Server is down, retrying in 3 Second...")
			time.Sleep(time.Second * 3)
		} else {
			success = true
		}
	}

	if (r.GetNodeID() == nil) || (len(r.GetNodeID()) == 0) {
		log.Error("File not found")
		return
	}
	Nodes := r.GetNodeID()

	// NOTE: modify the file n times
	a := *modify_file_count

	for success := false; !success; {
		// select a random node
		nodeID := findNode(Nodes)
		if nodeID == int32(*port) {
			// remove the node from the list
			for i, node := range Nodes {
				if node == nodeID {
					Nodes = append(Nodes[:i], Nodes[i+1:]...)
					break
				}
			}

			log.Info("Already have the file, retrying to find another node...")
			continue
		}

		if nodeID == -1 {
			log.Error("File not found")
			return
		}

		log.Info(fmt.Sprintf("Downloading from node %d...\n", nodeID))

		// download the file
		conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", nodeID), grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Error(err)
			return err
		}
		defer conn.Close()

		c_node := pb.NewNodeClient(conn)

		ctx_node, cancel_node := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel_node()
		r_node, err_node := c_node.Download(ctx_node, &pb.DownloadMessage{Filename: fileName, NodeID: int32(*port)})
		if err_node != nil {
			log.Info("Node is down, retrying to find another node...")

			for i, node := range Nodes {
				if node == nodeID {
					Nodes = append(Nodes[:i], Nodes[i+1:]...)
					break
				}
			}

			continue
		}

		toWait := r_node.GetLoad()
		if toWait == -1 {
			log.Error("Download failed!")
		}

		// mimic the latency
		fmt.Printf("Latency: %d ms\n", m[nodeID])
		time.Sleep(time.Duration(toWait) * time.Millisecond)

		// copy the file to my folder
		_, error := copy(fmt.Sprintf("share/%d/%d/%s", nodeID, *port, fileName), fmt.Sprintf("files/%d/%s", *port, fileName))

		if error != nil {
			log.Error(error)
			return error
		}
		if a > 0 {
			// slightly modify the file to make sure the MD5 is different
			f, err_ := os.OpenFile(fmt.Sprintf("files/%d/%s", *port, fileName), os.O_APPEND|os.O_WRONLY, 0644)
			if err_ != nil {
				log.Error(err_)
				return err_
			}
			defer f.Close()

			if _, err_ := f.WriteString(" "); err_ != nil {
				log.Error(err_)
				return err_
			}
			a--
		}

		if md5, _ := checksum.MD5sum(fmt.Sprintf("files/%d/%s", *port, fileName)); md5 != r.GetMd5() {
			log.Error("MD5 mismatch, retrying...")
			// remove the file
			os.Remove(fmt.Sprintf("files/%d/%s", *port, fileName))

			// remove the node from the list
			for i, node := range Nodes {
				if node == nodeID {
					Nodes = append(Nodes[:i], Nodes[i+1:]...)
					break
				}
			}

			continue
		}

		success = true
	}

	log.Info("Download success!")

	// update the list
	if err := updateList(fileName); err != nil {
		log.Error(err)
		return err
	}
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

	md5, _ := checksum.MD5sum(fmt.Sprintf("files/%d/%s", *port, filename))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UpdateList(ctx, &pb.UpdateMessage{NodeID: int32(*port), Filename: filename, Md5: md5})
	if err != nil {
		log.Info("Server is down!")
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
	fmt.Println("\t2. GetLoad")
	fmt.Println("\tq. Exit")
	fmt.Print("> ")
}

func startServer(addr string) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNodeServer(s, &server{})
	log.Info("Node listening at ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: ", err)
	}
}

func main() {
	flag.Parse()

	jsonFile, err := os.Open(*config_file)
	var wg sync.WaitGroup

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened json config")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &config)

	addr = fmt.Sprintf("localhost:%d", *server_port)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	m = make(map[int32]int)

	for _, node := range config.Nodes {
		if node.ID == *port {
			for _, laten := range node.Latency {
				m[int32(laten.Dest)] = laten.Lat
			}
		}
	}

	wg.Add(1)
	go func(addr string) {
		defer wg.Done()
		startServer(fmt.Sprintf(":%d", *port))
	}(fmt.Sprintf(":%d", *port))

	// create the share folder
	if err := os.MkdirAll(fmt.Sprintf("share/%d", *port), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// updateList
	files, err := ioutil.ReadDir(fmt.Sprintf("files/%d", *port))

	if err != nil {
		// if the folder does not exist, create it
		if os.IsNotExist(err) {
			if err := os.MkdirAll(fmt.Sprintf("files/%d", *port), os.ModePerm); err != nil {
				log.Fatal(err)
			}
		} else { // otherwise, log the error
			log.Fatal(err)
		}
	}

	for _, f := range files {
		if !f.IsDir() {
			if err := updateList(f.Name()); err != nil {
				// log.Error(err)
				break
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
			// _, err := fmt.Scanf("%s", &content)
			content, err = buf.ReadString('\n')
			// remove the newline character
			content = strings.TrimSuffix(content, "\n")
			fmt.Println(content)
			if err != nil {
				log.Error(err)
				break
			}
			handleDownload(content)
		case "2\n":
			var content string
			fmt.Print("Enter node ID: ")
			// _, err := fmt.Scanf("%s", &content)
			content, err = buf.ReadString('\n')
			// remove the newline character
			content = strings.TrimSuffix(content, "\n")
			if err != nil {
				log.Error(err)
				break
			}
			i, err := strconv.Atoi(content)

			if err != nil {
				log.Error(err)
				break
			}

			handleGetLoad(int32(i))
		default:
			fmt.Println("Invalid input!")
		}
	}
	wg.Wait()
}
