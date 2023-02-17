package test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,-"`
}

var token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzYzMDE5MjYsIm9yaWdfaWF0IjoxNjc2Mjk4MzI2LCJ1c2VySUQiOjEwfQ.lsnJdcUBpN_Io2nWWf7UrMHRprHcia_uOsNOqnpL5H0"

func TestStressRelation(t *testing.T) {
	t.Log("测试启动", t.Name())
	times := 10
	for i := 0; i < times; i++ {
		actionType := strconv.Itoa(i%2 + 1)
		relationAction(t, token, "11", actionType)
	}
	t.Log("pass")
}

func relationAction(t *testing.T, token string, toUserId string, actionType string) {

	url := fmt.Sprintf("http://127.0.0.1:8080/douyin/relation/action/?token=%v&to_user_id=%v&action_type=%v", token, toUserId, actionType)
	method := "POST"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("User-Agent", "Apifox/1.0.0 (https://www.apifox.cn)")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	resp := Response{}
	json.Unmarshal(body, &resp)
	t.Logf("输出: %v", resp)
	if resp.StatusCode != 0 {
		t.Logf("测试异常 to: %v action: %v msg: %v\n", toUserId, actionType, resp.StatusMsg)
		t.FailNow()
	}
}
