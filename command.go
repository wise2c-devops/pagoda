package main

import (
	"os/exec"
	"path"

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
	installChan  chan []string
	resetChan    chan []string
	ansibleFile  string
}

func NewCommands() *Commands {
	return &Commands{
		currentIndex: -1,
		stopChan:     make(chan bool),
		nextChan:     make(chan bool, 1), //the length must be one
		installChan:  make(chan []string),
		resetChan:    make(chan []string),
	}
}

func (c *Commands) Launch(w string) {
	for {
		select {
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
				m := make(map[string]interface{})
				m["data"] = "hehe"
				ansibleChan <- m
			} else {
				c.complete(failed)
			}
		case rec := <-c.installChan:
			c.received = rec
			c.ansibleFile = "install.ansible"
			c.nextChan <- true
		case rec := <-c.resetChan:
			c.received = rec
			c.ansibleFile = "reset.ansible"
			c.nextChan <- true
		}
	}
}

func (c *Commands) Install(recv []string) {
	c.installChan <- recv
}

func (c *Commands) Reset(recv []string) {
	c.resetChan <- recv
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
				c.nextChan <- false
			} else {
				c.nextChan <- true
			}
		}()

		return
	}
}

func (c *Commands) complete(code CompleteCode) {
	switch code {
	case finished:
		glog.V(3).Info("complete all install step")
	case stopped:
		if err := c.currentCmd.Process.Kill(); err != nil {
			glog.Errorf("stop install error: %v", err)
		}
	case failed:
		glog.V(3).Info("failed at a step")
	}

	c.currentIndex = -1
}
