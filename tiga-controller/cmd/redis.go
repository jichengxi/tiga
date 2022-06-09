package cmd

import (
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"tiga-controller/ControllerTypes"
	"tiga-controller/watcher/redis"
)

var redisCmd = &cobra.Command{
	Use:     "redis",
	Short:   "使用Redis作为配置中心",
	Version: "v0.0.1",
	Run: func(cmd *cobra.Command, args []string) {
		runRedisWatcher()
	},
}

var redisType string
var master string
var password string
var db int

func init() {
	redisCmd.Flags().StringVarP(&ControllerTypes.Port, "Port", "P", "10000", "http端口")
	redisCmd.Flags().StringVarP(&redisType, "redisType", "t", "alone", "集群策略, (alone, sentinel)")
	redisCmd.Flags().StringVarP(&ControllerTypes.Endpoint, "endpoint", "e", "", "redis地址:redis端口, []redis地址:redis端口")
	redisCmd.Flags().StringVarP(&master, "master", "m", "", "sentinel master名")
	redisCmd.Flags().StringVarP(&password, "password", "p", "", "sentinel 密码")
	redisCmd.Flags().IntVarP(&db, "db", "d", 0, "redis 数据库")
	rootCmd.AddCommand(redisCmd)
	err := redisCmd.MarkFlagRequired("endpoint")
	if err != nil {
		panic(err)
	}
}

func runRedisWatcher() {
	if redisType == "alone" {
		redisCli := redis.NewAloneClient(ControllerTypes.Endpoint, db)
		redisCli.Watcher()
	} else if redisType == "sentinel" {
		redisCli := redis.NewSentinelClient(ControllerTypes.Endpoint, master, password, db)
		redisCli.Watcher()
	} else {
		log.Error("redisType must be alone or sentinel, endpoint must be not empty")
	}
}
