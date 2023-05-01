# Project 3: "Serverless" File System
[GitHub Repo](https://github.com/ARui-tw/I2DS_serverless-file-system)
## Team members and contributions
- Group: 6
- Team members:
	- Eric Chen (chen8332@umn.edu)
	- Sai Vanamala (vanam017@umn.edu)
	- Rohan Shanbhag (shanb020@umn.edu)

### Member contributions:
We opted to work together as a team throughout the coding/development process by working together through Discord, meeting generally every night for an hour or two over the course of the project. Each team member shared equal responsibility for overseeing the code written, the design decisions/documentation, and the test cases. All three of us were working together on call to complete the different functions (like the Find(), Download() etc.) to ensure system functionality.

## Project Build and Compilation Instructions

In order to run the project:
### Update Go
Since the version on the lab machines is outdated, we need to update the Go version to the latest version. I've written a script to do this for you. Just run the following command:
```sh
. build-go.sh
```
**NOTE**: This script will modify the $PATH variable to point to the new Go installation. Once you exit the shell, you will need to go to the project's folder and run the following command again to update the $PATH variable. (I don't want to mess up with the $PATH variable in your shell profile, but you can add the following line to your shell profile if you want to make it permanent.)
```sh
export PATH=$PWD/BuildGo/go/bin:$PATH
```

### Run the server:
```sh
go run SFS_server/main.go -config_file config.json  -port 50051
```
Change the port number to the port number you want to use and the config_file to the config file you want to use.

### Run the client: 
```sh
go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50052
```
Change the port numbers to the port numbers you want to use and the config_file to the config file you want to use.

## Design Document

We decided upon the design of our Serverless File System for Project 3 with some specific design decisions as follows:
When a node makes a query for a file (with the Find function), we return a set of peers that are completely outlined by the exact filename specified.
This means that if a node attempts to Find("foo"), they will be greeted with an error message/a return indicating that their search is not refined enough
Instead, the user is prompted to indicate the file extension more specifically, such as Find("foo.txt") or Find("foo.c"), which would result in a successful Find query
Lastly, any successful Find queries will yield an output list of nodes and the corresponding file name requested, along with their checksum. If any nodes do not possess the requested file then they do not show up in the list of nodes possessing that file.
To compute the checksum we utilized Golang’s checksum functionality available [here](https://pkg.go.dev/github.com/codingsince1985/checksum#section-readme).
As specified by the writeup, when one node makes a download request to another, i.e. node A requests a download from node B, the requested file becomes a part of a "shared directory"
This means that that file is now accessible to both peer nodes A and B, but this shared directory is not visible to the other peers/nodes within the serverless file system.
Furthermore (since we already know the latency from the static configuration as outlined in the node configuration file), we ensure that the nodes download a file from the next best available node according to this file.
Peer Selection – we opted to call GetLoad to determine the current load of the peer (since we already know the latency from the static configuration as outlined in the node configuration file). We find the node by finding the smallest of the sum of the latency between the two nodes and load times average latency. Overall, the average file download time increase as the number of peers hosting a file increases. 

### Fault Tolerance

Downloaded files can be corrupted. Specifically, if a peer node attempts to download a file from another peer but the file gets corrupted during transmission, then we can identify the occurrence of this fault/failure by recomputing the checksum of the file upon receipt by the requesting node.
Tracking server crashes. Specifically, if the server crashes then we direct the server to obtain the up-to-date list of files that have been shared by essentially utilizing the UpdateList function for each peer/node within the system.
This allows each node to share an updated list of files that the server would need to keep track of.
We also ensure that the server maintains one copy of each new/unique file that is possessed by each of the nodes, so that no duplicate files are unnecessarily stored when the server recovers from a crash.
Peer Crashes (fail-stop). Specifically, if a peer node crashes while the system is running, the user can continue interacting with the system (with the other, alive nodes), while the crashed node attempts to reconnect to the system.
The crashed node attempts to reconnect by referencing the /machID, the tracking server, and any shared directories that the node was a part of prior to it crashing.

## Test Cases

Let the following be true across each of the following test cases:

Let `files/{NodeID}` initially contain the following:
- Node A: "foo.txt"; "bar.txt"; "func.c"
- Node B: "foo.txt"; "foo.c"
- Node C: "bar.c"; "func.c"

- Node A boots up with { "id": 50052 } and obtains its file list, then reports the file list, checksum and endpoint information (IP and Port) to the tracking server.
- Node B boots up with { "id": 50053 } and obtains its file list, then reports the file list,  checksum and endpoint information (IP and Port) to the tracking server.
- Node C boots up with { "id": 50054 } and obtains its file list, then reports the file list,  checksum and endpoint information (IP and Port) to the tracking server.
- `config.json`:
    ```json
    {
        "primary": 50051,
        "nodes": [
            {
                "id": 50052,
                "latency": [
                    {
                        "dest": 50053,
                        "latency": 1000
                    },
                    {
                        "dest": 50054,
                        "latency": 5000
                    }
                ]
            },
            {
                "id": 50053,
                "latency": [
                    {
                        "dest": 50052,
                        "latency": 1000
                    },
                    {
                        "dest": 50054,
                        "latency": 3000
                    }
                ]
            },
            {
                "id": 50054,
                "latency": [
                    {
                        "dest": 50052,
                        "latency": 5000
                    },
                    {
                        "dest": 50053,
                        "latency": 3000
                    }
                ]
            }
        ]
    }
    ```

### Test Case #1: Normal Condition 
This test case tests that the system is capable of finding the list of nodes which store a certain file (with the Find function) and can download that file if desired.
#### How to run:
- Server:
    ```sh
    go run SFS_server/main.go -config_file config.json -port 50051
    ```
- Nodes (in separate terminals):
    ```sh
    go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50052
    ```
    ```sh
    go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50053
    ```
    ```sh
    go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50054
    ```


In order for a node (i.e. Node A) to find the list of nodes storing "foo.txt", we have the following input:

- \<User enters\>: 1
- \<User enters\>: foo.c
- \<System Outputs\>:
    ```
    INFO[0103] Downloading from node 50053...               
    Latency: 1000 ms
    INFO[0104] Download success!                            
    ```
- \<User enters\>: q

The file `foo.c` will now appear in the `files/{NodeID}` directory of Node A.

- To clean up:
    ```sh
    rm files/50052/foo.c
    ```

### Test Case #2:
Corrupted Download File (Failure) Condition – This testcase observes how the serverless file system handles a file that has been corrupted during Download. We can emulate this by introducing a byte to the file not originally present, and upon computing the checksums and recognizing the difference between the initial and final values, the system returns an error message to the client highlighting the failure/corrupted download file.
#### How to run:
- Server:
    ```sh
    go run SFS_server/main.go -config_file config.json -port 50051
    ```
- Nodes (in separate terminals):
    ```sh
    go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50052
    ```
    ```sh
    go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50053
    ```
    ```sh
    go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50054 -modify_file_count 1
    ```

In order to test the fault tolerance of the system, we will modify the file `foo.txt` in Node C (NodeID: 50054) by adding a byte to the file. It will automatically be detected by the system and the file will be re-downloaded from another node.

- In node C (NodeID: 50054)
    - \<User enters\>: 1
    - \<User enters\>: "foo.txt"	//Enter filename to be corrupted
    - \<System Outputs\>:
        ```
        Enter the file name: foo.txt
        INFO[0054] Downloading from node 50052...               
        Latency: 5000 ms
        ERRO[0059] MD5 mismatch, retrying...                    
        INFO[0059] Downloading from node 50053...               
        Latency: 3000 ms
        INFO[0062] Download success!                            
        ```
    - \<User enters\>: q

### Test Case #3:
Tracking Server Crash Condition – This test case shows how the serverless file system handles a crash by the tracking server.
#### How to run:
- Server:
    ```sh
    go run SFS_server/main.go -config_file config.json -port 50051
    ```
- Nodes (in separate terminals):
    ```sh
    go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50052
    ```
    ```sh
    go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50053
    ```

Let us assume that the server crashes unexpectedly. We can emulate this by terminating the server process.

- In the server terminal:
    - \<User enters\>: `Ctrl + C`
    - \<System Outputs\>:
        ```
        signal: interrupt
        ```

- In the node A terminal:
    - \<User enters\>: 1
    - \<User enters\>: foo.c
    - \<System Outputs\>:
        ```
        INFO[0024] Server is down, retrying in 3 Second...
        ```

- In the server terminal:
    - \<User enters\>:
        ```sh
        go run SFS_server/main.go -config_file config.json -port 50051
        ```
    - \<System Outputs\>:
        ```
        Successfully Opened json config
        INFO[0000] Node listening at [::]:50051                 
        INFO[0000] Node 50052 list updated                      
        INFO[0000] Node 50053 list updated                      
        INFO[0000] Node 50054 is not on                         
        ```

- In the node A terminal:
    - \<System Outputs\>:
        ```
        INFO[0066] Server is down, retrying in 3 Second...
        INFO[0069] Downloading from node 50053...               
        Latency: 1000 ms
        INFO[0070] Download success! 
        ```

The server will fetch the file lists from the nodes that are on and update its own file list.

### Test Case #4:
Peer Crash Condition – This test case handles a peer/node failing, and what occurs following the node’s reboot.

#### How to run:
- Server:
    ```sh
    go run SFS_server/main.go -config_file config.json -port 50051
    ```
- Nodes (in separate terminals):
    ```sh
    go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50052
    ```
    ```sh
    go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50053
    ```
    ```sh
    go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50054
    ```
Let us assume that Node C (NodeID: 50054) crashes unexpectedly. We can emulate this by terminating the node process.

- In the node C terminal:
    - \<User enters\>: `Ctrl + C`
    - \<System Outputs\>:
        ```
        exit status 1
        ```
    - \<User enters\>:
        ```sh
        go run SFS_node/main.go -config_file config.json -server_port 50051 -port 50054
        ```
Now, the node will update the server with its new file list. The server will then update its own file list. The node will then be able to interact with the system as normal.

## Pledge

No-one sought out any on-line solutions, e.g. github for portions of this lab

Signed:
Rohan Shanbhag
Eric Chen
Sai Vanamala