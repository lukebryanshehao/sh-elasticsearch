/*
@Time : 2019/3/27 15:18 
@Author : lukebryan
@File : participle_service	sego分词器
@Software: GoLand
*/
package services

import (
	"context"
	"elasticsearch/conf"
	"elasticsearch/models"
	"elasticsearch/utils"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"log"
	"reflect"
)

func NewSearchService() SearchService {
	return &searchService{}
}

type searchService struct {
}

type SearchService interface {
	SearchCompanys(text string) (companys []models.Company)
	SearchUsers(text string) (users []models.User)
}

func (s searchService) SearchCompanys(text string) (companys []models.Company) {

	queryStrings := utils.Participle(text)
	fmt.Println("------分词前------")
	fmt.Println(text)
	fmt.Println("-----分词后-------")
	fmt.Println(queryStrings)

	var resResults []*elastic.SearchResult

	for i := range queryStrings{
		matchPhraseQuery := elastic.NewMatchPhraseQuery("Abbreviation", queryStrings[i])
		res, err := conf.GetClient().Search("company").Type("company").Query(matchPhraseQuery).Do(context.Background())
		if err != nil {
			log.Print(err)
		}else {
			resResults = append(resResults, res)
		}
		matchPhraseQuery2 := elastic.NewMatchPhraseQuery("Name", queryStrings[i])
		res2, err := conf.GetClient().Search("company").Type("company").Query(matchPhraseQuery2).Do(context.Background())
		if err != nil {
			log.Print(err)
		}else {
			resResults = append(resResults, res2)
		}
	}

	for i := range resResults{
		for _, item := range resResults[i].Each(reflect.TypeOf(models.Company{})) { //从搜索结果中取数据的方法
			t := item.(models.Company)
			exist := IsExistInArray(t,companys)
			if !exist {
				companys = append(companys, t)
			}

		}
	}

	return
}


func (s searchService) SearchUsers(text string) (users []models.User) {
	return
}

func IsExistInArray(element interface{},array []models.Company) bool {
	for i := range array{
		if element == array[i] {
			return true
		}
	}
	return false
}