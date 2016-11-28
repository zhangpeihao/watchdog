package alarm

import (
	"crypto/md5"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/golang/glog"
)

var (
	list map[string]*Alarm = make(map[string]*Alarm)
	lock sync.Mutex
)

type Alarm struct {
	RiseAt  int64  `json:"rise-at"`
	Type    string `json:"type"`
	Content string `json:"content"`
	ID      string `json:"id"`
}

type Alarms []Alarm

func (a Alarms) Len() int           { return len(a) }
func (a Alarms) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Alarms) Less(i, j int) bool { return a[i].RiseAt < a[j].RiseAt }

func Rise(t, content string) string {
	glog.Infoln("alarm::Rise " + t + "," + content)
	defer glog.Infoln("alarm::Rise end")
	lock.Lock()
	defer lock.Unlock()
	a := Alarm{
		RiseAt:  time.Now().UnixNano(),
		Type:    t,
		Content: content,
	}
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%d-%s-%s", a.RiseAt, a.Type, a.Content))
	a.ID = fmt.Sprintf("%X", h.Sum(nil))
	list[a.ID] = &a
	return a.ID
}

func Remove(id string) {
	glog.Infoln("alarm::Remove")
	defer glog.Infoln("alarm::Remove end")
	lock.Lock()
	defer lock.Unlock()
	if _, found := list[id]; found {
		delete(list, id)
	}
}

func Foreach(fn func(*Alarm)) {
	lock.Lock()
	defer lock.Unlock()
	for _, a := range list {
		fn(a)
	}
}
