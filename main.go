package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	config_file = flag.String("config_file", "config.json", "The server config file")
)

type Config struct {
	Nodes   []nodeInfo `json:"nodes"`
	Primary int        `json:"primary"`
}

type nodeInfo struct {
	ID      int
	latency []latencyInfo
}

type latencyInfo struct {
	dest int
	lat  int
}

func main() {
	flag.Parse()

	jsonFile, err := os.Open(*config_file)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened json config")
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var config Config

	json.Unmarshal(byteValue, &config)

	// start primary server
	cmd := exec.Command("go", "run", "SFS_server/main.go", "-config_file", *config_file, "-port", strconv.Itoa(config.Primary))
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	cmdReader, _ := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		return
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("\t > %s\n", scanner.Text())
		}
	}()

	cmd.Start()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		return
	}
	time.Sleep(1 * time.Second)

	// start child server
	for _, node := range config.Nodes {
		cmd := exec.Command("go", "run", "SFS_node/main.go", "-config_file", *config_file, "-server_port", strconv.Itoa(config.Primary), "-port", strconv.Itoa(node.ID))
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		cmdReader, _ := cmd.StderrPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
			return
		}

		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				fmt.Printf("\t > %s\n", scanner.Text())
			}
		}()

		cmd.Start()

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
			return
		}

		time.Sleep(1 * time.Second)
	}

	log.Info("All servers are started")

	// sleep forever
	<-make(chan int)
}
