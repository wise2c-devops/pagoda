package runtime

import (
	"fmt"
	"sort"
	"sync"

	"gitee.com/wisecloud/wise-deploy/database"
	"github.com/golang/glog"
)

type ByName []string

func (s ByName) Len() int {
	return len(s)
}
func (s ByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByName) Less(i, j int) bool {
	return ComponentMap[s[i]] < ComponentMap[s[j]]
}

type notifyChannel struct {
	name string
	c    chan *database.Notification
}

type ClusterStatus struct {
	Stages       []string `json:"stages"`
	CurrentStage string   `json:"currentStage"`
}

// ClusterRuntime all cluster in Runtime
// wschans should have cluster name as key if want to support multiple clusters
// be installing currently
type ClusterRuntime struct {
	cluster map[string]*ClusterStatus
	mux     *sync.Mutex

	notificationChan chan *database.Notification
	wsChans          map[string]chan *database.Notification
	registeChan      chan *notifyChannel
	unregisteChan    chan string
}

func NewClusterRuntime() *ClusterRuntime {
	return &ClusterRuntime{
		cluster:          make(map[string]*ClusterStatus),
		mux:              &sync.Mutex{},
		notificationChan: make(chan *database.Notification, 5),
		wsChans:          make(map[string]chan *database.Notification),
		registeChan:      make(chan *notifyChannel),
		unregisteChan:    make(chan string),
	}
}

func (cr *ClusterRuntime) ProcessCluster(c *database.Cluster) {
	//sorted components
	sc := make([]string, 0, len(c.Components))
	for _, cc := range c.Components {
		sc = append(sc, cc.Name)
	}

	sort.Sort(ByName(sc))
	clusterStatus := &ClusterStatus{
		Stages: sc,
	}

	database.Default().DeleteLogs(c.ID)

	cr.mux.Lock()
	defer cr.mux.Unlock()

	cr.cluster[c.ID] = clusterStatus
}

func (cr *ClusterRuntime) RmCluster(clusterID string) {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	if _, find := cr.cluster[clusterID]; find {
		delete(cr.cluster, clusterID)
	} else {
		glog.Errorf("can't find cluster %s", clusterID)
	}
}

func (cr *ClusterRuntime) RotateStage(clusterID, name string) {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	cr.cluster[clusterID].CurrentStage = name
}

func (cr *ClusterRuntime) RetrieveStatus(clusterID string) (*ClusterStatus, error) {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	if s, find := cr.cluster[clusterID]; find {
		return s, nil
	}

	return nil, fmt.Errorf("can't find cluster %s", clusterID)
}

func (cr *ClusterRuntime) Registe(name string) chan *database.Notification {
	c := make(chan *database.Notification)
	nc := &notifyChannel{
		name: name,
		c:    c,
	}
	select {
	case cr.registeChan <- nc:
	default:
		glog.Errorf("registe %s fail", name)
	}

	return c
}

func (cr *ClusterRuntime) Unregiste(name string) {
	select {
	case cr.unregisteChan <- name:
	default:
		glog.Errorf("unregiste %s fail", name)
	}
}

func (cr *ClusterRuntime) Notify(c *database.Cluster, n *database.Notification) {
	if err := database.Default().CreateLog(c.ID, n); err != nil {
		glog.Error(err)
	}

	select {
	case cr.notificationChan <- n:
	default:
		glog.Errorf("notify fail %s", n.Stage)
	}
	//TODO: set cluster status to Notification's stage
}

func (cr *ClusterRuntime) Run() {
	for {
		select {
		case n := <-cr.registeChan:
			cr.wsChans[n.name] = n.c
			glog.V(2).Infof("add an observer: %s", n.name)
		case name := <-cr.unregisteChan:
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
