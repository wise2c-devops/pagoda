package runtime

import (
	"fmt"
	"sort"
	"sync"

	"gitee.com/wisecloud/wise-deploy/database"

	"github.com/golang/glog"
)

// Notify - notify to runtime a ansible callback message
func Notify(n *database.Notification) {
	runtime.notify(n)
}

// Register - register a websocket client for receive notification
func Register(name string) chan *database.Notification {
	return runtime.register(name)
}

// Annul - annul a websocket client
func Annul(name string) {
	runtime.annul(name)
}

// RetrieveStatus - retrieve a cluster runtime status
func RetrieveStatus(clusterID string) (*RunningStatus, error) {
	return runtime.retrieveStatus(clusterID)
}

type byName []string

func (s byName) Len() int {
	return len(s)
}

func (s byName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byName) Less(i, j int) bool {
	return componentMap[s[i]] < componentMap[s[j]]
}

type notifyChannel struct {
	name string
	c    chan *database.Notification
}

// RunningStatus - record a cluster running status
type RunningStatus struct {
	Stages       []string `json:"stages"`
	CurrentStage string   `json:"currentStage"`
}

// Runtime all cluster in Runtime
// wschans should have cluster name as key if want to support multiple clusters
// be installing currently
type Runtime struct {
	runningStatus *RunningStatus
	clusterID     string
	mux           *sync.Mutex

	notificationChan chan *database.Notification
	wsChans          map[string]chan *database.Notification
	registerChan     chan *notifyChannel
	removeChan       chan string
}

func newRuntime() *Runtime {
	return &Runtime{
		mux:              &sync.Mutex{},
		notificationChan: make(chan *database.Notification, 5),
		wsChans:          make(map[string]chan *database.Notification),
		registerChan:     make(chan *notifyChannel),
		removeChan:       make(chan string),
	}
}

func (cr *Runtime) startOperate(c *database.Cluster) {
	cr.clusterID = c.ID
	database.Instance().DeleteLogs(c.ID)

	//sorted components
	sc := make([]string, 0, len(c.Components))
	for _, cc := range c.Components {
		sc = append(sc, cc.Name)
	}
	sort.Sort(byName(sc))

	cr.mux.Lock()
	defer cr.mux.Unlock()
	cr.runningStatus = &RunningStatus{
		Stages: sc,
	}
}

func (cr *Runtime) stopOperate() {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	cr.clusterID, cr.runningStatus = "", nil
}

func (cr *Runtime) rotateStage(clusterID, name string) {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	cr.runningStatus.CurrentStage = name
}

func (cr *Runtime) retrieveStatus(clusterID string) (*RunningStatus, error) {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	if cr.runningStatus != nil {
		return cr.runningStatus, nil
	}

	return nil, fmt.Errorf("no cluster in operating now")
}

func (cr *Runtime) register(name string) chan *database.Notification {
	c := make(chan *database.Notification)
	nc := &notifyChannel{
		name: name,
		c:    c,
	}
	select {
	case cr.registerChan <- nc:
	default:
		glog.Errorf("registe %s fail", name)
	}

	return c
}

func (cr *Runtime) annul(name string) {
	select {
	case cr.removeChan <- name:
	default:
		glog.Errorf("unregiste %s fail", name)
	}
}

func (cr *Runtime) notify(n *database.Notification) {
	if err := database.Instance().CreateLog(cr.clusterID, n); err != nil {
		glog.Error(err)
	}

	select {
	case cr.notificationChan <- n:
	default:
		glog.Errorf("notify fail %s", n.Stage)
	}
	//TODO: set cluster status to Notification's stage
}

func (cr *Runtime) run() {
	for {
		select {
		case n := <-cr.registerChan:
			cr.wsChans[n.name] = n.c
			glog.V(2).Infof("add an observer: %s", n.name)
		case name := <-cr.removeChan:
			c, find := cr.wsChans[name]
			if !find {
				glog.Warningf("can't find %s's chan", name)
				return
			}

			delete(cr.wsChans, name)
			close(c)
			glog.V(2).Infof("remove an observer: %s", name)
		case event := <-cr.notificationChan:
			for k, v := range cr.wsChans {
				select {
				case v <- event:
				default:
					glog.Errorf("send notify to %s fail", k)
				}
			}
		}
	}
}
