/*
@Time : 2019/3/25 16:00 
@Author : lukebryan
@File : main
@Software: GoLand
*/
package main

import (
	"elasticsearch/utils"
	"fmt"
	"github.com/spf13/cast"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"strings"
)

type conf struct {
	Tables map[string]interface{} `yaml:"Table"`
	ModelPath string `yaml:"ModelPath"`
}

var Conf = &conf{}

func main() {
	yamlFile, err := ioutil.ReadFile("table.yaml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(yamlFile, Conf)
	if err != nil {
		panic(err)
	}

	if !strings.HasSuffix(Conf.ModelPath,"/") {
		Conf.ModelPath += "/"
	}

	tableMaps := Conf.Tables
	for k,v := range tableMaps{

		field := v.(map[interface {}]interface {})["field"]
		imports := v.(map[interface {}]interface {})["import"]
		packag := v.(map[interface {}]interface {})["package"]

		goFileName := strings.ToLower(k)+".go"
		content := "package "+fmt.Sprintf("%v",packag)+"\r\n\r\n"

		if fmt.Sprintf("%v",packag) != "" {
			content += "import (\r\n"
			for _,im := range imports.([]interface {}){
				content += "	\"" + fmt.Sprintf("%v",im) + "\"\r\n"
			}
			content += ")\r\n\r\n"
		}
		content += "type "+k+" struct {\r\n"
		for _,fie := range field.([]interface {}){
			content += "	" + fmt.Sprintf("%v",fie) + "\r\n"
		}

		content += "}"
		fmt.Println("-------------")
		fmt.Println(goFileName)
		fmt.Println(content)
		fmt.Println("-------------")
		
		go BuildGoFile(Conf.ModelPath + goFileName,content)
	}

}

func BuildGoFile(file string,content string)  {
	path := file[0:strings.LastIndex(file,"/")]
	if exist,_ := utils.PathExists(path);!exist {
		utils.MkDirs(path)
	}
	go utils.CreateFile(file,content)
}