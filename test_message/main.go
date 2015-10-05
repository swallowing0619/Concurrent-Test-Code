package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

func info(v ...interface{}) {
	log.Println("[INFO]", v)
}

var c chan *Customer
var d chan Request

func main() {
	/*
	 * configuarion
	 */
	var temp int
	fmt.Printf("Enter send message number:\n")
	fmt.Scanln(&temp)

	if temp > 0 {
		// 缓冲数
		maxChanNumber := 1000

		// 并发数
		maxConcurrencyNumber := 8

		c = make(chan *Customer, maxChanNumber)
		d = make(chan Request, maxChanNumber)

		// 添加客户到channel
		go func() {
			j := 0
			for {
				for _, customerInfo := range getCustomerInfoList() {
					j++
					if request, err := makeRequest(Request{}, fmt.Sprintf("hello world客户, %d", j), customerInfo); err == nil {
						if j <= temp {
							d <- request
						} else {
							time.Sleep(85600 * time.Second)
						}
					} else {
						info(j, err)
					}
				}
			}
		}()

		// 消费客户
		var j int = 0
		for i := 0; i < maxConcurrencyNumber; i++ {
			go func() {
				for {
					request := <-d
					startTime := time.Now()
					status, content, err2 := request.Send()

					endTime := time.Now()
					if err2 != nil {
						request.Error = err2.Error()
					} else {
						request.Status = status
						request.Content = content
					}

					info("idx: ", j, ", 次数: ", request.Count, ", status: ", status, ", content: ", content, ", 耗时: ", endTime.Unix()-startTime.Unix(), ", error:", err2)
					j++
				}
			}()
		}

		for {
			time.Sleep(1 * time.Second)
		}
		/*
		   http.HandleFunc("/", URLAPI)
		   http.ListenAndServe(":8099", nil)
		*/
	} else {
		fmt.Printf("Please Enter int number!\n")
	}
}

// 提供http访问 API
func URLAPI(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	//read xml
	body, _ := ioutil.ReadAll(r.Body)
	var result returnMessage
	err := xml.Unmarshal(body, &result)
	if err == nil {
		platform := "0"
		openid := result.FromUserName
		content := result.Content
		if customer, err := getCustomer(&CustomerInfo{Platform: platform, Openid: openid}); err == nil {
			customer.Responses = append(customer.Responses, content)

			length := len(customer.Requests) - 1

			// todo:
			sendMessage := ""
			if content == "小圆太笨啦，无法识别你的消息。有单咨询请直接输入单号，无单咨询请输入1，退出请输入0" {
				sendMessage = "1" // todo
			} else if content == "您还没有评价，您的评价是我们进步的动力哦。请输入1,2,3,4,5进行评价" {
				sendMessage = "1" // todo
			} else if content == "拜托拜托，小圆生病了，正在玩命儿的医治。暂时不能进行人工服务，请先试用自助服务吧！" {
				sendMessage = "1"
			} else if matched, _ := regexp.MatchString("^Sorry", content); matched {
				sendMessage = "1"
			} else if content == "非正确的单号,退出输入0" {
				sendMessage = "1"
			} else if content == "会话结束" {
				// todo
			} else {
				sendMessage = fmt.Sprintf("%s, %d", content, time.Now().Unix())
			}

			/* "您已放弃召唤客服MM,现在您可以使用自助服务。" */
			if sendMessage != "" {
				if request, err := makeRequest(customer.Requests[length], sendMessage, customer.Info); err == nil {
					customer.Requests = append(customer.Requests, request)
					c <- customer
				}
			}
		}
	}
}
