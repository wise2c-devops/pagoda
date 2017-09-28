package main

import (
	"gitee.com/wisecloud/wise-deploy/database"
	"github.com/golang/glog"
)

type notifyChannel struct {
	name string
	c    chan *Notification
}

// ClusterRuntime all cluster in Runtime
// wschans should have cluser name as key if want to support multiple clusters
// be installed currently
type ClusterRuntime struct {
	cluster          map[string]*database.Cluster
	notificationChan chan *Notification
	wsChans          map[string]chan *Notification

	registeChan   chan *notifyChannel
	unregisteChan chan string
}

func NewClusterRuntime() *ClusterRuntime {
	return &ClusterRuntime{
		cluster:          make(map[string]*database.Cluster),
		notificationChan: make(chan *Notification),
		wsChans:          make(map[string]chan *Notification),
		registeChan:      make(chan *notifyChannel),
		unregisteChan:    make(chan string),
	}
}

func (cr *ClusterRuntime) Registe(name string) chan *Notification {
	c := make(chan *Notification)
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
		glog.Errorf("registe %s fail", name)
	}
}

func (cr *ClusterRuntime) Notify(n *Notification) {
	select {
	case cr.notificationChan <- n:
	default:
		glog.Errorf("notify fail %s", n.Stage)
	}
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
