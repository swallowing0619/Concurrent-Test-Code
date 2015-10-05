package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Params struct {
	ToUser     string
	Openid     string
	CreateTime time.Time
	MsgType    string
	Message    string
}

type Request struct {
	ID        int
	Params    Params
	Status    int       // 回复状态
	StartTime time.Time // 开始时间
	EndTime   time.Time // 结束时间
	Error     string    // 错误消息
	Count     int       // 重试次数
	Content   string    // 回复内容
}

type returnMessage struct {
	ToUserName   string
	FromUserName string
	CreateTime   string
	MsgType      string
	Content      string
	MsgId        string
}

//模拟post请求

var client = &http.Client{}

func (request *Request) Send() (status int, content string, err error) {
	params := request.Params
	sendMessage := getSendMessage(params.ToUser, params.Openid, params.CreateTime, params.MsgType, params.Message)

	// todo: 超时支持
	resp, err := client.Post("http://58.32.246.66/wx", "text/xml", bytes.NewBufferString(sendMessage))

	if err != nil {
		return 0, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		var Result returnMessage
		err := xml.Unmarshal(body, &Result)
		if err != nil {
			return resp.StatusCode, "", err
		} else {
			return resp.StatusCode, "", nil // todo: 内容呢？
		}
	} else {
		return resp.StatusCode, "", nil
	}
}

var maxretryNumber int = 3

func makeRetryRequest(request Request) (Request, error) {
	if request.Count >= maxretryNumber {
		return Request{}, fmt.Errorf("retry 3 time")
	}

	request.StartTime = time.Now()
	request.Count++
	return request, nil
}

var id int = 1

func makeRequest(req Request, content string, info CustomerInfo) (Request, error) {
	request := Request{
		Status:    0,
		StartTime: time.Now(),
		Error:     "",
		Count:     1,
	}
	if req.ID == 0 {
		request.ID = id
		id++
	} else {
		request.ID = req.ID
	}

	param := Params{
		ToUser:     info.Name,
		Openid:     info.Openid,
		CreateTime: request.StartTime,
		MsgType:    "text",
		Message:    content,
	}
	request.Params = param

	return request, nil
}

func getSendMessage(toUser string, openid string, createTime time.Time, msgType string, message string) string {
	return fmt.Sprintf(`<xml>
      <ToUserName><![CDATA[%s]]></ToUserName>
      <FromUserName><![CDATA[%s]]></FromUserName>
      <CreateTime>%d</CreateTime>
      <MsgType><![CDATA[%s]]></MsgType>
      <Content><![CDATA[%s]]></Content>
    </xml>`, toUser, openid, createTime.Unix(), msgType, fmt.Sprintf("%s", message))
}
