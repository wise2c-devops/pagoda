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

// Run - start command & runtime
func Run(w string) {
	go command.Launch(w)
	go runtime.run()
}

// StartOperate - start operate a cluster
func StartOperate(cluster *database.Cluster, config *LaunchParameters) error {
	command.deploySeed = playbook.NewDeploySeed(cluster, command.workDir)

	if err := playbook.PreparePlaybooks(command.workDir, command.deploySeed); err != nil {
		return err
	}

	return command.StartOperate(cluster, config)
}

// StopOperate - stop a cluser's operation
func StopOperate() {
	command.StopOperate()
}

type completeCode int

const (
	finished completeCode = iota
	stopped
	failed
)

const (
	installFile      = "install.ansible"
	resetFile        = "reset.ansible"
	installOperation = "install"
	resetOperation   = "reset"
)

const (
	docker       = "docker"
	registry     = "registry"
	etcd         = "etcd"
	mysql        = "mysql"
	loadbalancer = "loadbalancer"
	k8sMaster    = "k8smaster"
	k8sNode      = "k8snode"
	wisecloud    = "wisecloud"
)

var (
	componentMap = map[string]int{
		docker:       0,
		registry:     1,
		etcd:         2,
		mysql:        3,
		loadbalancer: 4,
		k8sMaster:    5,
		k8sNode:      6,
		wisecloud:    7,
	}

	runtime = newRuntime()
	command = newCommand(runtime)
)

// LaunchParameters - needed when start operate a cluser
type LaunchParameters struct {
	Operation  string            `json:"operation"`
	Components []string          `json:"components"`
	Cluster    *database.Cluster `json:"-"`
}

type commandT struct {
	received     []string
	currentIndex int
	currentCmd   *exec.Cmd

	cluster     *database.Cluster
	ansibleFile string
	workDir     string
	deploySeed  *playbook.DeploySeed2

	stopChan    chan bool
	nextChan    chan bool
	processChan chan bool
	startChan   chan *LaunchParameters

	runtime *Runtime
}

func newCommand(runtime *Runtime) *commandT {
	return &commandT{
		currentIndex: -1,
		stopChan:     make(chan bool),
		nextChan:     make(chan bool, 1), //the length must be one
		processChan:  make(chan bool, 1), //the length must be one
		startChan:    make(chan *LaunchParameters),
		runtime:      runtime,
	}
}

func (c *commandT) Launch(w string) {
	c.workDir = w
	for {
		select {
		case <-c.stopChan:
			c.complete(stopped)
		case next := <-c.nextChan:
			if next && c.cluster != nil {
				c.currentIndex++
				if c.currentIndex == len(c.received) {
					c.complete(finished)
					break
				}
				c.run(w)
				c.runtime.rotateStage(c.cluster.ID, c.received[c.currentIndex])
			} else {
				c.complete(failed)
			}
		case rec := <-c.startChan:
			c.start(rec)
		}
	}
}

func (c *commandT) acquire() error {
	select {
	case c.processChan <- true:
		return nil
	default:
		return fmt.Errorf("I'm processing a action")
	}
}

func (c *commandT) release() error {
	select {
	case <-c.processChan:
		return nil
	default:
		return fmt.Errorf("I have processed a action")
	}
}

func (c *commandT) start(config *LaunchParameters) {
	if config.Operation == installOperation {
		c.ansibleFile = installFile
	} else if config.Operation == resetOperation {
		c.ansibleFile = resetFile
	} else {
		glog.Errorf("get unexpected operation %s", config.Operation)
		return
	}

	sort.Sort(byName(config.Components))
	c.received = config.Components
	c.cluster = config.Cluster
	c.nextChan <- true
	c.runtime.startOperate(c.cluster)
}

func (c *commandT) StartOperate(cluster *database.Cluster, config *LaunchParameters) error {
	if err := c.acquire(); err != nil {
		return err
	}
	glog.V(3).Infof("begin to install cluster %s", cluster.Name)

	config.Cluster = cluster
	c.startChan <- config

	return nil
}

func (c *commandT) StopOperate() {
	c.stopChan <- true
}

func (c *commandT) run(w string) {
	step := c.received[c.currentIndex]
	cmd := exec.Command("ansible-playbook", c.ansibleFile)
	cmd.Dir = path.Join(
		w,
		step+playbook.PlaybookSuffix,
		map[string]*playbook.Component(*c.deploySeed)[step].Version,
	)
	c.currentCmd = cmd

	go func() {
		glog.V(3).Infof("start step %s", step)
		output, err := cmd.CombinedOutput()
		if err != nil {
			glog.V(3).Infof("step %s failed: %v with output: %s", step, err, output)
			c.nextChan <- false
		} else {
			glog.V(3).Infof("step %s compeleted", step)
			c.nextChan <- true
		}
	}()
}

func (c *commandT) complete(code completeCode) {
	defer c.release()
	switch code {
	case finished:
		if c.ansibleFile == installFile {
			c.cluster.State = database.Success
			glog.V(3).Info("complete all install step")
		} else if c.ansibleFile == resetFile {
			c.cluster.State = database.Initial
			glog.V(3).Info("complete all reset step")
		}
	case stopped:
		if c.cluster == nil {
			glog.Warning("receive a stop but I haven't start")
			return
		}
		if err := c.currentCmd.Process.Kill(); err != nil {
			c.cluster.State = database.Failed
			glog.Errorf("stop install error: %v", err)
		}
		c.cluster = nil
		c.currentCmd = nil
	case failed:
		c.cluster.State = database.Failed
		glog.V(3).Info("failed at a step")
	}

	err := database.Instance().UpdateCluster(c.cluster)
	if err != nil {
		glog.Errorf("update cluster %s error %v", c.cluster.ID, err)
		return
	}

	c.currentIndex = -1
	c.runtime.stopOperate(c.cluster.ID)
}
