package database

import (
	"flag"
	"fmt"
	"strconv"
	"time"

	"gitee.com/wisecloud/wise-deploy/cluster"

	"github.com/go-xorm/xorm"
	"github.com/golang/glog"
	_ "github.com/mattn/go-sqlite3"
)

type SQLEngine struct {
	xe *xorm.Engine
}

type EngineConfig struct {
	SQLType      string
	ShowSQL      bool
	ShowExecTime bool
}

func NewEngine(config *EngineConfig) (*SQLEngine, error) {
	e := &SQLEngine{}
	if config.SQLType == "sqlite3" {
		engine, err := xorm.NewEngine("sqlite3", "cluster.db")
		if err != nil {
			return nil, err
		}
		e.xe = engine
	}

	e.xe.ShowSQL(config.ShowSQL)
	e.xe.ShowExecTime(config.ShowExecTime)
	if err := e.xe.Sync2(new(cluster.Cluster)); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *SQLEngine) RetrieveClusters() (clusters []*cluster.Cluster, err error) {
	err = e.xe.Find(&clusters)
	return
}

func (e *SQLEngine) CreateCluster(c *cluster.Cluster) error {
	if _, err := e.xe.InsertOne(c); err != nil {
		return err
	}

	return nil
}

func (e *SQLEngine) DeleteCluster(id string) error {
	_, err := e.xe.Exec("delete from cluster where id = ?", id)
	return err
}

func (e *SQLEngine) UpdateCluster(c *cluster.Cluster) error {
	_, err := e.xe.Update(c)
	return err
}

func (e *SQLEngine) RetrieveCluster(id string) (c *cluster.Cluster, err error) {
	_, err = e.xe.Get(c)
	return
}

func main() {
	flag.Parse()

	var err error
	engine, err := xorm.NewEngine("sqlite3", "cluster.db")
	if err != nil {
		glog.Exit(err)
	}
	if v := flag.Lookup("v"); v != nil {
		i, _ := strconv.Atoi(v.Value.String())
		if i >= 3 {
			engine.ShowSQL(true)
		}
	}

	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	engine.SetDefaultCacher(cacher)

	err = engine.CreateTables(&cluster.Cluster{})
	if err != nil {
		fmt.Println(err)
		return
	}

	id := time.Now().Format("2006-01-02 15:04:05")
	_, err = engine.Insert(&cluster.Cluster{ID: id, Name: "xlw", Description: "gun"})
	if err != nil {
		fmt.Println(err)
		return
	}

	var users []cluster.Cluster
	err = engine.Find(&users)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("users:", users)

	var users2 []cluster.Cluster
	err = engine.Find(&users2)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("users2:", users2)

	var users3 []cluster.Cluster
	err = engine.Find(&users3)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("users3:", users3)

	user4 := new(cluster.Cluster)
	has, err := engine.Id("xxx").Get(user4)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("user4:", has, user4)

	user4.Name = "xiaolunwen"
	_, err = engine.Id("xxx").Update(user4)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("user4:", user4)

	user5 := new(cluster.Cluster)
	has, err = engine.Id("xxx").Get(user5)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("user5:", has, user5)

	user7 := new(cluster.Cluster)
	user7.Name = "d"
	user7.Description = ""
	user7.ID = "xxx"
	_, err = engine.Id("xxx").Delete(user7)
	if err != nil {
		fmt.Println(err)
		return
	}
}
