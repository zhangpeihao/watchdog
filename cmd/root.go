// Copyright © 2016 Zhang Peihao <zhangpeihao@gmail.com>

package cmd

import (
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"k8s.io/kubernetes/pkg/util/logs"
)

var cfgFile string
var cfgQuiet bool
var cfgWebUIHost string
var cfgWebUIPort int

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "Nginx status watchdog service",
	Long: `
Watchdog service

Watchdog service can check the status page of nginx, and report the warning alarm to influxdb.`,
}

func Execute() error {
	logs.InitLogs()
	defer logs.FlushLogs()
	glog.Infoln("root::Execute")
	defer glog.Infoln("root::Execute end")
	return RootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "/etc/watchdog/config.yaml", "config file")
	RootCmd.PersistentFlags().BoolVarP(&cfgQuiet, "quiet", "q", false, "quiet operation")

	RootCmd.PersistentFlags().StringVar(&cfgWebUIHost, "webui-host", "0.0.0.0", "The host bound for WebUI service.")
	RootCmd.PersistentFlags().IntVar(&cfgWebUIPort, "webui-port", 7080, "The port bound for WebUI service.")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/watchdog")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		glog.Infoln("Using config file:", viper.ConfigFileUsed())
	}
}
