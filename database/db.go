package main

import (
	"fmt"
	"os"

	"gitee.com/wisecloud/wise-deploy/cluster"

	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	f := "cache.db"
	os.Remove(f)

	Orm, err := xorm.NewEngine("sqlite3", f)
	if err != nil {
		fmt.Println(err)
		return
	}
	Orm.ShowSQL(true)
	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	Orm.SetDefaultCacher(cacher)

	err = Orm.CreateTables(&cluster.Cluster{})
	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = Orm.Insert(&cluster.Cluster{ID: "xxx", Name: "xlw", Description: "gun"})
	if err != nil {
		fmt.Println(err)
		return
	}

	var users []cluster.Cluster
	err = Orm.Find(&users)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("users:", users)

	var users2 []cluster.Cluster
	err = Orm.Find(&users2)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("users2:", users2)

	var users3 []cluster.Cluster
	err = Orm.Find(&users3)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("users3:", users3)

	user4 := new(cluster.Cluster)
	has, err := Orm.Id("xxx").Get(user4)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("user4:", has, user4)

	user4.Name = "xiaolunwen"
	_, err = Orm.Id("xxx").Update(user4)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("user4:", user4)

	user5 := new(cluster.Cluster)
	has, err = Orm.Id("xxx").Get(user5)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("user5:", has, user5)

	user7 := new(cluster.Cluster)
	user7.Name = "d"
	user7.Description = ""
	user7.ID = "xxx"
	user7.Hosts = make([]*cluster.Host, 0)
	user7.Components = make([]*cluster.Component, 0)
	_, err = Orm.Id("xxx").Delete(user7)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		user6 := new(cluster.Cluster)
		has, err = Orm.Id("xxx").Get(user6)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("user6:", has, user6)
	}
}
