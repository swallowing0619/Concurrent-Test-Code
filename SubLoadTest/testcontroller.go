package main

import (
	"encoding/json"
	"fmt"
	//gr "github.com/parnurzeal/gorequest"
	"time"
	o "yto.net.cn/kefu/models/objects"
)

func (t *TestSuite) TestSessionCreate(sub *o.Subscriber) {
	//@Todo :random nodeCode
	data := &o.ApiSessionCreateData{
		AccessIP: "127.0.0.1",
		Business: "test",
		NodeCode: "999999",
		OpenID:   sub.OpenID,
		Platform: sub.Platform,
	}

	req := &o.ApiRequest{
		Endpoint:       "session",
		EnterpriseCode: EntCode,
		RequestAt:      time.Now(),
		Action:         "session.create",
		Data:           data,
	}

	log.Debug("%s %s", t.Platform, sub.Platform)

	if t.Platform == sub.Platform {
		jstr, err := json.Marshal(req)
		if err != nil {
			log.Error("Error Marshall:%s", err.Error())
			return
		}
		log.Debug("query json :%s", jstr)

		resp, body, errs := t.Client.Post(fmt.Sprintf("%s/api/v1/session?access_token=%s", OpenUrl, t.Token)).Type("json").Send(fmt.Sprintf("%s", jstr)).End()

		if errs != nil {
			log.Error("Error doing session create:%+v", errs)
		}
		log.Debug("Got Resp:%+v,body:%s", resp, body)
	}
}
