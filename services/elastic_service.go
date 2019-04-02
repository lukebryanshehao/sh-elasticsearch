/*
@Time : 2019/3/26 10:11 
@Author : lukebryan
@File : elastic_service
@Software: GoLand
*/
package services

import (
	"elasticsearch/conf"
	"elasticsearch/utils"
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"log"
	"sync"
)

func NewElasticService() ElasticService {
	return &elasticService{}
}

type elasticService struct {
}

var CreateMap = make(map[string]string)

type ElasticService interface {
	Create(table string) (map[int]map[string]interface{},bool)
	Update() bool
	Delete() bool
}

func (e elasticService) Create(table string) (map[int]map[string]interface{},bool) {
	rows,_ := conf.GetDB().Table(table).Rows()
	//返回所有列
	columns,_ := rows.Columns()
	columnsTypes,_ := rows.ColumnTypes()
	columnsTypeMaps := make(map[string]string, len(columnsTypes))
	for i := range columnsTypes{
		fmt.Println("columnsTypes:  ",columnsTypes[i].Name(),columnsTypes[i].DatabaseTypeName())
		columnsTypeMaps[columnsTypes[i].Name()] = columnsTypes[i].DatabaseTypeName()
	}

	//这里表示一行所有列的值，用[]byte表示
	vals := make([][]byte, len(columns))
	//这里表示一行填充数据
	scans := make([]interface{}, len(columns))
	//这里scans引用vals，把数据填充到[]byte里
	for k := range vals {
		scans[k] = &vals[k]
	}

	i := 0
	result := make(map[int]map[string]interface{})
	for rows.Next() {
		//填充数据
		err := rows.Scan(scans...)
		if err != nil {
			log.Println(err)
		}
		//每行数据
		row := make(map[string]interface{})
		//把vals中的数据复制到row中
		for k, v := range vals {
			key := columns[k]
			//这里把[]byte数据转成string
			//row[key] = string(v)
			dataType := columnsTypeMaps[key]
			//"VARCHAR", "TEXT", "NVARCHAR", "DECIMAL", "BOOL", "INT", "BIGINT"
			switch dataType {

				case "VARCHAR":
					row[key] = string(v)
				case "TEXT":
					row[key] = string(v)
				case "NVARCHAR":
					row[key] = string(v)
				case "DECIMAL":
					row[key] = cast.ToString(string(v))
				case "BOOL":
					row[key] = cast.ToBool(string(v))
				case "INT":
					row[key] = cast.ToInt(string(v))
				case "BIGINT":
					row[key] = cast.ToUint(string(v))
				case "TIMESTAMP":
					row[key] = string(v)
				default:
					row[key] = string(v)
			}

		}
		//放入结果集
		result[i] = row
		i++
	}

	var wg sync.WaitGroup
	wg.Add(len(result))
	for k := range result {
		go func(maps map[string]interface{}) {
			defer wg.Done()
			byt,_ := json.Marshal(maps)
			//fmt.Println(string(byt))
			utils.Create(string(byt), table, table, cast.ToString(maps["id"]),0)
		}(result[k])
	}

	wg.Wait()

	//now := time.Now()
	//e1 := models.Company{uint(188),now,now,&now, "思倍捷", "湖南思倍捷", 4, 0,1,""}
	//utils.Create(e1, "company", "company", string(e1.ID))
	return result,true
}

func (e elasticService) Update() bool {
	maps := make(map[string]interface{})
	maps["Name"] = "思倍捷"
	return utils.Update(maps, "company", "company", "158")
}

func (e elasticService) Delete() bool {
	return utils.Delete("company", "company", "158")
}
