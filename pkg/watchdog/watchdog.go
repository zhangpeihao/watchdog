package watchdog

import (
	"time"

	"github.com/golang/glog"

	"k8s.io/kubernetes/pkg/util/wait"
)

type WatchJob func()

// Run watch jobs by interval setting
func RunWatchJob(jobs []WatchJob, interval time.Duration) {
	glog.V(3).Infoln("watchdog::RunWatchJob")
	defer glog.V(3).Infoln("watchdog::RunWatchJob end")
	wait.Until(func() {
		for _, job := range jobs {
			job()
		}
	}, interval, wait.NeverStop)
}
