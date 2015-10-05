package main

import (
	"fmt"
	r "github.com/garyburd/redigo/redis"
	gr "github.com/parnurzeal/gorequest"
	"strconv"
	"time"
	o "yto.net.cn/kefu/models/objects"
)

func AddSubList(subids []int64) error {
	conn := Mo.Redis.GetConn()
	defer conn.Close()

	conn.Send("MULTI")
	for i := 0; i < len(subids); i++ {
		conn.Send("SADD", SubScriberSet, fmt.Sprintf("%d", subids[i]))
	}
	_, err := conn.Do("EXEC")
	return err
}

func GetSubList() ([]string, error) {
	conn := Mo.Redis.GetConn()
	defer conn.Close()

	list, err := r.Strings(conn.Do("SMEMBERS", SubScriberSet))
	return list, err
}

//456u
func LoadTestSub(begin, limit int64) ([]*o.Subscriber, error) {
	list, err := GetSubList()
	subs := make([]*o.Subscriber, 0)

	if err != nil {
		log.Error("Error Fetching Test Sublist:%s", err.Error())
	}

	for i := 0; i < len(list); i++ {
		id, err := strconv.ParseInt(list[i], 10, 64)
		if err != nil {
			log.Error("Error Parsing id:%s", list[i])
			continue
		}

		sub, err := Mo.Redis.SubscriberGetByID(EntCode, id)
		subs = append(subs, sub)
	}
	return subs[begin:limit], nil
}

func InitTestData() {
	err, subs := FetchTestSub(0, 5000)
	if err != nil {
		log.Error("Error Fetching Test Subscriber:%s", err.Error())
	}
	subids := make([]int64, 0)

	for _, sub := range subs {
		subids = append(subids, sub.ID)
		//@Todo :automatic
		sub.ApplicationCode = GetAppCode(sub.Platform)
		sub.EnterpriseCode = EntCode
		if err := Mo.Redis.SubscriberSet(sub, -1); err != nil {
			log.Error("Error CaChing Subscribers:%s", err.Error())
		}
	}
	AddSubList(subids)
}

func GetAppCode(p o.Platform) string {
	pcmap := map[o.Platform]string{
		o.PlatformQQ:     "qq_yto",
		o.PlatformWechat: "wechat_yto",
		o.PlatformAlipay: "alipay_yto",
		o.PlatformSite:   "site_yto",
	}

	code, ok := pcmap[p]
	if !ok {
		return ""
	}
	return code
}

func GetAppToken(p o.Platform) string {
	ptmap := map[o.Platform]string{
		o.PlatformQQ:     "8d1600d6-f7df-11e4-934b-00262d0d05c7",
		o.PlatformWechat: "c132c165-3c2a-11e5-a8b4-00262d0d05c7",
		o.PlatformAlipay: "cc2f2b03-2e7c-11e5-a8b4-00262d0d05c7",
		o.PlatformSite:   "8e027c42-3c3a-11e5-a8b4-00262d0d05c7",
	}

	token, ok := ptmap[p]
	if !ok {
		return ""
	}
	return token
}

func FetchTestSub(start, end int64) (error, []*o.Subscriber) {
	if start >= end || start < 0 {
		return fmt.Errorf("Invalid Range!"), nil
	}
	var subs []*o.Subscriber

	_, err := Mo.MySQL.Select(&subs, "SELECT * from subscribers limit ?, ? ", start, end)
	if err != nil {
		log.Error("Error Fetching Test Subscriber")
		return err, nil
	}

	for i := 0; i < len(subs); i++ {
		loc, err := Mo.MySQL.GetLocationByID(subs[i].ID)
		if err != nil {
			log.Error("Error Fetching Location:%s", err.Error())
		}
		subs[i].Location = loc
		time.Sleep(time.Microsecond * 100)
		log.Debug("cnt:%d", i)
	}
	return nil, subs
}

func NewSuite(platform o.Platform) (error, *TestSuite) {
	t := new(TestSuite)
	t.Client = gr.New()

	app, err := Mo.ApplicationGetByCode(EntCode, GetAppCode(platform))
	if err != nil {
		log.Error("Error Fetching Application:%s", err.Error())
		return err, nil
	}

	t.Token = GetAppToken(platform)

	t.App = app
	t.Platform = platform
	return nil, t
}
