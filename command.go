package main

import (
	"fmt"
	"os/exec"
	"path"
	"sort"

	"gitee.com/wisecloud/wise-deploy/playbook"

	"github.com/golang/glog"
)

const (
	registry     = "registry"
	etcd         = "etcd"
	loadbalancer = "loadbalancer"
	mysql        = "mysql"
	k8sMaster    = "k8sMaster"
	k8sNode      = "k8sNode"
	wiseCloud    = "wiseCloud"
)

type Commands struct {
	supported    []string
	received     []string
	currentIndex int
	currentCmd   *exec.Cmd
	stopChan     chan bool
	nextChan     chan bool
	cmdChan      chan []string
}

func NewCommands() *Commands {
	return &Commands{
		supported: []string{
			registry,
			etcd,
			loadbalancer,
			mysql,
			k8sMaster,
			k8sNode,
			wiseCloud,
		},
		currentIndex: -1,
		stopChan:     make(chan bool),
		nextChan:     make(chan bool, 1),
		cmdChan:      make(chan []string),
	}
}

func (c *Commands) Launch(w string) {
	for {
		select {
		case <-c.stopChan:
			if err := c.currentCmd.Process.Kill(); err != nil {
				glog.Errorf("stop install error: %v", err)
			}
		case next := <-c.nextChan:
			if next {
				c.currentIndex++
				if c.currentIndex == len(c.supported) {
					glog.V(3).Info("complete all install step")
					c.currentIndex = -1
					break
				}
				for ; c.currentIndex < len(c.supported); c.currentIndex++ {
					step := c.supported[c.currentIndex]
					index := sort.SearchStrings(c.received, step)
					if index < len(c.received) && c.received[index] == step {
						cmd := exec.Command("ansible-playbook", "install.ansible")
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

						break
					}
				}
			} else {
				glog.V(3).Info("failed at a step")
				c.currentIndex = -1
			}
		case rec := <-c.cmdChan:
			c.received = rec
			sort.Strings(c.received)
			c.nextChan <- true
		}
	}
}

func (c *Commands) Start(recv []string) {
	c.cmdChan <- recv
}

func (c *Commands) Stop() {
	c.stopChan <- true
}

func (c *Commands) Run(rec []string) error {
	sort.Strings(rec)
	for _, i := range c.supported {
		if index := sort.SearchStrings(rec, i); index < len(rec) && rec[index] == i {
			cmd := exec.Command("ansible-playbook")
			cmd.Env = append(cmd.Env, "ANSIBLE_CONFIG=?")
			cmd.Env = append(cmd.Env, "ANSIBLE_host=?")

			go func() {
				err := cmd.Run()
				if err != nil {
					c.nextChan <- false
				} else {
					c.nextChan <- true
				}
			}()

			b := <-c.nextChan
			if !b {
				return fmt.Errorf("get error")
			}
		}
	}

	return nil
}
