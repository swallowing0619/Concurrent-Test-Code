package main

import (
	"encoding/json"
	"fmt"
	// "time"
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	// "unicode/utf8"

	"github.com/parnurzeal/gorequest"
	// "github.com/revel/revel"
	o "yto.net.cn/kefu/models/objects"
	// . "yto.net.cn/kefu/official_app/app/models/objects"
)

type OpenapiRequest struct {
	AppID       string //服务窗账号ID
	AppSecret   string
	Endpoint    string
	AccessToken string
}

type MessageStatus struct {
	Category           o.MessageCategory
	ServiceEmpCode     string
	SubscribeSessionID int64
	ServiceSessionID   int64
}

const (
	requestRetryCount = 5
	// timeout           = time.Second * time.Duration(10)
)

func getRetry(url, req string) (resp gorequest.Response, body string, errs []error) {
	for i := 0; i < requestRetryCount; i++ {
		resp, body, errs = gorequest.New().Get(url).
			Query(req).End()

		if errs == nil && resp.StatusCode == 200 {
			break
		}
	}

	return
}

//获取平台口令
func (r *OpenapiRequest) fetchToken() (string, error) {
	// revel.INFO.Println("OpenapiRequest.fetchToken run...")
	r.Endpoint = fmt.Sprintf("http://222.73.41.159:9001/api/v1/session?access_token=%s", r.AccessToken)
	url := fmt.Sprintf("%s/access", r.Endpoint)
	req := fmt.Sprintf("appid=%s&appsecret=%s", r.AppID, r.AppSecret)

	resp, body, errs := getRetry(url, req)
	if errs != nil {
		return "", errs[0]
	}

	// revel.INFO.Printf("response when requesting openapi: %#v\n", resp)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("OpenAPI has not started")
	}

	data := o.ApiResponseAccess{}
	err := json.Unmarshal([]byte(body), &data)
	if err != nil {
		// revel.ERROR.Printf("Unable to unserialize the token result(%s).", err.Error())
		return "", err
	}
	// revel.INFO.Printf("The token data is: %#v\n", data)
	fmt.Println("The token data is: %#v\n", data)
	return data.Data.AccessToken, nil
}

/*
获取平台token
*/
func getToken(platform string) string {
	switch platform {
	case "wechat":
		return "6b6487e6-4646-11e5-a8b4-00262d0d05c7"
	case "alipay":
		return "cc2f2b03-2e7c-11e5-a8b4-00262d0d05c7"
	case "qq":
		return "8d1600d6-f7df-11e4-934b-00262d0d05c7"
	case "site":
		return "ea7c1330-4d5b-11e5-a8b4-00262d0d05c7"
	default:
		return ""
	}
}

//构造openapi请求
func NewOpenapiRequest(appid, appsecret, endpoint string, token string) (*OpenapiRequest, error) {
	// if appid == "" || appsecret == "" || endpoint == "" || token == "" {
	// 	return nil, fmt.Errorf("Param error, appid, appsecret, entpoint and token are required\n")
	// }

	r := &OpenapiRequest{
		AppID:       appid,
		AppSecret:   appsecret,
		Endpoint:    endpoint,
		AccessToken: token,
	}

	return r, nil
}

/*
查询已存在的会话
*/
func (r *OpenapiRequest) SessionQuery(platform o.Platform, openid string) (*o.ApiResponseSessionQuery, error) {
	// revel.INFO.Println("OpenapiRequest.SessionQuery run...")

	// param := fmt.Sprintf(`access_token=%s`, r.AccessToken)

	req := o.ApiRequest{
		Action: "session.query",
		Data: o.ApiSessionQueryData{
			Platform: platform,
			OpenID:   openid,
		},
	}
	//api接口
	r.Endpoint = fmt.Sprintf("http://222.73.41.159:9001/api/v1/session?access_token=%s", r.AccessToken)

	// revel.INFO.Printf("request when session query: %#v\n", req)

	// resp, body, errs := postRetry(r.Endpoint, param, req)
	//转为json
	raw, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	//转为byte
	bf := bytes.NewBuffer(raw)
	//post请求
	resp, err := http.Post(r.Endpoint, "application/json", bf)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}
	bRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// revel.INFO.Printf("response when requesting openapi: %#v\n", resp)
	// if resp.StatusCode != 200 {
	// 	return nil, fmt.Errorf("OpenAPI has not started")
	// }

	// revel.INFO.Printf("The response body is: %s\n", body)

	result := &o.ApiResponseSessionQuery{}
	if err := json.Unmarshal(bRaw, result); err != nil {
		// revel.ERROR.Printf("Unable to unmarshal the result fron session.query(%s)\n", err.Error())
		return nil, err
	}

	// r.AccessToken = accessToken

	return result, nil
}

/*
发送消息请求
*/
func (r *OpenapiRequest) MessageRequest(platform o.Platform, openID string, content string, msgstatus *MessageStatus) (*o.ApiResponseMessageRequest, error) {
	var messageHolder interface{}

	messageHolder = &o.ApiMessageTextContainer{
		MsgType: "text",
		MsgDetail: o.ApiTextMessage{
			Content: content,
		},
	}

	r.Endpoint = fmt.Sprintf("http://222.73.41.159:9001/api/v1/message?access_token=%s", r.AccessToken)
	// API request
	req := o.ApiRequest{
		Action: "message.request",
		Data: &o.ApiMessageRequestData{
			Platform:           platform,
			OpenID:             openID,
			Offline:            false,
			Message:            messageHolder,
			Category:           msgstatus.Category,
			ServiceEmpCode:     msgstatus.ServiceEmpCode,
			SubscribeSessionID: msgstatus.SubscribeSessionID,
			ServiceSessionID:   msgstatus.ServiceSessionID,
		},
	}

	//转为json
	raw, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	bf := bytes.NewBuffer(raw)
	resp, err := http.Post(r.Endpoint, "application/json", bf)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("response is nil")
	}
	bRaw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// parse the response
	result := new(o.ApiResponseMessageRequest)
	if err := json.Unmarshal(bRaw, result); err != nil {
		return nil, err
	}

	return result, nil
}

/*
查询会话后发送消息请求
*/
func request(req *OpenapiRequest, platform o.Platform, openid string, content string) {
	resp, err := req.SessionQuery(platform, openid)

	if err != nil {
		fmt.Printf("session query err:%s\n", err.Error())
		os.Exit(0)
	}

	msgstatus := &MessageStatus{
		Category:           o.MessageCategoryManuReq,
		ServiceEmpCode:     resp.Data.ServiceEmpCode,
		SubscribeSessionID: resp.Data.SubscriberSessionID,
		ServiceSessionID:   resp.Data.ServiceSessionID,
	}
	respMsgReq, err := req.MessageRequest(platform, openid, content, msgstatus)
	if err != nil {
		fmt.Printf("message request err:%s\n", err.Error())
		os.Exit(0)
	}
	fmt.Printf("message request response:%#v\n", respMsgReq)
}
