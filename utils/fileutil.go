/*
@Time : 2019/3/25 16:10 
@Author : lukebryan
@File : fileutil
@Software: GoLand
*/
package utils

import (
	"fmt"
	"io"
	"log"
	"os"
)

/*
判断文件/文件夹是否存在
path	文件/文件夹路径
 */
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/*
复制文件
dstName		目标文件
srcName		源文件
 */
func CopyFile(dstName, srcName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return
	}
	defer dst.Close()
	return io.Copy(dst, src)
}

/*
写入文件
file	文件		C:\file\data\test.text
content	内容		Hello World!\r\n
 */
func CreateFile(file string,content string) bool {
	//写入文件
	f, error := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0766)
	if error != nil {
		log.Print(error)
		return false
	}
	defer f.Close()
	_,err := f.Write([]byte(content))
	if err != nil {
		log.Print(err)
		return false
	}

	return true
}

/*
创建文件夹
 */
func MkDirs(dir string) bool {
	exist, _ := PathExists(dir)

	if exist {
		fmt.Printf("has dir![%v]\n", dir)
		return true
	} else {
		fmt.Printf("no dir![%v]\n", dir)
		// 创建文件夹
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			fmt.Printf("mkdir failed![%v]\n", err)
			return false
		} else {
			fmt.Printf("mkdir success!\n")
			return true
		}
	}
}