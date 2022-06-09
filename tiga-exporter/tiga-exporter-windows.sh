set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64

go build -o tiga-exporter tiga-exporter.go
scp tiga-exporter root@192.168.1.21:/tmp/