/*
@Time : 2019/3/26 9:44 
@Author : lukebryan
@File : elasticutil
@Software: GoLand
*/
package utils

import (
	"context"
	"elasticsearch/conf"
	"log"
	"time"
)

//创建
func Create(body interface{},index string,typ string,id string,count int) bool {

	if count >0 {
		time.Sleep(time.Second * 3)
	}

	_, err := conf.GetClient().Index().
		Index(index).
		Type(typ).
		Id(id).
		BodyJson(body).
		Do(context.Background())
	if err != nil {
		Create(body,index,typ,id,count+1)
		if count > 5 {
			log.Print(err)
			return false
		}
	}
	return true
}

//删除
func Delete(index string,typ string,id string) bool {

	_, err := conf.GetClient().Delete().
		Index(index).
		Type(typ).
		Id(id).
		Do(context.Background())
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

//修改
func Update(maps map[string]interface{},index string,typ string,id string) bool {
	_, err := conf.GetClient().Update().
		Index(index).
		Type(typ).
		Id(id).
		Doc(maps).
		Do(context.Background())
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

