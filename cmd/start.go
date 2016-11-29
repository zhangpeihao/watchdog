// Copyright © 2016 Zhang Peihao <zhangpeihao@gmail.com>
//

package cmd

import (
	"fmt"
	"net/http"
	"time"

	//"github.com/VividCortex/godaemon"
	"github.com/golang/glog"
	"github.com/spf13/cobra"

	"github.com/zhangpeihao/watchdog/pkg/apiservice"
	"github.com/zhangpeihao/watchdog/pkg/watchdog"
)

var (
	cfgNginxStatusPages []string
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start watchdog",
	Long: `
Start watchdog services.
`,
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	//godaemon.MakeDaemon(&godaemon.DaemonAttr{})
	startWatch(cmd)
}

// Start watch
func startWatch(cmd *cobra.Command) {
	var pages []string
	var err error
	pages, err = cmd.PersistentFlags().GetStringSlice("nginx-status-page")
	if err != nil {
		glog.Errorln("Get nginx-status-page option error:", err.Error())
		return
	}
	interval, err := cmd.Flags().GetDuration("watch-frenquency")
	if err != nil {
		glog.Errorln("Get watch-frenquency from settings error:", err.Error())
		panic(err)
	}

	var jobs []watchdog.WatchJob
	for _, page := range pages {
		job, err := watchdog.NewNginxStatus(page)
		if err != nil {
			glog.Errorln("New nginx status page watchdog error:", err.Error())
			panic(err)
		}
		jobs = append(jobs, job)
	}

	go watchdog.RunWatchJob(jobs, interval)

	// Start Web UI & Web API
	var webui_root string
	webui_root, err = cmd.PersistentFlags().GetString("webui-root")
	if err != nil {
		glog.Errorln("Get webui-root option error:", err.Error())
		return
	}
	http.HandleFunc("/api/v1/", apiservice.HandleFunc)
	http.Handle("/webui/", http.StripPrefix("/webui/", http.FileServer(http.Dir(webui_root))))
	http.ListenAndServe(fmt.Sprintf("%s:%d", cfgWebUIHost, cfgWebUIPort), nil)
}

func init() {

	RootCmd.AddCommand(startCmd)

	startCmd.PersistentFlags().StringSlice("nginx-status-page", nil, "The nginx status page to watch.")
	startCmd.PersistentFlags().String("webui-root", "public", "The root folder of web UI.")

	startCmd.Flags().Duration("watch-frenquency", 10*time.Second, "The interval between every update of the nginx status page.")

}
