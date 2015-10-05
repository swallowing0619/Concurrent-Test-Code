package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
	//gr "github.com/parnurzeal/gorequest"
	//ws "golang.org/x/net/websocket"
	log "github.com/Sirupsen/logrus"
	ws "github.com/gorilla/websocket"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

func init() {
	log.SetLevel(log.WarnLevel)
}

func randInt(min int, max int) int {
	if max-min <= 0 {
		return min
	}
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

func main() {
	var temp int
	fmt.Printf("Enter test websocket number:\n")
	fmt.Scanln(&temp)
	if temp > 0 {
		var wg sync.WaitGroup
		var index = 0
		for _, empCode := range EmpCodes {
			index = index + 1
			if index > temp {
				break
			}

			wg.Add(1)

			go func(code string) {
				defer wg.Done()

				cnt := randInt(1, 30)
				time.Sleep(time.Second * time.Duration(cnt))

				err := SimulateKeepOnline(code)
				if err != nil {
					log.Errorf("[%s] has error: %s", code, err.Error())
					return
				}
			}(empCode)
		}

		wg.Wait()
	} else {
		fmt.Printf("Please Enter init number!\n")
	}

	// ticket, err := ConnectorAuth("00397589")
	// if err != nil {
	// 	fmt.Printf("err is %v\n", err)
	// }
	// fmt.Printf("ticket is %v\n", ticket)
	// err := SimulateKeepOnline("00397589")
	// if err != nil {
	// 	fmt.Printf("err is ", err)
	// }
}

func SimulateKeepOnline(empCode string) error {
	log.Debugf("[%s] Begin...", empCode)
	defer log.Debugf("[%s] Exit...", empCode)
	ticket, err := ConnectorAuth(empCode)
	if err != nil {
		return err
	}
	// Extract the cookie for later use
	log.Debugf("[%s] Establishing websocket connection...", empCode)

	// Establish websocket connection
	var wsHeader = http.Header(make(map[string][]string))
	// for _, cookie := range cookies {
	// 	wsHeader.Add("Cookie", cookie.String())
	// }
	cookie := http.Cookie{
		Name:  "WS_TICKET",
		Value: ticket,
	}
	wsHeader.Add("Cookie", cookie.String())

	// wsConfig, err := ws.NewConfig("ws://ws.hotpu.cn/employee/ws", "http://ws.hotpu.cn")
	// if err != nil {
	// 	return err
	// }
	// wsConfig.Header = wsHeader

	// wsConn, err := ws.DialConfig(wsConfig)
	// if err != nil {
	// 	return err
	// }
	wsConn, resp, err := ws.DefaultDialer.Dial("ws://ws.hotpu.cn/employee/ws", wsHeader)

	if err != nil {
		return err
	}

	wsConn.SetPingHandler(nil)

	log.Debugf("response:%#v\n", resp)
	// log.Printf("[%s] IsClientConn: %v", empCode, wsConn.IsClientConn())
	// log.Printf("[%s] IsServerConn: %v", empCode, wsConn.IsServerConn())

	log.Debugf("[%s] LocalAddr: %+v", empCode, wsConn.LocalAddr())
	log.Debugf("[%s] RemoteAddr: %+v", empCode, wsConn.RemoteAddr())

	log.Debugf("[%s] Receiving upstream messages...", empCode)

	for {
		var message string

		if _, raw, err := wsConn.ReadMessage(); err != nil {
			return err
			continue
		} else {
			message = string(raw)
		}
		//log.Debugf("[%s] received: %s", empCode, message)
		fmt.Printf("[%s] received: %s\n", empCode, message)
	}
}

type AuthRequest struct {
	EnterpriseCode string `json:"ent_code""`
	EmployeeCode   string `json:"emp_code"`
	Password       string `json:"password"`
}

var client = &http.Client{}

func ConnectorAuth(empcode string) (ticket string, err error) {
	aReq := getAuthRequest(empcode)
	// todo: 超时支持
	wsAuthAddr := "http://ws.hotpu.cn/employee/ws/auth"
	rawJson, err := json.Marshal(aReq)
	if err != nil {
		return "", err
	}
	buf := bytes.NewBuffer(rawJson)
	resp, err := client.Post(wsAuthAddr, "application/json", buf)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("request connector auth err:%s\n", "响应为空")
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("request connector auth err:响应码%d\n", resp.StatusCode)
	}
	body := resp.Body
	if body == nil {
		return "", fmt.Errorf("request connector auth err:%s\n", "响应为空")
	}
	log.Debugf("ConnectorAuth success for: %v\n", empcode)
	defer resp.Body.Close()
	raw, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	ticket, err = parseConnectorResponse(raw)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(ticket) == "" {
		return "", fmt.Errorf("获取的ticket不能为空")
	}
	return ticket, nil
}
func getAuthRequest(empcode string) *AuthRequest {
	aReq := &AuthRequest{
		EnterpriseCode: "yto_dev",
		EmployeeCode:   empcode,
		Password:       "1",
	}
	return aReq
}

func parseConnectorResponse(raw []byte) (string, error) {
	cResp := new(Response)
	if err := json.Unmarshal(raw, cResp); err != nil {
		return "", err
	}
	if cResp.StatusCode == CodeOK {
		m := cResp.Data.(map[string]interface{})
		if obj, ok := m["ticket"]; ok {
			ticket := obj.(string)
			return ticket, nil
		} else {
			return "", fmt.Errorf("websocket认证响应体没有ticket")
		}
	} else {
		return "", fmt.Errorf(cResp.StatusMessage)
	}
}

type Response struct {
	StatusCode    StatusCode  `json:"status_code"`
	StatusMessage string      `json:"status_message"`
	Data          interface{} `json:"data"`
}

type StatusCode int64

func (s StatusCode) String() string {
	str, ok := nameMap[s]
	if ok {
		return str
	}
	return fmt.Sprintf("未知错误(%d)", s)
}

var nameMap = map[StatusCode]string{
	CodeOK:             "OK",
	CodeServerInternal: "服务内部错误",
	CodeBadEntCode:     "无效的机构代码",
	CodeBadEmp:         "无效的工号代码",
	CodeBadPassword:    "密码错误",
	CodeInvalidFormat:  "请求格式错误",
}

var (
	CodeOK             = StatusCode(0)
	CodeServerInternal = StatusCode(50010)
	CodeBadEntCode     = StatusCode(50020)
	CodeBadEmp         = StatusCode(50021)
	CodeBadPassword    = StatusCode(50022)
	CodeInvalidFormat  = StatusCode(50023)
)
