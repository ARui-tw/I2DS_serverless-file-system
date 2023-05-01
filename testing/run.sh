go run ../SFS_server/main.go -config_file config.json -port 50051 &
go run ../SFS_node/main.go -config_file config.json -server_port 50051 -port 50052 < 50052.txt &
sleep 1
go run ../SFS_node/main.go -config_file config.json -server_port 50051 -port 50053 < 50053.txt &
go run ../SFS_node/main.go -config_file config.json -server_port 50051 -port 50054 < 50054.txt &