export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
go build -o test01 test01.go
scp test01 root@172.17.36.211:/tmp/