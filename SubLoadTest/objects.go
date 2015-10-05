package main

import (
	gr "github.com/parnurzeal/gorequest"
	o "yto.net.cn/kefu/models/objects"
)

type TestSuite struct {
	App      *o.Application
	Platform o.Platform
	Client   *gr.SuperAgent
	Token    string
}
