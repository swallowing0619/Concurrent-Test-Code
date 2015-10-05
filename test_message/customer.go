package main
import (
  "fmt"
  "time"
)

var customers map[string] *Customer = make( map[string] *Customer )

/*
 * 客户
 */
type Customer struct {
  Requests  []Request            // 请求列表
  Responses  []string             // 返回数据
  Step      int                  // 0
  StartTime time.Time            // 开始时间
  EndTime   time.Time            // 结束时间
  Info      CustomerInfo
}


// 生成customer对象
func makeCustomer(request Request, info CustomerInfo) (*Customer) {
  customer := &Customer{
    Step:   0,
    StartTime: time.Now(),
    Info: info,
  }
  customer.Requests = append(customer.Requests, request)

  return customer
}


func addCustomer(customer *Customer) {
  if err, key := customer.Info.GetKey(); err == nil {
    customers[key] = customer
  }
}

func getCustomer (customer *CustomerInfo) (*Customer, error) {
  if err, key := customer.GetKey(); err != nil {
    return nil, fmt.Errorf("get customer by platform and openid: " + err.Error())
  } else if v, ok := customers[key]; ok {
    return v, nil
  } 

  return  nil, fmt.Errorf("don't exists customer in customers")
}

func getCustomerList() map[string] *Customer {
  return customers
}

