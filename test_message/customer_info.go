package main

import (
	"fmt"
)

type CustomerInfo struct {
	Name     string
	Platform string
	Openid   string
}

func (cf *CustomerInfo) GetKey() (error, string) {
	if cf.Platform == "" || cf.Openid == "" {
		return fmt.Errorf("missing platform parameter or openid in customer info"), ""
	}
	return nil, cf.Platform + " - " + cf.Openid
}

func getCustomerInfoList() map[string]CustomerInfo {
	// todo: fetch data from mysql
	customerInfoList := [][]string{
		// []string{"老骥", "owFDXs1RM2Xvow7uok-nONnKVOEg"},
		[]string{"古月林夕", "owFDXszUEtoctweeTHhkR0S5iGqg"},

		/*       []string{"心凝形释", "owFDXs7efk2SNEymr0R5qNrs6dpg"},
		[]string{"xiaoxiao", "owFDXsxggHn7gFMWLHd3Vc1PcqYU"},
		[]string{"小弟王楚", "owFDXs8JTRgNEa49Qr6mLaKj4jjc"},
		[]string{"埋名", "owFDXs2DeqbAHD9GMysGeJGTh-jE"},
		[]string{"別問", "owFDXs6yfgicsSGc0iBpO1_jPVAs"},
		[]string{"徐航", "owFDXs1awOfiUP2gzzYEQDneh2hM"},
		[]string{"夕阳丶", "owFDXs8EABjLdUX72SGRvMIxodz4"},
		[]string{"易显维", "owFDXs7j76N0d-jlb_P4gdLnn-hk"},
		[]string{"gh", "owFDXs9V1KDsytB_88x6gboaxTPY"},
		[]string{"杜润", "owFDXsxpkZJT-rqnx3Y5kP3RCXQE"},
		[]string{"顾瑜", "owFDXs8hj7LpdA02ZukactJCVbwY"},
		[]string{"露丝儿", "owFDXs7Nkxb5gxXr5VHNH_IuXHCo"},
		[]string{"李昌振", "owFDXs4mt_sv5P1M9dkgNpI7Ablo"},
		[]string{"丹丹", "owFDXs5L6JiFK87_nqYzh4PaANPU"},
		[]string{"超人带你飞", "owFDXs8iBlmvW0vhzcRPaCP5rHNY"},*/
		// []string{"杨冬燕", "owFDXs0V2CfaSB8ION2l0o63NhCw"},
		/*     []string{"Edison", "owFDXs9QPZ-BQ8Mdm1MyLOD7x8YM"},
		[]string{"小兮", "owFDXs_-1dMQNH3lr2Y9f46VS9mU"},
		[]string{"妞", "owFDXs3CPbdbPxemZaQHUsOBzWNM"},
		[]string{"堕落天使", "owFDXs4qcvv5OQF_6852v5saE6dM"},
		[]string{"晟林爸爸", "owFDXs8vrXpQDYqpscKvhJ7uEAo8"},
		[]string{"spinach", "owFDXs_gJtgPlA_D86nUonhbxQtE"},
		[]string{"李欢", "owFDXsxW13-JojMAuz34AMXc9UUM"},
		[]string{"段晓", "owFDXs4yLD5Aa7ntHTD_FUx_uMr4"},
		[]string{"kz", "owFDXs2_6f4VhGqNO6D2zAQTGIJY"},
		[]string{"Humbert", "owFDXs-59P4f5WGw9agYDfxLcESI"},
		[]string{"也哉", "owFDXs1eb4JylUmFY04bEFo7I09Y"},
		*/
	}

	customerInfo := make(map[string]CustomerInfo)

	platform := "0"
	var c CustomerInfo
	var err error
	var key string
	for _, v := range customerInfoList {
		c = CustomerInfo{
			Name:     v[0],
			Platform: platform,
			Openid:   v[1],
		}

		if err, key = c.GetKey(); err == nil {
			customerInfo[key] = c
		}
	}
	return customerInfo
}
