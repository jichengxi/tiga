export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o test03 test03.go
scp test03 root@172.17.36.211:/tmp/