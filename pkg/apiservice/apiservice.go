package apiservice

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/golang/glog"

	"github.com/zhangpeihao/watchdog/pkg/alarm"
)

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	glog.Infoln("apiservice::HandleFunc()", r.URL.Path)
	defer glog.Infoln("apiservice::HandleFunc end")
	switch r.URL.Path {
	case "/api/v1/alarms":
		var alarms alarm.Alarms
		alarm.Foreach(func(a *alarm.Alarm) {
			alarms = append(alarms, alarm.Alarm{
				RiseAt:  a.RiseAt,
				Type:    a.Type,
				Content: a.Content,
				ID:      a.ID,
			})
		})
		sort.Sort(alarms)
		buf, err := json.Marshal(&ApiResponse{0, "OK", alarms})
		if err != nil {
			glog.Errorln("JSON marshal error:", err.Error())
			w.WriteHeader(500)
			return
		}
		w.Write(buf)
		return
	}
	w.WriteHeader(404)
}
