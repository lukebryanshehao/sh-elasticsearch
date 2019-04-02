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
	"encoding/json"
	"fmt"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func NewSearchService() SearchService {
	return &searchService{}
}



type searchService struct {
}

type SearchService interface {
	Search(text string,searchType string,searchFields []string) (resultList models.ResultList)
	SearchUsers(text string) (users []models.User)
}

func (s searchService) Search(text string,searchType string,searchFields []string) (resultList models.ResultList){

	queryStrings := utils.Participle(text)
	fmt.Println("------分词前------")
	fmt.Println(text)
	fmt.Println("-----分词后-------")
	fmt.Println(queryStrings)

	var resResults []*elastic.SearchResult

	for i := range queryStrings{
		for j := range searchFields{
			matchPhraseQuery := elastic.NewMatchPhraseQuery(searchFields[j], queryStrings[i])
			res, err := conf.GetClient().Search(searchType).Type(searchType).Query(matchPhraseQuery).Do(context.Background())
			if err != nil {
				log.Print("Bad Error :",err)
			}else {
				resResults = append(resResults, res)
			}
		}
	}
	var companyList []interface{}
	var userList []interface{}

	for i := range resResults{
		for _, hit := range resResults[i].Hits.Hits {
			// hit.Index contains the name of the index
			switch searchType {
				case "company":
					// Deserialize hit.Source into a LogOrder (could also be just a map[string]interface{}).

					var t models.Company
					err := json.Unmarshal(*hit.Source, &t)
					if err != nil {
						log.Println("Unmarshal Error!")
					}
					exist := IsExistInArray(t,companyList)
					if !exist {
						companyList = append(companyList, t)
					}
				case "user":
					var t models.User
					err := json.Unmarshal(*hit.Source, &t)
					if err != nil {
						log.Println("Unmarshal Error!")
					}
					exist := IsExistInArray(t,userList)
					if !exist {
						userList = append(userList, t)
					}
				default:
					log.Println("没有匹配到类型")
			}
		}

		//for _, item := range resResults[i].Each(reflect.TypeOf(models.Company{})) { //从搜索结果中取数据的方法
		//	t := item.(models.Company)
		//	fmt.Println(t)
		//	exist := IsExistInArray(t,companyList)
		//	if !exist {
		//		companyList = append(companyList, t)
		//	}
		//}
	}
	resultList.UserList = userList
	resultList.CompanyList = companyList
	return
}


func (s searchService) SearchUsers(text string) (users []models.User) {
	return
}

func IsExistInArray(element interface{},array []interface{}) bool {
	for i := range array{
		if element == array[i] {
			return true
		}
	}
	return false
}