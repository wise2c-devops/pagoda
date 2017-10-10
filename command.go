package main

import (
	"fmt"
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
	initHost     = "init"
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
		initHost,
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
	processChan  chan bool
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
		processChan:  make(chan bool, 1), //the length must be one
		installChan:  make(chan *database.Cluster),
		resetChan:    make(chan *database.Cluster),
		currentCmd:   nil,
	}
}

func (c *Commands) Launch(w string) {
	for {
		select {
		case n := <-ansibleChan:
			if c.currentIndex == -1 {
				glog.Error("received a improper notify")
				break
			}
			n.Stage = c.received[c.currentIndex]
			clusterRuntime.Notify(c.cluster, n)
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
				clusterRuntime.RotateStage(c.cluster.ID, c.received[c.currentIndex])
			} else {
				c.complete(failed)
			}
		case rec := <-c.installChan:
			c.cluster = rec
			c.ansibleFile = "install.ansible"
			c.nextChan <- true
			clusterRuntime.ProcessCluster(c.cluster)
		case rec := <-c.resetChan:
			c.cluster = rec
			c.ansibleFile = "clean.ansible"
			c.nextChan <- true
			clusterRuntime.ProcessCluster(c.cluster)
		}
	}
}

func (c *Commands) Process() error {
	select {
	case c.processChan <- true:
		return nil
	default:
		return fmt.Errorf("I'm processing a action")
	}
}

func (c *Commands) Install(cluster *database.Cluster) {
	glog.V(3).Infof("begin to install cluster %s", cluster.Name)
	c.installChan <- cluster
}

func (c *Commands) Reset(cluster *database.Cluster) {
	glog.V(3).Infof("begin to reset cluster %s", cluster.Name)
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
				glog.V(3).Infof("step %s failed ", step)
				c.nextChan <- false
			} else {
				glog.V(3).Infof("step %s compeleted", step)
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
		if c.cluster == nil {
			glog.Warning("receive a stop but I haven't start")
			return
		}
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
	select {
	case <-c.processChan:
	default:
	}
	glog.V(3).Info("finish a install/reset")
	c.currentIndex = -1
	clusterRuntime.RmCluster(c.cluster.ID)
}
