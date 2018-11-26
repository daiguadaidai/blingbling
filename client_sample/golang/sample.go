package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/daiguadaidai/blingbling/reviewer"
	"github.com/liudng/godump"
)

func main() {
	// 设置需要传入的参数
	params := make(map[string]interface{})
	params["Host"] = "10.10.10.21"
	params["Port"] = 3307
	params["Username"] = "root"
	params["Password"] = "root"
	params["Database"] = "employees"
	params["Sqls"] = "alter table employees add column age1 int not null; delete from employees WHERE id = 1;"

	// 变成Json格式
	jsonParams, err := json.Marshal(params)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 设置POST请求相关参数
	reader := bytes.NewReader(jsonParams)
	url := "http://10.10.10.55:18080/sqlReview"
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")

	// 请求
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	respData := new(reviewer.ResponseReviewData)
	err = json.NewDecoder(resp.Body).Decode(respData)
	if err != nil {
		fmt.Printf("json decode err: %v", err)
		return
	}

	godump.Dump(respData)
}
