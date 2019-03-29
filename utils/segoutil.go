/*
@Time : 2019/3/27 19:16 
@Author : lukebryan
@File : segoutil
@Software: GoLand
*/
package utils

import "github.com/huichen/sego"

var segmenter sego.Segmenter

func init() {
	// 载入词典
	segmenter.LoadDictionary("./dictionary/dictionary.txt,./dictionary/ctb8_word.txt,./dictionary/dict.txt,./dictionary/msra_word.txt,./dictionary/weibo_word.txt")

}

/*
分词(会加入自身词段用于精确搜索)
 */
func Participle(text string) (output []string) {
	//加入自身词段
	output = append(output, text)

	// 分词
	segments := segmenter.Segment([]byte(text))

	// 处理分词结果
	// 支持普通模式和搜索模式两种分词，见代码中SegmentsToString函数的注释。
	output = append(output,sego.SegmentsToSlice(segments, false)...)
	return
}
