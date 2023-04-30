mkdir BuildGo && cd BuildGo
wget -q https://go.dev/dl/go1.20.2.linux-amd64.tar.gz
tar -C . -xzf go1.20.2.linux-amd64.tar.gz
export PATH=$PWD/go/bin:$PATH
cd ..
go version