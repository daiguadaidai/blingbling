package common

import (
	"github.com/dlclark/regexp2"
	"strings"
)

/* 字符串匹配
Params:
    _str: 需要匹配的字符串
    _reg: 正则表达式
 */
func StrIsMatch(_str string, _reg string) bool {
	var matched bool = false

	re := regexp2.MustCompile(_reg, 0)
	if isMatch, _ := re.MatchString(_str); isMatch {
		matched = true
	}

	return matched
}

/* 通过切割符号切割字符串, 并保存在 map中
Params:
    _str: 需要分割的字符串
    _sep: 通过什么进行分割
 */
func SplitString2Map(_str string, _sep string) map[string]bool {
	itemMap := make(map[string]bool)

	items := strings.Split(_str, _sep)
	for _, item := range items {
		item = strings.ToLower(strings.TrimSpace(item))
		if item == "" {
			continue
		}

		itemMap[item] = true
	}

	return itemMap
}
