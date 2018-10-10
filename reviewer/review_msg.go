package reviewer

import (
	"encoding/json"
	"fmt"
	"strings"
)

// 定义返回代码
const (
	REVIEW_CODE_SUCCESS = iota
	REVIEW_CODE_WARNING
	REVIEW_CODE_ERROR
)

type ResponseReviewData struct {
	Code int
	MSG string
	ReviewMSGs []*ReviewMSG
}

// 将信息以json 字符串返回
func (this *ResponseReviewData) ToJson() string {
	responseReviewBytes, err := json.Marshal(this)
	if err != nil {
		return this.GetErrorJson(err)
	}

	return string(responseReviewBytes)
}

func (this *ResponseReviewData) GetErrorJson(_err error) string {
	rs := fmt.Sprintf("{Code: %v, MSG: %v, ReviewMSG: [%v]}",
		this.Code, _err, "%v")
	reviewMSGStr := make([]string, 0, 1)
	for _, reviewMSG := range this.ReviewMSGs {
		reviewMSGStr = append(reviewMSGStr, reviewMSG.String())
	}

	return fmt.Sprintf(rs, strings.Join(reviewMSGStr, ","))
}

type ReviewMSG struct {
	Sql string
	Code int
	MSG string
}

func (this *ReviewMSG) String() string {
	rs := fmt.Sprintf("{Code: %v, MSG: %v, Sql: %v}",
		this.Code, this.MSG, this.Sql)

	return rs
}


