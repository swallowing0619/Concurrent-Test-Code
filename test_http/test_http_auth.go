package main

import (
	"bytes"
	//"encoding/xml"
	"fmt"
	//"io/ioutil"
	"encoding/json"
	//"io/ioutil"
	"net/http"
	//"time"
	"sync"
)

type AuthRequest struct {
	EnterpriseCode string `json:"ent_code""` //wangdian daima
	EmployeeCode   string `json:"emp_code"`
	Password       string `json:"password"`
}

var client = &http.Client{} //kehuduan duixiang

func ConnectorAuth(empcode string, index int) (err error) {
	aReq := getAuthRequest(empcode) //kefu info
	// todo: 超时支持
	wsAuthAddr := "http://ws.hotpu.cn/employee/ws/auth" //address
	rawJson, err := json.Marshal(aReq)                  //struct turn into json
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(rawJson)                               //json turn into bytes
	resp, err := client.Post(wsAuthAddr, "application/json", buf) //send post requst
	if err != nil {
		return err
	}
	if resp == nil {
		return fmt.Errorf("request connector auth err:%s\n", "响应为空")
	}

	if resp.StatusCode != http.StatusOK { //statusCode
		return fmt.Errorf("request connector auth err:响应码%d\n", resp.StatusCode)
	}
	body := resp.Body //xiangying xiaoxiti
	if body == nil {
		return fmt.Errorf("request connector auth err:%s\n", "响应为空")
	}
	fmt.Printf("%v:ConnectorAuth success for: %v\n", index, empcode)
	defer resp.Body.Close() //hanshu return zhihou cai zhixing
	return nil
}
func getAuthRequest(empcode string) *AuthRequest {
	aReq := &AuthRequest{
		EnterpriseCode: "yto_dev",
		EmployeeCode:   empcode,
		Password:       "1",
	}
	return aReq
}
func main() {
	// for i, empCode := range EmpCodes {
	// 	err := ConnectorAuth(empCode)
	// 	if err != nil {
	// 		fmt.Printf("error", err)
	// 	}
	// }
	var temp int
	fmt.Printf("Enter test http number:\n")
	fmt.Scanln(&temp)
	if temp > 0 {
		var wg sync.WaitGroup
		var index = 0
		for _, empCode := range EmpCodes {
			index = index + 1
			if index > temp {
				break
			}
			fmt.Printf("index:%v\n", index)
			wg.Add(1)
			go func(code string) {
				defer wg.Done()

				err := ConnectorAuth(code, index)
				if err != nil {
					fmt.Printf("%v has error:%v\n", code, err)
					return
				}
			}(empCode)
		}

		wg.Wait()
	} else {
		fmt.Printf("Please Enter int number!\n")
	}

}
