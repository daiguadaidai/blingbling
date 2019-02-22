package common

import (
	"strings"

	"fmt"
	"github.com/dlclark/regexp2"
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

// 通过时间字符串转化成格式化字符串 1970-01-01 08:00:01.123 -> 2006-01-02 15:04:05.0000
func TimeFormatParse(timeStr string) (string, error) {
	items := strings.Split(timeStr, ".")
	if len(items) == 1 {
		return "2006-01-02 15:04:05", nil
	} else if len(items) == 2 {
		if len(items[1]) < 1 || len(items[1]) > 6 {
			return "", fmt.Errorf("%s 格式不能正常解析(小数点后面的取值个数范围为(1到6个)), "+
				"正确时间格式为 2006-01-02 15:04:05 / 2006-01-02 15:04:05.000000", timeStr)
		}
		dotStr := strings.Repeat("0", len(items[1]))
		return fmt.Sprintf("2006-01-02 15:04:05.%s", dotStr), nil
	}
	return "", fmt.Errorf("%s 格式不能正常解析, "+
		"正确时间格式为 2006-01-02 15:04:05 / 2006-01-02 15:04:05.000000", timeStr)
}

// interface转化为字符串
func InterfaceToStr(value interface{}) string {
	return fmt.Sprintf("%v", value)
}
