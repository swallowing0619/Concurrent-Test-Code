package main

import (
	// "fmt"
	// r "github.com/garyburd/redigo/redis"
	// gr "github.com/parnurzeal/gorequest"
	// "strconv"
	o "yto.net.cn/kefu/models/objects"
)

var (
	OpenUrl       string = "http://222.73.41.159:9001"
	EntCode       string = "yto_dev"
	SubScriberSet        = "subscriber_set"
)

func main() {
	//InitTestData()
	err, t := NewSuite(o.PlatformWechat)
	if err != nil {
		panic(err)
	}

	subs, err := LoadTestSub(1002, 1003)
	if err != nil {
		panic(err)
	}

	for _, sub := range subs {
		if sub != nil {
			t.TestSessionCreate(sub)
		}
	}
}

/*
	alipay:cc2f2b03-2e7c-11e5-a8b4-00262d0d05c7
	wechat:5c512919-3a58-11e5-a8b4-00262d0d05c7
	qq:8d1600d6-f7df-11e4-934b-00262d0d05c7
	site:b7186778-3b36-11e5-a8b4-00262d0d05c7
*/
