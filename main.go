/*
@Time : 2019/3/21 16:00 
@Author : lukebryan
@File : main
@Software: GoLand

go-mysql-elasticsearch
*/
package main

import (
	"context"
	"elasticsearch/conf"
	"elasticsearch/models"
	"elasticsearch/services"
	"encoding/json"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"reflect"
)

var esService = services.NewElasticService()
var searchService = services.NewSearchService()

type Car struct {
	Name string   `json:"name"`
	LPN  string   `json:"lpn"`
	Type       string      `json:"type"`	//大型车,中型车,小型车
	Brand     string   `json:"brand"`		//品牌:大众,丰田
	Driver string `json:"driver"`			//司机
}

func main() {
	defer func() {
		conf.GetDB().Close()

	}()
	//tableStrings := []string{"company","category","mst_expert"}
	//tableStrings := []string{"company"}
	//
	//for i := range tableStrings {
	//	results,success := esService.Create(tableStrings[i])
	//	if success {
	//		for i := range results{
	//			fmt.Println(results[i])
	//		}
	//	}
	//}



	//sql := fmt.Sprintf(`SELECT table_name FROM information_schema.tables WHERE
	//				table_name RLIKE "%s" AND table_schema = "%s";`, "company", "name")
	//fmt.Println(sql)


	text := "中"
	searchFields := make([]string,1)
	searchFields[0] = "abbreviation"
	//searchFields[1] = "name"
	resultList := searchService.Search(text,"company",searchFields)
	for i := range resultList.CompanyList{
		fmt.Println(resultList.CompanyList[i])
	}
	for i := range resultList.UserList{
		fmt.Println(resultList.UserList[i])
	}

	//c:=make(chan struct{})
	//<-c
}




/*下面是简单的CURD*/

//创建
func create() {

	//使用结构体
	e1 := Car{"Jane's car", "湘A:13ZD34", "大型车", "大众", "Jane"}
	put1, err := conf.GetClient().Index().
		Index("car").
		Type("bigcar").
		Id("1").
		BodyJson(e1).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put1.Id, put1.Index, put1.Type)

	//使用字符串
	e2 := `{"name":"John's car","LPN":"湘A:24ZD56","Type":"中型车","Brand":"捷豹","Driver":"John"}`
	put2, err := conf.GetClient().Index().
		Index("car").
		Type("bigcar").
		Id("2").
		BodyJson(e2).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index s%s, type %s\n", put2.Id, put2.Index, put2.Type)

	e3 := `{"name":"Peter's car","LPN":"湘A:35ZD78","Type":"小型车","Brand":"特斯拉","Driver":"Peter"}`
	put3, err := conf.GetClient().Index().
		Index("car").
		Type("bigcar").
		Id("3").
		BodyJson(e3).
		Do(context.Background())
	if err != nil {
		panic(err)
	}
	fmt.Printf("Indexed tweet %s to index %s, type %s\n", put3.Id, put3.Index, put3.Type)

}

//删除
func delete() {

	res, err := conf.GetClient().Delete().Index("car").
		Type("bigcar").
		Id("1").
		Do(context.Background())
	if err != nil {
		println(err.Error())
		return
	}
	fmt.Printf("delete result %s\n", res.Result)
}

//修改
func update() {
	res, err := conf.GetClient().Update().
		Index("car").
		Type("bigcar").
		Id("1").
		Doc(map[string]interface{}{"name": "Jane's car"}).
		Do(context.Background())
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("update age %s\n", res.Result)

}

//查找
func gets() {
	//通过id查找
	get1, err := conf.GetClient().Get().Index("megacorp").Type("car").Id("2").Do(context.Background())
	if err != nil {
		panic(err)
	}
	if get1.Found {
		fmt.Printf("Got document %s in version %d from index %s, type %s\n", get1.Id, get1.Version, get1.Index, get1.Type)
	}
}

//搜索
func query() {
	var res *elastic.SearchResult
	var err error
	//取所有
	res, err = conf.GetClient().Search("megacorp").Type("car").Do(context.Background())
	printcar(res, err)

	//字段相等
	q := elastic.NewQueryStringQuery("last_name:Smith")
	res, err = conf.GetClient().Search("megacorp").Type("car").Query(q).Do(context.Background())
	if err != nil {
		println(err.Error())
	}
	printcar(res, err)

	if res.Hits.TotalHits > 0 {
		fmt.Printf("Found a total of %d car \n", res.Hits.TotalHits)

		for _, hit := range res.Hits.Hits {

			var t Car
			err := json.Unmarshal(*hit.Source, &t) //另外一种取数据的方法
			if err != nil {
				fmt.Println("Deserialization failed")
			}

			//fmt.Printf("car name %s : %s\n", t.name, t.Driver)
		}
	} else {
		fmt.Printf("Found no car \n")
	}

	//条件查询
	//年龄大于30岁的
	boolQ := elastic.NewBoolQuery()
	boolQ.Must(elastic.NewMatchQuery("last_name", "smith"))
	boolQ.Filter(elastic.NewRangeQuery("age").Gt(30))
	res, err = conf.GetClient().Search("megacorp").Type("car").Query(q).Do(context.Background())
	printcar(res, err)

	//短语搜索 搜索about字段中有 rock climbing
	matchPhraseQuery := elastic.NewMatchPhraseQuery("about", "rock climbing")
	res, err = conf.GetClient().Search("car").Type("bigcar").Query(matchPhraseQuery).Do(context.Background())
	printcar(res, err)

	//分析 interests
	aggs := elastic.NewTermsAggregation().Field("interests")
	res, err = conf.GetClient().Search("megacorp").Type("car").Aggregation("all_interests", aggs).Do(context.Background())
	printcar(res, err)

}

//简单分页
func list(size,page int) {
	if size < 0 || page < 1 {
		fmt.Printf("param error")
		return
	}
	res,err := conf.GetClient().Search("megacorp").
		Type("car").
		Size(size).
		From((page-1)*size).
		Do(context.Background())
	printcar(res, err)

}

//打印查询到的car
func printcar(res *elastic.SearchResult, err error) {
	if err != nil {
		print(err.Error())
		return
	}
	var typ models.Company
	for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
		t := item.(models.Company)
		fmt.Printf("%#v\n", t)
	}
}
func printResult(typ interface{},res *elastic.SearchResult) {
	for _, item := range res.Each(reflect.TypeOf(typ)) { //从搜索结果中取数据的方法
		t := item
		fmt.Printf("%#v\n", t)
		structField, _ := reflect.TypeOf(typ).FieldByName("ID")
		fmt.Printf("%#v\n", structField)
	}
}

