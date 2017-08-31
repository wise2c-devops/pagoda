package main

import (
	"fmt"
	"os/exec"
	"sort"
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
	supported []string
	received  []string
	stopChan  chan bool
	nextChan  chan bool
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
		stopChan: make(chan bool),
		nextChan: make(chan bool),
	}
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
