/*
@Time : 2019/3/19 14:49 
@Author : lukebryan
@File : conf
@Software: GoLand
*/
package conf

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"gopkg.in/olivere/elastic.v5"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

var client *elastic.Client

var db *gorm.DB

func GetClient() *elastic.Client {
	return client
}
func GetDB() *gorm.DB {
	return db
}

type config struct {
	Port       string `json:"Port"`
	UserName       string `json:"UserName"`
	Password string `json:"Password"`
	Ip string `json:"Ip"`
	DBPort string `json:"DBPort"`
	DBName string `json:"DBName"`
	ESHost       string `json:"ESHost"`
}

var Config = &config{}

//初始化
func init() {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic("Client config read err")
	}
	err = json.Unmarshal(b, Config)
	if err != nil {
		panic(err)
	}

	//Elasticsearch
	errorlog := log.New(os.Stdout, "APP", log.LstdFlags)
	client, err = elastic.NewClient(elastic.SetErrorLog(errorlog), elastic.SetURL(Config.ESHost))
	if err != nil {
		panic(err)
	}
	info, code, err := client.Ping(Config.ESHost).Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := client.ElasticsearchVersion(Config.ESHost)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)

	//Database
	path := strings.Join([]string{Config.UserName, ":", Config.Password, "@tcp(", Config.Ip, ":", Config.DBPort, ")/", Config.DBName, "?charset=utf8&parseTime=true"}, "")
	db, err = gorm.Open("mysql", path)
	if err != nil {
		panic(err)
	}
	db.SingularTable(true)
	db.DB().SetConnMaxLifetime(1 * time.Second)
	db.DB().SetMaxIdleConns(20)   //最大打开的连接数
	db.DB().SetMaxOpenConns(2000) //设置最大闲置个数
	db.SingularTable(true)	//表生成结尾不带s
	// 启用Logger，显示详细日志
	db.LogMode(true)

}
