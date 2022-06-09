```
put /tiga-controller/daemonSet/pod {"a":1,"b":1,"c":1}
put /tiga-controller/daemonSet/node {"aa":1}

{"pod": {"a":1,"b":1,"c":1},"node":{"aa":1}}

go run main.go etcd -e 120.24.98.32:2379 -p 9090
```

