package runtime

import (
	"fmt"
	"os/exec"
	"path"
	"sort"

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
	ComponentMap = map[string]int{
		registry:     0,
		etcd:         1,
		mysql:        2,
		loadbalancer: 3,
		k8sMaster:    4,
		k8sNode:      5,
		wisecloud:    6,
	}
)

type LaunchParameters struct {
	Operation  string            `json:"operation"`
	Components []string          `json:"components"`
	Clster     *database.Cluster `json:"-"`
}

type Commands struct {
	received     []string
	currentIndex int
	currentCmd   *exec.Cmd
	stopChan     chan bool
	nextChan     chan bool
	processChan  chan bool
	startChan    chan *LaunchParameters
	ansibleFile  string
	Cluster      *database.Cluster

	runtime *ClusterRuntime
}

func NewCommands() *Commands {
	return &Commands{
		currentIndex: -1,
		stopChan:     make(chan bool),
		nextChan:     make(chan bool, 1), //the length must be one
		processChan:  make(chan bool, 1), //the length must be one
		startChan:    make(chan *LaunchParameters),
		currentCmd:   nil,
	}
}

func (c *Commands) Launch(w string, runtime *ClusterRuntime) {
	c.runtime = runtime

	for {
		select {
		// case n := <-ansibleChan:
		// 	if c.currentIndex == -1 {
		// 		glog.Error("received a improper notify")
		// 		break
		// 	}
		// 	n.Stage = c.received[c.currentIndex]
		// 	clusterRuntime.Notify(c.cluster, n)
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
				c.runtime.RotateStage(c.Cluster.ID, c.received[c.currentIndex])
			} else {
				c.complete(failed)
			}
		case rec := <-c.startChan:
			c.start(rec)
		}
	}
}

func (c *Commands) process() error {
	select {
	case c.processChan <- true:
		return nil
	default:
		return fmt.Errorf("I'm processing a action")
	}
}

func (c *Commands) unProcess() error {
	select {
	case <-c.processChan:
		return nil
	default:
		return fmt.Errorf("I have processed a action")
	}
}

func (c *Commands) start(config *LaunchParameters) {
	sort.Sort(ByName(config.Components))
	c.received = config.Components
	c.Cluster = config.Clster

	if config.Operation == "install" {
		c.ansibleFile = "install.ansible"
	} else {
		c.ansibleFile = "clean.ansible"
	}

	c.nextChan <- true
	c.runtime.ProcessCluster(c.Cluster)
}

func (c *Commands) Install(cluster *database.Cluster, config *LaunchParameters) error {
	if err := c.process(); err != nil {
		return err
	}
	glog.V(3).Infof("begin to install cluster %s", cluster.Name)

	config.Clster = cluster
	c.startChan <- config

	return nil
}

func (c *Commands) Stop() {
	c.stopChan <- true
}

func (c *Commands) run(w string) {
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
}

func (c *Commands) complete(code CompleteCode) {
	switch code {
	case finished:
		if c.ansibleFile == "install.ansible" {
			c.Cluster.State = database.Success
			glog.V(3).Info("complete all install step")
		} else if c.ansibleFile == "clean.ansible" {
			c.Cluster.State = database.Initial
			glog.V(3).Info("complete all reset step")
		}
	case stopped:
		if c.Cluster == nil {
			glog.Warning("receive a stop but I haven't start")
			return
		}
		c.Cluster.State = database.Failed
		if err := c.currentCmd.Process.Kill(); err != nil {
			glog.Errorf("stop install error: %v", err)
		}
	case failed:
		c.Cluster.State = database.Failed
		glog.V(3).Info("failed at a step")
	}

	err := database.Default().UpdateCluster(c.Cluster)
	if err != nil {
		glog.Errorf("update cluster %s error %v", c.Cluster.ID, err)
		return
	}
	c.unProcess()

	glog.V(3).Info("finish a install/reset")
	c.currentIndex = -1
	c.runtime.RmCluster(c.Cluster.ID)
}
