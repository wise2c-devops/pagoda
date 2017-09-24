package main

import (
	"os/exec"
	"path"

	"gitee.com/wisecloud/wise-deploy/database"
	"gitee.com/wisecloud/wise-deploy/playbook"

	"github.com/golang/glog"
)

type CompleteCode int

const (
	finished CompleteCode = iota
	stopped
	failed
)

const (
	registry     = "registry"
	etcd         = "etcd"
	mysql        = "mysql"
	loadbalancer = "loadbalancer"
	k8sMaster    = "k8smaster"
	k8sNode      = "k8snode"
	wisecloud    = "wisecloud"
)

var (
	step = []string{
		registry,
		etcd,
		mysql,
		loadbalancer,
		k8sMaster,
		k8sNode,
		wisecloud,
	}
)

type Commands struct {
	received     []string
	currentIndex int
	currentCmd   *exec.Cmd
	stopChan     chan bool
	nextChan     chan bool
	installChan  chan *database.Cluster
	resetChan    chan *database.Cluster
	ansibleFile  string
	cluster      *database.Cluster
}

func NewCommands() *Commands {
	return &Commands{
		currentIndex: -1,
		received:     step,
		stopChan:     make(chan bool),
		nextChan:     make(chan bool, 1), //the length must be one
		installChan:  make(chan *database.Cluster),
		resetChan:    make(chan *database.Cluster),
	}
}

func (c *Commands) Launch(w string) {
	for {
		select {
		case n := <-ansibleChan:
			n.Stage = c.received[c.currentIndex]
			n.State = n.Task.State
			statsChan <- n
		case <-c.stopChan:
			c.complete(stopped)
		case next := <-c.nextChan:
			if next {
				c.currentIndex++
				if c.currentIndex == len(c.received) {
					c.complete(finished)
					break
				}
				c.run(w)
			} else {
				c.complete(failed)
			}
		case rec := <-c.installChan:
			c.cluster = rec
			c.ansibleFile = "install.ansible"
			c.nextChan <- true
		case rec := <-c.resetChan:
			c.cluster = rec
			c.ansibleFile = "reset.ansible"
			c.nextChan <- true
		}
	}
}

func (c *Commands) Install(cluster *database.Cluster) {
	c.installChan <- cluster
}

func (c *Commands) Reset(cluster *database.Cluster) {
	c.resetChan <- cluster
}

func (c *Commands) Stop() {
	c.stopChan <- true
}

func (c *Commands) run(w string) {
	for ; c.currentIndex < len(c.received); c.currentIndex++ {
		step := c.received[c.currentIndex]
		cmd := exec.Command("ansible-playbook", c.ansibleFile)
		cmd.Dir = path.Join(w, step+playbook.PlaybookSuffix)
		c.currentCmd = cmd

		go func() {
			glog.V(3).Infof("start step %s", step)
			err := cmd.Run()
			if err != nil {
				glog.V(3).Infof("step %s compeleted", step)
				c.nextChan <- false
			} else {
				glog.V(3).Infof("step %s failed", step)
				c.nextChan <- true
			}
		}()

		return
	}
}

func (c *Commands) complete(code CompleteCode) {
	switch code {
	case finished:
		c.cluster.State = database.Success
		glog.V(3).Info("complete all install step")
	case stopped:
		c.cluster.State = database.Failed
		if err := c.currentCmd.Process.Kill(); err != nil {
			glog.Errorf("stop install error: %v", err)
		}
	case failed:
		c.cluster.State = database.Failed
		glog.V(3).Info("failed at a step")
	}

	err := database.Instance(sqlConfig).UpdateCluster(c.cluster)
	if err != nil {
		glog.Errorf("update cluster %s error %v", c.cluster.ID, err)
		return
	}
	c.currentIndex = -1
}
