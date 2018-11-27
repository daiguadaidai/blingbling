package handle

import (
	"encoding/json"
	"fmt"
	"github.com/daiguadaidai/blingbling/meta_parser"
	"github.com/daiguadaidai/blingbling/parser"
	"github.com/gorilla/schema"
	"github.com/juju/errors"
	"io/ioutil"
	"net/http"
)

type MetaSqls struct {
	Sqls string
}

type ResponseData struct {
	Status  bool
	Message string
	Data    interface{}
}

func ReturnSuccess(w http.ResponseWriter, data interface{}) {
	d := new(ResponseData)
	d.Status = true
	d.Data = data

	jsonBytes, err := json.Marshal(d)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(jsonBytes))
}

func ReturnError(w http.ResponseWriter, err error) {
	d := new(ResponseData)
	d.Status = false
	d.Message = err.Error()

	jsonBytes, err := json.Marshal(d)
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}
	fmt.Fprintf(w, string(jsonBytes))
}

func SqlMetaHandle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")

	param, err := GetSqls(r)
	if err != nil {
		ReturnError(w, err)
		return
	}

	mis, err := GetMetaInfo(param)
	if err != nil {
		ReturnError(w, err)
		return
	}

	ReturnSuccess(w, mis)
	return
}

/* 通过http请求的内容获取相关自定义审核参数
Params:
	_request: 请求
*/
func GetSqls(_request *http.Request) (*MetaSqls, error) {
	switch _request.Method {
	case "POST":
		return GetSqlsByPost(_request)
	case "GET":
		return GetSqlsByGet(_request)
	default:
		errMSG := fmt.Sprintf("错误请求类型: %v. 值允许使用 GET/POST 请求", _request.Method)
		return nil, errors.New(errMSG)
	}

	// 什么每个逻辑都会返回自己值所以不会走到这里, 所以返回值都是 nil, nil
	return nil, nil
}

// 通过post方法获取自定义审核参数
func GetSqlsByPost(_request *http.Request) (*MetaSqls, error) {
	bodyBytes, _ := ioutil.ReadAll(_request.Body)
	ms := new(MetaSqls)
	if len(bodyBytes) == 0 { // 没有个参数都使用默认值
		return ms, nil
	}

	err := json.Unmarshal(bodyBytes, ms)
	if err != nil {
		errMSG := fmt.Sprintf("POST请求, 不能正确解析给予的值: %v", err)
		return nil, errors.New(errMSG)
	}

	return ms, nil
}

// 通过get方法获取自定义审核参数
func GetSqlsByGet(_request *http.Request) (*MetaSqls, error) {
	_request.ParseForm()
	ms := new(MetaSqls)

	decoder := schema.NewDecoder()
	err := decoder.Decode(ms, _request.URL.Query())
	if err != nil {
		errMSG := fmt.Sprintf("GET请求, 不能正确解析URL参数: %v", err)
		return nil, errors.New(errMSG)
	}

	return ms, nil
}

/* 开始审核 SQL 语句
Params:
    _requestParam: http传来的自定义参数
Return:
	int: 审核状态码
	string: 审核相关信息, 如果成功是成功信息, 如果失败是失败信息
*/
func GetMetaInfo(param *MetaSqls) ([]*meta_parser.MetaInfo, error) {
	// 循环每一个sql语句进行解析, 并且生成相关审核信息
	mis := make([]*meta_parser.MetaInfo, 0, 1)

	// 解析SQL
	sqlParser := parser.New()
	stmtNodes, err := sqlParser.Parse(param.Sqls, "", "")
	if err != nil {
		return nil, err
	}

	for _, stmtNode := range stmtNodes {
		metaParser := meta_parser.NewMetaParser(stmtNode)
		if metaParser == nil {
			continue
		}

		mi, err := metaParser.MetaParse()
		if err != nil {
			return nil, err
		}
		mis = append(mis, mi)
	}

	return mis, nil
}
