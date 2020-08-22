package php

import (
	"dj-api/app/config"
	"dj-api/tools/ssl"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

//返回格式
type Rsp struct {
	Code int                    `json:"code"`
	Msg  string                 `json:"msg"`
	Data map[string]interface{} `json:"data"`
}

//get请求
func HttpGet(path string, data map[string]interface{}) (map[string]interface{}, error) {
	token, err := createToken(data)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", config.C.Third.Url+path+"?token="+token, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json;charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("path:%v, req: %v, rsp err: %v", path, data, resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	rsp := &Rsp{}
	err = json.Unmarshal(body, rsp)
	if err != nil {
		return nil, err
	}
	if rsp.Code == 200 {
		return rsp.Data, nil
	}
	return nil, fmt.Errorf("path:%v, req: %v, rsp err: %v", path, data, rsp.Msg)
}

//生成jwt token
func createToken(data map[string]interface{}) (string, error) {
	dataJson, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	jwt := ssl.NewJwt(dataJson)
	jwt.SetSecret(config.C.Third.JwtKey)
	token, err := jwt.Enable()
	if err != nil {
		return "", err
	}
	return token, nil
}
