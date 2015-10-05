package main

import (
	"fmt"
	"os"
	// "runtime"
	"sync"
	// "time"
	o "yto.net.cn/kefu/models/objects"
)

// //打印日志
// func info(v ...interface{}) {
// 	log.Println("[INFO]", v)
// }

// var c chan *Customer
// var d chan Request

func main() {

	var message_num int
	fmt.Printf("Enter send message number:\n")
	fmt.Scanln(&message_num)

	if message_num <= 0 {
		fmt.Printf("Please Enter init number!\n")
		os.Exit(0)
	}

	// 并发请求,遍历客户列表
	var wg sync.WaitGroup
	// channal := make(chan int)
	var j int
	// runtime.GOMAXPROCS(2) // 最多使用2个核
	for _, customerInfo := range getCustomerInfoList() {
		token := getToken(customerInfo.Platform)
		_platform := o.Platform(customerInfo.Platform)
		//构造openapi的请求
		req, _ := NewOpenapiRequest("", "", "", token)
		fmt.Printf("token=%s \n", token)
		//开启线程,发送消息
		j++
		wg.Add(1)
		go func() {
			fmt.Printf("go %d", j)

			for i := 0; i < message_num; i++ {
				fmt.Printf("send_message %d", i+1)
				content := fmt.Sprintf("it's a test! %d", i+1)
				//查询会话后发送消息请求
				request(req, _platform, customerInfo.Openid, content)
				// time.Sleep(1 * time.Second)
				// time.Sleep(1 * time.Second) // 停顿一秒
			}
			wg.Done()
			// channal <- 0
		}()
		wg.Wait()
	}

	// for n := 0; n <= j; n++ { // 等待所有消息发送完毕。
	// 	<-channal
	// }

}
