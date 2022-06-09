package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tiga-controller",
	Short: "Tiga控制器",
	Long:  "Tiga控制器，用于通过外部配置中心或kubernetes ConfigMap读取配置来更新tiga-exporter监听项",
}

func Execute() error {
	return rootCmd.Execute()
}
