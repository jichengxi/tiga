package cmd

//import (
//    "fmt"
//    "github.com/spf13/cobra"
//    "tiga-controller/watcher/etcd"
//)
//
//var etcdCmd = &cobra.Command{
//    Use:              "etcd",
//    Short:            "使用Etcd作为配置中心",
//    Version:          "etcd v3 client",
//    Run: func(cmd *cobra.Command, args []string)  {
//        fmt.Println(args)
//        fmt.Println(endpoint)
//        runEtcdWatcher()
//    },
//}
//
//var (
//    endpoint []string
//)
//
//func init() {
//    etcdCmd.Flags().IntVarP(&port, "port", "p", 8081, "etcd http端口")
//    etcdCmd.Flags().StringArrayVarP(&endpoint, "endpoints", "e", []string{}, "配置中心的IP列表")
//    rootCmd.AddCommand(etcdCmd)
//    err := etcdCmd.MarkFlagRequired("endpoints")
//    if err != nil {
//       panic(err)
//    }
//}
//
//func runEtcdWatcher() {
//    a := etcd.ClientSet{}
//    // "172.17.47.201:2379"
//    a.NewClient(endpoint)
//    defer a.Client.Close()
//    a.Watcher(port)
//}
