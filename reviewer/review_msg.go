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
	Code       int
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
	rs := fmt.Sprintf(`{"Code": %v, "ReviewMSGs": [%v]}`,
		this.Code, _err, "%v")
	reviewMSGStr := make([]string, 0, 1)
	for _, reviewMSG := range this.ReviewMSGs {
		reviewMSGStr = append(reviewMSGStr, reviewMSG.String())
	}

	return fmt.Sprintf(rs, strings.Join(reviewMSGStr, ","))
}

// 设置返回代码
func (this *ResponseReviewData) ResetCode() {
	for _, reviewMSG := range this.ReviewMSGs {
		if reviewMSG.HaveError {
			this.Code = REVIEW_CODE_ERROR
			return
		} else if reviewMSG.HaveWarning {
			this.Code = REVIEW_CODE_WARNING
			continue
		}
	}
}

type ReviewMSG struct {
	Sql         string
	HaveError   bool
	HaveWarning bool
	ErrorMSGs   []string
	WarningMSGs []string
}

func NewReivewMSG() *ReviewMSG {
	return &ReviewMSG{
		ErrorMSGs:   make([]string, 0, 1),
		WarningMSGs: make([]string, 0, 1),
	}
}

// 重新设置是否有错误和警告
func (this *ReviewMSG) ResetHaveErrorAndWarning() {
	if len(this.ErrorMSGs) > 0 {
		this.HaveError = true
	}
	if len(this.WarningMSGs) > 0 {
		this.HaveWarning = true
	}
}

/* 添加信息, 如果有错误则信息是错误信息, 如果没有错误且有信息, 则是警告信息
Params:
    _haveError: 是否有错误
    _msg: 相关信息
 */
func (this *ReviewMSG) AppendMSG(_haveError bool, _msg string) (haveMSG bool) {
	if _haveError {
		this.ErrorMSGs = append(this.ErrorMSGs, _msg)
	} else if _msg != "" {
		haveMSG = true
		this.WarningMSGs = append(this.WarningMSGs, _msg)
	}

	return
}

func (this *ReviewMSG) String() string {
	jsonBytes, err := json.Marshal(this)
	if err != nil {
		return fmt.Sprintf(`{"Sql": %v, "HaveError": true, "HaveWarning": false, "ErrorMSG": [%v], "WarningMSG": []}`,
			this.Sql,
			fmt.Sprintf("审核信息转化成json时出错 %v", err))
	}

	return string(jsonBytes)
}
