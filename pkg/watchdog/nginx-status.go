package watchdog

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/golang/glog"

	"github.com/zhangpeihao/watchdog/pkg/alarm"
)

type nginxStatus struct {
	page       string
	client     *http.Client
	reg        *regexp.Regexp
	errorHosts map[string]*NginxStatusResult
}

type NginxStatusResult struct {
	Name       string
	Host       string
	Status     string
	RiseCounts int
	FallCounts int
	CheckType  string
	RiseAt     int64
	AlarmId    string
}

// Create a new watch job for nginx status page
func NewNginxStatus(page string) (job WatchJob, err error) {
	glog.Info("watchdog::NewNginxStatus")
	defer glog.Info("watchdog::NewNginxStatus end")
	w := &nginxStatus{
		page:       page,
		client:     new(http.Client),
		reg:        regexp.MustCompile(`(?s:<tr bgcolor="#FF0000">(.*?)</tr>)`),
		errorHosts: make(map[string]*NginxStatusResult),
	}
	job = func() {
		w.Watch()
	}
	return
}

func (w *nginxStatus) Watch() {

	glog.Infoln("Start watch nginx status page:", w.page)
	resp, err := w.client.Get(w.page)
	if err != nil {
		glog.Warningln("Get", w.page, "error:", err.Error())
		return
	}
	defer resp.Body.Close()

	/*
		<tr bgcolor="#FF0000">
		    <td>0</td>
		    <td>arch-sms</td>
		    <td>10.6.80.210:8080</td>
		    <td>down</td>
		    <td>0</td>
		    <td>71516</td>
		    <td>http</td>
		    <td>0</td>
		</tr>
	*/

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Warningln("Read from response error:", err.Error())
		return
	}
	recoredHosts := make(map[string]struct{})
	for host, _ := range w.errorHosts {
		recoredHosts[host] = struct{}{}
	}
	for _, block := range w.reg.FindAllString(string(buf), -1) {
		// The 3rd line is service name
		lines := strings.Split(block, "\n")
		if len(lines) < 8 {
			continue
		}
		host := trimLine(lines[3])
		if s, found := w.errorHosts[host]; !found {
			// New error host
			s = &NginxStatusResult{
				Name:      trimLine(lines[2]),
				Host:      trimLine(lines[3]),
				Status:    trimLine(lines[4]),
				CheckType: trimLine(lines[7]),
				RiseAt:    time.Now().UnixNano(),
			}
			s.RiseCounts, err = strconv.Atoi(trimLine(lines[5]))
			s.FallCounts, err = strconv.Atoi(trimLine(lines[6]))
			glog.Infoln("Host[", s.Name, "-", host, "] error")
			w.errorHosts[host] = s
			s.AlarmId = alarm.Rise("Nginx", "Host["+s.Name+" - "+s.Host+"] error")
			// TODO: Alarm to monitor
		} else {
			delete(recoredHosts, host)
		}
	}

	for host, _ := range recoredHosts {
		if s, found := w.errorHosts[host]; !found {
			glog.Infoln("Host[", s.Name, "-", host, "] recovered")
			alarm.Remove(s.AlarmId)
			// TODO: Stop the alarm from monitor
			delete(w.errorHosts, host)
		}
	}
}

func trimLine(line string) string {
	name := strings.TrimSpace(line)
	name = strings.TrimPrefix(name, "<td>")
	name = strings.TrimSuffix(name, "</td>")
	return name
}
