export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

go build -o tiga-exporter tiga-exporter.go
scp tiga-exporter root@172.17.36.211:/tmp/